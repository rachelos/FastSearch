package index

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"gitee.com/rachel_os/fastsearch/filecache"
	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/storage"
	"gitee.com/rachel_os/fastsearch/searcher/words"
	"gitee.com/rachel_os/fastsearch/utils"
	// "github.com/Knetic/govaluate"
)

type Engine struct {
	IndexPath string  //索引文件存储目录
	Option    *Option //配置

	InvertedIndexStorages []*filecache.FileCache    //关键字和Id映射，倒排索引,key=id,value=[]words
	PositiveIndexStorages []*storage.LeveldbStorage //ID和key映射，用于计算相关度，一个id 对应多个key，正排索引
	DocStorages           []*storage.LeveldbStorage //文档仓
	tempStorage           *storage.LeveldbStorage   //临时索引仓
	sync.Mutex                                      //锁
	sync.WaitGroup                                  //等待
	addDocumentWorkerChan []chan *model.IndexDoc    //添加索引的通道
	IsDebug               bool                      //是否调试模式
	AllowDrop             bool                      //允许删除
	Tokenizer             *words.Tokenizer          //分词器
	DatabaseName          string                    //数据库名

	Shard         int   //分片数
	Timeout       int64 //超时时间,单位秒
	BufferNum     int   //分片缓冲数
	DocumentCount int64 //文档总数量
}
type Option struct {
	InvertedIndexName string //倒排索引
	PositiveIndexName string //正排索引
	DocIndexName      string //文档存储
}

// Init 初始化索引引擎
func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()

	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	if e.Timeout == 0 {
		e.Timeout = 10 * 3 // 默认30s
	}
	//-1代表没有初始化
	e.DocumentCount = -1
	//log.Println("数据存储目录：", e.IndexPath)
	log.Println("chain num:", e.Shard*e.BufferNum)
	e.addDocumentWorkerChan = make([]chan *model.IndexDoc, e.Shard)
	//初始化临时索引仓
	e.tempStorage, _ = storage.NewStorage(e.GetFilePath(fmt.Sprintf("temp_%s", e.Option.DocIndexName)), true, e.Timeout)
	go e.DoIndexTask()
	//初始化文件存储
	for shard := 0; shard < e.Shard; shard++ {

		//初始化chan
		worker := make(chan *model.IndexDoc, e.BufferNum)
		// worker := make(chan *model.IndexDoc)
		e.addDocumentWorkerChan[shard] = worker

		//初始化chan
		go e.DocumentWorkerExec(worker)

		s, err := storage.NewStorage(e.GetFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)), true, e.Timeout)

		if err != nil {
			panic(err)
		}
		e.DocStorages = append(e.DocStorages, s)

		//初始化Keys存储
		ks, kerr := filecache.NewStorage(e.GetFilePath(fmt.Sprintf("%s_%d/", e.Option.InvertedIndexName, shard)), true, e.Timeout)
		if kerr != nil {
			panic(err)
		}
		e.InvertedIndexStorages = append(e.InvertedIndexStorages, ks)

		//id和keys映射
		iks, ikerr := storage.NewStorage(e.GetFilePath(fmt.Sprintf("%s_%d", e.Option.PositiveIndexName, shard)), true, e.Timeout)
		if ikerr != nil {
			panic(ikerr)
		}
		e.PositiveIndexStorages = append(e.PositiveIndexStorages, iks)
	}

	go e.automaticGC()
	//log.Println("初始化完成")
}

// 自动保存索引，10秒钟检测一次
func (e *Engine) automaticGC() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		//定时GC
		runtime.GC()
		if e.IsDebug {
			log.Println("waiting:", e.GetQueue())
		}
	}
}

func (e *Engine) GetShardKeyByWord(word string) int {
	return int(utils.StringToInt(word) % uint32(e.Shard))
}

func (e *Engine) InitOption(option *Option) {

	if option == nil {
		//默认值
		option = e.GetOptions()
	}
	e.Option = option
	//shard默认值
	if e.Shard <= 0 {
		e.Shard = 10
	}
	if e.BufferNum <= 0 {
		e.BufferNum = 1000
	}
	//初始化其他的
	e.Init()
}

func (e *Engine) GetFilePath(fileName string) string {
	return e.IndexPath + string(os.PathSeparator) + fileName
}

func (e *Engine) GetOptions() *Option {
	return &Option{
		DocIndexName:      "docs",
		InvertedIndexName: "inverted_index",
		PositiveIndexName: "positive_index",
	}
}

// GetQueue 获取队列剩余
func (e *Engine) GetQueue() int {
	total := 0
	for _, v := range e.addDocumentWorkerChan {
		total += len(v)
	}
	return total
}

// GetIndexCount 获取索引数量
func (e *Engine) GetIndexCount() int64 {
	var size int64
	for i := 0; i < e.Shard; i++ {
		size += e.InvertedIndexStorages[i].GetCount()
	}
	return size
}
