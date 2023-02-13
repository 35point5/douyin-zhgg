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
		conn.Exec("DELETE FROM user_models")
		conn.Exec("DELETE FROM video_models")
		conn.Create(&domain.UserModel{Id: 0, Name: "user0", Password: "123456", FollowCount: 0, FollowerCount: 0})
		conn.Create(&domain.UserModel{Id: 1, Name: "user1", Password: "123456", FollowCount: 1, FollowerCount: 3})
		conn.Create(&domain.UserModel{Id: 2, Name: "user2", Password: "123456", FollowCount: 1, FollowerCount: 1})
		conn.Create(&domain.UserModel{Id: 3, Name: "user3", Password: "123456", FollowCount: 1, FollowerCount: 0})
		conn.Create(&domain.UserModel{Id: 4, Name: "user4", Password: "123456", FollowCount: 1, FollowerCount: 0})

		conn.Create(&domain.VideoModel{Id: 1, Uid: 1, PlayUrl: domainName + staticURL + "/bear.mp4",
			CoverUrl: domainName + staticURL + "/cover.jpg", FavoriteCount: 3, CommentCount: 2, Title: "video1", UpdatedTime: time.Now(),
		})
		conn.Create(&domain.VideoModel{Id: 2, Uid: 2, PlayUrl: domainName + staticURL + "/song1.mp4",
			CoverUrl: domainName + staticURL + "/song1.jpg", FavoriteCount: 1, CommentCount: 1, Title: "video2", UpdatedTime: time.Now(),
		})
		conn.Create(&domain.VideoModel{Id: 3, Uid: 2, PlayUrl: domainName + staticURL + "/song2.mp4",
			CoverUrl: domainName + staticURL + "/song2.jpg", FavoriteCount: 0, CommentCount: 0, Title: "video3", UpdatedTime: time.Now(),
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
	//res.FollowCount = m.GetFollowCnt(id)
	//res.FollowerCount = m.GetFollowerCnt(id)
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

//func (m *mysqlBasicRepository) GetFavoriteCnt(vid int64) int64 {
//	var temp []domain.FavoriteListModel
//	res := m.Mysql.Where("video_id = ?", vid).Find(&temp)
//	return res.RowsAffected
//}

// IsFollow returns whether id follows fid
func (m *mysqlBasicRepository) IsFollow(id int64, fid int64) bool {
	temp := domain.FollowListModel{
		UserID:   id,
		ToUserID: fid,
	}
	res := m.Mysql.First(&temp)
	if res.Error == nil && (temp.Status == 0 || temp.Status == 2) {
		return true
	}
	temp = domain.FollowListModel{
		UserID:   fid,
		ToUserID: id,
	}
	res = m.Mysql.First(&temp)
	if res.Error == nil && temp.Status == 2 {
		return true
	}
	return false
}

//func (m *mysqlBasicRepository) GetFollowCnt(id int64) int64 {
//	//var temp []domain.UserFollowModel
//	//res := m.Mysql.Where("user_id = ?", id).Find(&temp)
//	//return res.RowsAffected
//	return 0
//}

//func (m *mysqlBasicRepository) GetFollowerCnt(id int64) int64 {
//	//var temp []domain.UserFollowModel
//	//res := m.Mysql.Where("target_user_id = ?", id).Find(&temp)
//	//return res.RowsAffected
//	return 0
//}

func (m *mysqlBasicRepository) GetCommentCnt() {
	//TODO: 等comment接口实现
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
