package goauth

import (
	"encoding/json"
	"github/0studio/goauth/utils"
	"net/url"
	"time"
)

/*
爱贝海马助手
*/

func DoIAPPPAYAuth(appId string, sid string, now time.Time) (status int32) {
	value := url.Values{}
	value.Set("appid", appId)
	value.Set("logintoken", sid)
	jsonBytes, err := getIAppLoginResponse(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := IAppLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	if loginInfo.Name != "" {
		status = 0
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

func getIAppLoginResponse(v url.Values, now time.Time) (json []byte, err error) {
	return utils.PostFormHttpResponse("http://ipay.iapppay.com:8888/iapppay/tokencheck", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

type IAppLoginResp struct {
	Name   string `json:"loginname"`
	UserId string `json:"userid"`
}
