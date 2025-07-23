package filecache

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// FileCache 是一个用于缓存文件内容的结构体
type FileCache struct {
	sync.RWMutex
	Option
	cache map[string]*os.File
	data  map[string]map[string]interface{}
}
type Option struct {
	Timeout   int64
	UseCache  bool
	Path      string
	Extension string
}

// NewFileCache 创建一个新的 FileCache 实例
func NewStorage(path string, c bool, timeout int64) (*FileCache, error) {
	cache := &FileCache{
		data:  make(map[string]map[string]interface{}, 0),
		cache: make(map[string]*os.File),
	}
	opt := &Option{Path: path, Timeout: timeout, UseCache: c}
	if opt.Path == "" {
		pwd, _ := os.Getwd()
		cache.Path = pwd + "\\data\\"
	} else {
		cache.Path = opt.Path
	}
	if opt.Timeout >= 0 {
		cache.Timeout = opt.Timeout
	}
	cache.UseCache = opt.UseCache
	cache.Extension = ".txt"
	return cache, nil
}

// GetFile 从缓存中获取文件内容，如果不存在，则读取文件并缓存
func (e *FileCache) Add(key string, data string) error {

	if e.Has(key, data) {
		return errors.New(fmt.Sprintf("%s already exists：%s", key, data))
	}
	e.Lock()
	defer e.Unlock()
	if e.data[key] == nil {
		e.data[key] = make(map[string]interface{}, 0)
	}
	e.data[key][data] = ""
	err := e.write(key, data)
	return err
}
func (e *FileCache) Get(key string) (v map[string]interface{}, found bool) {
	e.RLock()
	defer e.RUnlock()
	v, err := e.read(key)
	if err != nil {
		return nil, false
	}
	return v, true
}
func (e *FileCache) Delete(key string, v string) error {
	delete(e.data[key], v)
	err := e.write(key, v)
	return err
}
func (e *FileCache) Set(key string, values map[string]any) error {
	e.Lock()
	defer e.Unlock()
	for k, _ := range values {
		e.Add(key, k)
	}
	return nil
}
func (e *FileCache) Has(key string, v string) bool {
	arr, err := e.read(key)
	if err == nil {
		_, ok := arr[v]
		return ok
	}
	return false
}
func (e *FileCache) Count(key string) int64 {
	return int64(len(e.data[key]))
}
func (e *FileCache) GetCount() int64 {
	return int64(len(e.data))
}

func (e *FileCache) Data() map[string]map[string]interface{} {
	return e.data
}

func Dir_Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}
func (e *FileCache) open(key string) *os.File {
	path := e.file_path(key)
	file, ok := e.cache[key]
	if ok {
		return file
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0655)
	// file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	e.cache[key] = file
	if e.Timeout > 0 {
		time.AfterFunc(time.Duration(e.Timeout), func() {
			// e.close(key)
		})
	}
	return file
}
func (e *FileCache) close(key string) {
	delete(e.cache, key)
	e.cache[key].Close()
}
func (e *FileCache) write(key string, newLine string) error {
	var (
		err error
	)
	file := e.open(key)
	// 移动到文件末尾
	writer := bufio.NewWriter(file)
	if _, err := file.Seek(0, os.SEEK_END); err != nil {
		return err
	}
	// 写入新内容
	_, err = fmt.Fprintln(writer, newLine)
	if err != nil {
		return err
	}
	return writer.Flush()
}
func (e *FileCache) read(key string) (map[string]interface{}, error) {
	arr := make(map[string]interface{}, 0)
	// file := e.open(key)
	file, err := os.Open(e.file_path(key))
	if err != nil {
		return arr, errors.New(fmt.Sprintf("Error reading file:%s", err))
	}
	defer file.Close()
	// 使用bufio.NewScanner创建一个新的Scanner
	scanner := bufio.NewScanner(file)

	// 逐行读取
	for scanner.Scan() {
		id := scanner.Text()
		arr[id] = ""
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		return arr, errors.New(fmt.Sprintf("Error reading file:%s", err))
	}
	return arr, nil
}
