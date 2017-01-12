package com

import (
  "github.com/Safircn/lib/conf"
  "github.com/Safircn/lib/config"
  "github.com/astaxie/beego/logs"
)

var (
  Logger *logs.BeeLogger
  Conf config.Configer
)

func init() {
  iniPath := initConf()
  initLog()
  Logger.Info("配置文件:%s", iniPath)
}

//log_path  conf log_path
func initLog() {
  Logger = logs.NewLogger(10000)
  Logger.EnableFuncCallDepth(true)
  if Conf.String("runmode") == "pro" {
    Logger.SetLogger("file", `{"filename":"` + Conf.DefaultString("log_path", "logs/app.log") + `"}`)
  } else {
    Logger.SetLogger("console", "")
  }
}

func initConf() string {
  var err error
  path, err := conf.FindConf("./")
  if err != nil {
    panic("无法找到配置文件")
  }
  Conf, err = config.NewConfig("ini", path)
  if err != nil {
    panic("配置文件初始化失败")
  }
  return path
}
