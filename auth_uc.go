package goauth

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goutils"
	"time"
)

type UCData struct {
	Sid string `json:"sid"`
}
type UCGame struct {
	GameId int `json:"gameId"`
}

type UCSend struct {
	Id   int64  `json:"id"`
	Data UCData `json:"data"`
	Game UCGame `json:"game"`
	Sign string `json:"sign"`
}

type UCRecv struct {
	Id    int64      `json:"id"`
	State UCState    `json:"state"`
	Data  UCDataRecv `json:"data"`
}

type UCState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (s UCState) IsSucc() bool {
	// 1 succ,10 参数错，11 未登录，99 sdk server error
	return s.Code == 1
}

type UCDataRecv struct {
	AccountId string `json:"accountId"`
	Creator   string `json:"creator"` // JY:九游，PP
	NickName  string `json:"nickname"`
}

func DoUCAuth(isSandbox bool, gameid int, appKey, sid string, now time.Time) (ucRecv UCRecv) {
	id := now.Unix()
	// service := "ucid.user.sidInfo"
	ucData := UCData{sid}
	ucGame := UCGame{gameid}

	sign := goutils.GetHexMd5(fmt.Sprintf("sid=%s%s", sid, appKey))

	ucSend := UCSend{id, ucData, ucGame, sign}
	content, _ := json.Marshal(ucSend)
	jsonBytes, err := getUCLoginResponse(isSandbox, content, now)
	if err != nil {
		return
	}
	loginInfo := UCRecv{}
	json.Unmarshal(jsonBytes, &loginInfo)
	return
}

const (
	UC_USER_VERIFY_URL         = "http://sdk.g.uc.cn/cp/account.verifySession"
	UC_USER_VERIFY_URL_SANDBOX = "http://sdk.test4.g.uc.cn/cp/account.verifySession"
)

// UC平台
func getUCLoginResponse(isSandbox bool, content []byte, now time.Time) (json []byte, err error) {
	if isSandbox {
		return goutils.PostHttpResponse(UC_USER_VERIFY_URL_SANDBOX, content, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
	}
	return goutils.PostHttpResponse(UC_USER_VERIFY_URL, content, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}
