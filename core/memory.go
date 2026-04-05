package core

import (
	"runtime"
)

type MemoryStats struct {
	Alloc      float64
	TotalAlloc float64
	Sys        float64
	NumGC      uint32
}

func GetMemoryStats() *MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &MemoryStats{
		Alloc:      float64(m.Alloc) / 1024 / 1024,
		TotalAlloc: float64(m.TotalAlloc) / 1024 / 1024,
		Sys:        float64(m.Sys) / 1024 / 1024,
		NumGC:      m.NumGC,
	}
}
