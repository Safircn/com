package order

import (
  "net/url"
  "github.com/pquerna/ffjson/ffjson"
  "strconv"
  "com/shadowVc"
)

type RespOrderWxCreate struct {
  OrderWxCreateResponse *struct {
    No             string `json:"no"`
    WechatPayJsapi struct {
                     AppID      string `json:"appId"`
                     TimeStamp  string `json:"timeStamp"`
                     NonceStr   string `json:"nonceStr"`
                     PackageStr string `json:"packageStr"`
                     SignType   string `json:"signType"`
                     PaySign    string `json:"paySign"`
                   } `json:"wechatPayJsapi"`
  } `json:"OrderWxCreateResponse"`
}

func OrderWxCreate(shadowvcSdk *shadowVc.ShadowVcSdk, orderDesc, orderType string, price int64, notifyUrl string,
clientIp, openId string, effective int) (*RespOrderWxCreate, error) {
  postData := make(url.Values)
  postData.Add("orderDesc", orderDesc)
  postData.Add("orderType", orderType)
  postData.Add("price", strconv.FormatInt(price, 10))
  postData.Add("notifyUrl", notifyUrl)
  postData.Add("clientIp", clientIp)
  postData.Add("openId", openId)
  postData.Add("effective", strconv.Itoa(effective))
  bt, err := shadowvcSdk.Send("sd.order.wx.create", postData)

  OrderWxNativeCreate(shadowvcSdk, orderDesc, orderType, price, notifyUrl, clientIp, effective)

  if err != nil {
    return nil, err
  }

  respSmsCodeSend := new(RespOrderWxCreate)
  err = ffjson.Unmarshal(bt, respSmsCodeSend)
  if err != nil {
    return nil, err
  }
  if respSmsCodeSend.OrderWxCreateResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respSmsCodeSend, err
}

type RespOrderWxNativeCreate struct {
  OrderWxNativeCreateResponse *struct {
    No      string `json:"no"`
    CodeURL string `json:"codeUrl"`
  } `json:"OrderWxNativeCreateResponse"`
}

func OrderWxNativeCreate(shadowvcSdk *shadowVc.ShadowVcSdk, orderDesc, orderType string, price int64, notifyUrl string,
clientIp string, effective int) (*RespOrderWxNativeCreate, error) {
  postData := make(url.Values)
  postData.Add("orderDesc", orderDesc)
  postData.Add("orderType", orderType)
  postData.Add("price", strconv.FormatInt(price, 10))
  postData.Add("notifyUrl", notifyUrl)
  postData.Add("clientIp", clientIp)
  postData.Add("effective", strconv.Itoa(effective))
  bt, err := shadowvcSdk.Send("sd.order.wx.native.create", postData)

  if err != nil {
    return nil, err
  }

  respOrderWxNativeCreate := new(RespOrderWxNativeCreate)
  err = ffjson.Unmarshal(bt, respOrderWxNativeCreate)
  if err != nil {
    return nil, err
  }
  if respOrderWxNativeCreate.OrderWxNativeCreateResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respOrderWxNativeCreate, err
}

type RespOrderWxGet struct {
  OrderWxGetResponse *struct {
    OrderWx struct {
              Appid        string `json:"appid"`
              No           string `json:"no"`
              OrderType    string `json:"orderType"`
              Price        int64  `json:"price"`
              PayType      string `json:"payType"`
              Status       string `json:"status"`
              PayStatus    string `json:"payStatus"`
              PayTime      string `json:"payTime"`
              RefundStatus string `json:"refundStatus"`
              RefundTime   string `json:"refundTime"`
              OrderDesc    string `json:"orderDesc"`
              ClientIP     string `json:"clientIp"`
              TradeType    string `json:"tradeType"`
              OpenID       string `json:"openId"`
              Created      string `json:"created"`
            } `json:"orderWx"`
  } `json:"OrderWxGetResponse"`
}

func OrderWxGet(shadowvcSdk *shadowVc.ShadowVcSdk, no string) (*RespOrderWxGet, error) {
  postData := make(url.Values)
  postData.Add("no", no)

  bt, err := shadowvcSdk.Send("sd.order.wx.get", postData)
  if err != nil {
    return nil, err
  }

  respOrderWxGet := new(RespOrderWxGet)
  err = ffjson.Unmarshal(bt, respOrderWxGet)
  if err != nil {
    return nil, err
  }
  if respOrderWxGet.OrderWxGetResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respOrderWxGet, err
}
