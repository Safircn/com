package wechatWeb

import (
  "bytes"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "math/rand"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "time"
  "com"
  "github.com/pquerna/ffjson/ffjson"
  "encoding/json"
  "compress/gzip"
)

func (w *Wechat) GetContacts() (err error) {
  name, resp := "webwxgetcontact", new(MemberResp)
  apiURI := fmt.Sprintf("%s/%s?pass_ticket=%s&skey=%s&r=%s", w.BaseUri, name, w.PassTicket, w.Skey, w.GetUnixTime())

  respBytes, err := w.send("GET", apiURI, nil)
  if err != nil {
    return err
  }
  err = ffjson.Unmarshal(respBytes, resp)
  if err != nil {
    return err
  }

  if !resp.IsSuccess() {
    return resp.Error()
  }

  w.MemberList = resp.MemberList
  w.TotalMember = resp.MemberCount
  for _, member := range w.MemberList {
    w.MemberMap[member.UserName] = member
    if member.UserName[:2] == "@@" {
      w.GroupMemberList = append(w.GroupMemberList, member) //群聊

    } else if member.VerifyFlag & 8 != 0 {
      w.PublicUserList = append(w.PublicUserList, member) //公众号
    } else if member.UserName[:1] == "@" {
      w.ContactList = append(w.ContactList, member)
    }
  }
  mb := Member{}
  mb.NickName = w.User.NickName
  mb.UserName = w.User.UserName
  w.MemberMap[w.User.UserName] = mb
  for _, user := range w.ChatSet {
    exist := false
    for _, initUser := range w.InitContactList {
      if user == initUser.UserName {
        exist = true
        break
      }
    }
    if !exist {
      value, ok := w.MemberMap[user]
      if ok {
        contact := User{
          UserName:  value.UserName,
          NickName:  value.NickName,
          Signature: value.Signature,
        }

        w.InitContactList = append(w.InitContactList, contact)
      }
    }

  }
  return
}

func (w *Wechat) getWechatRoomMember(roomID, userId string) (roomName, userName string, err error) {
  //apiUrl := fmt.Sprintf("%s/webwxbatchgetcontact?type=ex&r=%s&pass_ticket=%s", w.BaseUri, w.GetUnixTime(), w.PassTicket)
  //params := make(map[string]interface{})
  //params["Count"] = 1
  //params["List"] = []map[string]string{}
  //l := []map[string]string{}
  //params["List"] = append(l, map[string]string{
  //  "UserName":   roomID,
  //  "ChatRoomId": "",
  //})

  return "", "", nil
}

func (w *Wechat) getSyncMsg() (*SyncResp, error) {
  name := "webwxsync"

  urlStr := fmt.Sprintf("%s/%s?sid=%s&pass_ticket=%s&skey=%s", w.BaseUri, name, w.Sid, w.PassTicket, w.Skey)
  params := SyncParams{
    BaseRequest: json.RawMessage(w.RequestJson),
    SyncKey:     w.SyncKey,
    RR:          time.Now().Unix(),
  }
  data, err := ffjson.Marshal(params)
  if err != nil {
    return nil, err
  }
  syncResp := new(SyncResp)

  respBytes, err := w.send("POST", urlStr, bytes.NewReader(data))
  if err != nil {
    return nil, err
  }
  err = ffjson.Unmarshal(respBytes, syncResp)
  if err != nil {
    return nil, err
  }

  if !syncResp.IsSuccess() {
    return nil, syncResp.Error()
  }

  if syncResp.BaseResponse.Ret == 0 {
    w.SyncKey = syncResp.SyncKey
    w.setSyncKeyStr(syncResp.SyncKey)
  }
  return syncResp, nil
}

