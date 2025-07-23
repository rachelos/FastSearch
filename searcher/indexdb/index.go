package indexdb

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"gitee.com/rachel_os/fastsearch/searcher/arrays"
	"gitee.com/rachel_os/fastsearch/searcher/sorts"
	"gitee.com/rachel_os/fastsearch/searcher/storage"
	"gitee.com/rachel_os/fastsearch/utils"
)

func (e *Engine) FixIndex(id string, words []string, wg *sync.WaitGroup) {
	wg.Add(1)
	var (
		inserts            []string
		needUpdateInverted bool
	)
	// 检查是否需要更新倒排索引 words变更/id不存在
	t1 := utils.ExecTime(func() {
		inserts, needUpdateInverted = e.optimizeIndex(id, words)
	})
	// 将新增的word剔出单独处理，减少I/O操作
	t2 := utils.ExecTime(func() {
		if needUpdateInverted {
			for _, word := range inserts {
				e.addInvertedIndex(word, id)
			}
		}
	})
	t3 := utils.ExecTime(func() {
		e.addPositiveIndex(id, words)

	})
	utils.Printc("倒排索引 ID:%s,检查:%fms(%d个),新增%fms(%d个),正排序索引：%fms\n", utils.Green, id, t1, len(words), t2, len(inserts), t3)
	wg.Done()
}

// 添加倒排索引
func (e *Engine) addInvertedIndex(word string, id string) {
	shard := e.GetShardKeyByWord(word)
	s := e.InvertedIndexStorages[shard]
	//string作为key
	key := []byte(word)

	//存在
	//添加到列表
	buf, find := s.Get(key)
	ids := make([]string, 0)
	if find {
		utils.Decoder(buf, &ids)
	}

	utils.ExecTime(func() {
		if !arrays.Search(ids, id) {
			ids = append(ids, id)
		}
	})
	utils.ExecTime(func() {
		sort.Strings(ids)
		s.Set(key, utils.Encoder(ids))
	})
	// utils.Printc("关键词:%s,ID:%s,关键ID数量:%d,检测耗时：%fms,排序耗时:%fms\n", utils.Green, word, id, len(ids), t1, t2)
}

// 添加到临时索引仓
func (e *Engine) addTemp(id string, words []string) {
	e.tempStorage.Set([]byte(id), utils.Encoder(words))
}
func (e *Engine) TempCount() int64 {
	return e.tempStorage.GetCount()
}

// 处理索引数据
func (e *Engine) IndexTask() {
	wg := sync.WaitGroup{}
	// 获取临时索引仓数据
	tempStorage := e.tempStorage
	tempStorage.GetAll(0, 100, "asc", func(i storage.Item) (bool, storage.Item) {
		var words []string
		t := utils.ExecTime(func() {
			words = make([]string, 0)
			utils.Decoder(i.Value, &words)
			e.FixIndex(string(i.Key), ([]string)(words), &wg)
			tempStorage.Delete(i.Key)
			wg.Wait()
		})
		color := utils.Green
		if t > 2000 {
			color = utils.Red
		}
		utils.Printc("索引完成 ID:%s,关键字:%d:耗时：%fms，剩余索引:%d\n", color, (i.Key), len(words), t, tempStorage.GetCount())
		return true, i
	})

	time.AfterFunc(1*time.Second, e.IndexTask)
}
func (e *Engine) DoIndexTask() {
	e.IndexTask()
}

// 移除删去的词
func (e *Engine) optimizeIndex(id string, newWords []string) ([]string, bool) {
	// 判断id是否存在
	// e.Lock()
	// defer e.Unlock()

	// 计算差值
	removes, inserts, changed := e.GetDifference(id, newWords)
	if changed {
		if removes != nil && len(removes) > 0 {
			// 移除正排索引
			for _, word := range removes {
				e.RemoveIdInWordIndex(id, word)
			}
		}
	}
	return inserts, changed
}

func (e *Engine) RemoveIdInWordIndex(id string, word string) {
	// e.Lock()
	// defer e.Unlock()
	shard := e.GetShardKeyByWord(word)

	wordStorage := e.InvertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	buf, found := wordStorage.Get(key)
	if found {
		// ids := make([]uint32, 0)
		ids := make([]string, 0)
		utils.Decoder(buf, &ids)

		//移除
		index := arrays.FindString(ids, id)
		if index != -1 {
			ids = utils.DeleteStringArray(ids, index)
			if len(ids) == 0 {
				err := wordStorage.Delete(key)
				if err != nil {
					utils.Printc("删除索引词出错：%s\n", utils.Red, err.Error())
				}
			} else {
				wordStorage.Set(key, utils.Encoder(ids))
			}
		}
	}

}

