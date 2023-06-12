package mysql

import (
	"errors"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	VideoId int64  `json:"video_id" gorm:"index;not null" `
	Content string `json:"content" gorm:"type:varchar(255);not null"`
	UserId  int64  `json:"user_id" gorm:"not null" `
}

func (m *Comment) TableName() string {
	return "comment"
}

type CommentManager struct {
	db *gorm.DB
}

// create a new mysql manager, if table not exist, create it
func NewManager(db *gorm.DB) *CommentManager {
	m := db.Migrator()
	if !m.HasTable(&Comment{}) {
		if err := m.CreateTable(&Comment{}); err != nil {
			panic(err)
		}
	}
	return &CommentManager{
		db: db,
	}
}

// create a new comment
func (m *CommentManager) AddComment(videoId int64, userId int64, content string) (*Comment, error) {
	comment := Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: content,
	}
	err := m.db.Create(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// del a comment
func (m *CommentManager) DelComment(videoId int64, commentId int64) error {
	err := m.db.Where("video_id = ? and id = ?", videoId, commentId).Delete(&Comment{}).Error
	if err != nil {
		return err
	}
	return nil
}

// get comment list
func (m *CommentManager) GetCommentList(videoId int64) ([]*Comment, error) {
	var comments []*Comment
	err := m.db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return comments, nil
}
