package qiniu

import (
	"fmt"
	"github.com/Safircn/lib/id"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
	goIo "io"
	"net/http"
	"shadowvc/core"
	"time"
)

var defaultImg, imgHttp, bucket string

func QiniuInit() {
	imgHttp = core.Iniconf.String("qiniu::IMG_HTTP")
	bucket = core.Iniconf.String("qiniu::BUCKET")
	defaultImg = imgHttp + "/" + core.Iniconf.DefaultString("wechat::DEFAULT_IMG", "234057113.jpg")
	conf.ACCESS_KEY = core.Iniconf.String("qiniu::ACCESS_KEY")
	conf.SECRET_KEY = core.Iniconf.String("qiniu::SECRET_KEY")
}

func uptoken(bucketName string) string {
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
		//CallbackUrl: callbackUrl,
		//CallbackBody:callbackBody,
		//ReturnUrl:   returnUrl,
		//ReturnBody:  returnBody,
		//AsyncOps:    asyncOps,
		//EndUser:     endUser,
		//Expires:     expires,
	}
	return putPolicy.Token(nil)
}

/**
上传网络图片 返回 key
*/
func ImgUpdate(key string, data goIo.Reader) (err error) {

	// data io.Reader
	var ret io.PutRet
	var extra = &io.PutExtra{
	//Params:    params,
	//MimeType:  mieType,
	//Crc32:     crc32,
	//CheckCrc:  CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	uptoken := uptoken(bucket)
	err = io.Put(nil, &ret, uptoken, key, data, extra)
	if err != nil {
		//上传产生错误
		core.Log.Error("io.Put failed:%s", err)
		return
	}
	//上传成功，处理返回值
	return
}

func WxHeadImgUpdateTokey(key, imgHtml string) string {
	if imgHtml == "" {
		return defaultImg
	}
	var (
		err  error
		resp *http.Response
	)
	startTime := time.Now()
	for i := 0; i < 3; i++ {

		fmt.Println(time.Now().Sub(startTime).Seconds(), "get start")

		resp, err = http.Get(imgHtml)
		fmt.Println(time.Now().Sub(startTime).Seconds(), "get end")
		if err != nil || resp.StatusCode != 200 {
			if i < 2 {
				continue
			}
			return defaultImg
		}
		fmt.Println(time.Now().Sub(startTime).Seconds(), "put start")
		err = ImgUpdate(key, resp.Body)
		fmt.Println(time.Now().Sub(startTime).Seconds(), "put end")
		resp.Body.Close()
		if err != nil {
			if i < 2 {
				continue
			}
			return defaultImg
		}
		break
	}

	return imgHttp + "/" + key
}

func WxHeadImgUpdate(imgHtml string) string {
	return WxHeadImgUpdateTokey(id.NewId(), imgHtml)
}
