package r2m

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

//单点redis test git tool with vscode
var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:        "39.105.153.230:6379",
		PoolSize:    512,
		PoolTimeout: time.Second * time.Duration(5),
	})
}

func SaveCache(key, data string) {
	_, err := RedisClient.Set(key, data, time.Minute).Result()
	if err != nil {
		fmt.Println("redis save error", err)
		return
	}
	fmt.Printf("key[%s],data[%s] insert to redis \n", key, data)
}
