package tools

import (
	"math/rand"
	"time"
)

func RedisExpireTime() time.Duration {
    //[10~30) 分钟的失效时间
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return time.Duration((10 + r.Int63n(20)) * int64(time.Minute))
}
