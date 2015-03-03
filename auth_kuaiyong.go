package goauth

// 快用

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	log "github.com/cihub/seelog"
	"time"
)

func DoKuaiyongAuth(appId, appKey string, sessionId string, now time.Time) (int32, string, string) {
	RetResponse, id, name := getKUAIYONGLoginResponse(appId, appKey, sessionId, now)
	if RetResponse == 0 {
		return PB_STATUS_SUCC, id, name
	}
	return PB_ERRNO_AUTH_ERROR, id, name
}
func getKUAIYONGLoginResponse(appId, appKey string, token string, now time.Time) (int, string, string) {
	contentByte, err := goutils.GetHttpResponseAsJson(getKUAIYONGLoginUrl(appId, appKey, token), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		log.Error("auth_kuaiyong_error", err)
		return 1, "", ""
	}
	resp := KuaiyongResp{}
	json.Unmarshal(contentByte, &resp)
	return resp.Code, resp.Data.AccountId, resp.Data.AccountName
}
func getKUAIYONGLoginUrl(appId, appKey string, token string) string {
	signStr := fmt.Sprintf("%s%s", appKey, token)
	sign := goutils.GetHexMd5(signStr)
	urlStr := fmt.Sprintf("http://f_signin.bppstore.com/loginCheck.php?tokenKey=%s&sign=%s", token, sign)
	log.Debug(urlStr, "url")
	return urlStr
}

type KuaiyongResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data KyData `json:"data"`
}

type KyData struct {
	AccountId   string `json:"guid"`
	AccountName string `json:"username"`
}
