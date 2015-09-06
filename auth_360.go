package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	log "github.com/cihub/seelog"
	"net/url"
	"time"
)

const (
	//CONST_360_APPID  = "200487256"
	//CONST_360_APPKEY = "a454dc82091e051014d0fe130f56b3c5" // 200487256"
	// CONST_360_APPID         = "200487256"
	// CONST_360_APPID         = "202199751"
	// CONST_360_APPKEY        = "32c9e19bad9102823e06e15de7f43ea9"
	// CONST_360_APPKEY_SECRET = "f8f1e3fc55c8d0593f463fc314dfc65e"
	CONST_360_REDIRECT_URI = "oob"
	// private key
)
const (
	CONST_360_GRANT_TYPE_AUTHORIZATION_CODE = "authorization_code"
	CONST_360_GRANT_TYPE_REFRESH_TOKEN      = "refresh_token"
)

// https://openapi.360.cn/oauth2/access_token?grant_type=authorization_code&
// code=120653f48687763d6ddc486fdce6b51c383c7ee544e6e5eab&client_id=0fb2676d5007f123756d1c1b4b5968bc&
// client_secret=1234567890ab18384f562d7d3f.....&redirect_uri=oob
func Do360Auth(appId, appKey, appSecret string, accessToken string, AccountId string, AccountName string, now time.Time) (succ int32, newAccountId string, newAccountName string, sdkRec interface{}) {
	// accessTokenRec, err := get360AccessTokenRecord(appId, appKey, appSecret, now, authorizationCode)
	// if err != nil {
	// 	succ = PB_ERRNO_AUTH_ERROR
	// 	return
	// }

	succ, newAccountId, newAccountName = Get360UserInfo(accessToken, now)
	sdkRec = accessToken
	return
}

// func get360AccessTokenRecord(appId, appKey, appSecret string, now time.Time, authorizationCode string) (accessTokenRec AccessTokenRec360, err error) {
// 	jsonBytes, err := goutils.GetHttpResponseAsJson(get360AccessToken(appId, appKey, appSecret, authorizationCode), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
// 	if err != nil {
// 		log.Error(err)
// 		fmt.Println(err)
// 		return
// 	}

// 	accessTokenRec = AccessTokenRec360{}
// 	json.Unmarshal(jsonBytes, &accessTokenRec)

// 	if accessTokenRec.IsTokenAvailable() {
// 		accessTokenRec.ExpireTime = now.Add(time.Second * time.Duration(goutils.Str2Int(accessTokenRec.ExpiresIn, 0)))
// 	}

// 	return
// }
func (rec *AccessTokenRec360) Refresh360Token(appId, appKey, appSecret string, now time.Time) bool {
	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_360_APPKEY, CONST_360_LOGIN_ACT, AccountId, sessionID, CONST_360_APPKEY_SECRET))
	queryStr := fmt.Sprintf("grant_type=%s&refresh_token=%s&client_id=%s&clent_secret=%s&scope=basic",
		CONST_360_GRANT_TYPE_REFRESH_TOKEN,
		url.QueryEscape(rec.RefreshToken),
		url.QueryEscape(appKey),
		url.QueryEscape(appSecret),
		url.QueryEscape(CONST_360_REDIRECT_URI))
	urlStr := fmt.Sprintf("https://openapi.360.cn/oauth2/access_token?%s", queryStr)

	jsonBytes, err := goutils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		return false
	}

	accessTokenRec := AccessTokenRec360{}
	json.Unmarshal(jsonBytes, &accessTokenRec)
	if accessTokenRec.IsTokenAvailable() {
		rec.Access_token = accessTokenRec.Access_token
		rec.RefreshToken = accessTokenRec.RefreshToken
		rec.ExpireTime = now.Add(time.Second * time.Duration(goutils.Str2Int(accessTokenRec.ExpiresIn, 0)))
		return true
	}

	return false
}

