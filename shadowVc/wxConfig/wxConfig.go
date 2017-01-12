package wxConfig

import (
  "com/shadowVc"
  //"github.com/pquerna/ffjson/ffjson"
  "net/url"
  "github.com/pquerna/ffjson/ffjson"
)

type RespWxConfigTokenGet struct {
  WxConfigTokenGetResponse *struct {
    AccessToken string `json:"accessToken"`
  }
}

func WxConfigTokenGet(shadowvcSdk *shadowVc.ShadowVcSdk) (*RespWxConfigTokenGet, error) {
  postData := make(url.Values)
  bt, err := shadowvcSdk.Send("sd.wx.config.getToken", postData)
  if err != nil {
    return nil, err
  }
  respWxConfigTokenGet := new(RespWxConfigTokenGet)
  err = ffjson.Unmarshal(bt, respWxConfigTokenGet)
  if err != nil {
    return nil, err
  }
  if respWxConfigTokenGet.WxConfigTokenGetResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respWxConfigTokenGet, nil
}
