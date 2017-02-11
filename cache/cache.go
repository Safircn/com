package cache

import (
  "time"
  "sync"
)

func init() {
  caches = make(map[string]Cache)
}

var caches map[string]Cache
var _ Cache = (*CacheBox)(nil)

type CacheBox struct {
  name      string
  rwMutex   *sync.RWMutex
  cacheList map[string]*cacheRow
}

type cacheRow struct {
  expiredTime *time.Time
  data        interface{}
}

func (this *CacheBox) Set(key string, val interface{}, timeout int64) {
  if (timeout == 0) {
    return
  }
  this.rwMutex.Lock()
  expiredTime := time.Now().Add(time.Minute * time.Duration(timeout))
  this.cacheList[key] = &cacheRow{
    expiredTime:&expiredTime,
    data:val,
  }
  this.rwMutex.Unlock()
}

func (this *CacheBox) Delete(key string) {
  this.rwMutex.Lock()
  if _, flag := this.cacheList[key]; flag {
    delete(this.cacheList, key)
  }
  this.rwMutex.Unlock()
}

func (this *CacheBox) Get(key string) interface{} {
  this.rwMutex.RLock()
  defer this.rwMutex.RUnlock()
  if v, flag := this.cacheList[key]; flag {
    if v.expiredTime.After(time.Now()) {
      return v.data
    }
  }
  return nil
}

func (this *CacheBox) GetString(key string) string {
  if iv := this.Get(key); iv != nil {
    if vStr, flag := iv.(string); flag {
      return vStr
    }
  }
  return ""
}

func (this *CacheBox) GetInt(key string) int {
  if iv := this.Get(key); iv != nil {
    if vInt, flag := iv.(int); flag {
      return vInt
    }
  }
  return 0
}

func (this *CacheBox) GetInt64(key string) int64 {
  if iv := this.Get(key); iv != nil {
    if vInt64, flag := iv.(int64); flag {
      return vInt64
    }
  }
  return 0
}

func (this *CacheBox) IsExist(key string) bool {
  this.rwMutex.RLock()
  defer this.rwMutex.RUnlock()
  if v, flag := this.cacheList[key]; flag {
    if v.expiredTime.After(time.Now()) {
      return true
    }
  }
  return false
}

func (this *CacheBox) Flush() {
  this.rwMutex.Lock()
  this.cacheList = make(map[string]*cacheRow)
  this.rwMutex.Unlock()
}

type Cache interface {
  // Put puts value into cache with key and expire time.
  Set(key string, val interface{}, timeout int64)
  // Get gets cached value by given key.
  Get(key string) interface{}
  GetString(key string) string
  GetInt(key string) int
  GetInt64(key string) int64
  // Delete deletes cached value by given key.
  Delete(key string)
  // IsExist returns true if cached value exists.
  IsExist(key string) bool
  // Flush deletes all cached data.
  Flush()
}

func NewCache(name string) *CacheBox {
  newCache := &CacheBox{
    name:name,
    cacheList:make(map[string]*cacheRow),
    rwMutex:new(sync.RWMutex),
  }
  caches[name] = newCache
  return newCache
}

func Run(Interval int) {
  if len(caches) == 0 {
    return
  }
  if Interval == 0 {
    Interval = 1
  }
  var cacheBox *CacheBox
  for {
    time.Sleep(time.Minute * time.Duration(Interval))

    for _, v := range caches {
      cacheBox = v.(*CacheBox)
      cacheBox.rwMutex.Lock()
      for kk, vv := range cacheBox.cacheList {
        if vv.expiredTime.Sub(time.Now()) < -1 {
          delete(v.(*CacheBox).cacheList, kk)
        }
      }
      cacheBox.rwMutex.Unlock()
    }
  }
}

func RunTntervalTenM()  {
  Run(10)
}