package utils

import "strings"

// TruncateString 截断字符串，考虑中文字符
func TruncateString(s string, length int) string {
	r := []rune(s)
	if len(r) > length {
		return string(r[:length-3]) + "..."
	}
	padding := strings.Repeat(" ", length-len(r))
	return s + padding
}
