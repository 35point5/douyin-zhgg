package mysql

import (
	"douyin-service/domain"
	"errors"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type mysqlSocialRepository struct {
	Mysql *gorm.DB
}

func NewMysqlSocialRepository(conn *gorm.DB, debug bool) domain.SocialRepository {
	if err := conn.AutoMigrate(&domain.FollowListModel{}); err != nil {
		log.Fatal(err)
	}
	if debug {
		conn.Exec("DELETE FROM follow_list_models")
		conn.Save(&domain.FollowListModel{UserID: 1, ToUserID: 2, Status: 2})
		conn.Save(&domain.FollowListModel{UserID: 2, ToUserID: 1, Status: 2})
		conn.Save(&domain.FollowListModel{UserID: 3, ToUserID: 1, Status: 0})
		conn.Save(&domain.FollowListModel{UserID: 4, ToUserID: 1, Status: 0})
	}
	return &mysqlSocialRepository{conn}
}

func (m *mysqlSocialRepository) GetFollowListByUserId(id int64) ([]domain.FollowListModel, error) {
	var res []domain.FollowListModel
	if id == 0 {
		return res, nil
	}
	var userModel domain.UserModel
	userModel.Id = id
	judgeRes := m.Mysql.First(&userModel)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("user id is not exist!")
	}
	judgeRes = m.Mysql.Where("user_id = ? AND status != 1", id).Find(&res)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("follow list is null!")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	return res, nil
}

func (m *mysqlSocialRepository) GetFollowerListByUserId(id int64) ([]domain.FollowListModel, error) {
	var res []domain.FollowListModel
	if id == 0 {
		return res, nil
	}
	var userModel domain.UserModel
	userModel.Id = id
	judgeRes := m.Mysql.First(&userModel)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("user id is not exist!")
	}
	judgeRes = m.Mysql.Where("to_user_id = ? AND status != 1", id).Find(&res)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("follow list is null!")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	return res, nil
}

func (m *mysqlSocialRepository) GetFriendListByUserId(id int64) ([]domain.FollowListModel, error) {
	var res []domain.FollowListModel
	if id == 0 {
		return res, nil
	}
	var userModel domain.UserModel
	userModel.Id = id
	judgeRes := m.Mysql.First(&userModel)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("user id is not exist!")
	}
	judgeRes = m.Mysql.Where("user_id = ? AND status != 1", id).Find(&res)
	if errors.Is(judgeRes.Error, gorm.ErrRecordNotFound) {
		return res, errors.New("follow list is null!")
	} else if judgeRes.Error != nil {
		return res, judgeRes.Error
	}
	return res, nil
}

func (m *mysqlSocialRepository) FollowActionByUserId(user_id int64, to_user_id int64, action_type int32) (bool, error) {
	if user_id == 0 {
		return false, errors.New("user_id is zero")
	}
	if user_id == to_user_id {
		return false, errors.New("cannot follow yourself")
	}
	flm := domain.FollowListModel{
		UserID:   user_id,
		ToUserID: to_user_id,
		Status:   0,
	}
	um1 := domain.UserModel{
		Id: user_id,
	}
	um2 := domain.UserModel{
		Id: to_user_id,
	}
	res := m.Mysql.First(&um1)
	if res.Error != nil {
		return false, errors.New("cannot find user " + strconv.FormatInt(user_id, 10))
	}
	res = m.Mysql.First(&um2)
	if res.Error != nil {
		return false, errors.New("cannot find user " + strconv.FormatInt(to_user_id, 10))
	}
	if action_type == 1 {
		flm.Status = 0
		ret := m.Mysql.Save(&flm)
		if ret.Error != nil {
			return false, ret.Error
		}
		um1.FollowCount++
		um2.FollowerCount++
		ret = m.Mysql.Save(&um1)
		if ret.Error != nil {
			return false, ret.Error
		}
		ret = m.Mysql.Save(&um2)
		if ret.Error != nil {
			return false, ret.Error
		}

		// 查询 ToUserID 对 UserID 的关注状态
		flm.UserID = to_user_id
		flm.ToUserID = user_id
		err := m.Mysql.First(&flm).Error
		if errors.Is(err, gorm.ErrRecordNotFound) { // ToUserID 没有关注 UserID
			return true, nil
		}
		if flm.Status == 0 { // ToUserID 关注了 UserID
			flm.Status = 2 // 互相关注
			ret := m.Mysql.Save(&flm)
			if ret.Error != nil {
				return false, ret.Error
			}
		}
	} else if action_type == 2 {
		flm.Status = 1
		ret := m.Mysql.Save(&flm)
		if ret.Error != nil {
			return false, ret.Error
		}
		um1.FollowCount--
		um2.FollowerCount--
		ret = m.Mysql.Save(&um1)
		if ret.Error != nil {
			return false, ret.Error
		}
		ret = m.Mysql.Save(&um2)
		if ret.Error != nil {
			return false, ret.Error
		}

		// 查询 ToUserID 对 UserID 的关注状态
		flm.UserID = to_user_id
		flm.ToUserID = user_id
		err := m.Mysql.First(&flm).Error
		if errors.Is(err, gorm.ErrRecordNotFound) { // ToUserID 没有关注 UserID
			return true, nil
		}
		if flm.Status == 2 { // ToUserID 关注了 UserID
			flm.Status = 0 // 不再互相关注
			ret := m.Mysql.Save(&flm)
			if ret.Error != nil {
				return false, ret.Error
			}
		}
	} else {
		return false, errors.New("action_type must be 1 or 2 !")
	}
	return true, nil
}
