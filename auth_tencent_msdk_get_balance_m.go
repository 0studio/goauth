package goauth

import (
	// "fmt"
	"encoding/json"
	"github.com/0studio/goutils"
	"net/http"
	"net/url"
	"time"
)

func (sdk *TencentMSDK) getBalanceMUrl(now time.Time, extString string) (value string) {
	// 查询余额接口,获取用户游戏币余额
	if sdk.IsModeProduct() {
		return "http://msdk.qq.com/mpay/get_balance_m?"
		// return fmt.Sprintf("http://msdk.qq.com/mpay/get_balance_m?%s", sdk.getCommonUri(now, extString)) // 【现网】
	}
	return "http://msdktest.qq.com/mpay/get_balance_m?"
	// return fmt.Sprintf("http://msdktest.qq.com/mpay/get_balance_m?%s", sdk.getCommonUri(now, extString)) //  【沙箱】
}
func (sdk *TencentMSDK) getBalanceMCookie() (cookies []*http.Cookie) {
	// 用户账户类型，（手Q）session_id =“openid”；（微信）sessionId = "hy_gameid"
	if sdk.IsPlatfromQQ() {
		cookies = make([]*http.Cookie, 3)
		cookies[0] = &http.Cookie{Name: "session_id", Value: "openid"}
		cookies[1] = &http.Cookie{Name: "session_type", Value: "kp_actoken"}
		cookies[2] = &http.Cookie{Name: "org_loc", Value: url.QueryEscape("/mpay/get_balance_m")}
		return
	}
	if sdk.IsPlatfromWeiXin() {
		cookies = make([]*http.Cookie, 3)
		cookies[0] = &http.Cookie{Name: "session_id", Value: "hy_gameid"}
		cookies[1] = &http.Cookie{Name: "session_type", Value: "wc_actoken"}
		cookies[2] = &http.Cookie{Name: "org_loc", Value: url.QueryEscape("/mpay/get_balance_m")}
		return
	}

	return

}

type GetBalanceMResult struct {
	Status     int32   `json:"ret,omitempty"`         //
	Msg        string  `json:"msg,omitempty"`         //
	Balance    float64 `json:"balance,omitempty"`     //余额 	游戏币个数（包含了赠送游戏币）
	GenBalance float64 `json:"gen_balance,omitempty"` //赠送游戏币个数
	FirstSave  int32   `json:"first_save,omitempty"`  //是否满足首次充值，1：满足，0：不满足
	SaveAmt    float64 `json:"save_amt ,omitempty"`   //累计充值金额(单位：游戏币)
}

func (sdk *TencentMSDK) GetBalanceM(now time.Time, extString string) (balance float64, totalBalance float64, ok bool) {
	// 查询余额接口,获取用户游戏币余额
	urlStr := sdk.getBalanceMUrl(now, extString)
	cookies := sdk.getBalanceMCookie()

	uri := "/mpay/get_balance_m"
	params := make(map[string]string)
	params["openid"] = sdk.openId
	params["openkey"] = sdk.accessToken
	params["pay_token"] = sdk.getPayToken()
	params["appid"] = sdk.appId
	params["ts"] = goutils.Int322Str(int32(now.Unix()))
	params["pf"] = sdk.pf
	params["pfkey"] = sdk.pfKey
	params["zoneid"] = goutils.Uint642Str(sdk.serverId)
	params["format"] = "json"
	params["sig"] = snsSigCheck("GET", uri, sdk.appKey, params)

	resultData, err := goutils.GetHttpResponseWithCookie(urlStr, cookies, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		return
	}
	var result GetBalanceMResult
	json.Unmarshal(resultData, &result)
	if result.Status == 0 { // succ
		balance = result.Balance
		totalBalance = result.SaveAmt
		ok = true
		return
	}

	return
}
