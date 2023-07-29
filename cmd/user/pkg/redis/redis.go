package redis

import (
	"context"
	"strconv"

	"tiktok-demo/cmd/user/pkg/mysql"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/errno"

	"github.com/go-redis/redis/v8"
)

type Manager struct {
	client *redis.Client
}

// NewManager creates a redis manager.
func NewManager(client *redis.Client) *Manager {
	return &Manager{client: client}
}

// delete the userInfo from redis,double delete
func (m *Manager) DeleteRelation(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := "userInfo:" + strconv.FormatInt(uid, 10)
	tidStr := "userInfo:" + strconv.FormatInt(tid, 10)
	if err := m.client.Del(c, uidStr).Err(); err != nil {
		return false, err
	}
	if err := m.client.Del(c, tidStr).Err(); err != nil {
		return false, err
	}
	m.client.Expire(c, uidStr, consts.RedisExpireTime)
	m.client.Expire(c, tidStr, consts.RedisExpireTime)
	return true, nil
}

// get user info hash from redis
func (m *Manager) GetUserInfo(c context.Context, uid int64) (*mysql.User, error) {
	uidStr := "userInfo:" + strconv.FormatInt(uid, 10)
	cache, err := m.client.HGetAll(c, uidStr).Result()
	if err != nil {
		return nil, err
	}
	if len(cache) == 0 {
		return nil, errno.UserNotExistErr
	}
	FollowerNumUint64, _ := strconv.ParseUint(cache["FollowerCount"], 10, 64)
	FollowingNumUint64, _ := strconv.ParseUint(cache["FollowingCount"], 10, 64)
	Uid, _ := strconv.ParseInt(cache["Uid"], 10, 64)
	Name := cache["Name"]
	return &mysql.User{
		Uid:            Uid,
		Name:           Name,
		FollowingCount: FollowingNumUint64,
		FollowerCount:  FollowerNumUint64,
	}, nil
}

// insert user info into redis
func (m *Manager) InsertUserInfo(c context.Context, u *mysql.User) error {
	uidStr := "userInfo:" + strconv.FormatInt(u.Uid, 10)

	if err := m.client.HSet(c, uidStr, map[string]interface{}{
		"Uid":            u.Uid,
		"Name":           u.Name,
		"FollowingCount": u.FollowingCount,
		"FollowerCount":  u.FollowerCount,
	}).Err(); err != nil {
		return err
	}
	m.client.Expire(c, uidStr, consts.RedisExpireTime)
	return nil
}
