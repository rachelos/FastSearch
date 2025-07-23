package btree

import (
	"fmt"
	"testing"
	"time"
)

var start = time.Now()

func logTime(flag string, v ...any) {
	fmt.Printf("[%s]-%s->%v %s\n", time.Now().Format("2006-01-02 15:04:05"), flag, time.Since(start), v)
	start = time.Now()
}
func Test_Testing(*testing.T) {
	index := NewDataIndex()

	// ... 插入更多数据
	var max int = 10000000
	for io := 0; io < max; io++ {
		index.Insert(fmt.Sprintf("key-%d", io), DataItem(fmt.Sprintf("value-%d", io)))
	}
	// 插入数据
	index.Insert("apple", "iPhone")
	index.Insert("google", "Pixel")
	logTime("start")
	// 查找数据
	value, found := index.Get("apple")
	if found {
		fmt.Printf("Found: %s -> %s\n", "apple", value)
	} else {
		fmt.Println("Not found")
	}
	logTime("end")
}
