package goauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 安智
type ANZHILoginResp struct {
	SC   string `json:"sc"`
	ST   string `json:"st"`
	Time string `json:"time"`
	Msg  string `json:"msg"`
}

type ANZHIMsg struct {
	Uid string `jsont:"uid"`
}

func DoAnzhiAuth(appKey, appSecret string, token string, now time.Time) (status int32) {
	sid := token
	time := timeFormat()
	sign := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s%s", appKey, sid, appSecret)))
	value := url.Values{}
	value.Set("time", time)
	value.Set("appKey", appKey)
	value.Set("sid", sid)
	value.Set("sign", sign)
	jsonBytes, err := getANZHILoginResponse(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := ANZHILoginResp{}
	respBytes := []byte(strings.Replace(string(jsonBytes), "'", "\"", -1)) //返回的消息含',要replace成"
	err = json.Unmarshal(respBytes, &loginInfo)
	if loginInfo.SC == "1" {
		status = PB_STATUS_SUCC
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

// 安智平台
func getANZHILoginResponse(v url.Values, now time.Time) (json []byte, err error) {
	return goutils.PostFormHttpResponse("http://user.anzhi.com/web/api/sdk/third/1/queryislogin", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

func timeFormat() string {
	StampNano := "Jan _2 15:04:05.000"
	t := time.Now()
	year, month, day := t.Date()
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	nanoSecond := strings.SplitN(t.Format(StampNano), ".", -1)[1]
	return fmt.Sprintf("%s%s%s%s%s%s%s", format(year), format(int(month)), format(day), format(hour), format(minute), format(second), nanoSecond)
}

func format(i int) string {
	if i < 10 {
		return "0" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}
