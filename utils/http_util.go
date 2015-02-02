package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
func GetHttpResponseAsJson(urlStr string, now time.Time, timeout int) (data []byte, err error) {
	client := HttpWithTimeOut(now, timeout)
	response, err := client.Get(urlStr)
	if err != nil {
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return body, nil
}

func UrlEncode(urlStr string) string {
	return url.QueryEscape(urlStr)
}

func GetHexMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))

}

// func postHttpResponseAsJson(urlStr string, values url.Values) (data []byte, err error) {
// 	response, err := http.PostForm(urlStr, values)
// 	if err != nil {
// 		return
// 	}

// 	defer response.Body.Close()
// 	body, _ := ioutil.ReadAll(response.Body)
// 	return body, nil

// }
func PostHttpResponse(urlStr string, content []byte, now time.Time, timeout int) (data []byte, err error) {
	client := HttpWithTimeOut(now, timeout)
	response, err := client.Post(urlStr, "text/html", bytes.NewReader(content))
	if err != nil {
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return body, nil

}

func PostFormHttpResponse(urlStr string, v url.Values, now time.Time, timeout int) (data []byte, err error) {
	client := HttpWithTimeOut(now, timeout)
	respose, err := client.PostForm(urlStr, v)
	if err != nil {
		return
	}

	defer respose.Body.Close()
	body, _ := ioutil.ReadAll(respose.Body)
	return body, nil
}

func HttpWithTimeOut(now time.Time, timeoutMillSeconds int) http.Client {
	timeoutDur := time.Millisecond * time.Duration(timeoutMillSeconds)
	// 在拨号回调中，使用DialTimeout来支持连接超时，当连接成功后，利用SetDeadline来让连接支持读写超时。
	fun := func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeoutDur)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(now.Add(timeoutDur))
		return conn, nil
	}
	transport := &http.Transport{Dial: fun, ResponseHeaderTimeout: timeoutDur}

	client := http.Client{
		Transport: transport,
	}
	return client
}
