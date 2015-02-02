package goauth

import (
	"encoding/json"
	"github.com/0studio/goauth/utils"
	"net/url"
	"time"
)

/*
xy苹果助手 √
appid:
uid:
*/

// const (
// 	XYUID   = ""
// 	XYAppId = "100009"
// )

func DoXYAuth(appId, uid string, sid string, accountId string, now time.Time) (status int32) {
	value := url.Values{}
	value.Set("uid", accountId)
	value.Set("appid", appId)
	value.Set("token", sid)
	jsonBytes, err := getXYLoginResponse(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := XYLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	if loginInfo.Ret == 0 {
		status = 0
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

func getXYLoginResponse(v url.Values, now time.Time) (json []byte, err error) {
	return utils.PostFormHttpResponse("http://passport.xyzs.com/checkLogin.php", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

type XYLoginResp struct {
	Ret int    `json:"ret"`
	Err string `json:"error"`
}
