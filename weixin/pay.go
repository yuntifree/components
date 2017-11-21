package weixin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

const (
	orderURL  = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	packValue = "Sign=WXPay"
	succCode  = "SUCCESS"
)

//WxPay for weixin pay
type WxPay struct {
	MerID  string
	MerKey string
}

//UnifyOrderReq unify order request
type UnifyOrderReq struct {
	Appid          string `xml:"appid"`
	Body           string `xml:"body"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	NotifyURL      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	TotalFee       int64  `xml:"total_fee"`
	OutTradeNO     string `xml:"out_trade_no"`
	Sign           string `xml:"sign"`
	Openid         string `xml:"openid"`
}

//UnifyOrderResp unify order response
type UnifyOrderResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Openid     string `xml:"openid"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	TradeType  string `xml:"trade_type"`
	PrepayID   string `xml:"prepay_id"`
}

//SimpleResponse simple response to verify return_code
type SimpleResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

//NotifyRequest notify request
type NotifyRequest struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	Appid         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Openid        string `xml:"openid"`
	Sign          string `xml:"sign"`
	ResultCode    string `xml:"result_code"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	TotalFee      int64  `xml:"total_fee"`
	CashFee       int64  `xml:"cash_fee"`
	TransactionID string `xml:"transaction_id"`
	OutTradeNO    string `xml:"out_trade_no"`
	TimeEnd       string `xml:"time_end"`
	FeeType       string `xml:"fee_type"`
	IsSubscribe   string `xml:"is_subscribe"`
}

func checkRsp(body io.ReadCloser) bool {
	defer body.Close()
	rspbody, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("ReadAll resp failed:%v", err)
		return false
	}
	var rs SimpleResponse
	dec := xml.NewDecoder(bytes.NewReader(rspbody))
	err = dec.Decode(&rs)
	if err != nil {
		log.Printf("decode failed:%s %v", rspbody, err)
		return false
	}
	if rs.ReturnCode != succCode {
		log.Printf("fail response:%s", rspbody)
		return false
	}
	return true
}

//CalcSign calculate sign for map
func (p *WxPay) CalcSign(mReq map[string]interface{}) string {
	var sortedKeys []string
	for k := range mReq {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var signStr string
	for _, k := range sortedKeys {
		log.Printf("%v -- %v", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" && k != "sign" {
			signStr += k + "=" + value + "&"
		}
	}

	if p.MerKey != "" {
		signStr += "key=" + p.MerKey
	}

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStr))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

func (p *WxPay) calcReqSign(req UnifyOrderReq) string {
	m := make(map[string]interface{})
	m["appid"] = req.Appid
	m["body"] = req.Body
	m["mch_id"] = req.MchID
	m["notify_url"] = req.NotifyURL
	m["trade_type"] = req.TradeType
	m["spbill_create_ip"] = req.SpbillCreateIP
	m["total_fee"] = req.TotalFee
	m["out_trade_no"] = req.OutTradeNO
	m["nonce_str"] = req.NonceStr
	m["openid"] = req.Openid
	return p.CalcSign(m)
}

//UnifyPayRequest send unify order pay request
func (p *WxPay) UnifyPayRequest(req UnifyOrderReq) (*UnifyOrderResp, error) {
	req.Sign = p.calcReqSign(req)

	buf, err := xml.Marshal(req)
	if err != nil {
		log.Printf("UnifyPayRequest marshal failed:%v", err)
		return nil, err
	}

	reqStr := string(buf)
	reqStr = strings.Replace(reqStr, "XUnifyOrderReq", "xml", -1)
	log.Printf("reqStr:%s", reqStr)

	request, err := http.NewRequest("POST", orderURL, bytes.NewReader([]byte(reqStr)))
	if err != nil {
		log.Printf("UnifyPayRequest NewRequest failed:%v", err)
		return nil, err
	}
	request.Header.Set("Accept", "application/xml")
	request.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, err := c.Do(request)
	if err != nil {
		log.Printf("UnifyPayRequest request failed:%v", err)
		return nil, err
	}

	defer resp.Body.Close()
	dec := xml.NewDecoder(resp.Body)
	var res UnifyOrderResp
	err = dec.Decode(&res)
	if err != nil {
		log.Printf("UnifyPayRequest Unmarshal failed:%v", err)
		return nil, err
	}
	return &res, nil
}

//VerifyNotify verify notify sign
func (p *WxPay) VerifyNotify(req NotifyRequest) bool {
	vt := reflect.TypeOf(req)
	vv := reflect.ValueOf(req)
	m := make(map[string]interface{})

	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		name := f.Tag.Get("xml")
		m[name] = vv.FieldByName(f.Name).Interface()
		log.Printf("name:%s value:%v", name, vv.FieldByName(f.Name).Interface())
	}
	sign := p.CalcSign(m)
	if req.Sign != sign {
		return false
	}
	return true
}