// 计算差值
// @return []string: 需要删除的词
// @return bool    : words出现变更返回true，否则返回false
func (e *Engine) GetDifference(id string, newWords []string) ([]string, []string, bool) {
	shard := e.GetShard(id)
	wordStorage := e.PositiveIndexStorages[shard]
	key := []byte(id)
	buf, found := wordStorage.Get(key)
	if found {
		oldWords := make([]string, 0)
		utils.Decoder(buf, &oldWords)

		// 计算需要移除的
		removes := make([]string, 0)
		for _, word := range oldWords {
			// 旧的在新的里面不存在，就是需要移除的
			if !arrays.Search(newWords, word) {
				removes = append(removes, word)
			}
		}
		// 计算需要新增的
		inserts := make([]string, 0)
		for _, word := range newWords {
			if !arrays.Search(oldWords, word) {
				inserts = append(inserts, word)
			}
		}
		if len(removes) != 0 || len(inserts) != 0 {
			return removes, inserts, true
		}
		// 没有改变
		return removes, inserts, false
	}
	// id不存在，相当于insert
	return nil, newWords, true
}

// 添加正排索引 id=>keys id=>doc
func (e *Engine) addPositiveIndex(id string, keys []string) {
	// e.Lock()
	// defer e.Unlock()
	key := []byte(id)
	//设置到id和key的映射中
	shard := e.GetShard(id)
	//id和key的映射
	positiveIndexStorage := e.PositiveIndexStorages[shard]
	positiveIndexStorage.Set(key, utils.Encoder(keys))

}

// processKeySearch 实现了Engine结构体中处理关键字搜索的功能
//
// 参数：
// word string - 待搜索的关键字
// fastSort *sorts.FastSort - 用于排序搜索结果的快速排序对象
// wg *sync.WaitGroup - 等待组对象，用于等待该协程执行完毕
//
// 返回值：
// 无返回值
func (e *Engine) ProcessKeySearch(word string, fastSort *sorts.FastSort, wg *sync.WaitGroup) {
	defer wg.Done()

	shard := e.GetShardKeyByWord(word)
	//读取id
	invertedIndexStorage := e.InvertedIndexStorages[shard]
	key := []byte(word)

	buf, find := invertedIndexStorage.Get(key)
	if find {
		ids := make([]string, 0)
		//解码
		utils.Decoder(buf, &ids)
		fastSort.Add(&ids)
	}

}

// RemoveIndex 根据ID移除索引
func (e *Engine) RemoveIndex(id string) error {
	//移除
	// e.Lock()
	// defer e.Unlock()

	shard := e.GetShard(id)
	// key := utils.Uint32ToBytes(id)
	key := []byte(id)
	//关键字和Id映射
	//invertedIndexStorages []*storage.LeveldbStorage
	//ID和key映射，用于计算相关度，一个id 对应多个key
	ik := e.PositiveIndexStorages[shard]
	keysValue, found := ik.Get(key)
	if !found {
		return errors.New(fmt.Sprintf("没有找到id=%s", id))
	}

	keys := make([]string, 0)
	utils.Decoder(keysValue, &keys)

	//符合条件的key，要移除id
	for _, word := range keys {
		e.RemoveIdInWordIndex(id, word)
	}

	//删除id映射
	err := ik.Delete(key)
	if err != nil {
		return errors.New(err.Error())
	}

	//文档仓
	err = e.DocStorages[shard].Delete(key)
	if err != nil {
		return err
	}
	//减少数量
	e.DocumentCount--

	return nil
}

func (e *Engine) Close() {
	// e.Lock()
	// defer e.Unlock()
	for i := 0; i < e.Shard; i++ {
		e.InvertedIndexStorages[i].Close()
		e.PositiveIndexStorages[i].Close()
	}
	e.tempStorage.Close()
}

// Drop 删除
func (e *Engine) Drop() error {
	if !e.AllowDrop {
		return errors.New("不允许删除")
	}
	e.Lock()
	defer e.Unlock()

	e.Close()
	//删除文件
	if err := os.RemoveAll(e.IndexPath); err != nil {
		return err
	}
	e.DocStorages = make([]*storage.LeveldbStorage, 0)
	e.InvertedIndexStorages = make([]*storage.LeveldbStorage, 0)
	e.PositiveIndexStorages = make([]*storage.LeveldbStorage, 0)
	runtime.GC()
	return nil
}
