package sms

import (
	"com/shadowVc"
	"github.com/pquerna/ffjson/ffjson"
	"net/url"
	"strconv"
)

type RespSmsCodeSend struct {
	SmsCodeSendResponse *struct {
		SmsCodeId int64 `json:"smsCodeId"`
	}
}

func SendSmsCode(shadowvcSdk *shadowVc.ShadowVcSdk, mobile, label string, length, interval int, content string) (bool, error) {
	postData := make(url.Values)
	postData.Add("mobile", mobile)
	postData.Add("label", label)
	postData.Add("length", strconv.Itoa(length))
	postData.Add("interval", strconv.Itoa(interval))
	postData.Add("content", content)
	bt, err := shadowvcSdk.Send("sd.sms.code.send", postData)
	if err != nil {
		return false, err
	}
	respSmsCodeSend := new(RespSmsCodeSend)
	err = ffjson.Unmarshal(bt, respSmsCodeSend)
	if err != nil {
		return false, err
	}
	if respSmsCodeSend.SmsCodeSendResponse == nil {
		return false, shadowVc.NewRespError(bt)
	}
	return true, nil
}

type RespSmsCodeCheck struct {
	SmsCodeCheckResponse *struct {
		Result bool `json:"result"`
	}
}

func CheckSmsCode(shadowvcSdk *shadowVc.ShadowVcSdk, mobile, label, code string) (bool, error) {
	postData := make(url.Values)
	postData.Add("mobile", mobile)
	postData.Add("label", label)
	postData.Add("code", code)
	bt, err := shadowvcSdk.Send("sd.sms.code.check", postData)
	if err != nil {
		return false, err
	}
	respSmsCodeCheck := new(RespSmsCodeCheck)
	err = ffjson.Unmarshal(bt, respSmsCodeCheck)
	if err != nil {
		return false, err
	}
	if respSmsCodeCheck.SmsCodeCheckResponse == nil {
		return false, shadowVc.NewRespError(bt)
	}
	if respSmsCodeCheck.SmsCodeCheckResponse.Result {
		return true, nil
	}
	return false, nil
}

type RespSmsTemplateSend struct {
	SmsTemplateSendResponse *struct {
		Sid int64 `json:"sid"`
	}
}

func SmsTemplateSend(shadowvcSdk *shadowVc.ShadowVcSdk, mobile, content string) (int64, error) {
	postData := make(url.Values)
	postData.Add("mobile", mobile)
	postData.Add("content", content)
	bt, err := shadowvcSdk.Send("sd.sms.template.send", postData)
	if err != nil {
		return 0, err
	}
	respSmsTemplateSend := new(RespSmsTemplateSend)
	err = ffjson.Unmarshal(bt, respSmsTemplateSend)
	if err != nil {
		return 0, err
	}
	if respSmsTemplateSend.SmsTemplateSendResponse == nil {
		return 0, shadowVc.NewRespError(bt)
	}
	return respSmsTemplateSend.SmsTemplateSendResponse.Sid, nil
}


func SmsTemplateSendToRaw(shadowvcSdk *shadowVc.ShadowVcSdk, mobile, content string) (string, error) {
  postData := make(url.Values)
  postData.Add("mobile", mobile)
  postData.Add("content", content)
  bt, err := shadowvcSdk.Send("sd.sms.template.send", postData)
  if err != nil {
    return "", err
  }
  return string(bt),nil
}
