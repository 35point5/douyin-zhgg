package mysql

import (
	"douyin-service/domain"
	"errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"time"
)

type mysqlInteractRepository struct {
	Mysql *gorm.DB
}

func NewMysqlInteractRepository(conn *gorm.DB, debug bool) domain.InteractRepository {
	if err := conn.AutoMigrate(&domain.FavoriteListModel{}); err != nil {
		log.Fatal(err)
	}

	if debug {
		conn.Save(&domain.FavoriteListModel{1, 2, 1, time.Now()})
		conn.Save(&domain.FavoriteListModel{1, 3, 1, time.Now()})
		conn.Save(&domain.FavoriteListModel{1, 2, 1, time.Now()})
		conn.Save(&domain.FavoriteListModel{3, 3, 1, time.Now()})
	}
	return &mysqlInteractRepository{conn}
}

func (m *mysqlInteractRepository) GetVideoModelsById(id []int64) ([]domain.VideoModel, error) {
	var res []domain.VideoModel
	if len(id) == 0 {
		return res, nil
	}
	videoCount := viper.GetInt("video_limit")
	queryRes := m.Mysql.Where(id).Order("updated_time desc").Limit(videoCount).Find(&res)
	return res, queryRes.Error
}

func (m *mysqlInteractRepository) GetFavoriteListByUserId(id int64) ([]domain.FavoriteListModel, error) {

	var res []domain.FavoriteListModel
	if id == 0 {
		return res, nil
	}
	var userModel domain.UserModel
	userModel.Id = id
	judgeRes := m.Mysql.First(&userModel)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("user id is not exist!")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	videoCount := viper.GetInt("video_limit")
	judgeRes = m.Mysql.Where("user_id = ?", id).Order("updated_time desc").Limit(videoCount).Find(&res)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("favorite list is null !")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	return res, nil
}

func (m *mysqlInteractRepository) FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error) {

	var flm domain.FavoriteListModel
	if user_id == 0 {
		return false, errors.New("user_id is zero !")
	}
	flm.VideoID = video_id
	flm.UserID = user_id
	if action_type == 1 {
		flm.Status = 1
		ret := m.Mysql.Create(&flm)
		if ret.Error != nil {
			return false, ret.Error
		}
	} else if action_type == 2 {
		ret := m.Mysql.Delete(&flm)
		if ret.Error != nil {
			return false, ret.Error
		}
	} else {
		return false, errors.New("action_type must be 1 or 2 !")
	}
	return true, nil
}
