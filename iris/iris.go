package iris

import (
  "com"
  "github.com/iris-contrib/middleware/logger"
  "github.com/kataras/go-sessions/sessiondb/leveldb"
  "github.com/kataras/go-template/html"
  "github.com/kataras/iris"
  "time"
)

var (
  runmode = com.Conf.String("runmode")
)

type  Framework struct {
  *iris.Framework
}

func IrisInitFramework() *Framework {
  initValidator()
  framework := &Framework{
    iris.New(),
  }


  if runmode != "pro" {
    framework.Config.IsDevelopment = true

  }
  framework.Config.Gzip = true
  framework.Config.Charset = "UTF-8"
  framework.UseSessionDB(leveldb.New(leveldb.Config{
    Path: "dbpath",
    CleanTimeout:  iris.DefaultSessionGcDuration,
    MaxAge: time.Hour * 384,
  }))
  framework.Use(logger.New())
  framework.Use(recoveryHandleFunc)
  framework.UseTemplate(html.New(html.Config{
    Left:        "{{",
    Right:       "}}",
    Layout:      "",
    Funcs:       FuncMaps,
    LayoutFuncs: make(map[string]interface{}, 0),
  })).Directory("./app/view", ".html")

  return framework
}

type DemoStruct struct {
  Name string `form:"name" validate:"required"`
  City string `form:"city"`
  Sex  int    `form:"sex" validate:"required"`
}


func (f Framework) RunIris() {
  f.Listen(":" + com.Conf.DefaultString("port", "8000"))
}
