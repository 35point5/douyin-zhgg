package usecase

import (
	"douyin-service/domain"
)

type interactUsecase struct {
	basicRepo    domain.BasicRepository
	interactRepo domain.InteractRepository
}

func NewInteractUsecase(basicRepo domain.BasicRepository, interactRepo domain.InteractRepository) domain.InteractUsecase {
	return &interactUsecase{basicRepo, interactRepo}
}

func (u *interactUsecase) GetFavoriteListByUserId(id int64) ([]domain.Video, error) {
	fls, err := u.interactRepo.GetFavoriteListByUserId(id)
	if err != nil {
		return nil, err
	}
	var vids []int64
	for _, fl := range fls {
		vids = append(vids, fl.VideoID)
	}
	vms, err := u.interactRepo.GetVideoModelsById(vids)
	if err != nil {
		return nil, err
	}
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
		return res, nil
	}
	return res, nil
}

func (u *interactUsecase) FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error) {
	return u.interactRepo.FavoriteActionByUserId(user_id, video_id, action_type)
}
