package gin

import (
  "github.com/Safircn/com"
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/contrib/sessions"
)

var runmode = com.Conf.String("runmode")

func InitFramework() *gin.Engine {
  engine := gin.Default()
  initValidator()
  if runmode == "pro" {
    gin.SetMode(gin.ReleaseMode)
  }else{
    gin.SetMode(gin.DebugMode)
  }

  return engine
}

func UseSession(engine *gin.Engine) {
  store := sessions.NewCookieStore([]byte(com.Conf.DefaultString("gin::cookieKey", "41a84ae78fb46628319f97091026307d")))
  engine.Use(sessions.Sessions(com.Conf.DefaultString("gin::cookieName", "ecbox"), store))
}

type DemoStruct struct {
  Name string `form:"name" validate:"required"`
  City string `form:"city"`
  Sex  int    `form:"sex" validate:"required"`
}

func Run(engine *gin.Engine) {
  engine.Run(":" + com.Conf.DefaultString("port", "8000"))
}

