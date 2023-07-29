package redis

import (
	"context"
	"errors"
	"strconv"

	"tiktok-demo/shared/consts"
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

// update follow list to uid's follow list in redis using set
func (m *Manager) AddFollow(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := "followList:" + strconv.FormatInt(uid, 10)
	tidStr := "fansList:" + strconv.FormatInt(tid, 10)
	pipeline := m.client.TxPipeline()
	if m.client.Exists(c, uidStr).Val() != 0 {
		pipeline.SAdd(c, uidStr, tid)
		pipeline.Expire(c, uidStr, tools.RedisExpireTime())
	}
	if m.client.Exists(c, tidStr).Val() != 0 {
		pipeline.SAdd(c, tidStr, uid)
		pipeline.Expire(c, tidStr, tools.RedisExpireTime())
	}
	_, err := pipeline.Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// update redis when unfollow a user
func (m *Manager) UnFollow(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := "followList:" + strconv.FormatInt(uid, 10)
	tidStr := "fansList:" + strconv.FormatInt(tid, 10)
	pipeline := m.client.TxPipeline()
	// 1. update uid's follow list
	if m.client.Exists(c, uidStr).Val() != 0 {
		pipeline.SRem(c, uidStr, tid)
		pipeline.Expire(c, uidStr, tools.RedisExpireTime())
	}
	// 2. update tid's fans list
	if m.client.Exists(c, tidStr).Val() != 0 {
		pipeline.SRem(c, tidStr, uid)
		pipeline.Expire(c, tidStr, tools.RedisExpireTime())
	}
	_, err := pipeline.Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *Manager) SetFollow(c context.Context, uid int64, ids []int64) (bool, error) {
	uidStr := "followList:" + strconv.FormatInt(uid, 10)

	pipeline := m.client.TxPipeline()
	for _, id := range ids {
		pipeline.SAdd(c, uidStr, id)
	}
	pipeline.Expire(c, uidStr, tools.RedisExpireTime())
	_, err := pipeline.Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// update fans list to uid's fans list in redis using set
func (m *Manager) SetFans(c context.Context, uid int64, ids []int64) (bool, error) {
	uidStr := "fansList:" + strconv.FormatInt(uid, 10)
	// 1. only add fans list to redis when the number of fans is greater than the threshold
	if len(ids) < int(consts.RedisFansThreshold) {
		return true, nil
	}
	// 2. add fans list to redis
	pipeline := m.client.TxPipeline()
	for _, id := range ids {
		pipeline.SAdd(c, uidStr, id)
	}
	pipeline.Expire(c, uidStr, tools.RedisExpireTime())
	_, err := pipeline.Exec(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// query uid whether follow tid
func (m *Manager) QueryRelation(c context.Context, uid int64, tid int64) (bool, error) {
	// 检查uid的follow和tid的fans是否存在
	uidStr := "followList:" + strconv.FormatInt(uid, 10)
	tidStr := "fansList:" + strconv.FormatInt(tid, 10)
	if m.client.Exists(c, uidStr).Val() != 0 {
		if m.client.SIsMember(c, tidStr, tid).Val() {
			return true, nil
		}
	}
	if m.client.Exists(c, tidStr).Val() != 0 {
		if m.client.SIsMember(c, uidStr, uid).Val() {
			return true, nil
		}
	}
	return false, errors.New("no key")
}

// query follow list
func (m *Manager) QueryFollow(c context.Context, uid int64) ([]int64, error) {
	uidStr := "followList:" + strconv.FormatInt(uid, 10)
	if m.client.Exists(c, uidStr).Val() == 0 {
		return nil, nil
	}
	ids, err := m.client.MGet(c, uidStr).Result()
	if err != nil || len(ids)-1 == 0 {
		return nil, err
	}
	var res []int64
	for _, id := range ids {
		if id != nil {
			res = append(res, id.(int64))
		}
	}
	return res, nil
}

// query fans list
func (m *Manager) QueryFans(c context.Context, uid int64) ([]int64, error) {
	uidStr := "fansList:" + strconv.FormatInt(uid, 10)
	if m.client.Exists(c, uidStr).Val() == 0 {
		return nil, nil
	}
	ids, err := m.client.MGet(c, uidStr).Result()
	if err != nil || len(ids)-1 == 0 {
		return nil, err
	}
	var res []int64
	for _, id := range ids {
		if id != nil {
			res = append(res, id.(int64))
		}
	}
	return res, nil
}
