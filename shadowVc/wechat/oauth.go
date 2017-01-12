package wechat

import (
	"com/shadowVc"
	"github.com/pquerna/ffjson/ffjson"
	"net/url"
	"strconv"
)

type RespOauthUrl struct {
	WechatOauthUrlResponse *struct {
		Url string `json:"url"`
	}
}

func MakeOauthUrl(shadowvcSdk *shadowVc.ShadowVcSdk, redirectUrl string, isUserInfo int, state string) (*RespOauthUrl, error) {
	postData := make(url.Values)
	postData.Add("redirectUri", redirectUrl)
	postData.Add("state", state)
	postData.Add("isUserInfo", strconv.Itoa(isUserInfo))
	bt, err := shadowvcSdk.Send("sd.wx.oauth.url", postData)
	if err != nil {
		return nil, err
	}
	respOauthUrl := new(RespOauthUrl)
	err = ffjson.Unmarshal(bt, respOauthUrl)
	if err != nil {
		return nil, err
	}
	if respOauthUrl.WechatOauthUrlResponse == nil {
		return nil, shadowVc.NewRespError(bt)
	}
	return respOauthUrl, nil
}

type RespOauthGet struct {
	WechatOauthGetResponse *struct {
		Id         int64  `json:"id"`
		Name       string `json:"name"`
		HeadImg    string `json:"headImg"`
		OpenId     string `json:"openId"`
		NickName   string `json:"nickName"`
		HeadImgUrl string `json:"headImgUrl"`
		Sex        int    `json:"sex"`
		Province   string `json:"province"`
		City       string `json:"city"`
		Country    string `json:"country"`
		UnionId    string `json:"unionId"`
		IsHeed     bool   `json:"isHeed"`
	}
}

func GetOautMember(shadowvcSdk *shadowVc.ShadowVcSdk, code string) (*RespOauthGet, error) {
	postData := make(url.Values)
	postData.Add("code", code)
	bt, err := shadowvcSdk.Send("sd.wx.oauth.get", postData)
	if err != nil {
		return nil, err
	}
	respOauthGet := new(RespOauthGet)
	err = ffjson.Unmarshal(bt, respOauthGet)
	if err != nil {
		return nil, err
	}
	if respOauthGet.WechatOauthGetResponse == nil {
		return nil, shadowVc.NewRespError(bt)
	}
	return respOauthGet, nil
}
