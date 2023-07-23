package mysql

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	FromUserId int64  `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	ToUserId   int64  `gorm:"index:idx_userid_from;index:idx_userid_to;not null" json:"to_user_id"`
	Content    string `gorm:"type:varchar(255);not null" json:"content"`
}

func (m *Message) TableName() string {
	return "message"
}

type MessageManager struct {
	db *gorm.DB
}

// create a new mysql manager, if table not exist, create it
func NewManager(db *gorm.DB) *MessageManager {
	m := db.Migrator()
	if !m.HasTable(&Message{}) {
		if err := m.CreateTable(&Message{}); err != nil {
			panic(err)
		}
	}
	return &MessageManager{
		db: db,
	}
}

// create a new message
func (m *MessageManager) AddMessage(fromUserId int64, toUserId int64, content string) (*Message, error) {
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	message := Message{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Content:    content,
	}
	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &message, nil
}

// 根据两个用户的用户id获取给定时间后的聊天信息记录
func (m *MessageManager) GetMessageByUserId(userId int64, toUserId int64, timeStamp int64) ([]*Message, error) {
	var messages []*Message
	err := m.db.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userId, toUserId, toUserId, userId).Where("created_at > ?", time.Unix(timeStamp, 0)).Order("created_at asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// 根据消息id获取消息
func (m *MessageManager) GetMessageById(id int64) (*Message, error) {
	var message Message
	err := m.db.Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// 根据两个用户的用户id获取最新的一条聊天信息记录
func (m *MessageManager) GetLatestMessage(userId int64, toUserId int64) (*Message, error) {
	var message Message
	err := m.db.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userId, toUserId, toUserId, userId).Order("created_at desc").First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}
