package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"time"
)

func (sdk *TencentMSDK) getRefreshTokenUrl(now time.Time, extString string) string {
	if sdk.IsModeProduct() {
		return fmt.Sprintf("http://msdk.qq.com/auth/refresh_token?%s", sdk.getCommonUri(now, extString)) // 【现网】
	}
	return fmt.Sprintf("http://msdktest.qq.com/auth/refresh_token?%s", sdk.getCommonUri(now, extString)) //  【沙箱】
}
func (sdk *TencentMSDK) isTokenExired(now time.Time) bool {
	if sdk.accessToken != "" && now.Before(sdk.expireTime) {
		return false
	}
	return true
}

type msdkRefreshToken struct {
	AppId        string `json:"appid,omitempty"`        //
	RefreshToken string `json:"refreshToken,omitempty"` //
}
type msdkRefreshTokenResult struct {
	Status       int32  `json:"ret,omitempty"`          //
	Msg          string `json:"msg,omitempty"`          //
	ExpiresIn    int    `json:"expiresIn,omitempty"`    //
	AccessToken  string `json:"accessToken,omitempty"`  //
	RefreshToken string `json:"refreshToken,omitempty"` //
	OpenId       string `json:"openid,omitempty"`       //
	Scope        string `json:"scope,omitempty"`        //
}

func (sdk *TencentMSDK) isRefreshable() bool {
	if sdk.appId != "" && sdk.refreshToken != "" {
		return true
	}
	return false
}

func (sdk *TencentMSDK) RefreshWeiXinToken(now time.Time) {
	urlStr := sdk.getRefreshTokenUrl(now, "")
	req := msdkRefreshToken{
		AppId:        sdk.appId,
		RefreshToken: sdk.refreshToken,
	}
	data, _ := json.Marshal(&req)

	resultData, err := goutils.PostHttpResponse(urlStr, data, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		fmt.Println("weixin_refresh_token", err)
		return
	}
	var result msdkRefreshTokenResult
	json.Unmarshal(resultData, &result)
	if result.Status == 0 {
		sdk.accessToken = result.AccessToken
		sdk.refreshToken = result.RefreshToken
		sdk.expireTime = now.Add(time.Second * time.Duration(result.ExpiresIn))
		return
	}
	return
}

type msdkCheckToken struct {
	OpenId      string `json:"openid,omitempty"`      //
	AccessToken string `json:"accessToken,omitempty"` //
}
type msdkCheckTokenResult struct {
	Status int32  `json:"ret,omitempty"` //
	Msg    string `json:"msg,omitempty"` //
}

func (sdk *TencentMSDK) getCheckWeixinTokenUrl(now time.Time, extString string) string {
	if sdk.IsModeProduct() {
		return fmt.Sprintf("http://msdk.qq.com/auth/check_token?%s", sdk.getCommonUri(now, extString)) // 【现网】
	}
	return fmt.Sprintf("http://msdktest.qq.com/auth/check_token?%s", sdk.getCommonUri(now, extString)) //  【沙箱】
}

const (
	STATUS_FAIL = 1
)

func (sdk *TencentMSDK) CheckWeiXinToken(now time.Time) (status int32) {
	urlStr := sdk.getCheckWeixinTokenUrl(now, "")
	req := msdkCheckToken{
		OpenId:      sdk.openId,
		AccessToken: sdk.accessToken,
	}
	data, _ := json.Marshal(&req)

	resultData, err := goutils.PostHttpResponse(urlStr, data, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		fmt.Println("weixin_check_token_fail", err)
		return STATUS_FAIL
	}
	var result msdkCheckTokenResult
	err = json.Unmarshal(resultData, &result)
	if result.Status == 0 && err == nil {
		return 0
	}
	fmt.Println("weixin_check_token_fail", err, result)
	if result.Status != 0 {
		return result.Status
	}
	return STATUS_FAIL
}
