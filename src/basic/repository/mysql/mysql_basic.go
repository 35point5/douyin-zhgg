package mysql

import (
	"douyin-service/domain"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"time"
)

type mysqlBasicRepository struct {
	Mysql *gorm.DB
	//Redis *redis.Client
}

func NewMysqlBasicRepository(conn *gorm.DB, debug bool) domain.BasicRepository {
	if err := conn.AutoMigrate(&domain.VideoModel{}, &domain.UserModel{}); err != nil {
		log.Fatal(err)
	}
	if debug {
		domainName := viper.GetString("domain")
		staticURL := viper.GetString("static_url")
		conn.Create(&domain.VideoModel{
			Id:            0,
			Uid:           5,
			PlayUrl:       domainName + staticURL + "/bear.mp4",
			CoverUrl:      domainName + staticURL + "/cover.jpg",
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			UpdatedTime:   time.Now(),
		})
	}
	return &mysqlBasicRepository{conn}
}

func (m *mysqlBasicRepository) GetVideoByTime(t time.Time) []domain.VideoModel {
	var res []domain.VideoModel
	videoCount := viper.GetInt("video_limit")
	m.Mysql.Where("updated_time < ?", t).Order("updated_time desc").Limit(videoCount).Find(&res)
	return res
}

func (m *mysqlBasicRepository) GetUserById(id int64) domain.UserModel {
	var res domain.UserModel
	res.Id = id
	m.Mysql.First(&res)
	return res
}

//func (m *mysqlBasicRepository) SetToken(uid int64, token string) error {
//	return m.Redis.SetNX(token, uid, time.Hour).Err()
//}

func (m *mysqlBasicRepository) CreateUser(user domain.UserModel) (int64, error) {
	res := m.Mysql.Create(&user)
	return user.Id, res.Error
}
