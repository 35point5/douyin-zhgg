package mysql

import (
	"douyin-service/domain"
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type mysqlInteractRepository struct {
	Mysql *gorm.DB
}

func NewMysqlInteractRepository(conn *gorm.DB, debug bool) domain.InteractRepository {
	if err := conn.AutoMigrate(&domain.FavoriteListModel{}, &domain.CommentModel{}); err != nil {
		log.Fatal(err)
	}

	if debug {
		conn.Exec("DELETE FROM favorite_list_models")
		conn.Exec("DELETE FROM comment_models")
		conn.Save(&domain.FavoriteListModel{UserID: 1, VideoID: 2, Status: 1})
		conn.Save(&domain.FavoriteListModel{UserID: 2, VideoID: 1, Status: 1})
		conn.Save(&domain.FavoriteListModel{UserID: 3, VideoID: 1, Status: 1})
		conn.Save(&domain.FavoriteListModel{UserID: 4, VideoID: 1, Status: 1})

		conn.Save(&domain.CommentModel{ID: 1, UserID: 1, VideoID: 1, Status: 1, CommentText: "好看！"})
		conn.Save(&domain.CommentModel{ID: 2, UserID: 2, VideoID: 1, Status: 1, CommentText: "不好看！"})
		conn.Save(&domain.CommentModel{ID: 3, UserID: 3, VideoID: 2, Status: 1, CommentText: "太好看了！"})

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
	judgeRes = m.Mysql.Where("user_id = ?", id).Limit(videoCount).Find(&res)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("favorite list is null !")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	return res, nil
}

func (m *mysqlInteractRepository) FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error) {

	flm := domain.FavoriteListModel{
		UserID:  user_id,
		VideoID: video_id,
		Status:  1,
	}
	var vm domain.VideoModel
	m.Mysql.First(&vm, video_id)
	if user_id == 0 {
		return false, errors.New("user_id is zero !")
	}
	if action_type == 1 {
		flm.Status = 1
		ret := m.Mysql.Save(&flm)
		if ret.Error != nil {
			return false, ret.Error
		}
		vm.FavoriteCount++
		m.Mysql.Save(&vm)
	} else if action_type == 2 {
		ret := m.Mysql.Delete(&flm)
		if ret.Error != nil {
			return false, ret.Error
		}
		vm.FavoriteCount--
		m.Mysql.Save(&vm)
	} else {
		return false, errors.New("action_type must be 1 or 2 !")
	}
	return true, nil
}

func (m *mysqlInteractRepository) GetUser(user_id int64) (domain.UserModel, error) {
	var userModel domain.UserModel
	userModel.Id = user_id
	judgeRes := m.Mysql.First(&userModel)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return userModel, errors.New("user id is not exist!")
	} else if judgeRes.Error != nil {
		return userModel, judgeRes.Error
	}
	return userModel, nil
}

func (m *mysqlInteractRepository) CommentActionByUserId(user_id int64, video_id int64, content string, action_type int32, comment_id int64) error {
	return nil
}

func (m *mysqlInteractRepository) GetCommentListByVideoId(video_id int64) ([]domain.CommentModel, error) {
	var res []domain.CommentModel
	judgeRes := m.Mysql.Where("video_id = ?", video_id).Find(&res)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("favorite list is null !")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	return res, nil
}

func (m *mysqlInteractRepository) AddCommentByUserId(user_id int64, video_id int64, content string) (domain.CommentModel, error) {
	model := domain.CommentModel{
		UserID:      user_id,
		VideoID:     video_id,
		CommentText: content,
		CreateDate:  time.Now().Format("2006-01-02 15:04:05"),
	}
	ret := m.Mysql.Save(&model)
	if ret.Error != nil {
		return model, errors.New("评论失败！")
	} else {
		return model, nil
	}
}

func (m *mysqlInteractRepository) DeleteCommentById(user_id int64, comment_id int64) (domain.CommentModel, error) {
	var model domain.CommentModel
	ret := m.Mysql.First(&model, comment_id)
	if ret.Error != nil {
		return model, errors.New("评论不存在！")
	}
	ret = m.Mysql.Delete(&model)
	if ret.Error != nil {
		return model, errors.New("删除失败！")
	}
	return model, nil
}
