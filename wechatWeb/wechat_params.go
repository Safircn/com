package wechatWeb

import (
  "encoding/xml"
  "fmt"
  "os"
  "strings"
  "encoding/json"
)

type Config struct {
  AutoReply bool     `json:"auto_reply"`
  AutoSave  bool     `json:"auto_save"`
  ReplyMsgs []string `json:"reply_msgs"`
}

type MessageOut struct {
  ToUserName string
  Content    string
  Type       int
}

type Message struct {
  FromUserName         string
  PlayLength           int
  RecommendInfo        []string
  Content              string
  StatusNotifyUserName string
  StatusNotifyCode     int
  Status               int
  VoiceLength          int
  ToUserName           string
  ForwardFlag          int
  AppMsgType           int
  AppInfo              AppInfo
  Url                  string
  ImgStatus            int
  MsgType              int
  ImgHeight            int
  MediaId              string
  FileName             string
  FileSize             string
  FromUserNickName     string
  ToUserNickName       string
}

func (m Message) String() string {
  from := m.FromUserNickName
  to := m.ToUserNickName
  if from == "" {
    from = m.FromUserName
  }
  if to == "" {
    to = m.ToUserName
  }
  return from + "->" + to + ":" + m.Content + "\n"
}

type AppInfo struct {
  Type  int
  AppID string
}

type GetUUIDParams struct {
  AppID    string  `json:"appid"`
  Fun      string  `json:"fun"`
  Lang     string  `json:"lang"`
  UnixTime float64 `json:"-"`
}

type Response struct {
  BaseResponse *BaseResponse `json:"BaseResponse"`
}

type Request struct {
  BaseRequest   json.RawMessage

  MemberCount   int    `json:",omitempty"`
  MemberList    []User `json:",omitempty"`
  Topic         string `json:",omitempty"`
  Caller
  ChatRoomName  string `json:",omitempty"`
  DelMemberList string `json:",omitempty"`
  AddMemberList string `json:",omitempty"`
}
type Caller interface {
  IsSuccess() bool
  Error() error
}

type BaseRequest struct {
  XMLName    xml.Name `xml:"error" json:"-"`
  Ret        int      `xml:"ret" json:"-"`
  Message    string   `xml:"message" json:"-"`
  Skey       string   `xml:"skey" json:"Skey"`
  Wxsid      string   `xml:"wxsid" json:"Sid"`
  Wxuin      int64    `xml:"wxuin" json:"Uin"`
  PassTicket string   `xml:"pass_ticket" json:"-"`
  DeviceID   string   `xml:"-" json:"DeviceID"`
}

type BaseResponse struct {
  Ret    int
  ErrMsg string
}

type MsgResp struct {
  Response
}

type InitResp struct {
  Response
  User                User    `json:"User"`
  Count               int     `json:"Count"`
  ContactList         []User  `json:"ContactList"`
  SyncKey             SyncKey `json:"SyncKey"`
  ChatSet             string  `json:"ChatSet"`
  SKey                string  `json:"SKey"`
  ClientVersion       int     `json:"ClientVersion"`
  SystemTime          int     `json:"SystemTime"`
  GrayScale           int     `json:"GrayScale"`
  InviteStartCount    int     `json:"InviteStartCount"`
  MPSubscribeMsgCount int     `json:"MPSubscribeMsgCount"`
  //MPSubscribeMsgList  string  `json:"MPSubscribeMsgList"`
  ClickReportInterval int `json:"ClickReportInterval"`
}

type SyncKey struct {
  Count int      `json:"Count"`
  List  []KeyVal `json:"List"`
}

type KeyVal struct {
  Key int `json:"Key"`
  Val int `json:"Val"`
}

func (r *Response) IsSuccess() bool {
  return r.BaseResponse.Ret == StatusSuccess
}

func (r *Response) Error() error {
  return fmt.Errorf("message:[%s]", r.BaseResponse.ErrMsg)
}

