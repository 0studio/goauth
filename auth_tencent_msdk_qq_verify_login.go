package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"time"
)

func (sdk *TencentMSDK) getSig(now time.Time) (value string) {
	unixTime := int32(now.Unix())
	value = goutils.GetHexMd5(fmt.Sprintf("%s%d", sdk.appKey, unixTime))
	return
}

func (sdk *TencentMSDK) getCommonUri(now time.Time, extString string) (value string) { // extString 透传参数
	unixTime := int32(now.Unix())
	sig := sdk.getSig(now)
	if extString == "" {
		value = fmt.Sprintf("appid=%s&timestamp=%d&openid=%s&sig=%s&encode=1", sdk.appId, unixTime, sdk.openId, sig)
		return
	}
	value = fmt.Sprintf("appid=%s&timestamp=%d&openid=%s&sig=%s&encode=1&msdkExtInfo=%s", sdk.appId, unixTime, sdk.openId, sig, extString)
	return

}

type QQVerifyLogin struct {
	AppId   string `json:"appid,omitempty"`   //
	OpenId  string `json:"openid,omitempty"`  //
	OpenKey string `json:"openkey,omitempty"` //
	Userip  string `json:"userip,omitempty"`  //
}

type QQVerifyLoginResult struct {
	Status int32  `json:"ret,omitempty"` //
	Msg    string `json:"msg,omitempty"` //
}

func (sdk *TencentMSDK) getQQVerifyLoginUrl(now time.Time, extString string) (value string) { // extString 透传参数
	// 查询余额接口,获取用户游戏币余额
	if sdk.IsModeProduct() {
		// return "http://msdk.qq.com/auth/verify_login"
		return fmt.Sprintf("http://msdk.qq.com/auth/verify_login?%s", sdk.getCommonUri(now, extString)) // 【现网】
	}
	// return "http://msdktest.qq.com/auth/verify_login"
	return fmt.Sprintf("http://msdktest.qq.com/auth/verify_login?%s", sdk.getCommonUri(now, extString)) //  【沙箱】

}
func (sdk *TencentMSDK) QQVerifyLogin(now time.Time, ip string) (succ bool) { // extString 透传参数
	urlStr := sdk.getQQVerifyLoginUrl(now, "")
	req := QQVerifyLogin{
		AppId:   sdk.appId,
		OpenId:  sdk.openId,
		OpenKey: sdk.accessToken,
		Userip:  ip,
	}
	data, _ := json.Marshal(&req)
	resultData, err := goutils.PostHttpResponse(urlStr, data, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		fmt.Println("qqverifylogin", err)
		return
	}
	var result QQVerifyLoginResult
	json.Unmarshal(resultData, &result)
	if result.Status == 0 {
		return true
	}
	return false
}
