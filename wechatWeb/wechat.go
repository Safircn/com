package wechatWeb

import (
  "bytes"
  "crypto/tls"
  "encoding/xml"
  "errors"
  "fmt"
  //  "log"
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "strconv"
  "strings"
  "time"
  "encoding/base64"
  "com"
  "github.com/pquerna/ffjson/ffjson"
  "encoding/json"
)

type WechatStatus int

const (
  WechatStatusOffline WechatStatus = iota
  WechatStatusLogin
  WechatStatusOnline
)

func (w *Wechat) setWechatStatusLogin() {
  w.status = WechatStatusLogin
}
func (w *Wechat) setWechatStatusOnline() {
  w.offlineChan = make(chan bool, 1)
  w.status = WechatStatusOnline
}

func (w *Wechat) SetStatusOffline() {
  //关闭同步
  w.offlineChan <- true
  w.status = WechatStatusOffline
}

type Wechat struct {
  status          WechatStatus
  offlineChan     chan bool

  LoginCode       string

  User            User
  Root            string
  Uuid            string
  BaseUri         string
  RedirectedUri   string
  Uin             int64
  Sid             string
  Skey            string
  PassTicket      string
  DeviceId        string
  BaseRequest     map[string]string
  LowSyncKey      string
  SyncKeyStr      string
  SyncKey         SyncKey
  Users           []string
  InitContactList []User   //谈话的人
  MemberList      []Member //
  ContactList     []Member //好友
  GroupList       []string //群
  GroupMemberList []Member //群友
  PublicUserList  []Member //公众号
  SpecialUserList []Member //特殊账号

  AutoReplyMode   bool     //default false
  AutoOpen        bool
  Interactive     bool
  TotalMember     int
  TimeOut         int      // 同步时间间隔   default:20
  MediaCount      int      // -1
  SaveFolder      string
  QrImagePath     string
  Client          *http.Client

  MemberMap       map[string]Member
  ChatSet         []string

  RequestJson     []byte
  HostStr         string
  ServiceHost     *ServiceHost
}

func (w *Wechat) requestParamsJson() []byte {
  baseRequestJson, _ := ffjson.Marshal(BaseRequest{
    Skey:w.Skey,
    Wxsid:w.Sid,
    Wxuin:w.Uin,
    DeviceID:w.DeviceId,
  })
  return baseRequestJson
}

func NewWechat(deviceId string) *Wechat {
  jar, err := cookiejar.New(nil)
  if err != nil {
    return nil
  }

  transport := *(http.DefaultTransport.(*http.Transport))
  transport.ResponseHeaderTimeout = 1 * time.Minute
  transport.TLSClientConfig = &tls.Config{
    InsecureSkipVerify: true,
  }
  //  DeviceId:      "e123456789002237",
  return &Wechat{
    status:WechatStatusOffline,
    DeviceId:deviceId,
    AutoReplyMode: false,
    Interactive:   false,
    AutoOpen:      false,
    MediaCount:    -1,
    Client: &http.Client{
      Transport: &transport,
      Jar:       jar,
      Timeout:   1 * time.Minute,
    },
    MemberMap:   make(map[string]Member),
  }

}

func (w *Wechat) MakeUUIDToGetQR() (string, error) {
  if w.status != WechatStatusOffline {
    return "", errors.New("not Offline")
  }
  err := w.GetUUID()
  if err != nil {
    return "", err
  }
  var base64Qr string
  base64Qr, err = w.GetQR()
  if err != nil {
    return "", err
  }
  w.setWechatStatusLogin()
  return base64Qr, nil
}

func (w *Wechat) WaitForLogin() (err error) {
  code, tip := "", 1
  w.LoginCode = "408"
  for code != "200" {
    if w.status != WechatStatusLogin {
      return errors.New("not Login Status")
    }
    w.RedirectedUri, code, tip, err = w.waitToLogin(w.Uuid, tip)
    if err != nil {
      err = fmt.Errorf("二维码登陆失败：%s", err.Error())
      com.Logger.Error("WaitForLogin ERR:%v", err)
      return
    }
  }
  err = w.Login()
  if err != nil {
    com.Logger.Error("WaitForLogin ERR:%v", err)
  }
  w.status = WechatStatusOnline
  if err := w.StatusNotify(); err != nil {
    com.Logger.Error("开启状态栏通知失败:%v", err)
    return err
  }
  if err := w.GetContacts(); err != nil {
    com.Logger.Error("拉取联系人失败:%v\n", err)
    return err
  }
  //V  if err := w.TestCheck(); err != nil {
  //    com.Logger.Error("检查状态失败:%v\n", err)
  //    return err
  //  }
  go w.MsgChan() //接收消息
  return
}

