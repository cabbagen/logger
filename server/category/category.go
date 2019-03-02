package category

import (
	"time"
	"logger/server/cache"
)

type LoggerStruct struct {
	Appid         string
	Level         string
	IsCache       bool             // 是否通过 redis 缓存
	LimitTime     time.Duration    // 通过 redis 缓存时间
	Cacher        cache.Cacher
}

// 打印日志
type Logger interface {
	Log(message string)
}

// 处理缓存相关的方法
type HandleCacher interface {
	AppendToCache(message string)
	ListenSyncTask()
	ExecCacheRecord()
}