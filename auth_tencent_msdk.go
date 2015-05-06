package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"net/http"
	"net/url"
	"time"
)

const (
	TENCENT_PLATFORM_QQ     = "qq"
	TENCENT_PLATFORM_WEIXIN = "weixin"
)

type TencentMSDK struct {
	mode        string // dev or pro ,means develop or product
	openId      string
	platform    string // "qq" "weixin"
	pf          string // sdk pf 平台来源，平台-注册渠道-系统运行平台-安装渠道-业务自定义。 从MSDK的getpf接口获取例如： qq_m_qq-2001-android-2011-xxxx
	pfKey       string
	appId       string
	appKey      string
	accessToken string // openkey,从手Q登录态或者微信登录态中获取的access_token 的值
	payToken    string // Q登录时从手Q登录态中获取的pay_token的值,使用MSDK登录后获取到的eToken_QQ_Pay返回内容就是pay_token； 微信登录时特别注意该参数传空。
	serverId    uint64 // zoneid 	账户分区ID。应用如果没有分区：传zoneid=1
}

func NewTencentMSDK(mode string, tencentPlatform, appId, appKey, accessToken, payToken string, pf, pfKey string, serverId uint64, openId string) (sdk TencentMSDK) {
	sdk = TencentMSDK{
		mode:        mode,
		platform:    tencentPlatform,
		pf:          pf,
		pfKey:       pfKey,
		appId:       appId,
		appKey:      appKey,
		accessToken: accessToken,
		payToken:    payToken,
		openId:      openId,
	}

	return
}
func (sdk *TencentMSDK) IsPlatfromQQ() bool {
	return sdk.platform == TENCENT_PLATFORM_QQ
}
func (sdk *TencentMSDK) IsPlatfromWeiXin() bool {
	return sdk.platform == TENCENT_PLATFORM_WEIXIN
}

func (sdk *TencentMSDK) getPayToken() string {
	if sdk.IsPlatfromQQ() {
		return sdk.payToken
	}
	if sdk.IsPlatfromWeiXin() {
		return ""
	}
	return sdk.payToken
}
func (sdk *TencentMSDK) IsModeProduct() bool {
	return sdk.mode == "pro"
}
func (sdk *TencentMSDK) getSig(now time.Time) (value string) {
	unixTime := int32(now.Unix())
	value = goutils.GetHexMd5(fmt.Sprintf("%s%d", sdk.appKey, unixTime))
	return
}
func (sdk *TencentMSDK) getCommonUri(now time.Time, extString string) (value string) { // extString 透传参数
	unixTime := int32(now.Unix())
	sig := sdk.getSig(now)
	if extString == "" {
		value = fmt.Sprintf("appid=%s&timestamp=%d&sig=%s&encode=1", sdk.appId, unixTime, sig)
		return
	}
	value = fmt.Sprintf("appid=%s&timestamp=%d&sig=%s&encode=1&msdkExtInfo=%s", sdk.appId, unixTime, sig, extString)
	return

}
func (sdk *TencentMSDK) getBalanceMUrl(now time.Time, extString string) (value string) {
	// 查询余额接口,获取用户游戏币余额
	if sdk.IsModeProduct() {
		return fmt.Sprintf("http://msdk.qq.com/mpay/get_balance_m?%s", sdk.getCommonUri(now, extString)) // 【现网】
	}
	return fmt.Sprintf("http://msdktest.qq.com/mpay/get_balance_m?%s", sdk.getCommonUri(now, extString)) //  【沙箱】
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
