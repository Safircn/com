package com

func InArrayStr(arr []string, s string) bool {
  if len(arr) == 0 {
    return false
  }
  for _, v := range arr {
    if v == s {
      return true
    }
  }
  return false
}
