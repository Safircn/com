package iris

import (
  "github.com/kataras/iris"
  "com"
)

type (
  Rest struct {
    restMap    map[string][]iris.Handler

    restGroups []*RestGroup
  }
  RestGroup struct {
    handlers     []iris.HandlerFunc
    restHandlers map[string][]iris.HandlerFunc
  }

)

func RestNew() *Rest {
  rest := &Rest{
    restMap:make(map[string][]iris.Handler),
    restGroups:make([]*RestGroup, 1),
  }
  rest.restGroups[0] = &RestGroup{
    handlers: make([]iris.HandlerFunc, 0),
    restHandlers:make(map[string][]iris.HandlerFunc),
  }
  return rest
}

func (r *Rest) NewGroup(handleFuncs ... iris.HandlerFunc) *RestGroup {
  g := new(RestGroup)
  g.handlers = handleFuncs
  g.restHandlers = make(map[string][]iris.HandlerFunc)
  r.restGroups = append(r.restGroups, g)
  return g
}

func (r *Rest) Method(method string, handleFuncs ... iris.HandlerFunc) {
  r.restGroups[0].restHandlers[method] = handleFuncs
}

func (rg *RestGroup) Method(method string, handleFuncs ... iris.HandlerFunc) {
  rg.restHandlers[method] = handleFuncs
}

func (r *Rest) HandleFunc() iris.HandlerFunc {
  //rests
  var (
    handleFuncs []iris.HandlerFunc
  )
  for k, v := range r.restGroups {
    if len(v.restHandlers) > 0 {
      for kk, vv := range v.restHandlers {
        handleFuncs = make([]iris.HandlerFunc, 0)
        if k > 0 && len(v.handlers) > 0 {
          handleFuncs = append(handleFuncs, v.handlers...)
        }
        handleFuncs = append(handleFuncs, vv...)
        r.restMap[kk] = convertToHandlers(handleFuncs)
        com.Logger.Info("register Rest %s", kk)
      }
    }
  }
  controller := new(Controller)
  return func(ctx *iris.Context) {
    method := ctx.URLParam("method")
    if method == "" {
      controller.RespJsonMsg(ctx, 404, "method不能为空")
      return
    }
    if handlerFuncs, flag := r.restMap[method]; flag {
      com.Logger.Info("rest %s", method)
      ctx.Middleware = append(ctx.Middleware, handlerFuncs...)
      ctx.Next()
    } else {
      controller.RespJsonMsg(ctx, 404, "method方法名不存在")
      return
    }
  }
}

func convertToHandlers(handlersFn []iris.HandlerFunc) []iris.Handler {
  hlen := len(handlersFn)
  mlist := make([]iris.Handler, hlen)
  for i := 0; i < hlen; i++ {
    mlist[i] = iris.Handler(handlersFn[i])
  }
  return mlist
}