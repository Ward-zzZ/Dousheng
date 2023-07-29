package redis

import (
	"context"
	"errors"
	"strconv"

	"tiktok-demo/shared/tools"

	"github.com/go-redis/redis/v8"
)

type Manager struct {
	client *redis.Client
}

// NewManager creates a redis manager.
func NewManager(client *redis.Client) *Manager {
	return &Manager{client: client}
}

// Redis key design:
// User like list: set
// Video like count: int

func (m *Manager) SetUserLikeList(c context.Context, uid int64, vids []int64) error {
	uidStr := "userLikeList:" + strconv.FormatInt(uid, 10)
	pipeline := m.client.Pipeline()
	for _, vid := range vids {
		pipeline.SAdd(c, uidStr, vid)
	}
	pipeline.Expire(c, uidStr, tools.RedisExpireTime())
	_, err := pipeline.Exec(c)
	return err
}

// 先判断是否存在，如果不存在则跳过redis更新
func (m *Manager) AddUserLikeList(c context.Context, uid int64, vid int64) (bool, error) {
	uidStr := "userLikeList:" + strconv.FormatInt(uid, 10)
	// 先判断是否存在，再添加
	if m.client.Exists(c, uidStr).Val() == 0 {
		return false, errors.New("redis key not exist")
	}
	// 向关注列表中添加关注用户
	if err := m.client.SAdd(c, uidStr, vid).Err(); err != nil {
		return false, err
	}
	if err := m.client.Expire(c, uidStr, tools.RedisExpireTime()).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (m *Manager) DelUserLikeList(c context.Context, uid int64, vid int64) (bool, error) {
	uidStr := "userLikeList:" + strconv.FormatInt(uid, 10)
	if m.client.Exists(c, uidStr).Val() == 0 {
		return false, errors.New("redis key not exist")
	}
	if exit := m.client.SIsMember(c, uidStr, vid).Val(); exit {
		if err := m.client.SRem(c, uidStr, vid).Err(); err != nil {
			return false, err
		}
		if err := m.client.Expire(c, uidStr, tools.RedisExpireTime()).Err(); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("user like list not exist or video not in user like list")
}

func (m *Manager) QueryUserLike(c context.Context, uid int64, vid int64) (bool, error) {
	uidStr := "userLikeList:" + strconv.FormatInt(uid, 10)
	if m.client.Exists(c, uidStr).Val() == 0 {
		return false, errors.New("user like list not exist")
	}
	res, err := m.client.SIsMember(c, uidStr, vid).Result()
	if err != nil {
		return false, err
	}
	return res, nil
}

func (m *Manager) GetUserLikeList(c context.Context, uid int64) ([]int64, error) {
	uidStr := "userLikeList:" + strconv.FormatInt(uid, 10)
	if m.client.Exists(c, uidStr).Val() == 0 {
		return nil, errors.New("user like list not exist")
	}
	vals := m.client.SMembers(c, uidStr).Val()
	res := make([]int64, len(vals))
	for i, val := range vals {
		res[i], _ = strconv.ParseInt(val, 10, 64)
	}
	return res, nil
}

func (m *Manager) SetVideoLikeCount(c context.Context, vid int64, count int64) error {
	vidStr := "videoLikeCount:" + strconv.FormatInt(vid, 10)
	tx := m.client.TxPipeline()
	defer tx.Close()

	tx.Set(c, vidStr, count, 0)
	tx.Expire(c, vidStr, tools.RedisExpireTime())
	_, err := tx.Exec(c)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) AddVideoLikeCount(c context.Context, vid int64) error {
	vidStr := "videoLikeCount:" + strconv.FormatInt(vid, 10)
	if m.client.Exists(c, vidStr).Val() == 0 {
		return errors.New("video like count not exist")
	}
	if err := m.client.Incr(c, vidStr).Err(); err != nil {
		return err
	}
	if err := m.client.Expire(c, vidStr, tools.RedisExpireTime()).Err(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DelVideoLikeCount(c context.Context, vid int64) error {
	vidStr := "videoLikeCount:" + strconv.FormatInt(vid, 10)
	tx := m.client.TxPipeline()
	defer tx.Close()
	if m.client.Exists(c, vidStr).Val() != 0 {
		tx.Decr(c, vidStr)
		tx.Expire(c, vidStr, tools.RedisExpireTime())
		_, err := tx.Exec(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) GetVideoLikeCount(c context.Context, vid int64) (int64, error) {
	vidStr := "videoLikeCount:" + strconv.FormatInt(vid, 10)
	if m.client.Exists(c, vidStr).Val() == 0 {
		return 0, errors.New("video like count not exist")
	}
	return m.client.Get(c, vidStr).Int64()
}
