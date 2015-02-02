package goauth

//百度
import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goauth/utils"
	"strings"
	"time"
)

const ( //运营给的参数
// appId     = "3404678"
// appKey    = "b9mstdOGun5je3b5OwwXPGDi"
// appSecret = "gjqiYMGstfvsS2vVV44ItCvYY9m1G8ZT"
)

type SdkRecv struct {
	Code    string `json:"error_code"`
	Message string `json:"error_msg"`
}

func DoBaiduAuth(appId, appKey, appSecret string, sessionId string, AccountId string, now time.Time) (status int32) {
	//urlStr := "http://sdk.m.duoku.com/openapi/sdk/checksession?"
	uid := AccountId
	clientSecret := strings.ToLower(utils.GetHexMd5(fmt.Sprintf("%s%s%s%s%s", appId, appKey, uid, sessionId)))
	urlStr := fmt.Sprintf("http://sdk.m.duoku.com/openapi/sdk/checksession?appid=%s&appkey=%s&uid=%s&sessionid=%s&clientsecret=%s", appId, appKey, uid, sessionId, clientSecret)
	jsonBytes, err := getBaiduLoginResponse(urlStr, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	kv := SdkRecv{}
	json.Unmarshal(jsonBytes, &kv)
	if kv.Code == BAIDU_STATUS_SUCC {
		status = PB_STATUS_SUCC
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

const (
	BAIDU_STATUS_SUCC = "0"
	BAIDU_STATUS_FAIL = "1"
)

func getBaiduLoginResponse(urlStr string, now time.Time) (json []byte, err error) {
	return utils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}
