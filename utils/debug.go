package utils

import "fmt"

// DebugLevel 全局调试级别
var DebugLevel = 1

// SetDebugLevel 设置调试级别
func SetDebugLevel(level int) {
	DebugLevel = level
}

// DebugLog 输出调试信息
func DebugLog(level int, format string, args ...interface{}) {
	if level <= DebugLevel {
		fmt.Printf(format+"\n", args...)
	}
}
