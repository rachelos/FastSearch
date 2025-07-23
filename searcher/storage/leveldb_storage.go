package storage

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"gitee.com/rachel_os/fastsearch/leveldb"
	"gitee.com/rachel_os/fastsearch/leveldb/filter"
	"gitee.com/rachel_os/fastsearch/leveldb/iterator"
	"gitee.com/rachel_os/fastsearch/leveldb/opt"
	"gitee.com/rachel_os/fastsearch/searcher/exp"
	"gitee.com/rachel_os/fastsearch/utils"
)

// LeveldbStorage TODO 要支持事务
type LeveldbStorage struct {
	db          *leveldb.DB
	path        string
	mu          sync.RWMutex //加锁
	cmu         sync.RWMutex //加锁
	closed      bool
	cache       map[string][]byte
	CacheEnable bool
	timeout     int64
	retry_index int64
	lastTime    int64
	count       int64
}
type Item struct {
	Key   []byte
	Value []byte
	Index int64
	Score int64
}
type Items []Item

func (a Items) Len() int           { return len(a) }
func (a Items) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Items) Less(i, j int) bool { return a[i].Score < a[j].Score }

func (s *LeveldbStorage) NewItem(key []byte, value []byte) *Item {
	return &Item{Key: key, Value: value}
}

func (s *LeveldbStorage) AutoId() int64 {
	return s.count + 1
}
func (s *LeveldbStorage) autoOpenDB() {
	if s.isClosed() {
		s.ReOpen()
	}
	s.lastTime = time.Now().Unix()
}

// NewStorage 打开数据库
func NewStorage(path string, cache bool, timeout int64) (*LeveldbStorage, error) {
	db := &LeveldbStorage{
		path:        path,
		closed:      true,
		timeout:     timeout,
		CacheEnable: cache,
		cache:       make(map[string][]byte, 0),
		lastTime:    time.Now().Unix(),
	}
	if timeout > 0 {
		go db.task()
	}
	//计算总条数
	go db.compute()
	return db, nil
}

func (s *LeveldbStorage) task() {
	if s.timeout == -1 {
		//不检查
		return
	}
	for {
		if !s.isClosed() && time.Now().Unix()-s.lastTime > s.timeout {
			s.Close()
			//log.Println("leveldb storage timeout", s.path)
		}
		time.Sleep(time.Duration(5) * time.Second)

	}
}
func (s *LeveldbStorage) getCache(key []byte) ([]byte, bool) {
	if !s.CacheEnable {
		return nil, false
	}
	s.cmu.RLock()
	defer s.cmu.RUnlock()
	value, ok := s.cache[string(key)]
	return value, ok
}

func (s *LeveldbStorage) setCache(key []byte, value []byte) {
	if !s.CacheEnable {
		return
	}
	s.cmu.Lock()
	defer s.cmu.Unlock()
	s.cache[string(key)] = value
}
func (s *LeveldbStorage) delCache(key []byte) {
	if !s.CacheEnable {
		return
	}
	s.cmu.Lock()
	defer s.cmu.Unlock()
	delete(s.cache, string(key))
}

func openDB(path string) (*leveldb.DB, error) {

	////使用布隆过滤器
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(path, o)

	return db, err
}
func (s *LeveldbStorage) ReOpen() error {
	if !s.isClosed() {
		log.Println("db is not closed")
		return nil
	}
	db, err := openDB(s.path)
	if err != nil {
		utils.Printc("数据库还未初始化,PATH:%s,Error:%s\n\n", utils.Red, s.path, err.Error())
		// panic(err)
		return errors.New("数据库还未初始化")
	}
	s.db = db
	s.closed = false
	utils.Printc("open db:%s\n", utils.Green, s.path)
	return nil
}

func (s *LeveldbStorage) NGet(key []byte, cache bool) ([]byte, bool) {
	s.autoOpenDB()
	s.mu.RLock()
	defer s.mu.RUnlock()
	if cache {
		if value, ok := s.getCache(key); ok {
			return value, true
		}
	}

	buffer, err := s.db.Get(key, nil)

	if err != nil {
		return nil, false
	}
	s.setCache(key, buffer)
	return buffer, true
}
func (s *LeveldbStorage) Get(key []byte) ([]byte, bool) {
	return s.NGet(key, true)
}
func (s *LeveldbStorage) DB() *leveldb.DB {
	return s.db
}

