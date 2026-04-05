package core

import (
	"fmt"
	"time"

	config "github.com/anujkaushik1/GoDis/config"
)

var store = make(map[string]*Obj)

type Obj struct {
	Value     any
	ExpiresAt int64
}

func Set(key string, value any, durationSec ...int64) {

	if float64(config.App.Memory) <= GetMemoryStats().Alloc {
		fmt.Printf("Memory limit reached, preparing for eviction. Current memory allocated: %f\n", GetMemoryStats().Alloc)
		Evict()
		fmt.Println("Memory evicted.")
	}

	expiresAt := int64(-1)

	if len(durationSec) > 0 && durationSec[0] > 0 {
		expiresAt = time.Now().UnixMilli() + durationSec[0]*1000
	}

	obj := Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	store[key] = &obj
}

func Get(key string) *Obj {
	return store[key]
}

func Del(key string) bool {
	_, ok := store[key]
	if ok {
		delete(store, key)
		return true
	}
	return false
}
