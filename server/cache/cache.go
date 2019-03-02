package cache

// 定义缓存结构
type CacherStruct struct {
	Category       string
}

// 处理日志缓存
type Cacher interface {
	Push(key, message string) error;
	Flush(key string) ([]string, error);
}
