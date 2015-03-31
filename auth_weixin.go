package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"net/url"
	"time"
)

const (
	CONST_WeiXin_GRANT_TYPE_AUTHORIZATION_CODE = "authorization_code"
	CONST_WeiXin_GRANT_TYPE_REFRESH_TOKEN      = "refresh_token"
)

func NewWeiXinSDK(appId, appSecret, string, code string) (sdk WeiXinSDK) {
	sdk.AppId = appId
	sdk.AppSecret = appSecret
	sdk.Code = code
	return
}

// 	{
// "openid":"OPENID",
// "nickname":"NICKNAME",
// "sex":1,
// "province":"PROVINCE",
// "city":"CITY",
// "country":"COUNTRY",
// "headimgurl": "http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
// "privilege":[
// "PRIVILEGE1",
// "PRIVILEGE2"
// ],
// "unionid": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
// 	}

// 获取用户个人信息（UnionID机制）
// 此接口用于获取用户个人信息。开发者可通过OpenID来获取用户基本信息。
// 特别需要注意的是，如果开发者拥有多个移动应用、网站应用和公众帐号，
// 可通过获取用户基本信息中的unionid来区分用户的唯一性，因为只要是同一
// 个微信开放平台帐号下的移动应用、网站应用和公众帐号，用户的unionid是
// 唯一的。换句话说，同一用户，对同一个微信开放平台下的不同应用，
// unionid是相同的。

type UserInfoWeiXin struct {
	OpenId     string   `json:"openid,omitempty"`
	Unionid    string   `json:"unionid,omitempty"`
	NickName   string   `json:"nickname,omitempty"`
	Sex        int      `json:"sex,omitempty"`
	HeadImgUrl string   `json:"headimgurl,omitempty"`
	Province   string   `json:"province,omitempty"`
	City       string   `json:"city,omitempty"`
	Privilege  []string `json:"privilege,omitempty"`
}

type WeiXinSDK struct {
	AppId string
	// AppKey       string
	AppSecret    string
	Code         string // 用于获取access_token
	AccessToken  string `json:"access_token,omitempty"`  // 接口调用凭证
	ExpiresIn    int    `json:"expires_in,omitempty"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token,omitempty"` // 用户刷新access_token
	OpenId       string `json:"openid,omitempty"`        // 授权用户唯一标识
	Scope        string `json:"scope,omitempty"`         // 用户授权的作用域，使用逗号（,）分隔
	Unionid      string `json:"unionid,omitempty"`       // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。

	ErrorCode int    `json:"errcode,omitempty"`
	Error     string `json:"errmsg,omitempty"`
	// Scope      string `json:scope,omitempty`
	ExpireTime time.Time
	UserInfo   UserInfoWeiXin
}

func (sdk *WeiXinSDK) DoAuth(now time.Time) bool {
	return sdk.GetUserInfo(now)
}

func (sdk *WeiXinSDK) isTokenExired(now time.Time) bool {
	if sdk.AccessToken != "" && now.Before(sdk.ExpireTime) {
		return false
	}
	return true
}

func (sdk *WeiXinSDK) isRefreshable() bool {
	if sdk.AppId != "" && sdk.RefreshToken != "" {
		return true
	}
	return false
}
func (sdk *WeiXinSDK) GetAccessToken(now time.Time) (err error) {
	if !sdk.isTokenExired(now) { // 如果token 没过期， 直接返回，
		return
	}
	if sdk.isRefreshable() && sdk.RefreshWeiXinToken(now) {
		return
	}

	getAccessTokenurl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?grant_type=%s&code=%s&appid=%s&secret=%s",
		CONST_WeiXin_GRANT_TYPE_AUTHORIZATION_CODE,
		url.QueryEscape(sdk.Code),
		url.QueryEscape(sdk.AppId),
		url.QueryEscape(sdk.AppSecret))

	jsonBytes, err := goutils.GetHttpResponseAsJson(getAccessTokenurl, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, sdk)
	if err != nil {
		return
	}

	sdk.ExpireTime = now.Add(time.Second * time.Duration(sdk.ExpiresIn)).Add(-time.Minute)

	return
}

func (sdk *WeiXinSDK) RefreshWeiXinToken(now time.Time) bool {
	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_WeiXin_APPKEY, CONST_WeiXin_LOGIN_ACT, AccountId, sessionID, CONST_WeiXin_APPKEY_SECRET))
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?grant_type=%s&refresh_token=%s&appid=%s",
		CONST_WeiXin_GRANT_TYPE_REFRESH_TOKEN,
		url.QueryEscape(sdk.RefreshToken),
		url.QueryEscape(sdk.AppId))

	jsonBytes, err := goutils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		return false
	}

	response := WeiXinSDK{}
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return false
	}

	if sdk.ErrorCode == 0 {
		sdk.AccessToken = response.AccessToken
		sdk.OpenId = response.OpenId
		sdk.Scope = response.Scope
		sdk.ExpireTime = now.Add(time.Second * time.Duration(response.ExpiresIn)).Add(-time.Minute)
		return true

	}
	sdk.ErrorCode = response.ErrorCode
	sdk.Error = response.Error
	return false
}

func (sdk *WeiXinSDK) GetUserInfo(now time.Time) bool {
	var err error
	err = sdk.GetAccessToken(now)
	if err != nil {
		return false
	}

	return sdk.getWeiXinUserInfoRecord(now)
}
func (sdk *WeiXinSDK) getWeiXinUserInfoRecord(now time.Time) bool {
	jsonBytes, err := goutils.GetHttpResponseAsJson(getWeiXinUserInfoUrl(sdk.AccessToken, sdk.OpenId), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = json.Unmarshal(jsonBytes, &(sdk.UserInfo))
	if err != nil {
		return false
	}
	return true
}

func getWeiXinUserInfoUrl(accessToken string, openId string) string {
	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_WeiXin_APPKEY, CONST_WeiXin_LOGIN_ACT, AccountId, sessionID, CONST_WeiXin_APPKEY_SECRET))
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
		url.QueryEscape(accessToken), url.QueryEscape(openId))
	return urlStr
}

// func DoWeiXinAuth(appId, appKey, appSecret string, authorizationCode string, AccountId string, AccountName string, now time.Time) (succ int32, newAccountId string, newAccountName string, sdkRec interface{}) {
// 	sdk, err := getWeiXinAccessTokenRecord(appId, appKey, appSecret, now, authorizationCode)
// 	if err != nil {
// 		succ = PB_ERRNO_AUTH_ERROR
// 		return
// 	}

// 	if sdk.IsTokenAvailable() {
// 		succ, sdk.UserInfo = GetUserInfo(sdk.AccessToken, sdk.OpenId, now)
// 		// succ, newAccountId, newAccountName =
// 		sdkRec = sdk
// 		return
// 	} else {
// 		return PB_ERRNO_AUTH_ERROR, AccountId, AccountName, nil
// 	}
// }
