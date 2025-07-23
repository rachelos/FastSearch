package benchmark

import (
	"testing"

	"gitee.com/rachel_os/fastsearch/searcher/arrays"
)

func Benchmark(b *testing.B) {

	//测试两种方法的性能
	size := 10000
	arrayList := make([][]uint32, size)
	for i := 0; i < size; i++ {
		arrayList[i] = GetRandomUint32(1000)
	}

	b.Run("array", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var temp []uint32
			for _, nums := range arrayList {

				for _, num := range nums {
					if !arrays.BinarySearch(temp, num) {
						temp = append(temp, num)
					}
				}
			}
		}
	})

	stringList := make([][]string, size)
	for i := 0; i < size; i++ {
		stringList[i] = GetRandomString(1000)
	}
	b.Run("sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var temp []uint32
			for _, v := range arrayList {
				temp = append(temp, v...)
			}
			//去重
			var as []uint32
			for _, v := range temp {
				if !arrays.BinarySearch(as, v) {
					as = append(as, v)
				}
			}
		}
	})

	b.Run("array string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var temp []string
			for _, nums := range stringList {

				for _, num := range nums {
					if !arrays.BinarySearchString(temp, num) {
						temp = append(temp, num)
					}
				}
			}
		}
	})

	b.Run("sort string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var temp []string
			for _, v := range stringList {
				temp = append(temp, v...)
			}
			//去重
			var as []string
			for _, v := range temp {
				if !arrays.BinarySearchString(as, v) {
					as = append(as, v)
				}
			}
		}
	})
}
