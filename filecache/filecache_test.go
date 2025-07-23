package filecache

import (
	"fmt"
	"os"
	"testing"

	"gitee.com/rachel_os/fastsearch/utils"
)

var cache, err = NewStorage("./tmp/", true, 0)

func TestRead(t *testing.T) {
	// 示例使用 FileCache

	t2 := utils.ExecTime(func() {
		for k, a := range cache.Data() {
			fmt.Println(k, cache.Count(k), a)
		}
	})
	t1 := utils.ExecTime(func() {
		cache.Get("b")
		fmt.Println(cache.Count("b"))
	})
	fmt.Println("read1", t1, "read all", t2)
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
}

func TestAdd(t *testing.T) {
	cache.Add("a", "hello")
	for i := 0; i < 100000; i++ {
		s, _ := utils.GenerateRandomString(10)
		err := cache.Add("b", fmt.Sprintf("v-%v%v", s, i))
		if err != nil {
			fmt.Println("add err", err)
		}
	}
	cache.Add("c", "hello")
}
func BenchmarkFile(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Add("a", fmt.Sprintf("%v%v", "hello-", i))
		}
	})
	b.Run("read", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Get("a")
		}
	})
}
