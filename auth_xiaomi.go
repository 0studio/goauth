package goauth

//小米
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/0studio/goauth/utils"
	"strings"
	"time"
)

type XiaomiRecv struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errMsg"`
}

func DoXiaomiAuth(appId, appKey, appSecret string, sessionId string, AccountId string, now time.Time) (status int32) {
	mac := hmac.New(sha1.New, []byte(appSecret))
	sign := "appId=" + appId + "&session=" + sessionId + "&uid=" + AccountId
	mac.Write([]byte(sign))
	signature := mac.Sum(nil)
	c := hex.Dump(signature)
	d := strings.SplitN(c, "  ", -1)
	sig := strings.Replace(d[1]+d[2]+d[4], " ", "", -1)
	urlStr := fmt.Sprintf("http://mis.migc.xiaomi.com/api/biz/service/verifySession.do?appId=%s&session=%s&uid=%s&signature=%s", appId, sessionId, AccountId, sig)
	jsonBytes, err := getDangleLoginResponse(urlStr, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	kv := XiaomiRecv{}
	json.Unmarshal(jsonBytes, &kv)
	if kv.Code != 200 {
		return PB_ERRNO_AUTH_ERROR
	} else {
		return PB_STATUS_SUCC
	}

}

func getXiaomiLoginResponse(urlStr string, now time.Time) (json []byte, err error) {
	return utils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}
