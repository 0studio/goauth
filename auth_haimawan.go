package goauth

import (
	"fmt"
	"github.com/0studio/goutils"
	"net/url"
	"time"
)

//海马移动支付
// sdk 有ip 白名单机制 ，需要将你的ip 配置到白名单中， 才能登录成功 ,多个ip 用|分隔
// pay.haima.me

func DoHaimaWanAuth(appId string, token string, now time.Time) bool {
	value := url.Values{}
	value.Set("appid", appId)
	value.Set("t", token)
	jsonBytes, err := getHaimaWanAuthResp(value, now)
	if err != nil {
		return false
	}
	resp := string(jsonBytes)
	if resp == "success" {
		return true
	}
	fmt.Println("haima auth error:"+appId, token, resp)

	return false
}

// return success or fail
func getHaimaWanAuthResp(v url.Values, now time.Time) (json []byte, err error) {
	return goutils.PostFormHttpResponse("http://api.haimawan.com/index.php?m=api&a=validate_token", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}
