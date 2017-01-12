package company

import (
  "com/shadowVc"
  "strconv"
  "github.com/pquerna/ffjson/ffjson"
  "net/url"
)

type RespTransfers struct {
  CompanyTransfersResponse *struct {
    CompanyTransfers struct {
                       PaymentNo      string `json:"paymentNo"`
                       PartnerTradeNo string `json:"partnerTradeNo"`
                     } `json:"companyTransfers"`
  } `json:"CompanyTransfersResponse"`
}

func Transfers(shadowvcSdk *shadowVc.ShadowVcSdk, openid string, amount int64, userName, desc string) (*RespTransfers, error) {
  postData := make(url.Values)
  postData.Add("openid", openid)
  postData.Add("amount", strconv.FormatInt(amount, 10))
  if userName != "" {
    postData.Add("userName", openid)
  }
  postData.Add("desc", desc)

  bt, err := shadowvcSdk.Send("sd.wx.company.transfers", postData)

  if err != nil {
    return nil, err
  }

  respTransfers := new(RespTransfers)
  err = ffjson.Unmarshal(bt, respTransfers)
  if err != nil {
    return nil, err
  }
  if respTransfers.CompanyTransfersResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respTransfers, err
}

