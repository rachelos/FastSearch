package filecache

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
)

func insertSlash(input string) string {
	var builder strings.Builder
	length := len(input)
	var j = 1
	for i := 0; i < length; i++ {
		builder.WriteByte(input[i])
		if j == 2 {
			builder.WriteRune('/')
			j = 0
		}
		j++
	}
	return builder.String()
}

// 计算字符串的MD5哈希值
func md5String(s string) string {
	// 创建一个新的hash.Hash对象，这里使用的是md5算法
	hasher := md5.New()

	// 写入要计算哈希的数据
	_, err := io.WriteString(hasher, s)
	if err != nil {
		// 处理错误，这里只是简单地返回空字符串作为错误处理
		fmt.Println("Error writing string:", err)
		return ""
	}

	// Sum方法会返回计算出的哈希值（切片形式），其中Sum的第二个参数是可选的，
	// 用于提供一个预分配的切片以存储哈希值，以避免额外的内存分配。
	// 如果传入nil，则会分配一个新的切片。
	// 这里我们传入nil，因为我们不需要预先分配切片。
	hash := hasher.Sum(nil)

	// 将字节切片转换为十六进制表示的字符串
	return fmt.Sprintf("%x", hash)
}
func Sha1Hash(text string) string {
	// 创建一个新的hash.Hash对象，这里使用的是sha1.New()
	hasher := sha1.New()
	// 写入需要哈希的数据。这里我们将字符串转换为字节切片并写入
	io.WriteString(hasher, text)
	// Sum函数返回计算出的哈希值的字节切片。Sum(nil)表示将哈希值追加到一个空的切片中
	// 如果你不想追加到已有的切片，可以传递nil作为参数
	hashBytes := hasher.Sum(nil)
	// 将字节切片转换为十六进制字符串
	// 注意：strings.ToUpper()用于将小写字母转换为大写，如果你需要小写，可以去掉这个函数调用
	hashString := strings.ToUpper(fmt.Sprintf("%x", hashBytes))
	return hashString
}

func (e *FileCache) file_path(key string) string {
	p := md5String(key)
	path := e.Path + insertSlash(p[0:2])
	if !Dir_Exists(path) {
		os.MkdirAll(path, 0644)
	}
	data_path := path + p + e.Extension
	return data_path
}
