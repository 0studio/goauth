package utils

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
)

// 1、本demo签名验签使用的rsa私钥格式为pkcs8，因为爱贝提供的rsa私钥为pkcs1，所以需要使用cptools工具，将pkcs1私钥转为pkcs8格式
// var privateKey string = "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAKz0WssMzD9pwfHlEPy8+NFSnsX+CeZoogRyrzAdBkILTVCukOfJeaqS07GSpVgtSk9PcFk3LqY59znddga6Kf6HA6Tpr19T3Os1U3zNeU79X/nT6haw9T4nwRDptWQdSBZmWDkY9wvA28oB3tYSULxlN/S1CEXMjmtpqNw4asHBAgMBAAECgYBzNFj+A8pROxrrC9Ai6aU7mTMVY0Ao7+1r1RCIlezDNVAMvBrdqkCWtDK6h5oHgDONXLbTVoSGSPo62x9xH7Q0NDOn8/bhuK90pVxKzCCI5v6haAg44uqbpt7fZXTNEsnveXlSeAviEKOwLkvyLeFxwTZe3NQJH8K4OqQ1KzxK+QJBANmXzpVdDZp0nAOR34BQWXHHG5aPIP3//lnYCELJUXNB2/JYTN57dv5LlE5/Ckg0Bgak764A/CX62bKhe/b+FMsCQQDLe4F2qHGy7Sa81xatm66mEkG3u88g9qRARdEvgx9SW+F1xBt2k/bU2YI31hB8IYXzL8KW9NzDfQPihBBUFn4jAkEAzbrmq/pLPlo6mHV3qE5QA2+J+hRh0UYVKsVDKkJGLH98gepS45hArbawBne/NP1bJTUVGKP9w7sl0es01hbteQJATzLO/QQq3N15Cl8dMI07uN+6PG0Y/VeCLpH+DWQXuNKSOmgN2GVW2RmfmWP0Hpxdqn2YW3EKy/vIm02TnWbzyQJAXwujUR9u9s8BZI33kw3gQ7bvWVYt8yyiYzWD2Qrnyg08tN5o+JsjW3fEDWHm70jjZIc+l/5FaZ7H5NOYpnVcpA=="
// 签名
func RSASignWithMD5(privateKey string, content string) (sign string) {
	hash := md5.New()
	io.WriteString(hash, string(content))
	hashed := hash.Sum(nil)

	base64Data, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		fmt.Println("base64Decode_error", err)
		return
	}

	pk, err := x509.ParsePKCS8PrivateKey([]byte(base64Data))
	if err != nil {
		fmt.Println("read private key error", err)
		return
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, pk.(*rsa.PrivateKey), crypto.MD5, hashed)
	if err != nil {
		fmt.Println("rsa_sign_error", err)
		return
	}
	sign = base64.StdEncoding.EncodeToString(signature)
	return
}

// // 公钥
// var pubKey string = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCs9FrLDMw/acHx5RD8vPjRUp7F/gnmaKIEcq8wHQZCC01QrpDnyXmqktOxkqVYLUpPT3BZNy6mOfc53XYGuin+hwOk6a9fU9zrNVN8zXlO/V/50+oWsPU+J8EQ6bVkHUgWZlg5GPcLwNvKAd7WElC8ZTf0tQhFzI5raajcOGrBwQIDAQAB"
func VerifyRSASignWithMD5(publicKey string, content, sign string) bool {
	pubKeyData, err := base64.StdEncoding.DecodeString(publicKey)
	pk, err := x509.ParsePKIXPublicKey(pubKeyData)
	if err != nil {
		fmt.Println("read_public_key_error", err)
		return false
	}

	signData, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		fmt.Println("decode_base64_sign_error", err)

	}
	// var h crypto.Hash
	hash := md5.New()
	io.WriteString(hash, string(content))
	hashed := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(pk.(*rsa.PublicKey), crypto.MD5, hashed, signData)
	if err == nil {
		return true
	}
	fmt.Println("verify_sign_error", err)
	return false

}
