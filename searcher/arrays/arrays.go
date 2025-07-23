package arrays

const (
	LOW  = 0
	HIGH = 1
)

// BinarySearch 二分查找
func BinarySearch(arr []uint32, target uint32) bool {
	low := 0
	high := len(arr) - 1
	for low < high {
		mid := (low + high) >> 1
		if arr[mid] >= target {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return arr != nil && arr[low] == target
}

// BinarySearch 在已排序的字符串切片中查找目标字符串
// 如果找到，返回其索引；否则返回 -1
func BinarySearchString(slice []string, target string) bool {
	left, right := 0, len(slice)-1
	for left <= right {
		mid := left + (right-left)/2
		if slice[mid] == target {
			return true // 找到目标，返回索引
		} else if slice[mid] < target {
			left = mid + 1 // 目标在右半部分
		} else {
			right = mid - 1 // 目标在左半部分
		}
	}
	return false // 没有找到目标
}

// 去重
func RemoveDuplicate(arr []string) []string {
	resArr := make([]string, 0)
	tmpMap := make(map[string]interface{})
	for _, val := range arr {
		if _, ok := tmpMap[val]; !ok {
			resArr = append(resArr, val)
			tmpMap[val] = struct{}{}
		}
	}
	return resArr
}

func ArrayUint32Exists(arr []uint32, target uint32) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
func ArrayUint64Exists(arr []uint64, target uint64) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func ArrayStringExists(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// MergeArrayUint32 合并两个数组
func MergeArrayUint32(target []uint32, source []uint32) []uint32 {

	for _, val := range source {
		if !BinarySearch(target, val) {
			target = append(target, val)
		}
	}
	return target
}

func Find(arr []uint32, target uint32) int {
	for index, v := range arr {
		if v == target {
			return index
		}
	}
	return -1
}
func Find64(arr []uint64, target uint64) int {
	for index, v := range arr {
		if v == target {
			return index
		}
	}
	return -1
}
func FindString(arr []string, target string) int {
	for index, v := range arr {
		if v == target {
			return index
		}
	}
	return -1
}
