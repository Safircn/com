package gin

import (
  "html/template"
  "github.com/gin-gonic/gin/render"
  "github.com/gin-gonic/gin"
  "bytes"
  "com"
)

type (
  HTMLDebug struct {
    Files    []string
    Glob     string
    FuncMaps template.FuncMap
  }
)

func LoadHTMLGlob(engine *gin.Engine, pattern string) {
  templ := template.New("gin")
  if len(FuncMaps) > 0 {
    templ.Funcs(template.FuncMap(FuncMaps))
  }
  templ = template.Must(templ.ParseGlob(pattern))
  if gin.IsDebugging() {
    debugPrintLoadTemplate(templ)
    htmlDebug := HTMLDebug{Glob: pattern}
    if len(FuncMaps) > 0 {
      htmlDebug.FuncMaps = template.FuncMap(FuncMaps)
    }
    engine.HTMLRender = htmlDebug
  } else {
    engine.HTMLRender = render.HTMLProduction{Template: templ}
  }
}

func LoadHTMLFiles(engine *gin.Engine, files ...string) {
  if gin.IsDebugging() {
    htmlDebug := HTMLDebug{Files: files}
    if len(FuncMaps) > 0 {
      htmlDebug.FuncMaps = template.FuncMap(FuncMaps)
    }
    engine.HTMLRender = htmlDebug
  } else {
    templ := template.New("gin")
    if len(FuncMaps) > 0 {
      templ.Funcs(template.FuncMap(FuncMaps))
    }
    templ = template.Must(templ.ParseFiles(files...))
    if len(FuncMaps) > 0 {
      templ.Funcs(template.FuncMap(FuncMaps))
    }
    engine.HTMLRender = render.HTMLProduction{Template: templ}
  }
}

func debugPrintLoadTemplate(tmpl *template.Template) {
  if gin.IsDebugging() {
    var buf bytes.Buffer
    for _, tmpl := range tmpl.Templates() {
      buf.WriteString("\t- ")
      buf.WriteString(tmpl.Name())
      buf.WriteString("\n")
    }
    com.Logger.Info("Loaded HTML Templates (%d): \n%s\n", len(tmpl.Templates()), buf.String())
  }
}

func (r HTMLDebug) Instance(name string, data interface{}) render.Render {
  return render.HTML{
    Template: r.loadTemplate(),
    Name:     name,
    Data:     data,
  }
}
func (r HTMLDebug) loadTemplate() *template.Template {
  templ := template.New("gin")
  if len(FuncMaps) > 0 {
    templ.Funcs(template.FuncMap(FuncMaps))
  }
  if len(r.Files) > 0 {
    return template.Must(templ.ParseFiles(r.Files...))
  }
  if len(r.Glob) > 0 {
    return template.Must(templ.ParseGlob(r.Glob))
  }
  panic("the HTML debug render was created without files or glob pattern")
}