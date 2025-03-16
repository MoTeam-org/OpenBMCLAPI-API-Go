package utils

import "strings"

// ANSI 颜色转义码
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

// ColorText 为文本添加颜色
func ColorText(color, text string) string {
	return color + text + Reset
}

// PadString 返回固定宽度的字符串，考虑中文字符
func PadString(str string, width int, alignLeft bool) string {
	strRunes := []rune(str)
	strWidth := len(strRunes)

	if strWidth >= width {
		return string(strRunes[:width])
	}

	padding := strings.Repeat(" ", width-strWidth)
	if alignLeft {
		return str + padding
	}
	return padding + str
}
