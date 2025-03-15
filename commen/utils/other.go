package utils

import "strconv"

// 辅助函数，将字符串转换为 float64
func ParseFloat(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0 // 或者处理错误
	}
	return val
}

// 辅助函数，将字符串转换为 int
func ParseInt(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0 // 或者处理错误
	}
	return val
}
