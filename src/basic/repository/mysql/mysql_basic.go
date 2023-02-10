package mysql

import (
	"douyin-service/domain"
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
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
			Uid:           24,
			PlayUrl:       domainName + staticURL + "/bear.mp4",
			CoverUrl:      domainName + staticURL + "/cover.jpg",
			FavoriteCount: 0,
			CommentCount:  0,
			UpdatedTime:   time.Now(),
		})
	}
	return &mysqlBasicRepository{conn}
}

func (m *mysqlBasicRepository) GetVideoByTime(t time.Time) []domain.VideoModel {
	var res []domain.VideoModel
	videoCount := viper.GetInt("video_limit")
	m.Mysql.Where("updated_time <= ?", t).Order("updated_time desc").Limit(videoCount).Find(&res)
	return res
}

func (m *mysqlBasicRepository) GetUserById(id int64) domain.UserModel {
	var res domain.UserModel
	m.Mysql.First(&res, id)
	return res
}

func (m *mysqlBasicRepository) IsFavorite(uid int64, vid int64) bool {
	var temp domain.FavoriteListModel
	res := m.Mysql.Where("user_id = ? and video_id = ?", uid, vid).First(&temp)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false
	}
	if temp.Status != 1 {
		return false
	}
	return true
}

// IsFollow returns whether id follows fid
func (m *mysqlBasicRepository) IsFollow(id int64, fid int64) bool {
	//TODO: 等follow接口实现
	//temp := domain.UserFollowModel{
	//	UserId:       id,
	//	TargetUserId: fid,
	//}
	//res := m.Mysql.First(&temp)
	//if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	//	return false
	//}
	return true
}

func (m *mysqlBasicRepository) GetFollowCnt(id int64) int64 {
	//TODO: 等follow接口实现
	//var temp []domain.UserFollowModel
	//res := m.Mysql.Where("user_id = ?", id).Find(&temp)
	//return res.RowsAffected
	return 0
}

func (m *mysqlBasicRepository) GetFollowerCnt(id int64) int64 {
	//TODO: 等follow接口实现
	//var temp []domain.UserFollowModel
	//res := m.Mysql.Where("target_user_id = ?", id).Find(&temp)
	//return res.RowsAffected
	return 0
}

func (m *mysqlBasicRepository) GetUserByName(name string) domain.UserModel {
	var res domain.UserModel
	m.Mysql.Where("name = ?", name).First(&res)
	return res
}

//func (m *mysqlBasicRepository) SetToken(uid int64, token string) error {
//	return m.Redis.SetNX(token, uid, time.Hour).Err()
//}

func (m *mysqlBasicRepository) CreateUser(user domain.UserModel) (int64, error) {
	res := m.Mysql.Create(&user)
	return user.Id, res.Error
}
