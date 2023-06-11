package mysql

import (
	"errors"
	"tiktok-demo/shared/errno"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserId     int64 `json:"user_id" gorm:"not null;uniqueIndex:user_video"`
	VideoId    int64 `json:"video_id" gorm:"not null;uniqueIndex:user_video"`
	IsFavorite bool  `json:"is_favorite" gorm:"not null"`
}

func (m *Favorite) TableName() string {
	return "favorite"
}

type FavoriteManager struct {
	db *gorm.DB
}

// create a new mysql manager, if table not exist, create it
func NewManager(db *gorm.DB, salt string) *FavoriteManager {
	m := db.Migrator()
	if !m.HasTable(&Favorite{}) {
		if err := m.CreateTable(&Favorite{}); err != nil {
			panic(err)
		}
	}
	return &FavoriteManager{
		db: db,
	}
}

// get video's favorite user id list
func (m *FavoriteManager) GetFavoriteUserIdList(videoId int64) ([]*Favorite, error) {
	var favorites []*Favorite
	err := m.db.Where("video_id = ?", videoId).Find(&favorites).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		klog.Infof("get favorite user id list of video %d failed: %s", videoId, err.Error())
		return nil, err
	}
	return favorites, err
}

// get user's favorite video id list
func (m *FavoriteManager) GetFavoriteVideoIdList(userId int64) ([]*Favorite, error) {
	var favorites []*Favorite
	err := m.db.Where("user_id = ?", userId).Find(&favorites).Error
	if err != nil {
		klog.Infof("get favorite video id list of user %d failed: %s", userId, err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.FavoriteVideoListNotExistErr
		}
		return nil, err
	}
	return favorites, nil
}

// query a favorite
func (m *FavoriteManager) QueryFavorite(userId int64, videoId int64) (bool, error) {
	favorite := &Favorite{}
	err := m.db.Where("user_id = ? and video_id = ?", userId, videoId).First(favorite).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			klog.Infof("query favorite of user %d and video %d failed: %s", userId, videoId, err.Error())
			return false, nil
		}
		klog.Errorf("query favorite of user %d and video %d failed: %s", userId, videoId, err.Error())
		return false, err
	}
	return favorite.IsFavorite, nil
}

// favorite action
func (m *FavoriteManager) FavoriteAction(userId int64, videoId int64, isFavorite bool) error {
	var isExist bool
	var err error

	res := m.db.Where("user_id = ? and video_id = ?", userId, videoId).First(&Favorite{})
	if res.RowsAffected == 0 {
		isExist = false
	} else {
		if res.Error != nil {
			klog.Infof("query favorite data fail")
			return errors.New("query favorite data fail")
		}
		isExist = true
	}

	if !isExist {
		err = m.InsertFavorite(userId, videoId, isFavorite)
	} else {
		err = m.UpdateFavorite(userId, videoId, isFavorite)
	}

	if err != nil {
		klog.Infof("favorite action failed: %s", err.Error())
		return err
	}

	return nil
}

// insert a favorite
func (m *FavoriteManager) InsertFavorite(userId int64, videoId int64, isFavorite bool) error {
	favorite := &Favorite{
		UserId:     userId,
		VideoId:    videoId,
		IsFavorite: isFavorite,
	}
	err := m.db.Create(favorite).Error
	if err != nil {
		klog.Infof("insert favorite failed: %s", err.Error())
		return errors.New("insert favorite data fail")
	}
	return nil
}

// update a favorite
func (m *FavoriteManager) UpdateFavorite(userId int64, videoId int64, isFavorite bool) error {
	err := m.db.Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Update("is_favorite", isFavorite).Error
	if err != nil {
		klog.Infof("update favorite failed: %s", err.Error())
		return errors.New("update favorite data fail")
	}
	return nil
}

// get video's favorite count by video id
func (m *FavoriteManager) GetFavoriteCountByVideoId(videoId int64) (int64, error) {
	var count int64
	err := m.db.Model(&Favorite{}).Where("video_id = ? and is_favorite = ?", videoId, true).Count(&count).Error
	if err != nil {
		klog.Infof("get favorite count by video id %d failed: %s", videoId, err.Error())
		return 0, err
	}
	return count, nil
}
