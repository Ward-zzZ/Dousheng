package redis

import (
	"context"
	"strconv"
	"time"

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

// Key design:
// ZSet: key: vide_id, Member: comment_id, Score: timestamp
// Hash: key: comment_id, field: user_id,content

func (m *Manager) AddComment(videoId int64, commentId int64, userId int64, content string) error {
	commentIdStr := "commentInfo:" + strconv.FormatInt(commentId, 10)
	videoIdStr := "commentList:" + strconv.FormatInt(videoId, 10)

	// start transaction
	tx := m.client.TxPipeline()
	// 1. add comment to redis
	tx.HSet(context.Background(), commentIdStr, "user_id", userId, "content", content)
	tx.Expire(context.Background(), commentIdStr, tools.RedisExpireTime())
	// 2. add comment id to video's comment list
	if m.client.Exists(context.Background(), videoIdStr).Val() == 1 {
		member := &redis.Z{Score: float64(time.Now().Unix()), Member: strconv.FormatInt(commentId, 10)}
		tx.ZAdd(context.Background(), videoIdStr, member)
		// 3. set expire time
		tx.Expire(context.Background(), videoIdStr, tools.RedisExpireTime())
	}
	// exec transaction
	if _, err := tx.Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DelComment(videoId int64, commentId int64) (bool, error) {
	commentIdStr := "commentInfo:" + strconv.FormatInt(commentId, 10)
	videoIdStr := "commentList:" + strconv.FormatInt(videoId, 10)

	// start transaction
	tx := m.client.TxPipeline()
	// 1. del comment from video's comment list
	tx.Del(context.Background(), videoIdStr)
	// 2. del comment info
	tx.Del(context.Background(), commentIdStr)
	// exec transaction
	if _, err := tx.Exec(context.Background()); err != nil {
		// if err is redis.Nil, it means the comment is not exist
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *Manager) GetCommentList(videoId int64) (map[string]map[string]string, error) {
	videoIdStr := "commentList:" + strconv.FormatInt(videoId, 10)

	// Get the comments IDs with their scores from the ZSet starting from the specified offset and limited by the specified count
	commentIdsWithScores, err := m.client.ZRevRangeWithScores(context.Background(), videoIdStr, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Get the comment info for each comment ID
	comments := make(map[string]map[string]string)
	for _, commentIdWithScore := range commentIdsWithScores {
		commentId := commentIdWithScore.Member.(string)
		commentIdStr := "commentInfo:" + commentId

		commentInfo, err := m.client.HGetAll(context.Background(), commentIdStr).Result()
		if err != nil {
			return nil, err
		}
		commentInfo["create_time"] = strconv.FormatInt(int64(commentIdWithScore.Score), 10)
		comments[commentId] = commentInfo
	}

	return comments, nil
}

func (m *Manager) SetCommentList(videoId int64, commentIds []int64, userId []int64, content []string) error {
	videoIdStr := "commentList:" + strconv.FormatInt(videoId, 10)

	// start transaction
	tx := m.client.TxPipeline()
	// 先删除原来的
	tx.Del(context.Background(), videoIdStr)
	tx.Expire(context.Background(), videoIdStr, tools.RedisExpireTime())
	// 1. add comment to redis
	for i := 0; i < len(commentIds); i++ {
		commentIdStr := "commentInfo:" + strconv.FormatInt(commentIds[i], 10)
		tx.HSet(context.Background(), commentIdStr, "user_id", userId[i], "content", content[i])
		tx.Expire(context.Background(), commentIdStr, tools.RedisExpireTime())
	}
	// 2. add comment id to video's comment list
	for i := 0; i < len(commentIds); i++ {
		member := &redis.Z{Score: float64(time.Now().Unix()), Member: strconv.FormatInt(commentIds[i], 10)}
		tx.ZAdd(context.Background(), videoIdStr, member)
	}
	// 3. set expire time
	tx.Expire(context.Background(), videoIdStr, tools.RedisExpireTime())

	// exec transaction
	if _, err := tx.Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