type MemberResp struct {
  Response
  MemberCount  int
  ChatRoomName string
  MemberList   []Member
  Seq          int
}

func (this *Member) IsNormal(mySelf string) bool {
  return this.VerifyFlag & 8 == 0 && // 公众号/服务号
    !strings.HasPrefix(this.UserName, "@@") && // 群聊
    this.UserName != mySelf && // 自己
    !this.IsSpecail() // 特殊账号
}

func (this *Member) IsSpecail() bool {
  for i, count := 0, len(SpecialUsers); i < count; i++ {
    if this.UserName == SpecialUsers[i] {
      return true
    }
  }
  return false
}

type User struct {
  UserName          string `json:"UserName"`
  Uin               int64  `json:"Uin"`
  NickName          string `json:"NickName"`
  HeadImgUrl        string `json:"HeadImgUrl" xml:""`
  RemarkName        string `json:"RemarkName" xml:""`
  PYInitial         string `json:"PYInitial" xml:""`
  PYQuanPin         string `json:"PYQuanPin" xml:""`
  RemarkPYInitial   string `json:"RemarkPYInitial" xml:""`
  RemarkPYQuanPin   string `json:"RemarkPYQuanPin" xml:""`
  HideInputBarFlag  int    `json:"HideInputBarFlag" xml:""`
  StarFriend        int    `json:"StarFriend" xml:""`
  Sex               int    `json:"Sex" xml:""`
  Signature         string `json:"Signature" xml:""`
  AppAccountFlag    int    `json:"AppAccountFlag" xml:""`
  VerifyFlag        int    `json:"VerifyFlag" xml:""`
  ContactFlag       int    `json:"ContactFlag" xml:""`
  WebWxPluginSwitch int    `json:"WebWxPluginSwitch" xml:""`
  HeadImgFlag       int    `json:"HeadImgFlag" xml:""`
  SnsFlag           int    `json:"SnsFlag" xml:""`
}

type Member struct {
  Uin              int64
  UserName         string
  NickName         string
  HeadImgUrl       string
  ContactFlag      int
  MemberCount      int
  MemberList       []User
  RemarkName       string
  HideInputBarFlag int
  Sex              int
  Signature        string
  VerifyFlag       int
  OwnerUin         int
  PYInitial        string
  PYQuanPin        string
  RemarkPYInitial  string
  RemarkPYQuanPin  string
  StarFriend       int
  AppAccountFlag   int
  Statues          int
  AttrStatus       int
  Province         string
  City             string
  Alias            string
  SnsFlag          int
  UniFriend        int
  DisplayName      string
  ChatRoomId       int
  KeyWord          string
  EncryChatRoomId  string
}

type NotifyParams struct {
  BaseRequest  json.RawMessage
  Code         int
  FromUserName string
  ToUserName   string
  ClientMsgId  int
}

type SyncCheckResp struct {
  RetCode  int `json:"retcode"`
  Selector int `json:"selector"`
}

type SyncParams struct {
  BaseRequest json.RawMessage `json:"BaseRequest"`
  SyncKey     SyncKey     `json:"SyncKey"`
  RR          int64       `json:"rr"`
}

type SyncResp struct {
  Response
  SKey                   string `json:"SKey"`
  SyncKey                SyncKey       `json:"SyncKey"`
  SyncCheckKey           struct {
                           Count int `json:"Count"`
                           List  []struct {
                             Key int `json:"Key"`
                             Val int `json:"Val"`
                           } `json:"List"`
                         } `json:"SyncCheckKey"`

  AddMsgCount            int `json:"AddMsgCount"`
  AddMsgList             []AddMsgList `json:"AddMsgList"`
  ModContactCount        int `json:"ModContactCount"`
  ModContactList         []interface{} `json:"ModContactList"`
  DelContactCount        int `json:"DelContactCount"`
  DelContactList         []interface{} `json:"DelContactList"`
  ModChatRoomMemberCount int `json:"ModChatRoomMemberCount"`
  ModChatRoomMemberList  []interface{} `json:"ModChatRoomMemberList"`
  Profile                *Profile `json:"Profile"`
  ContinueFlag           int `json:"ContinueFlag"`
}

