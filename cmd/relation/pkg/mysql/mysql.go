package mysql

import (
	"errors"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	UserID   int64 `json:"user_id" gorm:"not null;uniqueIndex:user_touser" `    // fans id
	ToUserID int64 `json:"to_user_id" gorm:"not null;uniqueIndex:user_touser" ` // be followed id
}

func (m *Relation) TableName() string {
	return "relation"
}

type RelationManager struct {
	db   *gorm.DB
	salt string
}

// create a new mysql manager, if table not exist, create it
func NewManager(db *gorm.DB, salt string) *RelationManager {
	m := db.Migrator()
	if !m.HasTable(&Relation{}) {
		if err := m.CreateTable(&Relation{}); err != nil {
			panic(err)
		}
	}
	return &RelationManager{
		db:   db,
		salt: salt,
	}
}

// create a new relation
func (m *RelationManager) AddFollow(userId int64, toUserId int64) error {
	if userId == toUserId {
		return nil
	}
	// using transaction to ensure data consistency
	relation := Relation{
		UserID:   userId,
		ToUserID: toUserId,
	}
	existingRelation := Relation{}
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// 1. check if relation exist
		err := tx.Where("user_id = ? and to_user_id = ?", userId, toUserId).First(&existingRelation).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Infof("get relation failed: %s", err.Error())
			return err
		}
		// todo: 这里返回的err是否会影响到外层的err
		if existingRelation.UserID != 0 {
			klog.Infof("relation already exist")
			return err
		}
		// 2. create relation
		if err := tx.Create(&relation).Error; err != nil {
			klog.Infof("create relation failed: %s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		klog.Infof("create relation transaction failed: %s", err.Error())
		return err
	}
	return nil
}

// del a relation
func (m *RelationManager) DelFollow(userId int64, toUserId int64) error {
	if userId == toUserId {
		return nil
	}
	existingRelation := Relation{}
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// 1. check if relation exist
		err := tx.Where("user_id = ? and to_user_id = ?", userId, toUserId).First(&existingRelation).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			klog.Infof("get relation failed: %s", err.Error())
			return err
		}
		if existingRelation.UserID == 0 {
			klog.Infof("relation not exist")
			return err
		}
		// 2. delete relation, use Unscoped to delete soft deleted records
		if err := tx.Unscoped().Delete(&existingRelation).Error; err != nil {
			klog.Infof("delete relation failed: %s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		klog.Infof("delete relation transaction failed: %s", err.Error())
		return err
	}
	return nil
}

// query a relation
func (m *RelationManager) QueryRelation(userId int64, toUserId int64) (bool, error) {
	if userId == toUserId {
		return true, nil
	}
	err := m.db.Where("user_id = ? and to_user_id = ?", userId, toUserId).First(&Relation{}).Error
	if err != nil {
		klog.Infof("query relation failed: %s", err.Error())
		return false, err
	}
	return true, nil
}

// get user's fans
func (m *RelationManager) GetFansList(userId int64) ([]*Relation, error) {
	var relations []*Relation
	err := m.db.Where("to_user_id = ?", userId).Find(&relations).Error
	if err != nil {
		klog.Infof("get fans list failed: %s", err.Error())
		return nil, err
	}
	return relations, nil
}

// get user's follow
func (m *RelationManager) GetFollowList(userId int64) ([]*Relation, error) {
	var relations []*Relation
	err := m.db.Where("user_id = ?", userId).Find(&relations).Error
	if err != nil {
		klog.Infof("get follow list failed: %s", err.Error())
		return nil, err
	}
	return relations, nil
}

// get set of user's follow,using for quick query a user is fans or not
func (m *RelationManager) GetFollowSet(userId int64) (map[int64]struct{}, error) {
	followSet := make(map[int64]struct{})
	followList, err := m.GetFollowList(userId)
	if err != nil {
		klog.Infof("get follow list failed: %s", err.Error())
		return nil, err
	}
	for _, relation := range followList {
		followSet[relation.ToUserID] = struct{}{}
	}
	return followSet, nil
}
