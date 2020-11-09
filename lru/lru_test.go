package lru

import (
	"geecache"
	"testing"
)

type String string

func (s String) Len() int64  {
	return int64(len(s))
}

// get test -v lru_test.go lru.go
func TestGet(t *testing.T)  {
	lruCache := NewCache(1024, nil)

	lruCache.Set("key1", String("value1"))
	if v := lruCache.Get("key1"); v == nil || string(v.(String)) != "value1" {
		t.Errorf("except key1-value1, but got key1-%s", string(v.(String)))
	}

	if v := lruCache.Get("key2"); v != nil {
		t.Errorf("unexcept key key2")
	}
}

// 主要测试是否自动淘汰最近最少的 kv
func TestExpireOldest(t *testing.T)  {
	lruCache := NewCache(30, nil)

	// 4 + 28
	lruCache.Set("key1", String("this is value, cost 28 bytes"))
	// 4 + 28
	lruCache.Set("key2", String("this is value2, cost 29 bytes"))
	if v := lruCache.Get("key1"); v != nil {
		t.Errorf("key1 not expired")
	}
}

func TestDel(t *testing.T)  {
	lruCache := NewCache(1024, nil)
	lruCache.Set("key1", String("value1"))
	lruCache.Del("key1")
	if v := lruCache.Get("key1"); v != nil {
		t.Errorf("del key failed")
	}
}

func TestSet(t *testing.T)  {
	lruCache := NewCache(1024, nil)
	lruCache.Set("key1", String("value1"))
	memory := lruCache.Memory()
	lruCache.Set("key1", String("value1value1"))
	newMemory := lruCache.Memory()
	if newMemory - memory != 6 {
		t.Errorf("memory count has err")
	}
	if v := lruCache.Get("key1"); v == nil || string(v.(String)) != "value1value1" {
		t.Errorf("update value has err")
	}
}

func TestOnDeleted(t *testing.T)  {
	c := make(chan struct{}, 1)
	lruCache := NewCache(1024, func(key string, value geecache.Value) {
		// 删除 kv 的同时删除 key 的过期时间记录
		t.Logf("key:%s, value:%s has been deleted", key, string(value.(String)))
		c <- struct{}{}
	})
	lruCache.Set("key1", String("value1"))
	lruCache.Del("key1")
	select {
	case <- c:
	default:
		t.Errorf("func onDeleted has not been trigger")
	}
}
