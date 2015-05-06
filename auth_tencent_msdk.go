package goauth

import ()

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

func NewTencentMSDK(mode, tencentPlatform, appId, appKey, accessToken, payToken, pf, pfKey, openId string, serverId uint64) (sdk TencentMSDK) {
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

// func (sdk *TencentMSDK) getSig(now time.Time) (value string) {
// 	unixTime := int32(now.Unix())
// 	value = goutils.GetHexMd5(fmt.Sprintf("%s%d", sdk.appKey, unixTime))
// 	return
// }

// func (sdk *TencentMSDK) getCommonUri(now time.Time, extString string) (value string) { // extString 透传参数
// 	unixTime := int32(now.Unix())
// 	sig := sdk.getSig(now)
// 	if extString == "" {
// 		value = fmt.Sprintf("appid=%s&timestamp=%d&sig=%s&encode=1", sdk.appId, unixTime, sig)
// 		return
// 	}
// 	value = fmt.Sprintf("appid=%s&timestamp=%d&sig=%s&encode=1&msdkExtInfo=%s", sdk.appId, unixTime, sig, extString)
// 	return

// }