// // authorizationCode 换access token
// func get360AccessToken(appId, appKey, appSecret string, authorizationCode string) string {
// 	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_360_APPKEY, CONST_360_LOGIN_ACT, AccountId, sessionID, CONST_360_APPKEY_SECRET))
// 	queryStr := fmt.Sprintf("grant_type=%s&code=%s&client_id=%s&clent_secret=%s&redirect_uri=%s",
// 		CONST_360_GRANT_TYPE_AUTHORIZATION_CODE,
// 		url.QueryEscape(authorizationCode),
// 		url.QueryEscape(appKey),
// 		url.QueryEscape(appSecret),
// 		url.QueryEscape(CONST_360_REDIRECT_URI))
// 	urlStr := fmt.Sprintf("https://openapi.360.cn/oauth2/access_token?%s", queryStr)
// 	return urlStr
// }

// // {"error_code":"4000203","error":"client_id、client_secret不可用（OAuth2）"}
// {
// "access_token":"120652e586871bb6bbcd1c7b77818fb9c95d92f9e0b735873",
// "expires_in":"36000”,
// “scope”:”basic”, “refresh_token”:”12065961868762ec8ab911a3089a7ebdf11f8264d5836fd41”
// }
type AccessTokenRec360 struct {
	Error_code string `json:"error_code,omitempty"`
	Error      string `json:"error,omitempty"`
	// Scope      string `json:scope,omitempty`
	ExpiresIn string `json:"expires_in,omitempty"`

	Access_token string `json:"access_token,omitempty"`
	ExpireTime   time.Time
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (token *AccessTokenRec360) IsTokenAvailable() bool {
	if token.Access_token != "" {
		return true
	}
	return false
}

// id
// Y
// 360 用户 ID, 缺省返回
// name
// Y
// 360 用户名, 缺省返回
// avatar
// Y
// 360 用户头像, 缺省返回
// sex
// N
// 360 用户性别,仅在 fields 中包含时候才返回,返回值为:男,女或 者未知
// area
// N
// 360 用户地区,仅在 fields 中包含时候才返回
// nick
// N
// 用户昵称,无值时候返回空

type UserInfo360 struct {
	Id     string
	Name   string
	Avatar string
	Sex    string
	Nick   string
	Area   string
}

func Get360UserInfo(accessToken string, now time.Time) (status int32, accountid string, name string) {
	userInfo, err := get360UserInfoRecord(accessToken, now)
	if err != nil {
		log.Error("360auth", err)
		status = PB_ERRNO_AUTH_ERROR
		return
	}

	if userInfo.Id != "" {
		return PB_STATUS_SUCC, userInfo.Id, userInfo.Name
	}
	return PB_ERRNO_AUTH_ERROR, "", ""

}
func get360UserInfoRecord(accessToken string, now time.Time) (userInfo UserInfo360, err error) {
	jsonBytes, err := goutils.GetHttpResponseAsJson(get360UserInfoUrl(accessToken), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		return
	}
	userInfo = UserInfo360{}
	json.Unmarshal(jsonBytes, &userInfo)
	return
}

func get360UserInfoUrl(accessToken string) string {
	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_360_APPKEY, CONST_360_LOGIN_ACT, AccountId, sessionID, CONST_360_APPKEY_SECRET))

	urlStr := fmt.Sprintf("https://openapi.360.cn/user/me.json?access_token=%s&fields=id,name", url.QueryEscape(accessToken))
	return urlStr
}

// https://openapi.360.cn/user/me.json?access_token=12345678983b38aabcdef387453ac8133ac3263987654321&fields=id,name,avatar,sex,area
// {
// "id": "201459001",
// "name": "360U201459001",
// "avatar": "http://u1.qhimg.com/qhimg/quc/48_48/22/02/55/220255dq9816.3eceac.jpg?f=d140ae40ee93e8b 08ed6e9c53543903b",
// "sex": "未知"
// "area": "" }
