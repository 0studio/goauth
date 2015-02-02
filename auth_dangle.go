package goauth

//当乐
import (
	"encoding/json"
	"fmt"
	//"strings"
	"github.com/0studio/goauth/utils"
	"time"
)

type DangleFailRecv struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

type DangleSuccRecv struct {
	MemberId   int    `json:"memberId"`
	UserName   string `json:"username"`
	NickName   string `json:"nickname"`
	Gender     string `json:"gender"`
	Level      int    `json:"level"`
	AvatarUrl  string `json:"avatar_url"`
	CreateData int    `json:"created_date"`
	Token      string `json:"token"`
	ErrorCode  int    `json:"error_code"`
}

func DoDangleAuth(appId, appKey, appSecret string, sessionId string, AccountId string, now time.Time) (status int32) {
	//urlStr := "http://sdk.m.duoku.com/openapi/sdk/checksession?"
	mid := AccountId
	token := sessionId
	sig := utils.GetHexMd5(fmt.Sprintf("%s|%s", token, appKey))
	urlStr := fmt.Sprintf("http://connect.d.cn/open/member/info/?app_id=%s&mid=%s&token=%s&sig=%s", appId, mid, token, sig)
	jsonBytes, err := getDangleLoginResponse(urlStr, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	kv := DangleSuccRecv{}
	json.Unmarshal(jsonBytes, &kv)
	if kv.ErrorCode != 0 {
		return PB_ERRNO_AUTH_ERROR
	} else {
		return PB_STATUS_SUCC
	}
	/*if kv.Code == BAIDU_STATUS_SUCC {
		status = PB_STATUS_SUCC
		return
	}*/

}

func getDangleLoginResponse(urlStr string, now time.Time) (json []byte, err error) {
	return utils.GetHttpResponseAsJson(urlStr, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}
