package img

import (
  "com/shadowVc"
  "github.com/pquerna/ffjson/ffjson"
  "net/url"
  "strconv"
)

type RespImgUpload struct {
  ImgUploadResponse *struct {
    Img struct {
          CompanyID  int      `json:"companyId"`
          Format     string   `json:"format"`
          MemberID   int      `json:"memberId"`
          Name       string   `json:"name"`
          Size       int      `json:"size"`
          StorageKey string   `json:"storageKey"`
          Tags       []string `json:"tags"`
          URL        string   `json:"url"`
          WxOpenID   string   `json:"wxOpenId"`
        } `json:"img"`
  }
}

type (
  UploadImgConfig struct {
    CompanyId int64
    MemberId  int64
    WxOpenId  string
    Tag       []string
  }
)

func UploadImg(shadowvcSdk *shadowVc.ShadowVcSdk, storageType, fileName string, fileItem []byte, uploadImgConfig *UploadImgConfig) (*RespImgUpload, error) {
  postData := make(url.Values)
  postData.Add("storageType", storageType)
  if uploadImgConfig == nil {
    uploadImgConfig = &UploadImgConfig{
      0, 0, "", make([]string, 0),
    }
  }
  postData.Add("companyId", strconv.FormatInt(uploadImgConfig.CompanyId, 10))
  postData.Add("memberId", strconv.FormatInt(uploadImgConfig.MemberId, 10))
  postData.Add("wxOpenId", uploadImgConfig.WxOpenId)
  tagsBys, err := ffjson.Marshal(uploadImgConfig.Tag)
  if err != nil {
    return nil, err
  }
  postData.Add("tags", string(tagsBys))

  bt, err := shadowvcSdk.SendMultipart("sd.img.upload", postData, "fileItem", fileName, fileItem)
  if err != nil {
    return nil, err
  }
  respImgUpload := new(RespImgUpload)
  err = ffjson.Unmarshal(bt, respImgUpload)
  if err != nil {
    return nil, err
  }
  if respImgUpload.ImgUploadResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respImgUpload, nil
}

type RespImgGet struct {
  ImgGetResponse *struct {
    Img struct {
          CompanyID  int           `json:"companyId"`
          Format     string        `json:"format"`
          MemberID   int           `json:"memberId"`
          Name       string        `json:"name"`
          Size       int           `json:"size"`
          StorageKey string        `json:"storageKey"`
          Tags       []interface{} `json:"tags"`
          URL        string        `json:"url"`
          WxOpenID   string        `json:"wxOpenId"`
        } `json:"img"`
  } `json:"ImgGetResponse"`
}

func GetImg(shadowvcSdk *shadowVc.ShadowVcSdk, companyId int64, storageKey, pathVariable string) (*RespImgGet, error) {
  postData := make(url.Values)
  postData.Add("companyId", strconv.FormatInt(companyId, 10))
  postData.Add("storageKey", storageKey)
  postData.Add("pathVariable", pathVariable)
  bt, err := shadowvcSdk.Send("sd.img.get", postData)
  if err != nil {
    return nil, err
  }
  respImgGet := new(RespImgGet)
  err = ffjson.Unmarshal(bt, respImgGet)
  if err != nil {
    return nil, err
  }

  if respImgGet.ImgGetResponse == nil {
    return nil, shadowVc.NewRespError(bt)
  }
  return respImgGet, nil
}
