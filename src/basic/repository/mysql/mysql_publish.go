package mysql

import (
	"douyin-service/domain"
	"gorm.io/gorm"
	"log"
)

type mysqlPublishRepository struct {
	Mysql *gorm.DB
}

func NewMysqlPublishRepository(conn *gorm.DB, debug bool) domain.PublishRepository {
	if err := conn.AutoMigrate(&domain.VideoModel{}, &domain.UserModel{}); err != nil {
		log.Fatal(err)
	}
	return &mysqlPublishRepository{conn}
}

func (m *mysqlPublishRepository) AddVideo(v *domain.VideoModel) error {
	ret := m.Mysql.Save(v)
	return ret.Error
}

func (m *mysqlPublishRepository) LisVideoByUserId(userId int64) (videoList []domain.VideoModel, err error) {
	ret := m.Mysql.Where("user_id = ?", userId).Find(&videoList)
	err = ret.Error
	return
}
