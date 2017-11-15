package sms

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/yuntifree/components/httputil"
)

const (
	smsurl = "https://yun.tim.qq.com/v3/tlssmssvr/sendsms"
)

//Qcloud qcloud sms implement
type Qcloud struct {
	Appid  string
	Appkey string
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (q *Qcloud) genBody(phone string, code int) string {
	js, err := simplejson.NewJson([]byte(`{"tel":{"nationcode":"86"}, "type":"0","ext":"","extend":""}`))
	if err != nil {
		return ""
	}
	s := fmt.Sprintf("%06d", code)
	msg := "【东莞无线】欢迎使用东莞无线免费WiFi,您的验证码为:" + s
	js.Set("msg", msg)
	sig := getMD5Hash(q.Appkey + phone)
	js.Set("sig", sig)
	js.SetPath([]string{"tel", "phone"}, phone)
	data, err := js.Encode()
	if err != nil {
		return ""
	}

	return string(data[:])
}

//Send send verify code to phone
func (q *Qcloud) Send(phone string, code int) int {
	body := q.genBody(phone, code)
	fmt.Println(body)
	rand.Seed(42)
	url := smsurl + "?sdkappid=" + q.Appid + "&random=" + strconv.Itoa(rand.Int())
	fmt.Println(url)
	rspbody, err := httputil.Request(url, body)
	if err != nil {
		return -1
	}
	fmt.Println(string(rspbody))
	js, err := simplejson.NewJson([]byte(`{}`))
	err = js.UnmarshalJSON([]byte(rspbody))
	s, err := js.GetPath("result").String()
	if s != "0" {
		return -3
	}

	return 0
}
