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

// store id and username in redis using string
func (m *Manager) AddName(c context.Context, uid int64, uname string) error {
	if err := m.client.Set(c, strconv.FormatInt(uid, 10)+"_username", uname, consts.RedisExpireTime).Err(); err != nil {
		return err
	}
	return nil
}

// update FollowingNum and FollowerNum in redis using hash
func (m *Manager) UpdateFollow(c context.Context, uid int64, ufollowing uint64, ufollower uint64) error {
	uidStr := strconv.FormatInt(uid, 10) + "_follow_info"
	if err := m.client.HSet(c, uidStr, "FollowingNum", ufollowing).Err(); err != nil {
		return err
	}
	if err := m.client.HSet(c, uidStr, "FollowerNum", ufollower).Err(); err != nil {
		return err
	}
	m.client.Expire(c, uidStr, consts.RedisExpireTime)
	return nil
}

// delete the followingNum of uid and followerNum of tid in redis using hash
func (m *Manager) DeleteRelation(c context.Context, uid int64, tid int64) (bool, error) {
	uidStr := strconv.FormatInt(uid, 10) + "_follow_info"
	tidStr := strconv.FormatInt(tid, 10) + "_follow_info"
	if err := m.client.HDel(c, uidStr, "FollowingNum").Err(); err != nil {
		return false, err
	}
	if err := m.client.HDel(c, tidStr, "FollowerNum").Err(); err != nil {
		return false, err
	}
	m.client.Expire(c, uidStr, consts.RedisExpireTime)
	m.client.Expire(c, tidStr, consts.RedisExpireTime)
	return true, nil
}

// get user info(username,followingNum,followerNum) from redis by uid
func (m *Manager) GetUserInfo(c context.Context, uid int64) (*mysql.User, error) {
	var FollowingNum, FollowerNum, Name string
	uidStr := strconv.FormatInt(uid, 10)
	FollowingNum, _ = m.client.HGet(c, uidStr+"_follow_info", "FollowingNum").Result()
	FollowerNum, _ = m.client.HGet(c, uidStr+"_follow_info", "FollowerNum").Result()
	Name, _ = m.client.Get(c, uidStr+"_username").Result()
	if len(Name) > 0 && len(FollowingNum) > 0 && len(FollowerNum) > 0 {
		FollowerNumUint64, _ := strconv.ParseUint(FollowerNum, 10, 64)
		FollowingNumUint64, _ := strconv.ParseUint(FollowingNum, 10, 64)
		return &mysql.User{
			Uid:            uid,
			Name:           Name,
			FollowingCount: FollowingNumUint64,
			FollowerCount:  FollowerNumUint64,
		}, nil
	}
	return nil, errno.UserNotExistErr
}

// insert user info into redis
func (m *Manager) InsertUserInfo(c context.Context, u *mysql.User) error {
	if err := m.AddName(c, u.Uid, u.Name); err != nil {
		return err
	}
	if err := m.UpdateFollow(c, u.Uid, u.FollowingCount, u.FollowerCount); err != nil {
		return err
	}
	return nil
}
