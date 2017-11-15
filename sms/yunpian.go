package sms

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	simplejson "github.com/bitly/go-simplejson"
)

const (
	tplURL = "https://sms.yunpian.com/v2/sms/tpl_single_send.json"
)

//Yunpian yunpian sms implement
type Yunpian struct {
	Apikey string
	TplID  int
}

//Send send yunpian sms
func (y *Yunpian) Send(phone string, code int) int {
	tplValue := url.Values{"#code#": {fmt.Sprintf("%06d", code)}}.Encode()
	data := url.Values{"apikey": {y.Apikey}, "mobile": {phone},
		"tpl_id": {fmt.Sprintf("%d", y.TplID)}, "tpl_value": {tplValue}}
	err := sendData(data)
	if err != nil {
		return -1
	}
	return 0
}

func sendData(data url.Values) error {
	resp, err := http.PostForm(tplURL, data)
	if err != nil {
		log.Printf("SendYPSMS request failed:%v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("SendYPSMS read response failed:%v", err)
		return err
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Printf("SendYPSMS parse response failed:%v", err)
		return err
	}
	code, err := js.Get("code").Int()
	if err != nil {
		log.Printf("SendYPSMS get response code failed:%v", err)
		return err
	}
	if code != 0 {
		log.Printf("SendYPSMS illegal code:%s", string(body))
		return fmt.Errorf("illegal response code:%d", code)
	}
	return nil
}
