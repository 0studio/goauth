package goauth

import (
	"encoding/json"
	"net/url"
	"time"

	"strconv"

	"github.com/0studio/goutils"
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
	return goutils.PostFormHttpResponse("https://pay.i4.cn/member_third.action", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

type AiSiLoginResp struct {
	Status   int32  `json:"status"`
	UserName string `json:"username"`
	UserId   int32  `json:"userid"`
}

func (res AiSiLoginResp) GetUserIdAsString() string {
	return strconv.Itoa(int(res.UserId))
}
func (res AiSiLoginResp) IsSucc() bool {
	return res.Status == 0 && res.UserId != 0
}

func DoAiSiAuth2(token string, now time.Time) (loginInfo AiSiLoginResp) {
	value := url.Values{}
	value.Set("token", token)
	jsonBytes, err := getLoginAiSiResponse(value, now)
	if err != nil {
		return
	}
	loginInfo = AiSiLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	return
}
