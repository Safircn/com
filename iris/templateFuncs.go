package iris

import (
  "github.com/Safircn/com"
  "github.com/pquerna/ffjson/ffjson"
  "html/template"
  "strconv"
  "time"
  "os"
  "crypto/md5"
  "io"
  "encoding/hex"
  "fmt"
)

var fileUnix string
var staticList = make(map[string]string)

func init() {
  fileUnix = GetfileUnix()
}

func GetfileUnix() string {
  return "?unix=" + strconv.FormatInt(time.Now().Unix(), 16)
}

func GetStaticMD5(url string) string {
  if staticUrl, flag := staticList[url]; flag {
    return staticUrl
  }
  file, inerr := os.Open(com.Conf.DefaultString("template::staticPath", ".") + url)
  if inerr == nil {
    md5h := md5.New()
    io.Copy(md5h, file)
    newUrl := url + "?md5=" + hex.EncodeToString(md5h.Sum([]byte("")))
    staticList[url] = newUrl
    return newUrl
  }
  staticList[url] = url
  return url
}

var FuncMaps = map[string]interface{}{
  "script": func(url string) interface{} {
    if runmode != "pro" {
      url = com.Conf.DefaultString("template::devUrl", "http://localhost:4000") + GetStaticMD5(url)
    } else {
      url = GetStaticMD5(url)
    }
    return template.HTML(`<script src="` + url + `" defer="defer"></script>`)
  },
  "scriptAsync": func(url string) interface{} {
    if runmode != "pro" {
      url = com.Conf.DefaultString("template::devUrl", "http://localhost:4000") + GetStaticMD5(url)
    } else {
      url = GetStaticMD5(url)
    }
    return template.HTML(`<script>
    var asyncScript = document.createElement("script");
    asyncScript.src = '` + url + `';
    asyncScript.type = "text/javascript";
    asyncScript.async = "async";
    document.body.appendChild(asyncScript);
    </script>`)
  },
  "scriptAsyncs": func(urls ... string) interface{} {
    var urlStr, url string
    for _, v := range urls {
      if runmode != "pro" {
        url = com.Conf.DefaultString("template::devUrl", "http://localhost:4000") + GetStaticMD5(v)
      } else {
        url = GetStaticMD5(v)
      }
      urlStr += "'" + url + "',"
    }
    urlStr = urlStr[:(len(urlStr) - 1)]
    fmt.Println(urls)
    return template.HTML(`<script>
    window.onload = function(){
    var scriptList = [` + urlStr + `];
    for(var k in scriptList){
      var asyncScript = document.createElement("script");
      asyncScript.src = scriptList[k];
      asyncScript.type = "text/javascript";
      document.body.appendChild(asyncScript);
    }
    }
    </script>`)
  },
  "style": func(url string) interface{} {
    if runmode != "pro" {
      url = com.Conf.DefaultString("template::devUrl", "http://localhost:4000") + GetStaticMD5(url)
    } else {
      url = GetStaticMD5(url)
    }
    return template.HTML(`<link rel="stylesheet" href="` + url + `"/>`)
  },
  "reduxState": func(obj interface{}) template.HTML {
    bys, err := ffjson.Marshal(obj)
    if err != nil {
      return ""
    }
    text := "<script>var reduxState =" + string(bys) + "</script>"
    return template.HTML(text)
  },

  "jsonData": func(name string, obj interface{}) template.HTML {
    bys, err := ffjson.Marshal(obj)
    if err != nil {
      return ""
    }
    text := "<script>var " + name + " =" + string(bys) + "</script>"
    return template.HTML(text)
  },

}
