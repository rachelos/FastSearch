package benchmark

import (
	"math/rand"

	"gitee.com/rachel_os/fastsearch/utils"
)

func GetRandomUint32(n int) []uint32 {
	var array = make([]uint32, n)
	for i := 0; i < n; i++ {
		array[i] = rand.Uint32()
	}
	return array
}
func GetRandomString(n int) []string {
	var array = make([]string, n)
	for i := 0; i < n; i++ {
		array[i], _ = utils.GenerateRandomString(16)
	}
	return array
}
