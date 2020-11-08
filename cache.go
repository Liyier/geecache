package geecache

type Value interface {
	// 返回值所在内存大小
	Len() int64
}

type Cache interface {
	// todo 支持 ttl
	// 查找不到时， value 为 nil
	Get(key string) Value
	Set(key string, value Value)
}

type Entry struct {
	Key   string
	Value Value
}
