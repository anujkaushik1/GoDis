package core

import (
	config "github.com/anujkaushik1/GoDis/config"
)

func evictBatch(percent int) {
	count := len(store) * percent / 100
	if count < 1 {
		count = 1
	}

	deleted := 0
	for key := range store {
		delete(store, key)
		deleted++
		if deleted >= count {
			break
		}
	}
}

func Evict() {
	for float64(config.App.Memory) <= GetMemoryStats().Alloc {
		if len(store) == 0 {
			break
		}
		evictBatch(10) //10% keys will be deleted
	}
}
