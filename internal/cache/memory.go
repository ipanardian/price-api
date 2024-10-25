package cache

import (
	"time"

	memory "github.com/patrickmn/go-cache"
)

type MemoryCacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

var memClient *memory.Cache

func InitMemory(conf MemoryCacheConfig) {
	memClient = memory.New(conf.DefaultExpiration, conf.CleanupInterval)
}

func MemGet(k string) (interface{}, bool) {
	return memClient.Get(k)
}

func MemGetWithExpiration(k string) (interface{}, time.Time, bool) {
	return memClient.GetWithExpiration(k)
}

func MemSet(k string, x interface{}, d time.Duration) {
	memClient.Set(k, x, d)
}

func MemDelete(k string) {
	memClient.Delete(k)
}

func MemCLient() *memory.Cache {
	return memClient
}
