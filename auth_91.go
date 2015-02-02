package goauth

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github/0studio/goauth/utils"
	"net/url"
	"time"
)

const (
	CONST_91_LOGIN_ACT            = "4"
	CONST_91_QUERY_PAY_RESULT_ACT = "1" // 查询支付购买结果
)

// 错误码(0=失败,1=成功(SessionId 有效),2= AppId 无效,3= Act 无效,4=参数无效,5= Sign 无效,11=SessionId 无效)
const (
	LOGIN_91_ERROR_FAIL        = "0"
	LOGIN_91_ERROR_SUCC        = "1"
	LOGIN_91_ERROR_APPID       = "2"
	LOGIN_91_ERROR_ACT         = "3"
	LOGIN_91_ERROR_WRONG_PARAM = "4"
	LOGIN_91_ERROR_SIGN        = "5"
	LOGIN_91_ERROR_SESSIONID   = "6"
)

// http://service.sj.91.com/usercenter/AP.aspx?AppId=100814&Act=4&Uin=account&Sign=2910fd34a0e8128b66c790e0d1df3b42&SessionID=session
// {"ErrorCode":"0","ErrorDesc":"出错了"}

func Do91Auth(appId, appKey string, AccountId string, sessionID string, now time.Time) (status int32) {
	response := get91LoginResponse(appId, appKey, AccountId, sessionID, now)
	if response.ErrorCode == LOGIN_91_ERROR_SUCC {
		status = PB_STATUS_SUCC
		return
	} else if response.ErrorCode == LOGIN_91_ERROR_SIGN {
		status = PB_ERRNO_SDK_AUTH_WRONG_ACTION
		return
	}
	status = PB_ERRNO_AUTH_ERROR
	return
}

func get91LoginResponse(appId, appKey string, AccountId string, sessionID string, now time.Time) (response Response91Login) {
	jsonBytes, err := utils.GetHttpResponseAsJson(get91LoginUrl(appId, appKey, AccountId, sessionID), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	response = Response91Login{ErrorCode: LOGIN_91_ERROR_FAIL}
	if err != nil {
		log.Error("error_auth91", err)
		return
	}

	json.Unmarshal(jsonBytes, &response)
	return
}
func get91LoginUrl(appId, appKey string, AccountId string, sessionID string) string {
	Sign := utils.GetHexMd5(fmt.Sprintf("%s%s%s%s%s", appId, CONST_91_LOGIN_ACT, AccountId, sessionID, appKey))
	queryStr := fmt.Sprintf("AppId=%s&Act=4&Uin=%s&Sign=%s&SessionID=%s", url.QueryEscape(appId), url.QueryEscape(AccountId), Sign, url.QueryEscape(sessionID))
	urlStr := fmt.Sprintf("http://service.sj.91.com/usercenter/AP.aspx?%s", queryStr)
	return urlStr
}

type Response91Login struct {
	// {"ErrorCode":"0","ErrorDesc":"出错了"}
	ErrorCode string
	ErrorDesc string
}

// ////////////////////////////////////////////////////////////////////////////////
// // pay
// ////////////////////////////////////////////////////////////////////////////////
// func GetPayResult(CooOrderSerial string) (result PayResult91) { // 商户订单号
// 	jsonBytes := GetHttpResponseAsJson(getPayResultUrl(CooOrderSerial))
// 	result = PayResult91{}
// 	json.Unmarshal(jsonBytes, &result)
// 	return
// }

// // AppId=100010&Act=1&CooOrderSerial=4764f9ff47174e2287ded31673a73a50&Sign=b4a 67099d354dce29ab8a9b526344cb6
// func getPayResultUrl(CooOrderSerial string) string { // 商户订单号
// 	// 参数值与 AppKey 的 MD5 值String.Format("{0}{1}{2}{3}", AppId, Act, CooOrderSerial, AppKey).HashToMD5Hex();
// 	Sign := GetHexMd5(fmt.Sprintf("%s%s%s%s%s", CONST_91_APPID, CONST_91_QUERY_PAY_RESULT_ACT, CooOrderSerial, CONST_91_APPKEY))
// 	queryStr := fmt.Sprintf("AppId=%s&Act=1&CooOrderSerial=%s&Sign=%s", url.QueryEscape(CONST_91_APPID), url.QueryEscape(CONST_91_APPID), Sign)
// 	urlStr := fmt.Sprintf("http://service.sj.91.com/usercenter/AP.aspx?%s", queryStr)
// 	return urlStr
// }

// type PayResult91 struct {
// 	ConsumeStreamId string //  消费流水号,平台流水号
// 	CooOrderSerial  string // 商户订单号,购买时应用传入,原样返回给应用
// 	MerchantId      string // 商户 ID
// 	AppId           string // 应用 ID,必须对应游戏客户端中使用的 APPID
// 	ProductName     string // 应用名称
// 	Uin             string //  91 账号 ID,购买时应用传入,原样返回给应用
// 	GoodsId         string //  商品 ID,购买时应用传入,原样返回给应用
// 	GoodsInfo       string // 商品名称,购买时应用传入,原样返回给应用
// 	GoodsCount      string //  商品数量,购买时应用传入,原样返回给应用
// 	OriginalMoney   string // 原价(格式:0.00),购买时应用传入的单价*总数,总原价
// 	OrderMoney      string // 实际价格(格式:0.00),购买时应用传入的单价*总数,总实际 价格。(打个比方:原价 100,现价 50,可以说明原价和实际价 格的关系,如无特殊这两个价格应保持一致)
// 	Note            string // 即支付描述(客户端 API 参数中的 payDescription 字段) 购买时客户端应用通过 API 传入,原样返回给应用服务器 开发者可以利用该字段,定义自己的扩展数据。例如区分游戏 服务器。不超过 20 个字符,只支持英文或数字。 订单信息较多的可通过订单号将数据保存在游戏服务器端数据 库中。
// 	PayStatus       string //  支付状态:0=失败,1=成功,2=正在处理中 (仅当 ErrorCode=1,表示接口调用成功时,才需要检查此字段状态,开发商需要根据此参数状态发放物品)
// 	CreateTime      string // 创建时间(yyyy-MM-dd HH:mm:ss)
// 	Sign            string // 参数的 MD5 值,其中 AppKey 为 91SNS 平台分配的应用密钥 String.Format("{0}{1}{2}{3}{4}{5}{6}{7}{8}{9:0.00}{10: 0.00}{11}{12}{13:yyyy-MM-dd HH:mm:ss}{14}", ConsumeStreamId, CooOrderSerial, MerchantId, AppId, ProductName, Uin, GoodsId, GoodsInfo, GoodsCount, OriginalMoney, OrderMoney, Note, PayStatus, CreateTime, AppKey).HashToMD5Hex()
// 	ErrorCode       string //  错误码(0=失败,1=成功,2= AppId 无效,3= Act 无效,4=参 数无效,5= Sign 无效, 11=没有该订单)
// 	ErrorDesc       string //错误描述
// }
