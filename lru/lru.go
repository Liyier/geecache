package lru

import (
	"container/list"
)


type Cache struct {
	// 允许缓存组件占用的最大内存
	maxBytes int64
	// 当前缓存组件已经占用的内存
	nBytes   int64
	// 实现 lru 的双向链表
	ll       *list.List
	// 存储 kv 的 map
	provider    map[string]*list.Element
	// 当元素被删除时，触发的回调函数
	//OnDeleted func(key string, cache Value)
}



