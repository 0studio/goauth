package goauth

// 豌豆夹
import (
	"fmt"
	log "github.com/cihub/seelog"
	"github/0studio/goauth/utils"
	"net/url"
	"time"
)

// 直接返回true or false 的字符串

func DoWanDouJiaAuth(appId, appKey string, token string, AccountId string, now time.Time) int32 {
	RetResponse := getWanDouJiaLoginResponse(appId, appKey, AccountId, token, now)
	if RetResponse == "true" {
		return PB_STATUS_SUCC
	}
	return PB_ERRNO_AUTH_ERROR
}
func getWanDouJiaLoginResponse(appId, appKey string, AccountId string, token string, now time.Time) string {
	contentByte, err := utils.GetHttpResponseAsJson(getWanDouJiaLoginUrl(appId, appKey, AccountId, token), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		log.Error("auth_wandoujia_error", err)
		return "false"
	}

	return string(contentByte)
}
func getWanDouJiaLoginUrl(appId, appKey string, AccountId string, token string) string {
	urlStr := fmt.Sprintf("https://pay.wandoujia.com/api/uid/check?uid=%s&token=%s&appkey_id=%s",
		url.QueryEscape(AccountId), url.QueryEscape(token), url.QueryEscape(appKey))
	return urlStr
}