//同步守护goroutine
func (w *Wechat) SyncDaemon() {
  resp, err := w.SyncCheck()
  if err != nil {
    com.Logger.Error("SyncDaemon err:%v", err)
    w.SetStatusOffline()
  }
  com.Logger.Info("SyncDaemon info:%v", resp)
  switch resp.RetCode {
  case 1100, 1101, 1102:
    w.SetStatusOffline()
    break
  case 0:
    switch resp.Selector {
    case 0:

    case 2, 3: //有消息,未知
      msgs, err := w.getSyncMsg()

      if err != nil {
        com.Logger.Info("SyncDaemon err:%v", err)
      }
      com.Logger.Info("SyncDaemon info:%v", msgs)
      if msgs.AddMsgCount > 0 {

        for _, m := range msgs.AddMsgList {
          msg := new(Message)
          msg.MsgType = m.MsgType
          msg.FromUserName = m.FromUserName
          if nickNameFrom, ok := w.MemberMap[msg.FromUserName]; ok {
            msg.FromUserNickName = nickNameFrom.NickName
          }

          msg.ToUserName = m.ToUserName
          if nickNameTo, ok := w.MemberMap[msg.ToUserName]; ok {
            msg.ToUserNickName = nickNameTo.NickName
          }

          msg.Content = m.Content
          msg.Content = strings.Replace(msg.Content, "&lt;", "<", -1)
          msg.Content = strings.Replace(msg.Content, "&gt;", ">", -1)
          msg.Content = strings.Replace(msg.Content, " ", " ", 1)
          switch msg.MsgType {
          case 1:

            if msg.FromUserName[:2] == "@@" {
              //群消息，暂时不处理
              if msg.FromUserNickName == "" {
                contentSlice := strings.Split(msg.Content, ":<br/>")
                msg.Content = contentSlice[1]

              }
            } else {
            }
            if msg.ToUserNickName == "" {
              if user, ok := w.MemberMap[msg.ToUserName]; ok {
                msg.ToUserNickName = user.NickName
              }

            }
            if msg.FromUserNickName == "" {
              if user, ok := w.MemberMap[msg.FromUserNickName]; ok {
                msg.FromUserNickName = user.NickName
              }
            }
          case 3:
          //图片
          case 34:
          //语音
          case 47:
          //动画表情
          case 49:
          //链接
          case 51:
          //获取联系人信息成功
          case 62:
          //获得一段小视频
          case 10002:
          //撤回一条消息
          }
        }
      }

    case 4: //通讯录更新
      w.GetContacts()
    case 6: //可能是红包
    //w.Log.Println("请速去手机抢红包")
    case 7:
    //w.Log.Println("在手机上操作了微信")
    //w.Log.Println("无事件")
    default:
    }
  default:
  //w.Log.Printf("the resp:%+v", resp)
  }

}

func (w *Wechat) MsgChan() {
  if w.status != WechatStatusOnline {
    return
  }
  syncCheckTimer := time.NewTimer(MsgChanDuration)
  var (
    syncDaemonStart time.Time
    diffDuration time.Duration
  )
  forBack:
  for {
    select {
    case <-syncCheckTimer.C:
      syncDaemonStart = time.Now()
      w.SyncDaemon()
      diffDuration = MsgChanDuration - (time.Now().Sub(syncDaemonStart))
      if diffDuration < 0 {
        diffDuration = 0
      }
      syncCheckTimer.Reset(diffDuration)
    case <-w.offlineChan:
      syncCheckTimer.Stop()
      break forBack
    }
  }
}

func (w *Wechat) StatusNotify() (err error) {
  statusURL := w.BaseUri + fmt.Sprintf("/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", w.PassTicket)
  resp := new(NotifyResp)
  data, err := ffjson.Marshal(NotifyParams{
    BaseRequest:  json.RawMessage(w.RequestJson),
    Code:         3,
    FromUserName: w.User.UserName,
    ToUserName:   w.User.UserName,
    ClientMsgId:  w.GetUnixTime(),
  })

  respBytes, err := w.send("OIST", statusURL, bytes.NewReader(data))
  if err != nil {
    return err
  }
  err = ffjson.Unmarshal(respBytes, resp)
  if err != nil {
    return err
  }

  if !resp.IsSuccess() {
    return resp.Error()
  }

  return
}

