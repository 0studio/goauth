package goauth

// 拇指玩
import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	log "github.com/cihub/seelog"
	"net/url"
	"time"
)

// 游戏名字：超神学院
// APPKEY：534b9fc9c1652
// 签名校验串：534b9fcbc1aad(支付回调用到)

// const (
// )

const (
	MUZHIWAN_CODE_SUCC             = "1"
	MUZHIWAN_CODE_TOKEN_INVALIDATE = "2010"
	MUZHIWAN_CODE_APPKEY_ERROR     = "2001"
)

// 1：成功登录
// -1：用户名格式不正确（注册）
// -2：用户名或密码错误
// 1005：密码格式不正确
// 2001：appkey无效
// 2010：token无效
// 2011：token已过期

func DomuzhiwanAuth(appKey string, authorizationCode string, AccountId string, AccountName string, now time.Time) (succ int32, newAccountId string, newAccountName string, sdkRec interface{}) {
	accessTokenRec := getmuzhiwanAccessTokenRecord(appKey, authorizationCode, now)
	if accessTokenRec.Code == MUZHIWAN_CODE_SUCC {
		succ = PB_STATUS_SUCC
		newAccountId = accessTokenRec.User.Uid
		newAccountName = accessTokenRec.User.Username
		return
	}
	log.Errorf("muzhiwan auth error,accountid=%s,accountname=%s,token=%s code=%s (2011:token expire ,2010 invalidate token,2001:invalidate appkey,1005:password format wrong,-1 or -2 name or password wrong )", AccountId, AccountName, authorizationCode, accessTokenRec.Code)
	if accessTokenRec.Code == MUZHIWAN_CODE_TOKEN_INVALIDATE {
		return PB_ERRNO_AUTH_ERROR, AccountId, AccountName, nil
	}
	if accessTokenRec.Code == MUZHIWAN_CODE_APPKEY_ERROR {
		return PB_ERRNO_SDK_AUTH_WRONG_ACTION, AccountId, AccountName, nil
	}

	return PB_ERRNO_AUTH_ERROR, AccountId, AccountName, nil

}

func getmuzhiwanAccessTokenRecord(appKey string, authorizationCode string, now time.Time) (accessTokenRec AccessTokenRecmuzhiwan) {
	jsonBytes, err := goutils.GetHttpResponseAsJson(getmuzhiwanAccessToken(appKey, authorizationCode), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	accessTokenRec = AccessTokenRecmuzhiwan{}
	if err != nil {
		log.Error("auth_muzhiwan_error", err)
		return
	}
	log.Debug(string(jsonBytes))

	json.Unmarshal(jsonBytes, &accessTokenRec)

	if accessTokenRec.Code != MUZHIWAN_CODE_SUCC {
		accessErrorTokenRec := AccessTokenRecmuzhiwanErrorCode{}
		json.Unmarshal(jsonBytes, &accessErrorTokenRec)
		accessTokenRec.Code = goutils.Int2Str(accessErrorTokenRec.Code)
		accessTokenRec.Msg = accessErrorTokenRec.Msg
	}

	return
}

func getmuzhiwanAccessToken(appKey string, token string) string {
	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_muzhiwan_APPID, CONST_muzhiwan_LOGIN_ACT, AccountId, sessionID, CONST_MUZHIWAN_APPKEY))
	queryStr := fmt.Sprintf("token=%s&appkey=%s", url.QueryEscape(token), url.QueryEscape(appKey))
	urlStr := fmt.Sprintf("http://sdk.muzhiwan.com/oauth2/getuser.php?%s", queryStr)
	return urlStr
}

//变态拇指玩 code 正确时返回code:"1" , 错误时code:2011 ,一会字符串 一会数字
// {"code":2011,"msg":"Token is overtime"}
// str := `{"code":"1","msg":"","user":{"username":"1397468703318","uid":"4946302","sex":"0","mail":"","icon":"http:\/\/www.muzhiwan.com\/index.php?action=profile&opt=getPic&uid=4946302&size=78&ismobile=1"}}`

type AccessTokenRecmuzhiwan struct {
	Code string
	Msg  string
	User UserInfomuzhiwan `json:"user,omitempty"`
}

type AccessTokenRecmuzhiwanErrorCode struct {
	Code int
	Msg  string
}

type UserInfomuzhiwan struct {
	Uid      string
	Username string
	Avatar   string
	Sex      string
	Mail     string
	Icon     string
}
