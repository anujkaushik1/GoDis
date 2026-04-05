package core

import (
	"fmt"
	"os"
	"time"

	config "github.com/anujkaushik1/GoDis/config"
)

var store = make(map[string]*Obj)

type Obj struct {
	Value     any
	ExpiresAt int64
}

func Set(key string, value any, durationSec ...int64) {

	haha := false
	if float64(config.App.Memory) <= GetMemoryStats().Alloc {
		Evict()
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("EVICCCTEDDDDD ==== ")
		fmt.Println(float64(config.App.Memory))
		fmt.Println(GetMemoryStats().Alloc)
		fmt.Println(len(store))
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")

		haha = true

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

	if haha {
		fmt.Println("")
		fmt.Println("----------")
		fmt.Println("EVICCCTEDDDDD ==== ")
		fmt.Println(float64(config.App.Memory))
		fmt.Println(GetMemoryStats().Alloc)
		fmt.Println(len(store))
		fmt.Println("")
		fmt.Println("----------")
		os.Exit(1)
	}
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
