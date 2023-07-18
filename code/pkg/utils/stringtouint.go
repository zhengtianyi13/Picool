package util

import (
	"strconv"
)

// StringToUint 字符串转uint
func StringToUint(str string) uint {
	i, _ := strconv.Atoi(str)
	return uint(i)
}
