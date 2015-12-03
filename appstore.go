package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"strings"
	"time"
)

const (
	// 0普通定单，1沙盒测试定单
	PAY_ORDER_TYPE_COMMON  = 0
	PAY_ORDER_TYPE_SANDBOX = 1
)

const (
	PAY_APP_STORE_TIMEOUT = 10 * 1000 // 10s

)

func getAppStoreResponseTry5Times(url string, now time.Time, content []byte) (response []byte, err error) {
	for i := 0; i < 5; i++ {
		response, err = goutils.PostHttpResponse(url, content, now, PAY_APP_STORE_TIMEOUT)
		if err == nil {
			return
		}
	}
	return
}

type AppStorePayResponse struct {
	Status    int32
	OrderType int32            // 0普通定单，1沙盒测试定单
	Info      AppStoreRecvData `json:"receipt,omitempty"`
}

func (res AppStorePayResponse) IsSandbox() bool {
	return res.OrderType == PAY_ORDER_TYPE_SANDBOX
}
func (res AppStorePayResponse) IsSucc() bool {
	return res.Status == 0

}

type AppStoreRecvData struct {
	TransactionId string `json:"transaction_id,omitempty"`
	ProductId     string `json:"product_id,omitempty"`
}

// """
// {"status":21002, "exception":"java.lang.ClassCastException"}
//
// {"receipt":{"original_purchase_date_pst":"2012-11-12 03:10:48 America/Los_Angeles",
// "purchase_date_ms":"1352718648392",
// "unique_identifier":"98de2b8cb8b973773538c5c8743e1043677b9201",
// "original_transaction_id":"1000000058437063",
// "bvrs":"50000",
// "transaction_id":"1000000058437063",
// "quantity":"1",
// "unique_vendor_identifier":"FB22FBAA-CFC9-4ABD-9522-BCECF78C2866",
// "item_id":"577373064",
// "product_id":"com.yyshtech.gold_6",
// "purchase_date":"2012-11-12 11:10:48 Etc/GMT",
// "original_purchase_date":"2012-11-12 11:10:48 Etc/GMT",
// "purchase_date_pst":"2012-11-12 03:10:48 America/Los_Angeles",
// "bid":"com.yyshtech.zhajinhua", "original_purchase_date_ms":"1352718648392"}, "status":0}
// product_id: com.yysh.goldcount.6

// How do I verify my receipt (iOS)?
// Always verify your receipt first with the production URL; proceed to verify with the sandbox URL if you receive a 21007 status code. Following this approach ensures that you do not have to switch

func GetAppStoreResponse(content []byte, now time.Time) (jsonData AppStorePayResponse, err error) {
	response, err := getAppStoreResponseTry5Times("https://buy.itunes.apple.com/verifyReceipt", now, content)
	// response, err := getAppStoreResponseTry5Times("https://sandbox.itunes.apple.com/verifyReceipt", content)

	if err != nil {
		return
	}
	jsonData = AppStorePayResponse{}
	json.Unmarshal(response, &jsonData)
	jsonData.OrderType = PAY_ORDER_TYPE_COMMON
	if jsonData.Status == 21007 {
		jsonData.OrderType = PAY_ORDER_TYPE_SANDBOX
		response, err = getAppStoreResponseTry5Times("https://sandbox.itunes.apple.com/verifyReceipt", now, content)
		if err != nil {
			return
		}
		json.Unmarshal(response, &jsonData)

	}

	return
}

type ClientAppStoreSend struct {
	TransactionId string `json:"transaction_id,omitempty"`
	ReceiptData   string `json:"receipt-data,omitempty"`
}

func AppStorePay(content string, now time.Time) (jsonData AppStorePayResponse) {
	var err error
	send := ClientAppStoreSend{}
	err = json.Unmarshal([]byte(content), &send)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonData.Info.TransactionId = send.TransactionId
	if send.ReceiptData == "" {
		return
	}

	if len(send.TransactionId) >= 19 || len(send.TransactionId) <= 8 {
		return
	}
	if strings.Contains(send.TransactionId, "-") {
		return

	}
	jsonData, err = GetAppStoreResponse([]byte(content), now)
	if err != nil {
		jsonData.Info.TransactionId = send.TransactionId
		jsonData.Status = 1 // err
	}

	return
}
