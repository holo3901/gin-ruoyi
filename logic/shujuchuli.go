package logic

import (
	"fmt"
	"strconv"
	"strings"
)

// 数组转string
func Join(sep []int, sp string) string {
	sarr := make([]string, len(sep))
	for i, v := range sep {
		sarr[i] = fmt.Sprint(v)
	}
	return strings.Join(sarr, fmt.Sprint(sp))
}

// string转数组int
func Split(data string) []int {
	var sa = strings.Split(data, ", ")
	var sarr []int
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		var v1, _ = strconv.Atoi(v)
		sarr = append(sarr, v1)
	}
	return sarr
}

// string转数组string
func SplitStr(data string) []string {
	var sa = strings.Split(data, ", ")
	var sarr []string
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		sarr = append(sarr, v)
	}
	return sarr
}