func (s *LeveldbStorage) GetAll(start int64, end int64, order string, filter func(Item) (bool, Item)) (int64, Items) {
	return s.GetList(start, end, order, filter, false)
}
func (s *LeveldbStorage) GetList(start int64, end int64, order string, filter func(Item) (bool, Item), is_sanp bool) (int64, Items) {
	s.autoOpenDB()
	var iter iterator.Iterator
	if end > s.count || end == 0 {
		end = s.count
	}
	count := s.count
	if is_sanp {
		snap, _ := s.db.GetSnapshot()
		defer snap.Release()
		iter = snap.NewIterator(nil, nil)
	} else {
		iter = s.db.NewIterator(nil, nil)
	}
	temps := make(Items, 0)
	found := int64(0)
	for iter.Next() {
		found++
		if found <= start {
			continue
		}
		if found > end {
			break
		}
		key := iter.Key()
		value, _ := s.Get(key)
		item := Item{
			Key:   key,
			Value: value,
			Index: found,
			Score: found,
		}
		if filter != nil {
			val, item := filter(item)
			if val {
				temps = append(temps, item)
			} else {
				found--
				count--
			}
		} else {
			temps = append(temps, item)
		}

	}
	sort.Sort(Items(temps))
	if strings.ToLower(order) == "desc" {
		sort.Sort(sort.Reverse(Items(temps)))
	}
	iter.Release()
	page_size := int(end - start)
	if page_size == 0 {
		return count, temps
	}
	cur_page := int(math.Ceil(float64(int(start)/int(page_size)))) + 1
	if count < s.count {
		c_count := int64(page_size * cur_page)
		if len(temps) == page_size {
			count = c_count + 1
		} else {
			c_count = int64(int(c_count) - (page_size - int(count)))
		}
	}
	return count, temps
}

type SearchOption struct {
	FilterExp string
	ScoreExp  string
	Start     int64
	End       int64
	Order     string
}

// callbacl传入要转换的对你
func (s *LeveldbStorage) Search(callback func(Item) (v any), option *SearchOption) (int64, Items) {
	return s.GetAll(option.Start, option.End, option.Order, func(item Item) (bool, Item) {
		var (
			filter_exp_str = option.FilterExp
			score_exp_str  = option.ScoreExp
		)

		if filter_exp_str == "" {
			return true, item
		}
		obj := callback(item)
		if obj == nil {
			return true, item
		}
		params := Obj2Map(obj)
		filterExp, _ := exp.NewEvaluableExpression(filter_exp_str)
		val, _ := filterExp.Evaluate((params))

		if score_exp_str != "" {
			scoreExp, _ := exp.NewEvaluableExpression(score_exp_str)
			score, _ := scoreExp.Evaluate((params))
			item.Score = int64(score.(int))
		}
		return val.(bool), item
	})
}

func (s *LeveldbStorage) Has(key []byte) bool {
	s.autoOpenDB()
	has, err := s.db.Has(key, nil)
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		return false
	}
	return has
}

func (s *LeveldbStorage) Set(key []byte, value []byte) error {
	s.autoOpenDB()
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.Has(key) {
		s.count++
	} else {
		// s.db.Delete(key, nil)
	}
	err := s.db.Put(key, value, nil)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		return err
	}
	// 更新缓存
	s.setCache(key, value)
	return nil
}

func (s *LeveldbStorage) BatchSet(kv *Items) ([]string, error) {
	s.autoOpenDB()
	s.mu.Lock()
	var (
		batch *leveldb.Batch
		ids   []string
	)
	batch = new(leveldb.Batch)
	for _, item := range *kv {
		if !s.Has(item.Key) {
			s.count++
		}
		batch.Put(item.Key, item.Value)
		ids = append(ids, string(item.Key))
		s.setCache(item.Key, item.Value)
	}
	err := s.db.Write(batch, nil)
	defer s.mu.Unlock()
	if err != nil {
		// panic(err)
		fmt.Println(err.Error())
		return ids, err
	}
	return ids, nil
}

// Delete 删除
func (s *LeveldbStorage) Delete(key []byte) error {
	s.autoOpenDB()
	err := s.db.Delete(key, nil)
	if err != nil {
		return err
	}
	s.delCache(key)
	s.count--
	return nil
}

// Close 关闭
func (s *LeveldbStorage) Close() error {
	if s.isClosed() {
		return nil
	}
	s.mu.Lock()
	err := s.db.Close()
	if err != nil {
		return err
	}
	s.closed = true
	s.mu.Unlock()
	return nil
}

func (s *LeveldbStorage) isClosed() bool {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	if s.db == nil {
		return true
	}
	return s.closed
}

func (s *LeveldbStorage) compute() {
	var count int64
	s.autoOpenDB()
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		count++
	}
	iter.Release()
	s.count = count
}

func (s *LeveldbStorage) GetCount() int64 {
	return s.count
}
