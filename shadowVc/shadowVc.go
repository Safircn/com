package shadowVc

import (
	"com"
	"bytes"
	"github.com/Safircn/lib/md5"
	"github.com/pquerna/ffjson/ffjson"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	V           = "2.0"
	FORMAT      = "json"
	SIGN_METHOD = "md5"
	SDK_VERSION = "ecbox-sdk-golang-20160601"
)

type ShadowVcSdk struct {
	serviceUrl string
	appKey     string
	appSecret  string

	httpClient *http.Client
}

func NewShadowVcSdk(serverUrl, appKey, appSecret string) *ShadowVcSdk {
	return &ShadowVcSdk{
		serviceUrl: serverUrl,
		appKey:     appKey,
		appSecret:  appSecret,
		httpClient: new(http.Client),
	}
}

func (sd *ShadowVcSdk) GetAppKey() string {
	return sd.appKey
}

func (sd *ShadowVcSdk) GetAppSecret() string {
	return sd.appSecret
}

var requestPool = sync.Pool{
	New: func() interface{} {
		res := &http.Request{
			Method:     "POST",
			URL:        nil,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       nil,
			Host:       "",
		}
		return res
	},
}

func (sd *ShadowVcSdk) Send(method string, postForm url.Values) (respBody []byte, err error) {
	//SERVER_URL
	request := requestPool.Get().(*http.Request)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	defer func() {
		request.Body.Close()
		requestPool.Put(request)
	}()

	urlvalue := url.Values{}
	urlvalue.Set("appKey", sd.appKey)
	urlvalue.Set("method", method)
	urlvalue.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	urlvalue.Set("v", V)
	urlvalue.Set("sdkVersion", SDK_VERSION)
	urlvalue.Set("format", FORMAT)
	urlvalue.Set("signMethod", SIGN_METHOD)

	valueForm := url.Values{}
	for k, v := range urlvalue {
		valueForm[k] = v
	}
	for k, v := range postForm {
		valueForm[k] = v
	}

	strs := sort.StringSlice{}
	for k, _ := range valueForm {
		strs = append(strs, k)
	}
	strs.Sort()

	makeString := sd.appSecret
	for _, v := range strs {
		if "@" != Substr(valueForm[v][0], 0, 1) {
			makeString += v + strings.Join(valueForm[v], ",")
		}
	}
	makeString += sd.appSecret

	urlvalue.Set("sign", strings.ToUpper(md5.Md5(makeString)))
	urlSu, err := url.Parse(sd.serviceUrl)
	if err != nil {
		return
	}

	request.Host = urlSu.Host
	urlSu.RawQuery = urlvalue.Encode()
	request.URL = urlSu

	urlvalue = url.Values{}
	rc := strings.NewReader(postForm.Encode())
	request.ContentLength = int64(rc.Len())
	request.Body = ioutil.NopCloser(rc)

	resp, err := sd.httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (sd *ShadowVcSdk) SendMultipart(method string, postForm url.Values, fieldname, filename string, fileItem []byte) (respBody []byte, err error) {
	//SERVER_URL
	request := requestPool.Get().(*http.Request)
	defer func() {
		request.Body.Close()
		requestPool.Put(request)
	}()
	buf := new(bytes.Buffer)
	defer buf.Reset()
	multWrier := multipart.NewWriter(buf)

	multWrier.CreateFormFile(fieldname, filename)
	buf.Write(fileItem)

	urlvalue := url.Values{}
	urlvalue.Set("appKey", sd.appKey)
	urlvalue.Set("method", method)
	urlvalue.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	urlvalue.Set("v", V)
	urlvalue.Set("sdkVersion", SDK_VERSION)
	urlvalue.Set("format", FORMAT)
	urlvalue.Set("signMethod", SIGN_METHOD)

	valueForm := url.Values{}
	for k, v := range urlvalue {
		valueForm[k] = v
	}
	for k, v := range postForm {
		valueForm[k] = v
		multWrier.WriteField(k, v[0])
	}

	strs := sort.StringSlice{}
	for k, _ := range valueForm {
		strs = append(strs, k)
	}
	strs.Sort()

	makeString := sd.appSecret
	for _, v := range strs {
		if "@" != Substr(valueForm[v][0], 0, 1) {
			makeString += v + strings.Join(valueForm[v], ",")
		}
	}
	makeString += sd.appSecret

	urlvalue.Set("sign", strings.ToUpper(md5.Md5(makeString)))
	urlSu, err := url.Parse(sd.serviceUrl)
	if err != nil {
		return
	}

	request.Host = urlSu.Host
	urlSu.RawQuery = urlvalue.Encode()
	request.URL = urlSu

	urlvalue = url.Values{}

	request.Header.Set("Content-Type", multWrier.FormDataContentType())
	multWrier.Close()

	request.ContentLength = int64(buf.Len())
	rc, ok := io.Reader(buf).(io.ReadCloser)
	if !ok && buf != nil {
		rc = ioutil.NopCloser(buf)
	}

	request.Body = rc

	resp, err := sd.httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func NewRespError(bt []byte) error {
	respError := new(RespError)
  com.Logger.Error("%s", string(bt))
	err := ffjson.Unmarshal(bt, respError)
	if err != nil {
		return err
	}
	return respError
}

type RespError struct {
	ErrorResponse struct {
		ErrorCode int    `json:"errorCode"`
		Msg       string `json:"msg"`
		SubCode   int    `json:"subCode"`
		SubMsg    string `json:"subMsg"`
	}
}

func (err *RespError) Error() string {
	return err.ErrorResponse.Msg + " " + err.ErrorResponse.SubMsg
}
