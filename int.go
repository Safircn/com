package com

import (
  "strings"
  "strconv"
)

func SafeSplitInt(s, sep string) []int {
  sl := strings.Split(s, sep)
  slleng := len(sl)
  if slleng == 1 && sl[0] == "" {
    return make([]int, 0)
  }
  liInt := make([]int, slleng)
  for k, v := range sl {
    liInt[k], _ = strconv.Atoi(v)
  }
  return liInt
}

func JoinInt(is []int, sep string) string {
  if len(is) == 0 {
    return ""
  }
  var joinS string
  for _, v := range is {
    joinS += strconv.Itoa(v) + sep
  }
  return joinS[0:(len(joinS) - len(sep))]
}

func SafeSplitInt64(s, sep string) []int64 {
  sl := strings.Split(s, sep)
  slleng := len(sl)
  if slleng == 1 && sl[0] == "" {
    return make([]int64, 0)
  }
  liInt := make([]int64, slleng)
  for k, v := range sl {
    liInt[k], _ = strconv.ParseInt(v, 10, 64)
  }
  return liInt
}

func JoinInt64(is []int64, sep string) string {
  if len(is) == 0 {
    return ""
  }
  var joinS string
  for _, v := range is {
    joinS += strconv.FormatInt(v, 10) + sep
  }
  return joinS[0:(len(joinS) - len(sep))]
}