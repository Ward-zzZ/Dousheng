package initialize

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"tiktok-demo/cmd/relation/config"
	"tiktok-demo/shared/consts"
)

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalServerConfig.RedisInfo.Host, config.GlobalServerConfig.RedisInfo.Port),
		Password: config.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisRelationClientDB,
	})
}
