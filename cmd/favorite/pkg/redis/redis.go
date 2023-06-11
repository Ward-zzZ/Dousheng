package redis

import (
	"context"
	"errors"
	"strconv"

	"tiktok-demo/shared/consts"

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

func (m *Manager) AddUserLikeList(c context.Context, uid int64, vid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_user_like_list"
	if err := m.client.SAdd(c, uidStr, vid).Err(); err != nil {
		return false, err
	}
	if err := m.client.Expire(c, uidStr, consts.RedisExpireTime).Err(); err != nil {
		return false, err
	}
	if err := m.AddVideoLikeCount(c, vid); err != nil {
		return false, err
	}
	return true, nil
}

func (m *Manager) DelUserLikeList(c context.Context, uid int64, vid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_user_like_list"
	if exit := m.client.SIsMember(c, uidStr, vid).Val(); exit {
		if err := m.client.SRem(c, uidStr, vid).Err(); err != nil {
			return false, err
		}
		if err := m.client.Expire(c, uidStr, consts.RedisExpireTime).Err(); err != nil {
			return false, err
		}
		if err := m.DelVideoLikeCount(c, vid); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("user like list not exist or video not in user like list")
}

func (m *Manager) AddVideoLikeCount(c context.Context, vid int64) error {
	vidStr := strconv.FormatInt(vid, 10) + "_video_like_count"
	if err := m.client.Incr(c, vidStr).Err(); err != nil {
		return err
	}
	if err := m.client.Expire(c, vidStr, consts.RedisExpireTime).Err(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DelVideoLikeCount(c context.Context, vid int64) error {
	vidStr := strconv.FormatInt(vid, 10) + "_video_like_count"
	if err := m.client.Decr(c, vidStr).Err(); err != nil {
		return err
	}
	if err := m.client.Expire(c, vidStr, consts.RedisExpireTime).Err(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) QueryUserLike(c context.Context, uid int64, vid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_user_like_list"
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
	uidStr := strconv.FormatInt(uid, 10) + "_user_like_list"
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

func (m *Manager) GetVideoLikeCount(c context.Context, vid int64) (int64, error) {
	vidStr := strconv.FormatInt(vid, 10) + "_video_like_count"
	if m.client.Exists(c, vidStr).Val() == 0 {
		return 0, errors.New("video like count not exist")
	}
	return m.client.Get(c, vidStr).Int64()
}
