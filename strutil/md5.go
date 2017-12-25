package strutil

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

//Sha1 gen hex sha1 of content
func Sha1(content string) string {
	hash := sha1.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

//MD5 gen hex md5 of content
func MD5(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

//HmacSha1 gen HMAC-SHA1 of content with key
func HmacSha1(content, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(content))
	hash := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash[:])
}
