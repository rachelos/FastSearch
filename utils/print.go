package utils

import (
	"fmt"
	"time"
)

var DEBUG = true

// 定义一些ANSI转义码的颜色常量
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

func Println(a ...interface{}) {
	if DEBUG {
		t := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s]", t)
		fmt.Println(a...)
	}
}
func Printf(format string, a ...interface{}) {
	if DEBUG {
		t := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s]", t)
		fmt.Printf(format, a...)
	}
}
func Printc(format string, color string, a ...interface{}) {
	if DEBUG {
		format = color + format + Reset
		t := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s]", t)
		fmt.Printf(format, a...)
	}
}
