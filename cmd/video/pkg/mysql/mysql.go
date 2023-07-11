package mysql

import (
	"tiktok-demo/shared/consts"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"tiktok-demo/shared/errno"
)

type Video struct {
	gorm.Model
	AuthorId      int64  `gorm:"index:idx_authorid;not null"`
	VideoURL      string `gorm:"type:varchar(200);not null"`
	CoverURL      string `gorm:"type:varchar(200);not null"`
	FavoriteCount int64  `gorm:"type:int;default:0;not null"`
	CommentCount  int64  `gorm:"type:int;default:0;not null"`
	// IsFavorite    bool   `gorm:"type:tinyint(1);default:0;not null"`
	Title         string `gorm:"type:varchar(200);not null"`
}

func (v *Video) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(consts.VideoSnowflakeNode)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	if v.ID == 0 {
		v.ID = uint(sf.Generate().Int64())
	}
	return nil
}

func (v *Video) TableName() string {
	return "video"
}

type Manager struct {
	db *gorm.DB
}

// create a new mysql manager, if table not exist, create it
func NewManager(db *gorm.DB) *Manager {
	m := db.Migrator()
	if !m.HasTable(&Video{}) {
		if err := m.CreateTable(&Video{}); err != nil {
			panic(err)
		}
	}
	return &Manager{
		db: db,
	}
}

func (m *Manager) PublishVideo(v *Video) error {
	return m.db.Model(&Video{}).Create(v).Error
}

func (m *Manager) GetVideoByTime(lastTime int64, limit int) ([]*Video, error) {
	videos := make([]*Video, limit)
	if lastTime == 0 {
		lastTime = int64(time.Now().UnixMilli())
	}
	res := m.db.Model(&Video{}).Where("created_at < ?", time.UnixMilli(lastTime)).Order("created_at desc").Limit(limit).Find(&videos)
	if res.RowsAffected == 0 {
		return nil, errno.VideoNotExistErr
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return videos, nil
}

func (m *Manager) GetVideoByUserId(userId int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := m.db.Model(&Video{}).Where(&Video{AuthorId: userId}).Find(&videos).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.VideoNotExistErr
		} else {
			return nil, err
		}
	}
	return videos, nil
}

func (m *Manager) GetVideoById(videoId int64) (*Video, error) {
	video := &Video{}
	err := m.db.Model(&Video{}).Where("id = ?", videoId).First(video).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.VideoNotExistErr
		} else {
			return nil, err
		}
	}
	return video, nil
}
