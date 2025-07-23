package index

import (
	"fmt"
	"strings"
	"sync"

	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/utils"
)

// AddDocument 分词索引
func (e *Engine) AddDocument(index *model.IndexDoc) {
	//等待初始化完成
	e.Wait()
	var (
		id         string
		text       string
		splitWords []string
	)

	time1 := utils.ExecTime(func() {
		//是否对文档内容进行索引
		if index.Cut_Document {
			s := utils.Map2String(index.Document, false)
			text = fmt.Sprintf("%s\n%s\n%v\n%s\n%s", index.Title, index.Text, s, index.Tags, index.Flag)
		} else {
			text = index.Text
		}
		splitWords = e.Tokenizer.Cut(text)
		//增加KV关键词
		if index.Keys != nil {
			splitWords = append(splitWords, utils.Map2Strings(index.Keys, true)...)
		}
		//增加SEO关键词
		if index.Seo != nil {
			splitWords = append(splitWords, utils.Map2Strings(index.Seo, false)...)
		}
		splitWords = append(splitWords, index.Id)
		id = index.Id
	})

	time2 := utils.ExecTime(func() {
		e.addTemp(id, splitWords)
	})
	utils.Println("倒排索引", index.Id, index.Title, time2, "ms")

	// TODO: 是否需要更新正排索引 - 检测document变更
	time3 := utils.ExecTime(func() {
		e.SaveDocument(index, splitWords)
	})
	utils.Println("正排索引", index.Id, index.Title, time3, "ms")
	utils.Printc("添加文档完成!!!\nID:%s\n标题:%s\n分词索引:%fms\n关键词:%d\n索引文本：%d\n\n", utils.Green, index.Id, index.Title, time1, len(splitWords), len(index.Text))
}

// GetDocumentCount 获取文档数量
func (e *Engine) GetDocumentCount() int64 {
	if e.DocumentCount == -1 {
		var count int64
		//使用多线程加速统计
		wg := sync.WaitGroup{}
		wg.Add(e.Shard)
		//这里的统计可能会出现数据错误，因为没加锁
		for i := 0; i < e.Shard; i++ {
			go func(i int) {
				count += e.DocStorages[i].GetCount()
				wg.Done()
			}(i)
		}
		wg.Wait()
		e.DocumentCount = count
		return count + 1
	}

	return e.DocumentCount + 1
}
func (e *Engine) SaveDocument(index *model.IndexDoc, keys []string) {
	key := []byte(index.Id)
	shard := e.GetShard(index.Id)
	docStorage := e.DocStorages[shard]
	doc := &model.StorageIndexDoc{
		IndexDoc: index,
		Keys:     keys,
	}
	//存储id和key以及文档的映射
	if !docStorage.Has(key) {
		e.DocumentCount++
	}
	docStorage.Set(key, utils.Encoder(doc))
}

func (e *Engine) Get(item model.SliceItem, request *model.SearchRequest, wordMap *map[string]bool) (model.ResponseDoc, error) {
	var (
		doc model.ResponseDoc
	)
	buf := e.GetDocById(string(item.Id))
	doc.Score = item.Score
	if buf != nil {
		//gob解析
		storageDoc := new(model.StorageIndexDoc)
		utils.Decoder(buf, &storageDoc)
		doc.Document = storageDoc.Document
		doc.Keys = storageDoc.Keys
		text := utils.RemoveHTMLTags(storageDoc.Text)
		title := utils.RemoveHTMLTags(storageDoc.Title)
		//处理关键词高亮
		highlight := request.Highlight
		if highlight != nil {
			//全部小写
			// text = strings.ToLower(text)
			//还可以优化，只替换击中的词
			for _, key := range storageDoc.Keys {
				if ok := (*wordMap)[key]; ok {
					key = utils.RemoveHTMLTags(key)
					text = strings.ReplaceAll(text, key, fmt.Sprintf("%s%s%s", highlight.PreTag, key, highlight.PostTag))
					title = strings.ReplaceAll(title, key, fmt.Sprintf("%s%s%s", highlight.PreTag, key, highlight.PostTag))
				}
			}
			//放置原始文本
			doc.OriginalText = storageDoc.Text
			doc.OriginalTitle = storageDoc.Title
		}
		doc.Text = text
		doc.Tags = storageDoc.Tags
		doc.Time = storageDoc.Time
		doc.Flag = storageDoc.Flag
		doc.Title = title
		doc.Id = string(item.Id)
	}
	return doc, nil
}
func (e *Engine) GetDocument(item model.SliceItem, doc *model.ResponseDoc, request *model.SearchRequest, wordMap *map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	*doc, _ = e.Get(item, request, wordMap)
}

func (e *Engine) IndexDocument(doc *model.IndexDoc) error {
	doc.Time = utils.Now()
	//数量增加
	e.addDocumentWorkerChan[e.GetShard(doc.Id)] <- doc
	return nil
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec(worker chan *model.IndexDoc) {
	for {
		doc := <-worker
		e.AddDocument(doc)
	}
}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id string) []byte {
	shard := e.GetShard(id)
	key := []byte(id)
	buf, found := e.DocStorages[shard].Get(key)
	if found {
		return buf
	}

	return nil
}

// getShardKey 计算索引分布在哪个文件块
func (e *Engine) GetShard(id string) int {
	// return int(id % uint32(e.Shard))
	sid := e.GetShardKeyByWord(id)
	// fmt.Println("sid:", sid)
	return sid
}
