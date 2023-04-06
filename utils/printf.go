package utils

import (
	"fmt"
	"os"
)

func ErrorFmtPrintf(format string, a ...interface{}) {
	// 设置红色字体
	const redColor = "\x1b[31m"
	// 重置样式
	const resetColor = "\x1b[0m"

	// 使用 fmt.Sprintf 格式化错误信息
	errorMsg := fmt.Sprintf(format, a...)
	// 打印红色的错误提示
	fmt.Printf(redColor + "error：" + errorMsg + resetColor + "\n")

	// 退出程序
	os.Exit(1)
}

func DefaultFmtPrintf(format string, a ...interface{}) {
	// 设置绿色字体
	const greenColor = "\x1b[32m"
	// 重置样式
	const resetColor = "\x1b[0m"

	// 使用 fmt.Sprintf 格式化信息
	msg := fmt.Sprintf(format, a...)
	// 打印绿色输出
	fmt.Printf(greenColor + msg + resetColor + "\n")
}
