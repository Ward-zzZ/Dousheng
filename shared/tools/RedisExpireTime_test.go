package tools_test

import (
    "testing"
    "time"

    "tiktok-demo/shared/tools"

)

func TestRedisExpireTime(t *testing.T) {
    // 调用被测试函数
    expireTime := tools.RedisExpireTime()

    // 检查返回值是否符合预期
    if expireTime < 10*time.Minute || expireTime >= 30*time.Minute {
        t.Errorf("RedisExpireTime() = %v, want [10m, 30m)", expireTime)
    }
}
