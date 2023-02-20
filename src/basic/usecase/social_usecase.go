package usecase

import "douyin-service/domain"

type socialUsecase struct {
	basicRepo  domain.BasicRepository
	socialRepo domain.SocialRepository
}

func NewSocialUsecase(basicRepo domain.BasicRepository, socialRepo domain.SocialRepository) domain.SocialUsecase {
	return &socialUsecase{basicRepo, socialRepo}
}

func (u *socialUsecase) GetFollowListByUserId(id int64) ([]domain.User, error) {
	fls, err := u.socialRepo.GetFollowListByUserId(id)
	if err != nil {
		return nil, err
	}
	var uids []int64
	for _, fl := range fls {
		uids = append(uids, fl.ToUserID)
	}

	var res []domain.User
	for _, uid := range uids {
		user := u.basicRepo.GetUserById(uid)
		res = append(res, domain.User{
			Id:            uid,
			Name:          user.Name,
			FollowCount:   user.FollowerCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true,
		})
	}
	return res, nil
}

func (u *socialUsecase) GetFollowerListByUserId(id int64) ([]domain.User, error) {
	fls, err := u.socialRepo.GetFollowerListByUserId(id)
	if err != nil {
		return nil, err
	}
	var res []domain.User
	for _, fl := range fls {
		uid := fl.UserID
		um := u.basicRepo.GetUserById(uid)
		user := domain.User{
			Id:            uid,
			Name:          um.Name,
			FollowCount:   um.FollowerCount,
			FollowerCount: um.FollowerCount,
		}
		if fl.Status == 2 {
			user.IsFollow = true
		} else {
			user.IsFollow = false
		}
		res = append(res, user)
	}

	return res, nil
}

func (u *socialUsecase) GetFriendListByUserId(id int64) ([]domain.User, error) {
	fls, err := u.socialRepo.GetFriendListByUserId(id)
	if err != nil {
		return nil, err
	}
	var uids []int64
	for _, fl := range fls {
		uids = append(uids, fl.ToUserID)
	}

	var res []domain.User
	for _, uid := range uids {
		user := u.basicRepo.GetUserById(uid)
		res = append(res, domain.User{
			Id:            uid,
			Name:          user.Name,
			FollowCount:   user.FollowerCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true,
		})
	}
	return res, nil
}

func (u *socialUsecase) FollowActionByUserId(user_id int64, to_user_id int64, action_type int32) (bool, error) {
	return u.socialRepo.FollowActionByUserId(user_id, to_user_id, action_type)
}
