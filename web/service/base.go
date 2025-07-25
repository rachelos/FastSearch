package service

import (
	"os"
	"runtime"

	"gitee.com/rachel_os/fastsearch/global"
	"gitee.com/rachel_os/fastsearch/searcher"
	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/system"
)

// Base 基础管理
type Base struct {
	Container *searcher.Container
	Callback  func() map[string]interface{}
}

func NewBase() *Base {
	return &Base{
		Container: global.Container,
		Callback:  Callback,
	}
}

// Query 查询
func (b *Base) Query(request *model.SearchRequest) (any, error) {
	return b.Container.GetDataBase(request.Database).MultiSearch(request)
}

// GC 释放GC
func (b *Base) GC() {
	runtime.GC()
}

// Status 获取服务器状态
func (b *Base) Status() map[string]interface{} {
	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)

	s := b.Callback()

	r := map[string]interface{}{
		"memory":  system.GetMemStat(),
		"cpu":     system.GetCPUStatus(),
		"disk":    system.GetDiskStat(),
		"system":  s,
		"version": global.VERSION,
	}
	return r
}

// Restart 重启服务
func (b *Base) Restart() {
	// TODD 未实现
	os.Exit(0)
}
