package qqmap

import (
  "net/http"
  "net/url"
  "strconv"
  "fmt"
  "strings"
  "com"
  "io/ioutil"
  "github.com/pquerna/ffjson/ffjson"
  "errors"
  "github.com/Safircn/lib/md5"
  "sort"
)

const (
  apiUrl = "http://apis.map.qq.com"
  gcoderUrl = "/ws/geocoder/v1"
)

type QQMap struct {
  key        string
  sk         string
  httpClient *http.Client
}

func (q *QQMap) GetKey() string {
  return q.key
}

func NewQQMap(key, sk string, httpClient *http.Client) *QQMap {
  if httpClient == nil {
    httpClient = new(http.Client)
  }
  return &QQMap{
    key:key,
    httpClient:httpClient,
    sk:sk,
  }
}

type GencoderInfoCond struct {
  CoordType  int
  PoiOptions *GencoderInfoCondPoi
}

type GencoderInfoCondPoi struct {
  AddressFormat string
  Radius        int
  PageSize      int
  PageIndex     int
  Policy        int
  Category      []string
}

func (q *QQMap) caculateAKSN(urlStr string, params url.Values) string {
  var (
    encodeStr string
    encodeStrEscape string
    prefix string
    valueEscape string
  )
  keys := make([]string, 0, len(params))
  for k := range params {
    keys = append(keys, k)
  }
  sort.Strings(keys)
  for _, k := range keys {
    vs := params[k]
    prefix = url.QueryEscape(k)
    for _, v := range vs {
      if encodeStr != "" {
        encodeStr += "&"
        encodeStrEscape += "%26"
      }
      encodeStr += prefix + "="
      encodeStrEscape += prefix + "%3D"
      valueEscape = url.QueryEscape(v)
      encodeStr += valueEscape
      encodeStrEscape += valueEscape
    }
  }

  return urlStr + "?" + encodeStr + "&sn=" + md5.Md5(url.QueryEscape(urlStr + "?") + encodeStrEscape + q.sk)
}

func (q *QQMap) get(url string, params url.Values, respDate interface{}) error {

  resp, err := q.httpClient.Get(apiUrl + q.caculateAKSN(url, params))
  if err != nil {
    com.Logger.Error("url:%s,err:%s", url, err.Error())
    return err
  }
  if resp.StatusCode != 200 {
    com.Logger.Error("url:%s,err:%s", url, "status != 200")
    return errors.New("status != 200")
  }

  bys, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    com.Logger.Error("url:%s,err:%s", url, "body EOF")
    return errors.New("Body EOF")
  }

  err = ffjson.Unmarshal(bys, respDate)
  if err != nil {
    com.Logger.Error("url:%s,err:%s", url, err.Error())
    return err
  }

  return nil
}

//http://lbs.qq.com/webservice_v1/guide-gcoder.html
func (q *QQMap) Gcoder(lat, lng float64, gencoderInfoCond *GencoderInfoCond) (*RespGcoder, error) {
  getForm := make(url.Values)
  getForm.Set("key", q.key)
  getForm.Set("location", strconv.FormatFloat(lat, 'f', -1, 64) + "," + strconv.FormatFloat(lng, 'f', -1, 64))
  if gencoderInfoCond != nil {
    getForm.Set("coord_type", strconv.Itoa(gencoderInfoCond.CoordType))
    if gencoderInfoCond.PoiOptions != nil {
      getForm.Set("get_poi", "1")
      poiOptionsStr := fmt.Sprintf("poi_options=address_format=%s;radius=%d;page_size=%d;page_index=%d;policy=%d",
        gencoderInfoCond.PoiOptions.AddressFormat, gencoderInfoCond.PoiOptions.Radius,
        gencoderInfoCond.PoiOptions.PageSize, gencoderInfoCond.PoiOptions.PageIndex, gencoderInfoCond.PoiOptions.Policy)
      if len(gencoderInfoCond.PoiOptions.Category) > 0 {
        poiOptionsStr += ";category=" + strings.Join(gencoderInfoCond.PoiOptions.Category, ",")
      }
      getForm.Set("poi_options", poiOptionsStr)
    }
  }
  respGcoder := new(RespGcoder)
  err := q.get(gcoderUrl, getForm, respGcoder)
  if err != nil {
    return nil, err
  }
  if respGcoder.Status != 0 {
    return nil, errors.New(respGcoder.Message)
  }
  return respGcoder, nil
}

