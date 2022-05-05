package tools

import "strconv"

func StringToInt64(str string) int64 {
	result, _ := strconv.ParseInt(str, 10, 64)
	return result
}

func Int64ToString(num int64) string {
	result := strconv.FormatInt(num, 10)
	return result
}
