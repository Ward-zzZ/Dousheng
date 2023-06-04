package redis

import (
	"context"
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

// part of uid's follow list in redis using set, solve hot key problem
func (m *Manager) AddRelation(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_relation_list"
	if err := m.client.SAdd(c, uidStr, tid).Err(); err != nil {
		return false, err
	}
	if err := m.client.Expire(c, uidStr, consts.RedisExpireTime).Err(); err != nil {
		return false, err
	}
	return true, nil
}

// update follow list to uid's follow list in redis using set
func (m *Manager) AddFollow(c context.Context, uid int64, ids []int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_follow_list"
	for _, id := range ids {
		if err := m.client.SAdd(c, uidStr, id).Err(); err != nil {
			return false, err
		}
	}
	if err := m.client.Expire(c, uidStr, consts.RedisExpireTime).Err(); err != nil {
		return false, err
	}
	return true, nil
}

// update fans list to uid's fans list in redis using set
func (m *Manager) AddFans(c context.Context, uid int64, ids []int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_fans_list"
	// 1. only add fans list to redis when the number of fans is greater than the threshold
	if len(ids) < int(consts.RedisFansThreshold) {
		return true, nil
	}
	// 2. add fans list to redis
	for _, id := range ids {
		if err := m.client.SAdd(c, uidStr, id).Err(); err != nil {
			return false, err
		}
	}
	if err := m.client.Expire(c, uidStr, consts.RedisExpireTime).Err(); err != nil {
		return false, err
	}
	return true, nil
}

// update redis when unfollow a user
func (m *Manager) UnFollow(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_follow_list"
	tidStr := strconv.FormatInt(tid, 10) + "_fans_list"
	// 1. update uid's follow list
	if cnt, _ := m.client.SCard(c, uidStr).Result(); cnt > 0 {
		if err := m.client.SRem(c, uidStr, tid).Err(); err != nil {
			return false, err
		}
		if err := m.client.Expire(c, uidStr, consts.RedisExpireTime).Err(); err != nil {
			return false, err
		}
	}
	// 2. update tid's fans list
	if cnt, _ := m.client.SCard(c, tidStr).Result(); cnt > 0 {
		if err := m.client.SRem(c, tidStr, uid).Err(); err != nil {
			return false, err
		}
		if err := m.client.Expire(c, tidStr, consts.RedisExpireTime).Err(); err != nil {
			return false, err
		}
	}
	// 3. delete uid's relation list
	uidStr = strconv.FormatInt(uid, 10) + "_relation_list"
	if cnt, _ := m.client.SCard(c, uidStr).Result(); cnt > 0 {
		m.client.SRem(c, uidStr, tid)
		m.client.Expire(c, uidStr, consts.RedisExpireTime)
	}
	return true, nil
}

// query a relation
func (m *Manager) QueryRelation(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_relation_list"
	if cnt, _ := m.client.SCard(c, uidStr).Result(); cnt > 0 {
		if ok, _ := m.client.SIsMember(c, uidStr, tid).Result(); ok {
			m.client.Expire(c, uidStr, consts.RedisExpireTime)
			return true, nil
		}
	}
	return false, nil
}

// query follow list
func (m *Manager) QueryFollow(c context.Context, uid int64) ([]int64, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_follow_list"
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
	uidStr := strconv.FormatInt(uid, 10) + "_fans_list"
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
