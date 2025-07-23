package arrays

import (
	"sort"
)

func Search(str_array []string, target string) bool {
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func SearchBySort(str_array []string, target string) bool {
	sort.Strings(str_array)
	return Search(str_array, target)
}
