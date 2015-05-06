package goauth

// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/0studio/goutils"
// 	"github.com/vincent-petithory/dataurl"
// 	"strings"
// 	"time"
// )

// // 疑问 待确定， ， 文档中两处用到timestamp 一处string 一处int ,尚未确定

// // 对传入参数所有键值对的 value 进行 urlencode 转码(注意!进行
// // urlencode 时要将空 格转化为%20 而不是+)后重新拼接成字符串 string2。
// // golang 默认 url/QueryEscape() 会把空格转成+
// // 所以， 改用 dataurl.EscapeString(来转换
// const (
// 	CONST_WeiXin_GRANT_TYPE_AUTHORIZATION_CODE = "authorization_code"
// 	CONST_WeiXin_GRANT_TYPE_REFRESH_TOKEN      = "refresh_token"
// )

// // 注意:appSecret、paySignKey、partnerKey 是验证商户唯一性的安全标识,
// // 请妥善保管。 对于 appSecret 和 paySignKey 的区别,可以这样认
// // 为:appSecret 是 API 使用时的登录密码, 会在网络中传播的;而
// // paySignKey 是在所有支付相关数据传输时用于加密并进行身份校验 的密钥,
// // 仅保留在第三方后台和微信后台,不会在网络中传播,而且 paySignKey 仅用
// // 于支 付请求。

// // appKey 即paySignKey
// func NewWeiXinSDK(appId, appKey, appSecret, partnerId, partnerKey, code string, notifyUrl string) (sdk WeiXinSDK) {
// 	sdk.AppId = appId
// 	sdk.AppKey = appKey
// 	sdk.AppSecret = appSecret
// 	sdk.PartnerId = partnerId
// 	sdk.PartnerKey = partnerKey
// 	sdk.Code = code
// 	sdk.NotifyUrl = notifyUrl
// 	return
// }

// // 	{
// // "openid":"OPENID",
// // "nickname":"NICKNAME",
// // "sex":1,
// // "province":"PROVINCE",
// // "city":"CITY",
// // "country":"COUNTRY",
// // "headimgurl": "http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
// // "privilege":["PRIVILEGE1", "PRIVILEGE2"],
// // "unionid": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
// // 	}

// // 获取用户个人信息（UnionID机制）
// // 此接口用于获取用户个人信息。开发者可通过OpenID来获取用户基本信息。
// // 特别需要注意的是，如果开发者拥有多个移动应用、网站应用和公众帐号，
// // 可通过获取用户基本信息中的unionid来区分用户的唯一性，因为只要是同一
// // 个微信开放平台帐号下的移动应用、网站应用和公众帐号，用户的unionid是
// // 唯一的。换句话说，同一用户，对同一个微信开放平台下的不同应用，
// // unionid是相同的。

// type UserInfoWeiXin struct {
// 	OpenId     string   `json:"openid,omitempty"`
// 	Unionid    string   `json:"unionid,omitempty"`
// 	NickName   string   `json:"nickname,omitempty"`
// 	Sex        int      `json:"sex,omitempty"`
// 	HeadImgUrl string   `json:"headimgurl,omitempty"`
// 	Province   string   `json:"province,omitempty"`
// 	City       string   `json:"city,omitempty"`
// 	Privilege  []string `json:"privilege,omitempty"`
// }

// type WeiXinSDK struct {
// 	AppId      string
// 	AppKey     string //即paySignKey
// 	PartnerId  string // 注册时分配的财付通商户号 partnerId;
// 	PartnerKey string
// 	NotifyUrl  string // for pay
// 	// AppKey       string
// 	AppSecret    string
// 	Code         string // 用于获取access_token
// 	AccessToken  string `json:"access_token,omitempty"`  // 接口调用凭证
// 	ExpiresIn    int    `json:"expires_in,omitempty"`    // access_token接口调用凭证超时时间，单位（秒）
// 	RefreshToken string `json:"refresh_token,omitempty"` // 用户刷新access_token
// 	OpenId       string `json:"openid,omitempty"`        // 授权用户唯一标识
// 	Scope        string `json:"scope,omitempty"`         // 用户授权的作用域，使用逗号（,）分隔
// 	Unionid      string `json:"unionid,omitempty"`       // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。

// 	ErrorCode int    `json:"errcode,omitempty"`
// 	Error     string `json:"errmsg,omitempty"`
// 	// Scope      string `json:scope,omitempty`
// 	ExpireTime time.Time
// 	UserInfo   UserInfoWeiXin
// }

// func (sdk *WeiXinSDK) DoAuth(now time.Time) bool {
// 	return sdk.GetUserInfo(now)
// }

