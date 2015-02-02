package goauth

import (
	"encoding/json"
	"fmt"
	"github/0studio/goauth/utils"
	"time"
)

// CONST_UC_APPKEY    = ""
// CONST_UC_APPSECRET = ""
// const_UC_CPID      = 20087
// const_UC_GAMEID    = 119474
// const_UC_SERVERID  = 1333

type UCData struct {
	Sid string `json:"sid"`
}
type UCGame struct {
	CpId      int    `json:"cpId"`
	GameId    int    `json:"gameId"`
	ChannelId string `json:"channelId"`
	ServerId  int    `json:"serverId"`
}

type UCSend struct {
	Id      int    `json:"id"`
	Service string `json:"service"`
	Data    UCData `json:"data"`
	Game    UCGame `json:"game"`
	Sign    string `json:"sign"`
}

type UCRecv struct {
	Id    int        `json:"id"`
	State UCState    `json:"state"`
	Data  UCDataRecv `json:"data"`
}

type UCState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type UCDataRecv struct {
	UCid     int    `json:"ucid"`
	NickName string `json:"nickname"`
}

func DoUCAuth(appKey, appSecret string, cpId, gameid, serverid int, sid string, now time.Time) (status int32) {
	id := int(now.Unix())
	service := "ucid.user.sidInfo"
	ucData := UCData{sid}
	ucGame := UCGame{cpId, gameid, "2", serverid}

	sign := utils.GetHexMd5(fmt.Sprintf("%s%s=%s%s", cpId, "sid", sid, appKey))

	ucSend := UCSend{id, service, ucData, ucGame, sign}
	content, _ := json.Marshal(ucSend)
	jsonBytes, err := getUCLoginResponse(content, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR
	}
	loginInfo := UCRecv{}
	json.Unmarshal(jsonBytes, &loginInfo)
	if loginInfo.State.Code == 1 {
		status = PB_STATUS_SUCC
		return
	}
	return PB_ERRNO_AUTH_ERROR
}

// UC平台
func getUCLoginResponse(content []byte, now time.Time) (json []byte, err error) {
	return utils.PostHttpResponse("http://sdk.g.uc.cn/ss/", content, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}
