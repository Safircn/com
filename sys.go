package com

import (
  "runtime"
  "fmt"
  "time"
)

var SysStatus struct {
  Uptime       string
  NumGoroutine int

                      // General statistics.
  MemAllocated string // bytes allocated and still in use
  MemTotal     string // bytes allocated (even if freed)
  MemSys       string // bytes obtained from system (sum of XxxSys below)
  Lookups      uint64 // number of pointer lookups
  MemMallocs   uint64 // number of mallocs
  MemFrees     uint64 // number of frees

                      // Main allocation heap statistics.
  HeapAlloc    string // bytes allocated and still in use
  HeapSys      string // bytes obtained from system
  HeapIdle     string // bytes in idle spans
  HeapInuse    string // bytes in non-idle span
  HeapReleased string // bytes released to the OS
  HeapObjects  uint64 // total number of allocated objects

                      // Low-level fixed-size structure allocator statistics.
                      //	Inuse is bytes used now.
                      //	Sys is bytes obtained from system.
  StackInuse  string // bootstrap stacks
  StackSys    string
  MSpanInuse  string // mspan structures
  MSpanSys    string
  MCacheInuse string // mcache structures
  MCacheSys   string
  BuckHashSys string // profiling bucket hash table
  GCSys       string // GC metadata
  OtherSys    string // other system allocations

                      // Garbage collector statistics.
  NextGC       string // next run in HeapAlloc time (bytes)
  LastGC       string // last run in absolute time (ns)
  PauseTotalNs string
  PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
  NumGC        uint32
}

var (
  startTime = time.Now()
)

func UpdateSystemStatus() {
  SysStatus.Uptime = TimeSincePro(startTime)

  m := new(runtime.MemStats)
  runtime.ReadMemStats(m)
  SysStatus.NumGoroutine = runtime.NumGoroutine()

  SysStatus.MemAllocated = FileSize(int64(m.Alloc))
  SysStatus.MemTotal = FileSize(int64(m.TotalAlloc))
  SysStatus.MemSys = FileSize(int64(m.Sys))
  SysStatus.Lookups = m.Lookups
  SysStatus.MemMallocs = m.Mallocs
  SysStatus.MemFrees = m.Frees

  SysStatus.HeapAlloc = FileSize(int64(m.HeapAlloc))
  SysStatus.HeapSys = FileSize(int64(m.HeapSys))
  SysStatus.HeapIdle = FileSize(int64(m.HeapIdle))
  SysStatus.HeapInuse = FileSize(int64(m.HeapInuse))
  SysStatus.HeapReleased = FileSize(int64(m.HeapReleased))
  SysStatus.HeapObjects = m.HeapObjects

  SysStatus.StackInuse = FileSize(int64(m.StackInuse))
  SysStatus.StackSys = FileSize(int64(m.StackSys))
  SysStatus.MSpanInuse = FileSize(int64(m.MSpanInuse))
  SysStatus.MSpanSys = FileSize(int64(m.MSpanSys))
  SysStatus.MCacheInuse = FileSize(int64(m.MCacheInuse))
  SysStatus.MCacheSys = FileSize(int64(m.MCacheSys))
  SysStatus.BuckHashSys = FileSize(int64(m.BuckHashSys))
  SysStatus.GCSys = FileSize(int64(m.GCSys))
  SysStatus.OtherSys = FileSize(int64(m.OtherSys))

  SysStatus.NextGC = FileSize(int64(m.NextGC))
  SysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
  SysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
  SysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
  SysStatus.NumGC = m.NumGC
}