package goauth

import (
	"encoding/json"
	"github/0studio/goauth/utils"
	"time"
)

func DoPPAuth(token string, now time.Time) (status int32, AccountId string, AccountName string, sdkRec interface{}) {
	jsonBytes, err := getLoginResponse(token, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR, "", "", nil
	}
	// var kv = make(map[string]interface{})
	var loginInfo = &PPLoginInfo{Status: PP_STATUS_FAIL}
	// 默认 25PP返回的 json格式不全  ，需要加前后{}
	jsonFormat := "{" + string(jsonBytes) + "}"
	json.Unmarshal([]byte(jsonFormat), loginInfo)
	statusStr := loginInfo.Status
	if statusStr == PP_STATUS_SUCC {
		status = PB_STATUS_SUCC
		AccountId = utils.Int2Str(loginInfo.Userid)
		AccountName = loginInfo.Username
		return
	}
	return PB_ERRNO_AUTH_ERROR, "", "", nil
}

// 25PP平台
func getLoginResponse(token string, now time.Time) (json []byte, err error) {
	return utils.PostHttpResponse("http://passport_i.25pp.com:8080/index?tunnel-command=2852126756", []byte(token), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

const (
	PP_STATUS_SUCC = 0
	PP_STATUS_FAIL = 1
)

type PPLoginInfo struct {
	Status   int32
	Userid   int `json:"userid,omitempty"`
	Username string
}
