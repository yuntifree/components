package weixin

import (
	"fmt"
	"log"
	"net/url"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/yuntifree/components/httputil"
)

const (
	wxAuthURL  = "https://open.weixin.qq.com/connect/oauth2/authorize"
	wxTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
)

//WxInfo weixin info
type WxInfo struct {
	Appid  string
	Appkey string
}

//UserInfo user info
type UserInfo struct {
	Openid string
	Token  string
}

//GenRedirect generate redirect url
func (w WxInfo) GenRedirect(redirect string) string {
	return fmt.Sprintf(`%s?appid=%s&redirect_uri=%s&response_type=code
	&scope=snsapi_userinfo&state=list#wechat_redirect`, wxAuthURL,
		w.Appid, url.QueryEscape(redirect))
}

//GetCodeToken use code to get user info
func (w *WxInfo) GetCodeToken(code string) (*UserInfo, error) {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wxTokenURL, w.Appid, w.Appkey, code)
	log.Printf("url:%s", url)
	res, err := httputil.Request(url, "")
	if err != nil {
		log.Printf("fetch url %s failed:%v", url, err)
		return nil, err
	}

	log.Printf("GetCodeToken resp:%s", res)
	js, err := simplejson.NewJson([]byte(res))
	if err != nil {
		log.Printf("parse resp failed:%v", err)
		return nil, err
	}

	openid, err := js.Get("openid").String()
	if err != nil {
		log.Printf("get openid failed:%v", err)
		return nil, err
	}

	token, err := js.Get("access_token").String()
	if err != nil {
		log.Printf("get access_token failed:%v", err)
		return nil, err
	}
	log.Printf("openid:%s token:%s", openid, token)

	var info UserInfo
	info.Openid = openid
	info.Token = token
	return &info, nil
}
