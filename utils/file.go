package utils

import (
	"fmt"
	"strings"
)

func ValidateDirectoryName(dirName string) error {
	// 检查文件夹名称长度是否为空或者太短/太长
	if len(dirName) == 0 || len(dirName) < 2 || len(dirName) > 255 {
		return fmt.Errorf("directory name is too short or too long")
	}
	// 检查文件夹名称中是否包含非法字符
	for _, r := range dirName {
		if strings.ContainsRune(`<>$#%{}|\^~[]'";`, r) {
			return fmt.Errorf("directory name contains invalid character: %s", string(r))
		}
	}
	return nil
}
