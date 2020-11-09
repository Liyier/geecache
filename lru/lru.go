package lru

import (
	"container/list"
	"geecache"
)

// todo 增加并发安全控制
type Cache struct {
	// 允许缓存组件占用的最大内存
	maxBytes int64
	// 当前缓存组件已经占用的内存
	nBytes int64
	// 实现 lru 的双向链表, 队尾是最近访问的元素
	ll *list.List
	// 存储 kv 的 map
	provider map[string]*list.Element
	// 当元素被删除时，触发的回调函数
	OnDeleted func(key string, cache geecache.Value)
}

func NewCache(maxBytes int64, onDeleted func(key string, value geecache.Value)) geecache.Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		provider:  make(map[string]*list.Element),
		OnDeleted: onDeleted,
	}
}

func (c *Cache) Get(key string) geecache.Value {
	if e := c.get(key); e != nil {
		kv := e.Value.(*geecache.Entry)
		return kv.Value
	}
	return nil
}

func (c *Cache) Set(key string, value geecache.Value) {
	if e := c.get(key); e != nil {
		// 更新当前使用内存
		kv := e.Value.(*geecache.Entry)
		c.nBytes += value.Len() - kv.Value.Len()
		// 更新值
		kv.Value = value
	} else {
		e := c.ll.PushBack(&geecache.Entry{Key: key, Value: value})
		c.provider[key] = e
		c.nBytes += value.Len() + int64(len(key))
	}
	c.expireOldest()
}

func (c *Cache) Del(key string) {
	if e := c.get(key); e != nil {
		c.del(e)
	}
}

func (c *Cache) Size() int {
	return c.ll.Len()
}

func (c *Cache) Memory() int64 {
	return c.nBytes
}

func (c *Cache) expireOldest() {
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		// 淘汰队首元素
		e := c.ll.Front()
		c.del(e)
	}
}

func (c *Cache) del(e *list.Element) {
	kv := e.Value.(*geecache.Entry)
	key := kv.Key
	c.ll.Remove(e)
	delete(c.provider, key)
	c.nBytes -= kv.Value.Len() + int64(len(key))
	if c.OnDeleted != nil {
		c.OnDeleted(key, kv.Value)
	}
}

func (c *Cache) get(key string) *list.Element {
	if e, ok := c.provider[key]; ok {
		c.ll.MoveToBack(e)
		return e
	} else {
		return nil
	}
}
