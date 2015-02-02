package goauth

//itools
import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goauth/utils"
	log "github.com/cihub/seelog"
	"net/url"
	"time"
)

/*
成功: {"status":"success"} 失败: {"status":"fail"}
*/

func DoIToolsAuth(appId, appKey string, sessionId string, now time.Time) int32 {
	RetResponse := getIToolsLoginResponse(appId, appKey, sessionId, now)
	if RetResponse == "success" {
		return PB_STATUS_SUCC
	}
	return PB_ERRNO_AUTH_ERROR
}
func getIToolsLoginResponse(appId, appKey string, token string, now time.Time) string {
	contentByte, err := utils.GetHttpResponseAsJson(getIToolsLoginUrl(appId, appKey, token), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		log.Error("auth_itools_error", err)
		return "false"
	}
	var resp map[string]string
	json.Unmarshal(contentByte, &resp)
	return resp["status"]
}
func getIToolsLoginUrl(appId, appKey string, token string) string {
	signStr := fmt.Sprintf("appid=%s&sessionid=%s", appId, token)
	sign := utils.GetHexMd5(signStr)
	urlStr := fmt.Sprintf("https://pay.itools.cn/?r=auth/verify&appid=%s&sessionid=%s&sign=%s", appId, url.QueryEscape(token), sign)
	return urlStr
}