type Profile struct {
  BitFlag           int `json:"BitFlag"`
  UserName          struct {
                      Buff string `json:"Buff"`
                    } `json:"UserName"`
  NickName          struct {
                      Buff string `json:"Buff"`
                    } `json:"NickName"`
  BindUin           int `json:"BindUin"`
  BindEmail         struct {
                      Buff string `json:"Buff"`
                    } `json:"BindEmail"`
  BindMobile        struct {
                      Buff string `json:"Buff"`
                    } `json:"BindMobile"`
  Status            int `json:"Status"`
  Sex               int `json:"Sex"`
  PersonalCard      int `json:"PersonalCard"`
  Alias             string `json:"Alias"`
  HeadImgUpdateFlag int `json:"HeadImgUpdateFlag"`
  HeadImgURL        string `json:"HeadImgUrl"`
  Signature         string `json:"Signature"`
}

type AddMsgList struct {
  MsgID                string `json:"MsgId"`
  FromUserName         string `json:"FromUserName"`
  ToUserName           string `json:"ToUserName"`
  MsgType              int `json:"MsgType"`
  Content              string `json:"Content"`
  Status               int `json:"Status"`
  ImgStatus            int `json:"ImgStatus"`
  CreateTime           int `json:"CreateTime"`
  VoiceLength          int `json:"VoiceLength"`
  PlayLength           int `json:"PlayLength"`
  FileName             string `json:"FileName"`
  FileSize             string `json:"FileSize"`
  MediaID              string `json:"MediaId"`
  URL                  string `json:"Url"`
  AppMsgType           int `json:"AppMsgType"`
  StatusNotifyCode     int `json:"StatusNotifyCode"`
  StatusNotifyUserName string `json:"StatusNotifyUserName"`
  RecommendInfo        struct {
                         UserName   string `json:"UserName"`
                         NickName   string `json:"NickName"`
                         QQNum      int `json:"QQNum"`
                         Province   string `json:"Province"`
                         City       string `json:"City"`
                         Content    string `json:"Content"`
                         Signature  string `json:"Signature"`
                         Alias      string `json:"Alias"`
                         Scene      int `json:"Scene"`
                         VerifyFlag int `json:"VerifyFlag"`
                         AttrStatus int `json:"AttrStatus"`
                         Sex        int `json:"Sex"`
                         Ticket     string `json:"Ticket"`
                         OpCode     int `json:"OpCode"`
                       } `json:"RecommendInfo"`
  ForwardFlag          int `json:"ForwardFlag"`
  AppInfo              struct {
                         AppID string `json:"AppID"`
                         Type  int `json:"Type"`
                       } `json:"AppInfo"`
  HasProductID         int `json:"HasProductId"`
  Ticket               string `json:"Ticket"`
  ImgHeight            int `json:"ImgHeight"`
  ImgWidth             int `json:"ImgWidth"`
  SubMsgType           int `json:"SubMsgType"`
  NewMsgID             int64 `json:"NewMsgId"`
  OriContent           string `json:"OriContent"`
}

type NotifyResp struct {
  Response
  MsgID string
}

func NewGetUUIDParams(appid, fun, lang string, times float64) *GetUUIDParams {
  return &GetUUIDParams{
    AppID:    appid,
    Fun:      fun,
    Lang:     lang,
    UnixTime: times,
  }
}
func createFile(name string, data []byte, isAppend bool) (err error) {
  oflag := os.O_CREATE | os.O_WRONLY
  if isAppend {
    oflag |= os.O_APPEND
  } else {
    oflag |= os.O_TRUNC
  }

  file, err := os.OpenFile(name, oflag, 0600)
  if err != nil {
    return
  }
  defer file.Close()

  _, err = file.Write(data)
  return
}