// func (sdk *WeiXinSDK) isTokenExired(now time.Time) bool {
// 	if sdk.AccessToken != "" && now.Before(sdk.ExpireTime) {
// 		return false
// 	}
// 	return true
// }

// func (sdk *WeiXinSDK) isRefreshable() bool {
// 	if sdk.AppId != "" && sdk.RefreshToken != "" {
// 		return true
// 	}
// 	return false
// }
// func (sdk *WeiXinSDK) GetAccessToken(now time.Time) (err error) {
// 	if !sdk.isTokenExired(now) { // 如果token 没过期， 直接返回，
// 		return
// 	}
// 	if sdk.isRefreshable() && sdk.RefreshWeiXinToken(now) {
// 		return
// 	}

// 	getAccessTokenurl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?grant_type=%s&code=%s&appid=%s&secret=%s",
// 		CONST_WeiXin_GRANT_TYPE_AUTHORIZATION_CODE,
// 		dataurl.EscapeString(sdk.Code),
// 		dataurl.EscapeString(sdk.AppId),
// 		dataurl.EscapeString(sdk.AppSecret))

// 	jsonBytes, err := goutils.GetHttpResponseAsJson(getAccessTokenurl, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(jsonBytes, sdk)
// 	if err != nil {
// 		return
// 	}

// 	sdk.ExpireTime = now.Add(time.Second * time.Duration(sdk.ExpiresIn)).Add(-time.Minute)

// 	return
// }

// func (sdk *WeiXinSDK) RefreshWeiXinToken(now time.Time) bool {
// 	// Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_WeiXin_APPKEY, CONST_WeiXin_LOGIN_ACT, AccountId, sessionID, CONST_WeiXin_APPKEY_SECRET))
// 	urlStr := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?grant_type=%s&refresh_token=%s&appid=%s",
// 		CONST_WeiXin_GRANT_TYPE_REFRESH_TOKEN,
// 		dataurl.EscapeString(sdk.RefreshToken),
// 		dataurl.EscapeString(sdk.AppId))

// 	jsonBytes, err := goutils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
// 	if err != nil {
// 		return false
// 	}

// 	response := WeiXinSDK{}
// 	err = json.Unmarshal(jsonBytes, &response)
// 	if err != nil {
// 		return false
// 	}

// 	if sdk.ErrorCode == 0 {
// 		sdk.AccessToken = response.AccessToken
// 		sdk.OpenId = response.OpenId
// 		sdk.Scope = response.Scope
// 		sdk.ExpireTime = now.Add(time.Second * time.Duration(response.ExpiresIn)).Add(-time.Minute)
// 		return true

// 	}
// 	sdk.ErrorCode = response.ErrorCode
// 	sdk.Error = response.Error
// 	return false
// }

// func (sdk *WeiXinSDK) GetUserInfo(now time.Time) bool {
// 	var err error
// 	err = sdk.GetAccessToken(now)
// 	if err != nil {
// 		return false
// 	}

// 	return sdk.getUserInfoRecord(now)
// }
// func (sdk *WeiXinSDK) getUserInfoRecord(now time.Time) bool {
// 	urlStr := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", dataurl.EscapeString(sdk.AccessToken), dataurl.EscapeString(sdk.OpenId))
// 	jsonBytes, err := goutils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	err = json.Unmarshal(jsonBytes, &(sdk.UserInfo))
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// //  生成预支付订单
// func (sdk *WeiXinSDK) GetPrePay(traceId, productDesc, attach, noncestr, clientIP string, totalFee int32, now time.Time) (clientReqJson string, ok bool) {
// 	var err error
// 	err = sdk.GetAccessToken(now)
// 	if err != nil {
// 		ok = false
// 		return
// 	}

// 	urlStr := fmt.Sprintf("https://api.weixin.qq.com/pay/genprepay?access_token=%s", dataurl.EscapeString(sdk.AccessToken))
// 	jsonData := sdk.genPrePayPostData(traceId, productDesc, attach, noncestr, clientIP, totalFee, now)
// 	responseData, err := goutils.PostHttpResponse(urlStr, jsonData, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
// 	if err != nil {
// 		ok = false
// 		return
// 	}
// 	prePayResp := prePayResponse{}
// 	json.Unmarshal(responseData, &prePayResp)
// 	if prePayResp.ErrMsg != "Success" {
// 		ok = false
// 		return
// 	}
// 	return sdk.getPrePay(traceId, noncestr, prePayResp.Prepayid, now)
// }
// func (sdk *WeiXinSDK) getPrePay(traceId, noncestr, prePayId string, now time.Time) (clientReqJson string, ok bool) {
// 	prePayData := prePayPostData{

