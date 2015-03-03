package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
//游戏工厂
gamekey：645110998e14a6926e0b2cb13bca01e7
secretkey：0AxyGZEmWPAtKJXZGLXl5SG492xMeZt7
*/

func DoYouXiGongChangAuth(appKey, appSecret string, sessionid string, acountName string, now time.Time) (status int32) {
	token := sessionid
	cp := acountName
	value := url.Values{}
	value.Set("token", token)
	value.Set("cp", cp)
	timeStamp := strconv.Itoa(int(now.Unix()))
	value.Set("timestamp", timeStamp)
	value.Set("gamekey", appKey)
	sign := getYouXiGongChangLoginSign(appSecret, value)
	value.Set("_sign", sign)
	//sign

	jsonBytes, err := getLoginYouXiGongChangResp(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := YouXiGongChangLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	if loginInfo.Result == "0" {
		status = 0
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

func getLoginYouXiGongChangResp(v url.Values, now time.Time) (json []byte, err error) {
	return goutils.PostFormHttpResponse("http://anyapi.mobile.youxigongchang.com/foreign/oauth/verification.php", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

func getYouXiGongChangLoginSign(appSecret string, v url.Values) (sign string) {
	keys := make([]string, len(v))
	i := 0
	for k, _ := range v {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	data := make([]string, len(keys))
	for i, key := range keys {
		data[i] = fmt.Sprintf("%s=%s", key, goutils.UrlEncode(v.Get(key)))
	}
	str := strings.Join(data, "&")
	strMd5 := goutils.GetHexMd5(str)
	sign = goutils.GetHexMd5(strMd5 + appSecret)
	return
}

type YouXiGongChangLoginResp struct {
	Result     string `json:"result"`
	ResultDesc string `json:"result_desc"`
}
