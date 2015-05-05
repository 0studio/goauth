package goauth

import (
	"github.com/stretchr/testify/assert"
	"testing"
	// "time"
)

func TestPrePayPackage(t *testing.T) {
	pack := prePayPackage{}
	pack.BankType = "WX"
	pack.ProductDesc = "支付测试"
	pack.FeeType = "1"
	pack.InputCharSet = "UTF-8"
	pack.NotifyUrl = "http://weixin.qq.com"
	pack.TradeNo = "7240b65810859cbf2a8d9f76a638c0a3"
	pack.PartnerId = "1900000109"
	pack.ClientIP = "196.168.1.1"
	pack.TotalFee = "1"
	pack.PaternerKey = "8934e7d15453e97507ef794cf7b0519d"
	assert.Equal(t, "bank_type=WX&body=支付测试&fee_type=1&input_charset=UTF-8&notify_url=http://weixin.qq.com&out_trade_no=7240b65810859cbf2a8d9f76a638c0a3&partner=1900000109&spbill_create_ip=196.168.1.1&total_fee=1",
		pack.toString1())
	assert.Equal(t, "7F77B507B755B3262884291517E380F8", pack.getSign())
	assert.Equal(t, "bank_type=WX&body=%E6%94%AF%E4%BB%98%E6%B5%8B%E8%AF%95&fee_type=1&input_charset=UTF-8&notify_url=http%3A%2F%2Fweixin.qq.com&out_trade_no=7240b65810859cbf2a8d9f76a638c0a3&partner=1900000109&spbill_create_ip=196.168.1.1&total_fee=1&sign=7F77B507B755B3262884291517E380F8", pack.toString())
}

func TestPrePaySign(t *testing.T) {
	pack := prePayPackage{}
	pack.BankType = "WX"
	pack.ProductDesc = "支付测试"
	pack.FeeType = "1"
	pack.InputCharSet = "UTF-8"
	pack.NotifyUrl = "http://weixin.qq.com"
	pack.TradeNo = "7240b65810859cbf2a8d9f76a638c0a3"
	pack.PartnerId = "1900000109"
	pack.ClientIP = "196.168.1.1"
	pack.TotalFee = "1"
	pack.PaternerKey = "8934e7d15453e97507ef794cf7b0519d"

	data := prePayIdData{
		AppId:      "wxd930ea5d5a258f4f",
		appKey:     "L8LrMqqeGRxST5reouB0K66CaYAWpqhAVsq7ggKkxHCOastWksvuX1uvmvQclxaHoYd3ElNBrNO2DHnnzgfVG9Qs473M3DTOZug5er46FhuGofumV8H2FVR9qkjSlC5K",
		TradeId:    "test_1399514976",
		Noncestr:   "e7d161ac8d8a76529d39d9f5b4249ccb",
		Timestamp:  1399514976,
		Package:    pack.toString(),
		SignMethod: "sha1",
	}
	assert.Equal(t, "appid=wxd930ea5d5a258f4f&appkey=L8LrMqqeGRxST5reouB0K66CaYAWpqhAVsq7ggKkxHCOastWksvuX1uvmvQclxaHoYd3ElNBrNO2DHnnzgfVG9Qs473M3DTOZug5er46FhuGofumV8H2FVR9qkjSlC5K&noncestr=e7d161ac8d8a76529d39d9f5b4249ccb&package=bank_type=WX&body=%E6%94%AF%E4%BB%98%E6%B5%8B%E8%AF%95&fee_type=1&input_charset=UTF-8&notify_url=http%3A%2F%2Fweixin.qq.com&out_trade_no=7240b65810859cbf2a8d9f76a638c0a3&partner=1900000109&spbill_create_ip=196.168.1.1&total_fee=1&sign=7F77B507B755B3262884291517E380F8&timestamp=1399514976&traceid=test_1399514976", data.toStrint1())
	assert.Equal(t, "8893870b9004ead28691b60db97a8d2c80dbfdc6", data.toAppSignature())
}

func TestGetPrePay(t *testing.T) {
	// partnerKey := "8934e7d15453e97507ef794cf7b0519d"
	prePayData := prePayPostData{

		AppId:     "wxd930ea5d5a258f4f",
		appKey:    "L8LrMqqeGRxST5reouB0K66CaYAWpqhAVsq7ggKkxHCOastWksvuX1uvmvQclxaHoYd3ElNBrNO2DHnnzgfVG9Qs473M3DTOZug5er46FhuGofumV8H2FVR9qkjSlC5K",
		TradeId:   "test_1399514976",
		Noncestr:  "e7d161ac8d8a76529d39d9f5b4249ccb",
		Timestamp: 1399514976,
		Package:   "Sign=WXpay",
		PartnerId: "1900000109",
		Prepayid:  "1101000000140429eb40476f8896f4c9",
	}
	// 这个测试不过 ，可能是因为  实际并不一定等于7ffecb600d7157c5aa49810d2d8f28bc2811827b
	assert.Equal(t, "7ffecb600d7157c5aa49810d2d8f28bc2811827b", prePayData.getSign())
	// sdk := NewWeiXinSDK(prePayData.AppId, prePayData.appKey, "", prePayData.PartnerId, partnerKey, "", "http://")
	// 	jsonData, ok := sdk.getPrePay(prePayData.TradeId, prePayData.Noncestr, prePayData.Prepayid, time.Unix(int64(1399514976), 0))
	// 	assert.True(t, ok)
	// 	assert.Equal(t, `{
	// "appid":"wxd930ea5d5a258f4f", "noncestr":"e7d161ac8d8a76529d39d9f5b4249ccb", "package":"Sign=WXpay";
	// "partnerid":"1900000109" "prepayid":"1101000000140429eb40476f8896f4c9", "sign":"7ffecb600d7157c5aa49810d2d8f28bc2811827b", "timestamp":"1399514976"
	// }`, jsonData)
	// 	fmt.Println(jsonData)
}
