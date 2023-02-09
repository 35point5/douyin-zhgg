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

func (u *basicUsecase) GetVideoByTime(t time.Time) ([]domain.Video, time.Time) {
	vms := u.basicRepo.GetVideoByTime(t)
	res := make([]domain.Video, 0, len(vms))
	for _, vm := range vms {
		user := u.basicRepo.GetUserById(vm.Uid)
		res = append(res, domain.Video{
			Id:            vm.Id,
			Author:        user,
			PlayUrl:       vm.PlayUrl,
			CoverUrl:      vm.CoverUrl,
			FavoriteCount: vm.FavoriteCount,
			CommentCount:  vm.CommentCount,
			IsFavorite:    vm.IsFavorite,
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
		IsFollow:      false,
	})
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (u *basicUsecase) UserRequest(userauth domain.UserAuth) domain.UserModel {
	//log.Println("user_req id", userauth.Id)
	return u.basicRepo.GetUserById(userauth.Id)
}

func (u *basicUsecase) UserLogin(user domain.UserRegisterRequest) domain.UserModel {
	return u.basicRepo.GetUserByName(user.Username)
}
