package goauth

import (
	"bytes"
	"fmt"
	"github.com/0studio/goutils"
	"github.com/vincent-petithory/dataurl"
	"sort"
)

// http://wiki.open.qq.com/wiki/%E8%85%BE%E8%AE%AF%E5%BC%80%E6%94%BE%E5%B9%B3%E5%8F%B0%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%94%E7%94%A8%E7%AD%BE%E5%90%8D%E5%8F%82%E6%95%B0sig%E7%9A%84%E8%AF%B4%E6%98%8E
// http://wiki.mg.open.qq.com/index.php?title=Sig%E8%AE%A1%E7%AE%97%E4%BA%8B%E4%BE%8B
// method:GET or POST
// URI不含host，URI示例：/v3/user/get_info）
func snsSigCheck(method, uri, appkey string, params map[string]string) (sign string) {
	encodedURI := dataurl.EscapeString(uri)
	sortParams := sortParams(params)
	unSignedStr := fmt.Sprintf("%s&%s&%s", method, encodedURI, dataurl.EscapeString(sortParams.JoinParam(false)))
	hashKey := appkey + "&"
	sign = goutils.HmacSha1Base64([]byte(hashKey), []byte(unSignedStr))
	return
}

type kv struct {
	key   string
	value string
}
type kvList []kv

func (list kvList) JoinParam(encodeValue bool) (ret string) {
	var buf bytes.Buffer

	for idx, entity := range list {
		buf.WriteString(entity.key)
		buf.WriteString("=")
		if encodeValue {
			buf.WriteString(dataurl.EscapeString(entity.value))
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
