package redis

import (
	"context"
	"strconv"
	"time"

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

// Key design:
// ZSet: key: vide_id, Member: comment_id, Score: timestamp
// Hash: key: comment_id, field: user_id,content

func (m *Manager) AddComment(videoId int64, commentId int64, userId int64, content string) error {
	commentIdStr := "comment:" + strconv.FormatInt(commentId, 10) + "_comment_info"
	videoIdStr := "video:" + strconv.FormatInt(videoId, 10) + "_comment_list"

	// start transaction
	tx := m.client.TxPipeline()
	// 1. add comment to redis
	tx.HSet(context.Background(), commentIdStr, "user_id", userId, "content", content)
	// 2. add comment id to video's comment list
	member := &redis.Z{Score: float64(time.Now().Unix()), Member: strconv.FormatInt(commentId, 10)}
	tx.ZAdd(context.Background(), videoIdStr, member)
	// 3. set expire time
	tx.Expire(context.Background(), videoIdStr, consts.RedisExpireTime)
	tx.Expire(context.Background(), commentIdStr, consts.RedisExpireTime)
	// exec transaction
	if _, err := tx.Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DelComment(videoId int64, commentId int64) (bool, error) {
	commentIdStr := "comment:" + strconv.FormatInt(commentId, 10) + "_comment_info"
	videoIdStr := "video:" + strconv.FormatInt(videoId, 10) + "_comment_list"

	// start transaction
	tx := m.client.TxPipeline()
	// 1. del comment from video's comment list
	tx.ZRem(context.Background(), videoIdStr, commentId)
	// 2. del comment info
	tx.Del(context.Background(), commentIdStr)
	// 3. set expire time
	tx.Expire(context.Background(), videoIdStr, consts.RedisExpireTime)
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
    videoIdStr := "video:" + strconv.FormatInt(videoId, 10) + "_comment_list"

    // Get the comments IDs with their scores from the ZSet starting from the specified offset and limited by the specified count
    commentIdsWithScores, err := m.client.ZRevRangeWithScores(context.Background(), videoIdStr, 0,-1).Result()
    if err != nil {
        return nil, err
    }

    // Get the comment info for each comment ID
    comments := make(map[string]map[string]string)
    for _, commentIdWithScore := range commentIdsWithScores {
        commentId := commentIdWithScore.Member.(string)
        commentIdStr := "comment:" + commentId + "_comment_info"

        commentInfo, err := m.client.HGetAll(context.Background(), commentIdStr).Result()
        if err != nil {
            return nil, err
        }
        commentInfo["create_time"] = strconv.FormatInt(int64(commentIdWithScore.Score), 10)
        comments[commentId] = commentInfo
    }

    return comments, nil
}
