package utils

import (
	"log"
	"strconv"
)

// 辅助函数，将字符串转换为 float64
func ParseFloat(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("parsefloat failed, err:", err)
		return 0.0, err // 或者处理错误
	}
	return val, nil
}

// 辅助函数，将字符串转换为 int
func ParseInt(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println("parse int failed, err:", err)
		return -1, err // 或者处理错误
	}
	return val, nil
}

func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "是":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "否":
		return false, nil
	}
	return false, nil
}
