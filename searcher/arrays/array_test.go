package arrays

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

func Test_Testing(t *testing.T) {
	var m []string
	var max int = 10000000
	for io := 0; io < max; io++ {
		m = append(m, fmt.Sprintf("sn-%d", io))
	}
	logTime("Init", max)
	m = append(m, fmt.Sprintf("sn-%d", max/2))
	m = append(m, fmt.Sprintf("sn-%d", max/4))
	found := fmt.Sprintf("sn-%d", max-1)

	logTime("Search")
	pos := Search(m, found)
	if pos {
		fmt.Println("找到1", pos)
	}
	logTime("Search")
}
