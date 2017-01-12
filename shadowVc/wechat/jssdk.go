package wechat

import (
	"com/shadowVc"
	"github.com/pquerna/ffjson/ffjson"
	"net/url"
)

type RespJssdkGet struct {
	WechatJssdkGetResponse *struct {
		AppId     string `json:"appId"`
		Timestamp string `json:"timestamp"`
		NonceStr  string `json:"nonceStr"`
		Signature string `json:"signature"`
	}
}

func GetJssdk(shadowvcSdk *shadowVc.ShadowVcSdk, urlStr string) (*RespJssdkGet, error) {
	postData := make(url.Values)
	postData.Add("url", urlStr)
	bt, err := shadowvcSdk.Send("sd.wx.jssdk.get", postData)
	if err != nil {
		return nil, err
	}
	respJssdkGet := new(RespJssdkGet)
	err = ffjson.Unmarshal(bt, respJssdkGet)
	if err != nil {
		return nil, err
	}
	if respJssdkGet.WechatJssdkGetResponse == nil {
		return nil, shadowVc.NewRespError(bt)
	}
	return respJssdkGet, nil
}
