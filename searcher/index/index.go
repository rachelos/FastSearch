package index

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"gitee.com/rachel_os/fastsearch/searcher/sorts"
	"gitee.com/rachel_os/fastsearch/searcher/storage"
	"gitee.com/rachel_os/fastsearch/utils"
)

func (e *Engine) FixIndex(id string, words []string, wg *sync.WaitGroup) {
	wg.Add(1)
	// 将新增的word剔出单独处理，减少I/O操作
	t2 := utils.ExecTime(func() {
		_wg := sync.WaitGroup{}
		for _, word := range words {
			_wg.Add(1)
			go e.addInvertedIndex(word, id, &_wg)
		}
		_wg.Wait()
	})
	t3 := utils.ExecTime(func() {
		e.addPositiveIndex(id, words)

	})
	utils.Printc("倒排索引 ID:%s,检查:(%d个),新增%fms(%d个),正排序索引：%fms\n", utils.Green, id, len(words), len(words), t2, t3)
	wg.Done()
}

// 添加倒排索引
func (e *Engine) addInvertedIndex(word string, id string, wg *sync.WaitGroup) {
	shard := e.GetShardKeyByWord(word)
	s := e.InvertedIndexStorages[shard]
	//string作为key
	key := word

	s.Add(key, id)
	wg.Done()
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

func (e *Engine) RemoveIdInWordIndex(id string, word string) {
	// e.Lock()
	// defer e.Unlock()
	shard := e.GetShardKeyByWord(word)

	wordStorage := e.InvertedIndexStorages[shard]

	//string作为key
	key := string(word)
	wordStorage.Delete(key, id)

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
	// key := []byte(word)

	id, find := invertedIndexStorage.Get(word)
	ids := utils.Map2Array(id)
	if find {
		//解码
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
		// e.InvertedIndexStorages[i].Close()
		// e.PositiveIndexStorages[i].Close()
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
	runtime.GC()
	return nil
}
