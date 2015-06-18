package goauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/0studio/goutils"
	"net/url"
	"time"
)

/*
爱贝海马助手
*/

func DoIAPPPAYAuth(appId string, sid string, now time.Time) (status int32) {
	value := url.Values{}
	value.Set("appid", appId)
	value.Set("logintoken", sid)
	jsonBytes, err := getIAppLoginResponse(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := IAppLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	if loginInfo.Name != "" {
		status = 0
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

func getIAppLoginResponse(v url.Values, now time.Time) (json []byte, err error) {
	return goutils.PostFormHttpResponse("http://ipay.iapppay.com:8888/iapppay/tokencheck", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

type IAppLoginResp struct {
	Name   string `json:"loginname"`
	UserId string `json:"userid"`
}

// 2.4	支付结果查询
// http://ipay.iapppay.com:9999/payapi/queryresult
func IPayQueryResult(appid, cporderid, rsaPrivateKey, rsaPublicKey string, now time.Time) (ret IPayQueryResultResponse, err error) {

	reqJsonStruct := iPayQueryResultRequest{
		AppId:     appid,
		CPorderid: cporderid,
	}
	reqJsonBytes, _ := json.Marshal(reqJsonStruct)
	transdata := string(reqJsonBytes)
	reqParams := make(map[string]string)
	reqParams["transdata"] = transdata
	reqParams["sign"] = goutils.RSAPack1SignWithMD5(rsaPrivateKey, transdata)
	reqParams["signtype"] = "RSA"
	postContent := JoinParam(reqParams, true, true)
	response, err := goutils.PostHttpResponse("http://ipay.iapppay.com:9999/payapi/queryresult", []byte(postContent), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		fmt.Println(err)
		return

	}
	values, err := url.ParseQuery(string(response))
	if err != nil {
		return
	}
	data := values.Get("transdata")
	signRet := values.Get("sign")
	err = json.Unmarshal([]byte(data), &ret)
	if err != nil {
		return
	}
	if !goutils.VerifyRSASignWithMD5(rsaPublicKey, data, signRet) {
		return IPayQueryResultResponse{}, errors.New("sign_verify_fail")
	}
	return
}

type iPayQueryResultRequest struct {
	AppId     string `json:"appid"`
	CPorderid string `json:"cporderid"`
}
type IPayQueryResultResponse struct {
	Code      string  `json:"code,omitempty"`
	Result    int32   `json:"result,omitempty"`    // 0–交易成功； 2–待支付
	CPorderid string  `json:"cporderid,omitempty"` // 商户订单号
	TransId   string  `json:"transid,omitempty"`   // 计费支付平台的交易流水号
	Appuserid string  `json:"appuserid,omitempty"` // 用户在商户应用的唯一标识
	Waresid   int32   `json:"waresid,omitempty"`   // 平台为应用内需计费商品分配的编码
	Money     float64 `json:"money,omitempty"`     // 本次交易的金额
}

func (res IPayQueryResultResponse) IsOk() bool {
	if res.Code == "" && res.Result == 0 && res.CPorderid != "" {
		return true
	}
	return false

}
