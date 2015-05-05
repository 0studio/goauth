package goauth

import (
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

// openid 	《Android SDK 公共参数解释说明》
// openkey 	《Android SDK 公共参数解释说明》
// pay_token 	《Android SDK 公共参数解释说明》
// appid 	《Android SDK 公共参数解释说明》
// ts 	《Android SDK 公共参数解释说明》
// sig 	《Android SDK 公共参数解释说明》
// pf 	《Android SDK 公共参数解释说明》
// pfkey 	《Android SDK 公共参数解释说明》
// zoneid 	《Android SDK 公共参数解释说明》
// format 	《Android SDK 公共参数解释说明》

// type GetBalanceM struct {
// 	OpenId    string `json:"openid,omitempty"`    //
// 	OpenKey   string `json:"openkey,omitempty"`   //
// 	PayToken  string `json:"pay_token,omitempty"` //
// 	AppId     string `json:"appid,omitempty"`     //
// 	Timestamp int32  `json:"ts,omitempty"`        //时间戳,为自 1970 年 1 月 1 日 00:00(时区:东八区)至当前时间的秒数
// 	Sig       string `json:"sig,omitempty"`       //
// 	Pf        string `json:"pf,omitempty"`        //
// 	PfKey     string `json:"pfkey,omitempty"`     //
// 	ServerId  string `json:"zoneid,omitempty"`    //
// 	Format    string `json:"format,omitempty"`    //
// }

// func (m GetBalanceM) getSign(method string, appKey string, now time.Time) string {
// }

func (sdk *TencentMSDK) GetBalanceM(now time.Time, extString string) (value float64) {
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

	goutils.GetHttpResponseWithCookie(urlStr, cookies, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	return
}