// 		AppId:     sdk.AppId,
// 		appKey:    sdk.AppKey,
// 		TradeId:   traceId,
// 		Noncestr:  noncestr,
// 		Timestamp: int32(now.Unix()),
// 		Package:   "Sign=WXpay",
// 		PartnerId: sdk.PartnerId,
// 		Prepayid:  prePayId,
// 	}
// 	prePayData.Sign = prePayData.getSign()
// 	jsonData, _ := json.Marshal(prePayData)

// 	clientReqJson = string(jsonData)
// 	ok = true
// 	return

// }

// type prePayPostData struct {
// 	AppId     string `json:"appid,omitempty"`     //
// 	appKey    string `json:"appkey,omitempty"`    //
// 	TradeId   string `json:"traceid,omitempty"`   //:由开发者自定义,可用于订单的查询与跟踪,建议根据支付用户信息生成此 id
// 	Noncestr  string `json:"noncestr,omitempty"`  //32 位内的随机串,防重发
// 	Timestamp int32  `json:"timestamp,omitempty"` //时间戳,为自 1970 年 1 月 1 日 00:00(时区:东八区)至当前时间的秒数
// 	Package   string `json:"package,omitempty"`   //
// 	PartnerId string `json:"partnerid,omitempty"` //
// 	Prepayid  string `json:"prepayid,omitempty"`  //
// 	Sign      string `json:"sign,omitempty"`      //
// }

// func (data prePayPostData) toString1() string {
// 	return fmt.Sprintf("appid=%s&appkey=%s&noncestr=%s&package=%s&partnerid=%s&prepayid=%s&timestamp=%d",
// 		data.AppId, data.appKey, data.Noncestr, data.Package, data.PartnerId, data.Prepayid, data.Timestamp)
// }
// func (data prePayPostData) getSign() string {
// 	return goutils.GetSha1(data.toString1())
// }

// type prePayResponse struct {
// 	Prepayid string `json:"prepayid,omitempty"` //
// 	Errcode  int    `json:"errcode,omitempty"`  //
// 	ErrMsg   string `json:"errmsg,omitempty"`   //
// }

// // {
// // "appid":"wxd930ea5d5a258f4f",
// // "traceid":"test_1399514976",
// // "noncestr":"e7d161ac8d8a76529d39d9f5b4249ccb ",
// // "timestamp":1399514976,
// // "package":"bank_type=WX&body=%E6%94%AF%E4%BB%98%E6%B5%8B%E8%AF%95&fee_type=1&input_charset=UTF-8&notify_url=http%3A%2F%2Fweixin.qq.com&out_trade_ no=7240b65810859cbf2a8d9f76a638c0a3&partner=1900000109&spbill_create_ip=196.168.1.1& total_fee=1&sign=7F77B507B755B3262884291517E380F8",
// // "sign_method":"sha1",
// // "app_signature":"7f77b507b755b3262884291517e380f8"
// // }
// // 生成如上json 串
// // totalFee 单位 分
// const (
// 	WEIXIN_FEE_TYPE  = "1" // rmb
// 	WEIXIN_BANK_TYPE = "WX"
// 	WEIXIN_CHARSET   = "UTF-8" // GBK or UTF-8
// )

// func (sdk *WeiXinSDK) genPrePayPostData(traceId, productDesc, attach, noncestr, clientIP string, totalFee int32, now time.Time) []byte {
// 	data := prePayIdData{
// 		AppId:      sdk.AppId,
// 		TradeId:    traceId,
// 		Noncestr:   noncestr,
// 		Timestamp:  int32(now.Unix()),
// 		Package:    sdk.genPrePayPackage(traceId, productDesc, sdk.NotifyUrl, attach, clientIP, totalFee),
// 		SignMethod: "sha1",
// 	}
// 	data.AppSignature = data.toAppSignature()
// 	jsonData, _ := json.Marshal(&data)

// 	return jsonData
// }

// type prePayIdData struct {
// 	AppId        string `json:"appid,omitempty"`       //
// 	appKey       string `json:"appkey,omitempty"`      //
// 	TradeId      string `json:"traceid,omitempty"`     //:由开发者自定义,可用于订单的查询与跟踪,建议根据支付用户信息生成此 id
// 	Noncestr     string `json:"noncestr,omitempty"`    //32 位内的随机串,防重发
// 	Timestamp    int32  `json:"timestamp,omitempty"`   //订单详情
// 	Package      string `json:"package,omitempty"`     //时间戳,为自 1970 年 1 月 1 日 00:00(时区:东八区)至当前时间的秒数
// 	SignMethod   string `json:"sign_method,omitempty"` //签名算法 default sha1
// 	AppSignature string `json:"app_signature,omitempty"`
// }

