package order

import (
  "net/url"
  "github.com/pquerna/ffjson/ffjson"
  "strconv"
  "com/shadowVc"
)

type RespOrderWxRefund struct {
  OrderWxRefundResponse *struct {
    OrderWxRefund struct {
                  } `json:"orderWxRefund"`
  } `json:"OrderWxRefundResponse"`
}

func OrderWxRefund(shadowvcSdk *shadowVc.ShadowVcSdk, no string, price int64) (*RespOrderWxRefund, error) {
  postData := make(url.Values)
  postData.Add("no", no)
  postData.Add("price", strconv.FormatInt(price, 10))
  bt, err := shadowvcSdk.Send("sd.order.wx.refund", postData)

  if err != nil {
    return nil, err
  }


  respOrderWxRefund := new(RespOrderWxRefund)
  err = ffjson.Unmarshal(bt, respOrderWxRefund)
  if err != nil {
    return nil, err
  }
  if respOrderWxRefund.OrderWxRefundResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respOrderWxRefund, err
}
