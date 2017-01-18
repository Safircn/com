package wechatWeb

import (
  "time"
  "regexp"
)

var (
  UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36"
  StatusSuccess = 0

  SpecialUsers = []string{
    "newsapp", "fmessage", "filehelper", "weibo", "qqmail",
    "tmessage", "qmessage", "qqsync", "floatbottle", "lbsapp",
    "shakeapp", "medianote", "qqfriend", "readerapp", "blogapp",
    "facebookapp", "masssendapp", "meishiapp", "feedsapp", "voip",
    "blogappweixin", "weixin", "brandsessionholder", "weixinreminder", "wxid_novlwrv3lqwv11",
    "gh_22b87fa7cb3c", "officialaccounts", "notification_messages", "wxitil", "userexperience_alarm",
  }

  SyncHosts = []string{
    //"webpush.wx.qq.com",
    "webpush.weixin.qq.com",
    "webpush.wx2.qq.com",
    "webpush2.wx.qq.com",
    "webpush.wx.qq.com",
  }
  SaveSubFolders = map[string]string{
    "webwxgeticon":                                  "icons",
    "webwxgetheadimg":                               "headimgs",
    "webwxgetmsgimg":                                "msgimgs",
    "webwxgetvideo":                                 "videos",
    "webwxgetvoice":                                 "voices",
    "_showQRCodeImg":                                "qrcodes",
  }
  AppID = "wx782c26e4c19acffb"
  Lang = "zh_CN"
  LastCheckTs = time.Now()
  LoginUrl = "https://login.weixin.qq.com/jslogin"
  QrUrl = "https://login.weixin.qq.com/qrcode/"
  APIKEY = "391ad66ebad2477b908dce8e79f101e7"
  TUringUserId = "abc123"

  HeaderAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
  HeaderAcceptLanguage = "zh-CN,zh;q=0.8"
  HeaderCacheControl = "no-cache"
  HeaderAcceptEncoding = "gzip"
  HeaderConnection = "keep-alive"
  HeaderPragma = "no-cache"
  HeaderUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2970.0 Safari/537.36"

  bytes200 []byte = []byte("200")

  MsgChanDuration = time.Second * 3

  ServiceHosts []ServiceHost = []ServiceHost{
    {
      Host:"weixin.qq.com",
      LoginHost:"login.weixin.qq.com",
      FileHost:"file.wx.qq.com",
      PushHost:"webpush.weixin.qq.com",
    },

    {
      Host:"wx.qq.com",
      LoginHost:"login.wx.qq.com",
      FileHost:"file.wx.qq.com",
      PushHost:"webpush.wx.qq.com",
    },
    {
      Host:"wx2.qq.com",
      LoginHost:"login.wx2.qq.com",
      FileHost:"file.wx2.qq.com",
      PushHost:"webpush.wx2.qq.com",
    },
    {
      Host:"wx8.qq.com",
      LoginHost:"login.wx8.qq.com",
      FileHost:"file.wx8.qq.com",
      PushHost:"webpush.wx8.qq.com",
    },
    {
      Host:"web.wechat.com",
      LoginHost:"login.web.wechat.com",
      FileHost:"file.web.wechat.com",
      PushHost:"webpush.web.wechat.com",
    },
    {
      Host:"web2.wechat.com",
      LoginHost:"login.web2.wechat.com",
      FileHost:"file.web2.wechat.com",
      PushHost:"webpush.web2.wechat.com",
    },
  }
)

func GetServiceHost(host string) *ServiceHost {
 for _,v := range ServiceHosts {
   if v.Host == host {
     return &v
   }
 }
  return nil
}

type ServiceHost struct {
  Host      string
  LoginHost string
  FileHost  string
  PushHost  string
}

var (
  regexpUUID *regexp.Regexp = regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+?)"`)
  regexpRedirect *regexp.Regexp = regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)

  regexpGerWindowCode  *regexp.Regexp = regexp.MustCompile(`window.code=(\d+);`)
  regexpLoginRedirect  *regexp.Regexp = regexp.MustCompile(`window.redirect_uri="(\S+?)"`)
)


