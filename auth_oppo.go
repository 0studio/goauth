package goauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/0studio/goauth/utils"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
oppo
游戏名称：别挡我超神
游戏ID：1690
游戏key：22j8XTPa3iv4wgsCoCK0wwgwo
游戏secret：187e6ec09c058306fAD90D30F3fd3384
*/

func DoOPPOAuth(appId, appKey, appSecret string, token string, now time.Time) (status int32, AccountId string, AccountName string, sdkRec interface{}) {
	oauthToken := strings.Split(strings.Split(token, "&")[0], "=")[1]
	oauthTokenSecret := strings.Split(strings.Split(token, "&")[1], "=")[1]
	oauthNonce := strconv.FormatInt(rand.Int63(), 10)
	oauthSignMethod := "HMAC-SHA1"
	oauthTime := strconv.FormatInt(time.Now().Unix(), 10)
	oauthVersion := "1.0"
	key := appSecret + "&" + oauthTokenSecret
	urlParamString := fmt.Sprintf("oauth_consumer_key=%s&oauth_nonce=%s&oauth_signature_method=%s&oauth_timestamp=%s&oauth_token=%s&oauth_version=%s",
		appKey, oauthNonce, oauthSignMethod, oauthTime, oauthToken, oauthVersion)
	baseString := "POST&" + utils.UrlEncode("http://thapi.nearme.com.cn/account/GetUserInfoByGame") + "&" + utils.UrlEncode(urlParamString)
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(baseString))
	oauthSignature := mac.Sum(nil)
	value := url.Values{}
	value.Set("oauth_consumer_key", appKey)
	value.Set("oauth_nonce", oauthNonce)
	value.Set("oauth_signature_method", oauthSignMethod)
	value.Set("oauth_timestamp", oauthTime)
	value.Set("oauth_version", oauthVersion)
	value.Set("oauth_signature", string(oauthSignature))
	value.Set("oauth_token", oauthToken)

	jsonBytes, err := getOPPOLoginResponse(value, now)
	if err != nil {
		return PB_ERRNO_AUTH_ERROR, "", "", nil
	}
	loginInfo := OPPOLoginResp{}
	json.Unmarshal(jsonBytes, &loginInfo)
	fmt.Println(loginInfo, "oppop auth ............")
	if loginInfo.BriefUser.Id != "" {
		status = PB_STATUS_SUCC
		AccountId = loginInfo.BriefUser.Id
		AccountName = loginInfo.BriefUser.UserName
		return
	}
	return PB_ERRNO_AUTH_ERROR, "", "", nil
}

// OPPO平台
func getOPPOLoginResponse(v url.Values, now time.Time) (json []byte, err error) {
	return utils.PostFormHttpResponse("http://thapi.nearme.com.cn/account/GetUserInfoByGame", v, now, DEFAULT_AUTH_HTTP_REQUEST_TIMEOUT)
}

type OPPOLoginResp struct {
	BriefUser BriefUser `json:"BriefUser"`
}

type BriefUser struct {
	Id                string `json:"id"`
	Sex               string `json:"sex"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	Name              string `json:"name"`
	UserName          string `json:"userName"`
	EmailStatus       string `json:"emailStatus"`
	MobileStatus      string `json:"mobileStatus"`
	Status            string `json:"status"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
	GameBalance       string `json:"gameBalance"`
}
