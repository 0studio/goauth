package goauth

import (
	"bytes"
	"github.com/vincent-petithory/dataurl"
	"net/url"
	"sort"
)

// 对传入参数所有键值对的 value 进行 urlencode 转码(注意!进行
// urlencode 时要将空 格转化为%20 而不是+)后重新拼接成字符串 string2。
// golang 默认 url/QueryEscape() 会把空格转成+
// 所以， 改用 dataurl.EscapeString(来转换

type kv struct {
	key   string
	value string
}
type kvList []kv

func JoinParam(params map[string]string, encodeValue, useStandardEncode bool) (ret string) {
	var buf bytes.Buffer
	var idx int
	for key, value := range params {
		buf.WriteString(key)
		buf.WriteString("=")
		if encodeValue {
			if useStandardEncode {
				buf.WriteString(url.QueryEscape(value))
			} else {
				buf.WriteString(dataurl.EscapeString(value))
			}

		} else {
			buf.WriteString(value)
		}

		if idx != len(params)-1 {
			buf.WriteString("&")
		}
		idx++
	}
	return buf.String()
}

func (list kvList) JoinParam(encodeValue, useStandardEncode bool) (ret string) {
	var buf bytes.Buffer
	for idx, entity := range list {
		buf.WriteString(entity.key)
		buf.WriteString("=")
		if encodeValue {
			if useStandardEncode {
				buf.WriteString(url.QueryEscape(entity.value))
			} else {
				buf.WriteString(dataurl.EscapeString(entity.value))
			}

		} else {
			buf.WriteString(entity.value)
		}

		if idx != len(list)-1 {
			buf.WriteString("&")
		}
	}
	return buf.String()
}

func (l kvList) Sort() {
	sort.Sort(l)
}

// 实现sort 接口
func (l kvList) Len() int {
	return len(l)
}
func (l kvList) Less(i, j int) bool {
	return l[i].key < l[j].key //
}
func (l kvList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func sortParams(params map[string]string) (l kvList) {
	l = make(kvList, len(params))
	var i int32 = 0
	for key, value := range params {
		l[i] = kv{key, value}
		i++
	}
	l.Sort()
	return
}
