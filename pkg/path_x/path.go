package path_x

import (
	"path/filepath"
	"strings"
)

// PathValid 路径合法性检查
func PathValid(path string) bool {
	// 统一转换成 Unix 路径分隔符
	path = filepath.ToSlash(path)
	// path是否以./开头
	if strings.HasSuffix(path, "./") {
		// 去掉./
		path = strings.TrimPrefix(path, "./")
	}
	// 使用filepath.Clean函数规范化路径
	cleanedPath := filepath.Clean(path)

	// 检查规范化前后的路径是否相同
	if filepath.ToSlash(cleanedPath) != path {
		// 路径不合法，包含了路径遍历漏洞
		return false
	}

	// 路径合法
	return true
}
