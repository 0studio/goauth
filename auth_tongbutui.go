package goauth

// 同步推
import (
	"fmt"
	"github.com/0studio/goauth/utils"
	log "github.com/cihub/seelog"
	"net/url"
	"time"
)

// 直接返回1,0,-1
//1: Token 有效
//0:失效
//-1:格式有错 的字符串

//AppId：140616
//AppKey：coe2LYNlxI@Uh5ubRodLAYl9xIUsh5Eb

func DoTongbutuiAuth(token string, now time.Time) int32 {
	RetResponse := getTongbutuiLoginResponse(token, now)
	if RetResponse == "1" {
		return PB_STATUS_SUCC
	}
	return PB_ERRNO_AUTH_ERROR
}
func getTongbutuiLoginResponse(token string, now time.Time) string {
	contentByte, err := utils.GetHttpResponseAsJson(getTongbutuiLoginUrl(token), now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	if err != nil {
		log.Error("auth_tongbutui_error", err)
		return "false"
	}

	return string(contentByte)
}
func getTongbutuiLoginUrl(token string) string {
	urlStr := fmt.Sprintf("http://tgi.tongbu.com/check.aspx?k=%s", url.QueryEscape(token))
	return urlStr
}
