package iris

import (
  "com"
  "github.com/kataras/iris"
  "gopkg.in/go-playground/validator.v9"
  "reflect"
)

type Controller struct{}

type Req struct {
  Code int         `json:"code"`
  Msg  string      `json:"msg"`
  Data interface{} `json:"data"`
}

type BindJsonHandlerFunc func(*iris.Context, interface{})

func (Controller) BindJsonFunc(req interface{}, handleFunc BindJsonHandlerFunc) iris.HandlerFunc {
  if reflect.TypeOf(req).Kind() != reflect.Struct {
    com.Logger.Error("BindJson:类型有误")
  }
  reqType := reflect.TypeOf(req)
  return func(ctx *iris.Context) {
    req = reflect.New(reqType).Interface()
    var err error
    err = ctx.ReadJSON(req)
    if err != nil {
      ctx.JSON(200, Req{
        100,
        "json 格式有误:" + err.Error(),
        nil,
      })
      return
    }
    err = ValidateStruct(req)
    if err != nil {
      errMsg := err.Error()
      if validationErrors, isValidErr := err.(validator.ValidationErrors); isValidErr {
        errMsg = errorMsg(req, validationErrors)
      }
      ctx.JSON(200, Req{
        100,
        errMsg,
        nil,
      })
      return
    }
    handleFunc(ctx, req)
  }
}

func (Controller) RespJson(ctx *iris.Context, code int, msg string, obj interface{}) {
  ctx.JSON(200, Req{
    Code:code,
    Msg:msg,
    Data:obj,
  })
}

func (c *Controller) RespJsonMsg(ctx *iris.Context, code int, msg string) {
  c.RespJson(ctx, code, msg, nil)
}

type IReqErrorMsg interface {
  ErrorMsg() map[string]map[string]string
}

func errorMsg(v interface{}, vailErr validator.ValidationErrors) string {
  if len(vailErr) == 0 {
    return ""
  }
  if errorMsgST, flag := v.(IReqErrorMsg); flag {
    errorMap := errorMsgST.ErrorMsg()

    for _, v := range vailErr {
      if mapV, flag := errorMap[v.StructField()]; !flag {
        return vailErr.Error()
      } else {
        if errMsg, flag := mapV[v.ActualTag()]; !flag {
          if errMsg, flag := mapV["default"]; flag {
            return errMsg
          }
          return vailErr.Error()
        } else {
          return errMsg
        }
      }
    }
  }
  return vailErr.Error()
}