func (w *Wechat) GetContactsInBatch() (err error) {
  resp := new(MemberResp)
  apiUrl := fmt.Sprintf("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r=%s&pass_ticket=%s", w.GetUnixTime(), w.PassTicket)

  respBytes, err := w.send("GET", apiUrl, nil)
  if err != nil {
    return err
  }
  err = ffjson.Unmarshal(respBytes, resp)
  if err != nil {
    return err
  }

  if !resp.IsSuccess() {
    return resp.Error()
  }

  return
}
func (w *Wechat) SyncCheck() (*SyncCheckResp, error) {
  params := url.Values{}
  curTime := strconv.FormatInt(time.Now().Unix(), 10)
  params.Set("r", curTime)
  params.Set("sid", w.Sid)
  params.Set("uin", strconv.FormatInt(w.Uin, 10))
  params.Set("skey", w.Skey)
  params.Set("deviceid", w.DeviceId)
  params.Set("synckey", w.SyncKeyStr)
  params.Set("_", curTime)
  checkUrl := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", w.ServiceHost.PushHost)
  Url, err := url.Parse(checkUrl)
  if err != nil {
    return nil, err
  }
  Url.RawQuery = params.Encode()

  bodyBytes, err := w.send("GET", Url.String(), nil)
  if err != nil {
    return nil, err
  }
  resp := new(SyncCheckResp)
  pmSub := regexpRedirect.FindStringSubmatch(string(bodyBytes))
  if len(pmSub) == 0 {
    com.Logger.Error("regex err body:%s", string(bodyBytes))
    return nil, errors.New("regex error in window.redirect_uri")
  }
  resp.RetCode, err = strconv.Atoi(pmSub[1])
  resp.Selector, err = strconv.Atoi(pmSub[2])
  return resp, nil
}

func (w *Wechat) SendMsg(toUserName, message string, isFile bool) (err error) {
  resp := new(MsgResp)

  apiUrl := fmt.Sprintf("%s/webwxsendmsg?pass_ticket=%s", w.BaseUri, w.PassTicket)
  clientMsgId := strconv.Itoa(w.GetUnixTime()) + "0" + strconv.Itoa(rand.Int())[3:6]
  params := make(map[string]interface{})
  params["BaseRequest"] = w.BaseRequest
  msg := make(map[string]interface{})
  msg["Type"] = 1
  msg["Content"] = message
  msg["FromUserName"] = w.User.UserName
  msg["LocalID"] = clientMsgId
  msg["ClientMsgId"] = clientMsgId
  msg["ToUserName"] = toUserName
  params["Msg"] = msg
  data, err := ffjson.Marshal(params)
  if err != nil {
    return err
    //w.Log.Printf("json.Marshal(%v):%v\n", params, err)
  }

  respBytes, err := w.send("POST", apiUrl, bytes.NewReader(data))
  if err != nil {
    return err
  }
  err = ffjson.Unmarshal(respBytes, resp)
  if err != nil {
    return err
  }

  if !resp.IsSuccess() {
    return resp.Error()
  }

  return
}

func (w *Wechat) GetGroupName(id string) (name string) {
  return
}

func (w *Wechat) SendMsgToAll(word string) (err error) {

  return
}

func (w *Wechat) SendImage(name, fileName string) (err error) {

  return
}

func (w *Wechat) AddMember(name string) (err error) {

  return
}

func (w *Wechat) CreateRoom(name string) (err error) {

  return
}

func (w *Wechat) PullMsg() {
  return
}

func (w *Wechat) newResquest(method, urlStr string, body io.Reader) (*http.Request, error) {
  res, err := http.NewRequest(method, urlStr, body)
  if err != nil {
    return nil, err
  }
  res.Header.Set("Accept", HeaderAccept)
  res.Header.Set("Accept-Encoding", HeaderAcceptEncoding)
  res.Header.Set("Accept-Language", HeaderAcceptLanguage)
  res.Header.Set("Cache-Control", HeaderCacheControl)
  res.Header.Set("Connection", HeaderConnection)
  res.Header.Set("Pragma", HeaderPragma)
  res.Header.Set("User-Agent", HeaderUserAgent)

  if method == "POST" {
    res.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  }
  return res, err
}

func (w *Wechat) send(method, urlStr string, body io.Reader) ([]byte, error) {
  req, err := w.newResquest(method, urlStr, body)
  if err != nil {
    return nil, err
  }
  resp, err := w.Client.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode != 200 {
    return nil, errors.New("send response code not 200")
  }
  if strings.Index(resp.Header.Get("Content-Encoding"), "gzip") > -1 {
    gzipReader, err := gzip.NewReader(resp.Body)
    if err != nil {
      return nil, err
    }
    bodyBytes, err := ioutil.ReadAll(gzipReader)
    if err != nil {
      return nil, err
    }
    return bodyBytes, nil
  }
  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  return bodyBytes, nil
}

func (w *Wechat) SetCookies() {
  //w.Client.Jar.SetCookies(u, cookies)

}

func (w *Wechat) GetUnixTime() int {
  return int(time.Now().Unix())
}
