package cache


import (
  "testing"
  . "github.com/smartystreets/goconvey/convey"
  "fmt"
  "time"
)


func Test_Run(t *testing.T) {
  Convey("Get Run", t, func() {

    demo1 := NewCache("demo1")
    demo2 := NewCache("demo2")
    go Run(2)

    demo1.Set("demo1","111111111111",1)

    demo2.Set("demo2","111111111111",1)

    fmt.Println(demo1.Get("demo1"),demo1.IsExist("demo1"))
    fmt.Println(demo2.Get("demo2"))
    demo2.Delete("demo2")
    fmt.Println(demo2.Get("demo2"))
    demo1.Flush()
    fmt.Println(demo1.Get("demo1"))
    fmt.Println(demo1.Get("demo1"),demo1.IsExist("demo1"))
time.Sleep(1*time.Minute)

    fmt.Println(demo1.Get("demo1"))
    fmt.Println(demo2.Get("demo2"))


  })
}