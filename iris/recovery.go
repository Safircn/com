package iris

import (
  "runtime"
  "github.com/kataras/iris"
)

var recoveryHandleFunc = iris.HandlerFunc(func(ctx *iris.Context) {
  defer func() {
    if err := recover(); err != nil {
      var i = 0
      for {
        _, file, line, ok := runtime.Caller(i)
        if !ok {
          break
        }
        ctx.Log("Recovery file %s lie %d\n", file, line)
        i++
      }
      ctx.Log("Recovery from panic\n%s", err)
      ctx.Panic()
    }
  }()
  ctx.Next()
})