func (w *Wechat) waitToLogin(uuid string, tip int) (redirectUri, code string, rt int, err error) {
  loginUri := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=%d&uuid=%s&_=%s", tip, uuid, time.Now().Unix())
  rt = tip

  respBytes, err := w.send("GET", loginUri, nil)
  if err != nil {
    return
  }

  pm := regexpGerWindowCode.FindStringSubmatch(string(respBytes))
  if len(pm) != 0 {
    code = pm[1]
  } else {
    err = errors.New("can't find the code")
    return
  }

  rt = 0
  w.LoginCode = code
  switch code {
  case "201":
  //扫描成功，请在手机上点击确认登陆"
  case "200":
    pmSub := regexpLoginRedirect.FindStringSubmatch(string(respBytes))
    if len(pmSub) != 0 {
      redirectUri = pmSub[1]
    } else {
      err = errors.New("regex error in window.redirect_uri")
      return
    }
    redirectUri += "&fun=new"
  case "408":
  case "0":
    err = errors.New("超时了，请重启程序")
    w.status = WechatStatusOffline
  default:
    err = errors.New("其它错误，请重启")
    w.status = WechatStatusOffline
  }
  return
}

//https://login.weixin.qq.com/I/IfPNLdzWzQ==
func (w *Wechat) GetQR() (string, error) {

  if w.Uuid == "" {
    return "", errors.New("no this uuid")
  }
  params := url.Values{}
  params.Set("t", "webwx")
  params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))

  respBytes, err := w.send("POST", QrUrl + w.Uuid, strings.NewReader(params.Encode()))
  if err != nil {
    return "", err
  }

  return "data:image/png;base64," + base64.StdEncoding.EncodeToString(respBytes), nil
}

func (w *Wechat) SetSynKey() {

}

func (w *Wechat) GetUUID() (err error) {
  params := url.Values{}
  params.Set("appid", AppID)
  params.Set("fun", "new")
  params.Set("lang", "zh_CN")
  params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))
  dateBytes, err := w.send("POST", LoginUrl, strings.NewReader(params.Encode()))
  if err != nil {
    return err
  }

  pm := regexpUUID.FindSubmatch(dateBytes)

  if len(pm) > 0 {
    code := pm[1]
    if !bytes.Equal(code, bytes200) {
      return errors.New("the status error")
    }
    w.Uuid = string(pm[2])
    return nil
  }
  return errors.New("get uuid failed")
}

func (w *Wechat) Login() (err error) {
  respBytes, err := w.send("GET", w.RedirectedUri, nil)
  if err != nil {
    return err
  }

  baseRequest := new(BaseRequest)
  if err = xml.Unmarshal(respBytes, baseRequest); err != nil {
    return err
  }

  w.DeviceId = w.DeviceId
  w.Skey = baseRequest.Skey
  w.Sid = baseRequest.Wxsid
  w.Uin = baseRequest.Wxuin
  w.PassTicket = baseRequest.PassTicket

  w.RequestJson = w.requestParamsJson()

  //XMLName    xml.Name `xml:"error" json:"-"`
  //Ret        int      `xml:"ret" json:"-"`
  //Message    string   `xml:"message" json:"-"`
  //Skey       string   `xml:"skey" json:"Skey"`
  //Wxsid      string   `xml:"wxsid" json:"Sid"`
  //Wxuin      int64    `xml:"wxuin" json:"Uin"`
  //PassTicket string   `xml:"pass_ticket" json:"-"`
  //DeviceID   string   `xml:"-" json:"DeviceID"`

  index := strings.LastIndex(w.RedirectedUri, "/")
  if index == -1 {
    index = len(w.RedirectedUri)
  }
  w.BaseUri = w.RedirectedUri[:index]

  uriTmp := w.BaseUri[8:]
  w.HostStr = uriTmp[:strings.Index(uriTmp, "/")]
  serviceHost := GetServiceHost(w.HostStr)
  if serviceHost == nil {
    com.Logger.Error("serviceHost nil host:%s", serviceHost)
    return errors.New("serviceHost nil")
  }

  w.ServiceHost = serviceHost

  err = w.WebWXInit()
  if err != nil {
    return err
  }

  return
}

func (w *Wechat) WebWXInit() error {
  newResp := new(InitResp)
  apiUri := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%d", w.BaseUri, w.PassTicket, w.Skey, int(time.Now().Unix()))

  initWxRes, _ := ffjson.Marshal(Request{
    BaseRequest: json.RawMessage(w.RequestJson),
  })

  respBytes, err := w.send("POST", apiUri, bytes.NewReader(initWxRes))
  if err != nil {
    return err
  }
  err = ffjson.Unmarshal(respBytes, newResp)
  if err != nil {
    return err
  }

  for _, contact := range newResp.ContactList {
    w.InitContactList = append(w.InitContactList, contact)
  }
  w.ChatSet = strings.Split(newResp.ChatSet, ",")
  w.User = newResp.User
  w.SyncKey = newResp.SyncKey
  w.setSyncKeyStr(newResp.SyncKey)

  return nil
}

func (w *Wechat) setSyncKeyStr(syncKey SyncKey) {
  var syncKeyStr string
  for _, item := range w.SyncKey.List {
    syncKeyStr += "|" + strconv.Itoa(item.Key) + "_" + strconv.Itoa(item.Val)
  }
  w.SyncKeyStr = syncKeyStr
}
