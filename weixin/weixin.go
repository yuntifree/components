package weixin

import (
	"fmt"
	"net/url"
)

const (
	wxAuthURL = "https://open.weixin.qq.com/connect/oauth2/authorize"
)

//WxInfo weixin info
type WxInfo struct {
	Appid  string
	Appkey string
}

//GenRedirect generate redirect url
func (w WxInfo) GenRedirect(redirect string) string {
	return fmt.Sprintf(`%s?appid=%s&redirect_uri=%s&response_type=code
	&scope=snsapi_userinfo&state=list#wechat_redirect`, wxAuthURL,
		w.Appid, url.QueryEscape(redirect))
}
