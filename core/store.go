package core

import "time"

var store = make(map[string]*Obj)

type Obj struct {
	Value     any
	ExpiresAt int64
}

func Set(key string, value any, durationSec ...int64) {
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
