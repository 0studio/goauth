package goauth

import (
	"encoding/json"
	"github/0studio/goauth/utils"
	"net/url"
	"time"
)

/*
	爱思助手sdk
*/

func DoAiSiAuth(token string, now time.Time) (status int32) {
	value := url.Values{}
	value.Set("token", token)
	jsonBytes, err := getLoginAiSiResponse(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := AiSiLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	if loginInfo.Status == 0 {
		status = 0
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

func getLoginAiSiResponse(v url.Values, now time.Time) (json []byte, err error) {
	return utils.PostFormHttpResponse("https://pay.i4.cn/member_third.action", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

type AiSiLoginResp struct {
	Status   int32  `json:"status"`
	UserName string `json:"username"`
	UserId   int32  `json:"userid"`
}
