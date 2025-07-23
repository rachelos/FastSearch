package dbcache

import (
	"testing"
	"time"

	"gitee.com/rachel_os/fastsearch/utils"
)

func TestDbCache(t *testing.T) {
	// username:password@tcp(localhost:3306)/dbname
	db := NewStorage("mysql", "root:123654@tcp(127.0.0.1:3306)/test")
	defer db.Close()
	for i := 0; i < 100; i++ {
		k, _ := utils.GenerateRandomString(10)
		// db.Add(k, k)
		db.Add("a", k)
	}
	db.Add("a", "a")
	time.Sleep(time.Second * 10) // 等待1秒
}
