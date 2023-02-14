package usecase

import (
	"douyin-service/domain"
	"errors"
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
			Id: vm.Id,
			Author: domain.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false, // TODO 应该去从数据库查
			},
			PlayUrl:       vm.PlayUrl,
			CoverUrl:      vm.CoverUrl,
			FavoriteCount: vm.FavoriteCount,
			CommentCount:  vm.CommentCount,
			IsFavorite:    false, // TODO
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

func (u *interactUsecase) GetCommentListByVideoId(video_id int64, user_id int64) ([]domain.Comment, error) {
	models, err := u.interactRepo.GetCommentListByVideoId(video_id)
	if err != nil {
		return nil, err
	}
	var comments []domain.Comment
	for _, c := range models {
		user_model, err := u.interactRepo.GetUser(c.UserID)
		if err != nil {
			continue
		}
		user := domain.User{
			Id:            user_model.Id,
			Name:          user_model.Name,
			FollowCount:   user_model.FollowCount,
			FollowerCount: user_model.FollowerCount,
			IsFollow:      u.basicRepo.IsFollow(user_id, user_model.Id),
		}
		comment := domain.Comment{
			Id:         c.ID,
			User:       user,
			Content:    c.CommentText,
			CreateDate: c.CreateDate,
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (u *interactUsecase) CommentAction(user_id int64, video_id int64, content string, comment_id int64, action_type int) (domain.Comment, error) {
	var comment_model domain.CommentModel
	var err error
	switch action_type {
	// add comment
	case 1:
		comment_model, err = u.interactRepo.AddCommentByUserId(user_id, video_id, content)
	// delete comment
	case 2:
		comment_model, err = u.interactRepo.DeleteCommentById(user_id, comment_id)
	// unknown
	default:
		return domain.Comment{}, errors.New("未知评论操作！")
	}
	if err != nil {
		return domain.Comment{}, err
	}
	user_model, err := u.interactRepo.GetUser(comment_model.UserID)
	if err != nil {
		return domain.Comment{}, errors.New("用户信息缺失！")
	}
	user := domain.User{
		Id:            user_model.Id,
		Name:          user_model.Name,
		FollowCount:   user_model.FollowCount,
		FollowerCount: user_model.FollowerCount,
	}
	comment := domain.Comment{
		Id:         comment_model.ID,
		User:       user,
		Content:    comment_model.CommentText,
		CreateDate: comment_model.CreateDate,
	}
	return comment, nil
}
