package com

type TmpCacheCreateI64 func(int64) interface{}
type ITmpCacheInt64 interface {
  Get(int64) interface{}
}
type TmpCacheInt64 struct {
  list       map[int64]interface{}
  createFunc TmpCacheCreateI64
}

type TmpCacheCreateI func(int) interface{}
type ITmpCacheInt interface {
  Get(int) interface{}
}
type TmpCacheInt  struct {
  list       map[int]interface{}
  createFunc TmpCacheCreateI
}

type TmpCacheCreateStr func(string) interface{}
type ITmpCacheString interface {
  Get(string) interface{}
}
type TmpCacheString struct {
  list       map[string]interface{}
  createFunc TmpCacheCreateStr
}

func TmpCacheInt64New(c TmpCacheCreateI64) ITmpCacheInt64 {
  ti64 := new(TmpCacheInt64)
  ti64.list = make(map[int64]interface{})
  ti64.createFunc = c
  return ti64
}

func (t *TmpCacheInt64) Get(id int64) interface{} {
  if v, flag := t.list[id]; flag {
    return v
  }
  data := t.createFunc(id)
  t.list[id] = data
  return data
}

func TmpCacheIntNew(c TmpCacheCreateI) ITmpCacheInt {
  ti := new(TmpCacheInt)
  ti.list = make(map[int]interface{})
  ti.createFunc = c
  return ti
}

func (t *TmpCacheInt) Get(id int) interface{} {
  if v, flag := t.list[id]; flag {
    return v
  }
  data := t.createFunc(id)
  t.list[id] = data
  return data
}

func TmpCacheStringNew(c TmpCacheCreateStr) ITmpCacheString {
  tStr := new(TmpCacheString)
  tStr.list = make(map[string]interface{})
  tStr.createFunc = c
  return tStr
}

func (t *TmpCacheString) Get(code string) interface{} {
  if v, flag := t.list[code]; flag {
    return v
  }
  data := t.createFunc(code)
  t.list[code] = data
  return data
}
