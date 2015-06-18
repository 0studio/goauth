package goauth

import (
	"fmt"
	"github.com/0studio/goutils"
	"github.com/vincent-petithory/dataurl"
)

// http://wiki.open.qq.com/wiki/%E8%85%BE%E8%AE%AF%E5%BC%80%E6%94%BE%E5%B9%B3%E5%8F%B0%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%94%E7%94%A8%E7%AD%BE%E5%90%8D%E5%8F%82%E6%95%B0sig%E7%9A%84%E8%AF%B4%E6%98%8E
// http://wiki.mg.open.qq.com/index.php?title=Sig%E8%AE%A1%E7%AE%97%E4%BA%8B%E4%BE%8B
// method:GET or POST
// URI不含host，URI示例：/v3/user/get_info）
func snsSigCheck(method, uri, appkey string, params map[string]string) (sign string) {
	encodedURI := dataurl.EscapeString(uri)
	sortParams := sortParams(params)
	unSignedStr := fmt.Sprintf("%s&%s&%s", method, encodedURI, dataurl.EscapeString(sortParams.JoinParam(false, false)))
	hashKey := appkey + "&"
	sign = goutils.HmacSha1Base64([]byte(hashKey), []byte(unSignedStr))
	return
}