type RespGcoder struct {
  Status    int `json:"status"`
  Message   string `json:"message"`
  RequestID string `json:"request_id"`
  Result    *struct {
    Location           struct {
                         Lat float64 `json:"lat"`
                         Lng float64 `json:"lng"`
                       } `json:"location"`
    Address            string `json:"address"`
    FormattedAddresses struct {
                         Recommend string `json:"recommend"`
                         Rough     string `json:"rough"`
                       } `json:"formatted_addresses"`
    AddressComponent   struct {
                         Nation       string `json:"nation"`
                         Province     string `json:"province"`
                         City         string `json:"city"`
                         District     string `json:"district"`
                         Street       string `json:"street"`
                         StreetNumber string `json:"street_number"`
                       } `json:"address_component"`
    AdInfo             struct {
                         Adcode   string `json:"adcode"`
                         Name     string `json:"name"`
                         Location struct {
                                    Lat float64 `json:"lat"`
                                    Lng float64 `json:"lng"`
                                  } `json:"location"`
                         Nation   string `json:"nation"`
                         Province string `json:"province"`
                         City     string `json:"city"`
                         District string `json:"district"`
                       } `json:"ad_info"`
    AddressReference   struct {
                         BusinessArea AddressReferenceRow `json:"business_area"`
                         FamousArea   AddressReferenceRow `json:"famous_area"`
                         Crossroad    AddressReferenceRow `json:"crossroad"`
                         Village      AddressReferenceRow`json:"village"`
                         Town         AddressReferenceRow `json:"town"`
                         Street       AddressReferenceRow `json:"street"`
                         LandmarkL1   AddressReferenceRow `json:"landmark_l1"`
                         LandmarkL2   AddressReferenceRow `json:"landmark_l2"`
                       } `json:"address_reference"`
    PoiCount           int `json:"poi_count"`
    Pois               *[]struct {
      ID       string `json:"id"`
      Title    string `json:"title"`
      Address  string `json:"address"`
      Category string `json:"category"`
      Location struct {
                 Lat float64 `json:"lat"`
                 Lng float64 `json:"lng"`
               } `json:"location"`
      AdInfo   struct {
                 Adcode   string `json:"adcode"`
                 Province string `json:"province"`
                 City     string `json:"city"`
                 District string `json:"district"`
               } `json:"ad_info"`
      Distance float64 `json:"_distance"`
      DirDesc  string `json:"_dir_desc"`
    } `json:"pois"`
  } `json:"result"`
}
type AddressReferenceRow struct {
  Title    string `json:"title"`
  Location struct {
             Lat float64 `json:"lat"`
             Lng float64 `json:"lng"`
           } `json:"location"`
  Distance float64 `json:"_distance"`
  DirDesc  string `json:"_dir_desc"`
}


//http://lbs.qq.com/webservice_v1/guide-geocoder.html
func (q *QQMap) Geocoder(address string) (*RespGeocoder, error) {
  getForm := make(url.Values)
  getForm.Set("key", q.key)
  getForm.Set("address", url.QueryEscape(address))
  respGeocoder := new(RespGeocoder)
  err := q.get(gcoderUrl, getForm, respGeocoder)
  if err != nil {
    return nil, err
  }
  if respGeocoder.Status != 0 {
    return nil, errors.New(respGeocoder.Message)
  }
  return respGeocoder, nil
}

type RespGeocoder struct {
  Status  int `json:"status"`
  Message string `json:"message"`
  Result  *struct {
    Title             string `json:"title"`
    Location          struct {
                        Lng float64 `json:"lng"`
                        Lat float64 `json:"lat"`
                      } `json:"location"`
    AddressComponents struct {
                        Province     string `json:"province"`
                        City         string `json:"city"`
                        District     string `json:"district"`
                        Street       string `json:"street"`
                        StreetNumber string `json:"street_number"`
                      } `json:"address_components"`
    Similarity        float64 `json:"similarity"`
    Deviation         int `json:"deviation"`
    Reliability       int `json:"reliability"`
  } `json:"result"`
}