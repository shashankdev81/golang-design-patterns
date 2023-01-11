package main

import (
	"testing"
	"time"
)

type cacheItem struct {
	data      string
	expire_at time.Time
}

type stringCache struct {
	m   map[string]cacheItem
	exp time.Duration
}

func (sc *stringCache) Get(key string) (string, bool) {
	if item, ok := sc.m[key]; !ok {
		return "", false
	} else {
		return item.data, true
	}
}

func (sc *stringCache) Put(key, data string) {
	sc.m[key] = cacheItem{
		data:      data,
		expire_at: time.Now().Add(sc.exp),
	}
}

func NewStringCache(d time.Duration) *stringCache {
	return &stringCache{
		m:   make(map[string]cacheItem),
		exp: d,
	}
}


func TestStringCache(t *testing.T) {
	cache := NewStringCache(time.Minute)

	ch := make(chan struct{})

	go func() {
		cache.Put("here", "this")
		close(ch)
	}()

	_, _ = cache.Get("here")

	<-ch
}