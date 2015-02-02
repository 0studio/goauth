package utils

import (
	"strconv"
)

func Str2Int(str string, defaultvalue int) (value int) {
	value, err := strconv.Atoi(str)
	if err != nil {
		value = defaultvalue
	}
	return
}
func Int2Str(v int) string {
	return strconv.Itoa(v)
}
