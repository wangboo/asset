package asset

import (
	"os"
	"time"
)

var (
	cache = map[string]*BuildResult{}
)

//编译结果
type BuildResult struct {
	size     int64
	data     []byte
	filePath string
	buildAt  time.Time
}

// 查询缓存
func findInCache(filePath string, buildFunc func(filePath string) []byte) []byte {
	if value, ok := cache[filePath]; ok {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			// 文件不存在了
			delete(cache, filePath)
			return []byte("//find not found")
		}
		if value.size == fileInfo.Size() {
			return value.data
		}
		// 文件改变了
		data := buildFunc(filePath)
		value.data = data
		return data
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return []byte("//find not found")
	}
	data := buildFunc(filePath)
	rst := &BuildResult{size: fileInfo.Size(), data: data, filePath: filePath, buildAt: time.Now()}
	cache[filePath] = rst
	return data
}
