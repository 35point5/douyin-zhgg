package usecase

import (
	"douyin-service/domain"
	"time"
)

type basicUsecase struct {
	basicRepo domain.BasicRepository
}

func NewBasicUsecase(basicRepo domain.BasicRepository) domain.BasicUsecase {
	return &basicUsecase{basicRepo}
}

// GetUserInfo
// fid: whether fid follows uid
func (u *basicUsecase) GetUserInfo(uid int64, fid int64) domain.User {
	um := u.basicRepo.GetUserById(uid)
	return domain.User{
		Id:            uid,
		Name:          um.Name,
		FollowCount:   um.FollowCount,
		FollowerCount: um.FollowerCount,
		IsFollow:      u.basicRepo.IsFollow(fid, uid),
	}
}

func (u *basicUsecase) GetVideoByTime(t time.Time, uid int64) ([]domain.Video, time.Time) {
	vms := u.basicRepo.GetVideoByTime(t)
	res := make([]domain.Video, 0, len(vms))
	for _, vm := range vms {
		res = append(res, domain.Video{
			Id:            vm.Id,
			Author:        u.GetUserInfo(vm.Uid, uid),
			PlayUrl:       vm.PlayUrl,
			CoverUrl:      vm.CoverUrl,
			FavoriteCount: vm.FavoriteCount,
			CommentCount:  vm.CommentCount,
			IsFavorite:    u.basicRepo.IsFavorite(uid, vm.Id),
		})
	}
	if len(vms) == 0 {
		return res, time.Unix(0, 0)
	}
	return res, vms[len(vms)-1].UpdatedTime
}

func (u *basicUsecase) UserRegister(user domain.UserRegisterRequest) (int64, error) {
	uid, err := u.basicRepo.CreateUser(domain.UserModel{
		Id:            0,
		Name:          user.Username,
		Password:      user.Password,
		FollowCount:   0,
		FollowerCount: 0,
	})
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (u *basicUsecase) UserRequest(userauth domain.UserAuth) domain.User {
	//log.Println("user_req id", userauth.Id)
	return u.GetUserInfo(userauth.Id, 0)
}

func (u *basicUsecase) UserLogin(user domain.UserRegisterRequest) domain.UserModel {
	return u.basicRepo.GetUserByName(user.Username)
}
