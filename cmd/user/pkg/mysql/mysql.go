package mysql

import (
	// "time"

	"tiktok-demo/cmd/user/pkg/md5"
	"tiktok-demo/shared/consts"
	"tiktok-demo/shared/errno"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	Uid            int64  `gorm:"column:uid;primary_key" json:"uid"`                                      // 用户的唯一标识
	Name           string `gorm:"column:name;type:varchar(40);NOT NULL" json:"name"`                      // 用户名称（唯一索引，不重复）
	Password       string `gorm:"column:password;type:varchar(257);NOT NULL" json:"password"`             // 密码(MD5加盐加密)
	FollowingCount uint64 `gorm:"column:following_count;type:bigint(20) unsigned" json:"following_count"` // 关注总数
	FollowerCount  uint64 `gorm:"column:follower_count;type:bigint(20) unsigned" json:"follower_count"`   // 粉丝总数
}

func (m *User) TableName() string {
	return "user"
}

// generate an id for user
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	// if user already has an id, do not generate a new one
	if u.Uid != 0 {
		return nil
	}
	sf, err := snowflake.NewNode(consts.UserSnowflakeNode)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	u.Uid = sf.Generate().Int64()
	return nil
}

type UserManager struct {
	db   *gorm.DB
	salt string
}

// create a new mysql manager, if table not exist, create it
func NewUserManager(db *gorm.DB, salt string) *UserManager {
	m := db.Migrator()
	if !m.HasTable(&User{}) {
		if err := m.CreateTable(&User{}); err != nil {
			panic(err)
		}
	}
	return &UserManager{
		db:   db,
		salt: salt,
	}
}

// create a new user
func (m *UserManager) CreateUser(user *User) (*User, error) {
	user.Password = md5.Md5Crypt(user.Password, m.salt)
	err := m.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

// get user by id
func (m *UserManager) GetUserByID(id int64) (*User, error) {
	var user User
	err := m.db.Where("uid = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.UserNotExistErr
		}
		return nil, err
	}
	return &user, nil
}

// get user by name
func (m *UserManager) GetUserByName(name string) (*User, error) {
	var user User
	err := m.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.UserNotExistErr
		}
		return nil, err
	}
	return &user, nil
}

// get  users list by id
func (m *UserManager) GetUsersByID(ids []int64) ([]*User, error) {
	var users []*User
	for _, id := range ids {
		user, err := m.GetUserByID(id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// follow a user
func (m *UserManager) FollowUser(userID int64, followUserID int64) error {
	// use transaction to ensure data consistency
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// check if user exist
		var user User
		if err := tx.Where("uid = ?", userID).First(&user).Error; err != nil {
			klog.Info("get user by id failed: %s", err.Error())
			return err
		}
		// check if follow user exist
		var followUser User
		if err := tx.Where("uid = ?", followUserID).First(&followUser).Error; err != nil {
			klog.Info("get follow user by uid failed: %s", err.Error())
			return err
		}
		// add following count
		if err := tx.Model(&user).Update("following_count", gorm.Expr("following_count + ?", 1)).Error; err != nil {
			klog.Info("update user following count failed: %s", err.Error())
			return err
		}
		// add follower count
		if err := tx.Model(&followUser).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
			klog.Info("update follow user follower count failed: %s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		klog.Info("follow transaction failed: %s", err.Error())
		return err
	}
	return nil
}

// unfollow a user
func (m *UserManager) UnFollowUser(userID int64, followUserID int64) error {
	// use transaction to ensure data consistency
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// check if user exist
		var user User
		if err := tx.Where("uid = ?", userID).First(&user).Error; err != nil {
			klog.Info("get user by uid failed: %s", err.Error())
			return err
		}
		// check if follow user exist
		var followUser User
		if err := tx.Where("uid = ?", followUserID).First(&followUser).Error; err != nil {
			klog.Info("get follow user by uid failed: %s", err.Error())
			return err
		}
		// minus following count
		if err := tx.Model(&user).Update("following_count", gorm.Expr("following_count - ?", 1)).Error; err != nil {
			klog.Info("update user following count failed: %s", err.Error())
			return err
		}
		// minus follower count
		if err := tx.Model(&followUser).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			klog.Info("update follow user follower count failed: %s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		klog.Info("unfollow transaction failed: %s", err.Error())
		return err
	}
	return nil
}
