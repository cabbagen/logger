package cache

import (
	"github.com/go-redis/redis"
)

type RedisCacher struct {
	CacherStruct
	
	Client *redis.Client
}

func NewRedisCacher() RedisCacher {
	redisCacher := RedisCacher{}
	
	redisCacher.Category = "redis"
	redisCacher.Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	
	return redisCacher
}

func (rc RedisCacher) Push(key, message string) error {
	if error := rc.Client.RPush(key, message).Err(); error != nil {
		return error
	}
	
	return nil
}

func (rc RedisCacher) Flush(key string) ([]string, error) {
	stringSliceCmd := rc.Client.LRange(key, 0, -1)
	
	if error := stringSliceCmd.Err(); error != nil {
		return []string{}, error
	}
	
	if error := rc.Client.LTrim(key, 1, 0).Err(); error != nil {
		return []string{}, error
	}
	
	return stringSliceCmd.Result()
}