// func (data prePayIdData) toStrint1() string {
// 	return fmt.Sprintf("appid=%s&appkey=%s&noncestr=%s&package=%s&timestamp=%d&traceid=%s",
// 		data.AppId, data.appKey, data.Noncestr, data.Package, data.Timestamp, data.TradeId)
// }
// func (data prePayIdData) toAppSignature() string {
// 	return goutils.GetSha1(data.toStrint1())
// }
// func (sdk *WeiXinSDK) genPrePayPackage(traceId, productDesc, notifyUrl, attach, clientIP string, totalFee int32) string {
// 	pack := prePayPackage{
// 		PartnerId:    sdk.PartnerId,
// 		PaternerKey:  sdk.PartnerKey,
// 		Attach:       attach,
// 		BankType:     WEIXIN_BANK_TYPE,
// 		ProductDesc:  productDesc,
// 		FeeType:      WEIXIN_FEE_TYPE,
// 		InputCharSet: WEIXIN_CHARSET,
// 		NotifyUrl:    notifyUrl,
// 		TradeNo:      traceId,
// 		ClientIP:     clientIP,
// 		TotalFee:     goutils.Int322Str(totalFee),
// 	}
// 	return pack.toString()
// }

// type prePayPackage struct {
// 	PaternerKey  string
// 	Attach       string `json:"attach,omitempty"`           //附加数据,原样返回;
// 	BankType     string `json:"bank_type,omitempty"`        //银行通道类型,固定为"WX";
// 	ProductDesc  string `json:"body,omitempty"`             //:商品描述
// 	FeeType      string `json:"fee_type,omitempty"`         //:取值:1(人民币),暂只支持 1;
// 	InputCharSet string `json:"input_charset,omitempty"`    //取值范围:"GBK"、"UTF-8",默认:"GBK"
// 	NotifyUrl    string `json:"notify_url,omitempty"`       //通知 URL
// 	TradeNo      string `json:"out_trade_no,omitempty"`     //商户系统内部的订单号,32 个字符内、可包含字 母,确保在商户系统唯一,详见注意事项,第 5 项
// 	PartnerId    string `json:"partner,omitempty"`          //注册时分配的财付通商户号 partnerId;
// 	ClientIP     string `json:"spbill_create_ip,omitempty"` //通知 URL
// 	TotalFee     string `json:"total_fee,omitempty"`        //订单总金额,单位为分;
// }

// // a.对所有传入参数按照字段名的 ASCII 码从小到大排序(字典序)后,使用 URL 键值 对的格式(即 key1=value1&key2=value2...)拼接成字符串 string1;
// func (pack prePayPackage) toString1() (str string) {
// 	if pack.Attach != "" {
// 		str = fmt.Sprintf("attach=%s&", pack.Attach)
// 	}
// 	str = fmt.Sprintf("%sbank_type=%s&body=%s&fee_type=%s&input_charset=%s&notify_url=%s&out_trade_no=%s&partner=%s&spbill_create_ip=%s&total_fee=%s",
// 		str,
// 		WEIXIN_BANK_TYPE, pack.ProductDesc, WEIXIN_FEE_TYPE, pack.InputCharSet, pack.NotifyUrl,
// 		pack.TradeNo, pack.PartnerId, pack.ClientIP, pack.TotalFee)

// 	return
// }
// func (pack prePayPackage) getSign() string {
// 	tmpStr := fmt.Sprintf("%s&key=%s", pack.toString1(), pack.PaternerKey)
// 	return strings.ToUpper(goutils.GetHexMd5(tmpStr))
// }
// func (pack prePayPackage) toString() (str string) {
// 	if pack.Attach != "" {
// 		str = fmt.Sprintf("attach=%s&", pack.Attach)
// 	}

// 	str = fmt.Sprintf("%sbank_type=%s&body=%s&fee_type=%s&input_charset=%s&notify_url=%s&out_trade_no=%s&partner=%s&spbill_create_ip=%s&total_fee=%s&sign=%s",
// 		str,
// 		WEIXIN_BANK_TYPE, dataurl.EscapeString(pack.ProductDesc), dataurl.EscapeString(pack.FeeType), dataurl.EscapeString(pack.InputCharSet), dataurl.EscapeString(pack.NotifyUrl),
// 		dataurl.EscapeString(pack.TradeNo), dataurl.EscapeString(pack.PartnerId), dataurl.EscapeString(pack.ClientIP), dataurl.EscapeString(pack.TotalFee),
// 		pack.getSign())
// 	return
// }
