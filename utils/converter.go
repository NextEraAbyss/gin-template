package utils

import (
	"strconv"
)

// ParseInt 将字符串解析为整数
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseInt64 将字符串解析为int64
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseFloat64 将字符串解析为float64
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseBool 将字符串解析为bool
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
