package main

type Value interface {
	// 返回值所在内存大小
	Len() int64
}

type Cache interface {
	// todo 支持 ttl
	Get(key string) (Value, bool)
	Set(key string, value Value)
}

