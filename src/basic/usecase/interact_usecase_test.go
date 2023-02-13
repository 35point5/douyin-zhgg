package usecase_test

import (
	ucase "douyin-service/basic/usecase"
	"douyin-service/domain"
	mocks "douyin-service/domain/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_interactUsecase_GetFavoriteListByUserId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDoerInteractRepo := mocks.NewMockInteractRepository(mockCtrl)
	mockDoerBasicRepo := mocks.NewMockBasicRepository(mockCtrl)

	mockUser := domain.User{
		Id: 123, Name: "xiaomin", FollowCount: 100, FollowerCount: 1000, IsFollow: false,
	}
	mockUserModel := domain.UserModel{
		Id: 123, Name: "xiaomin", Password: "password", FollowCount: 100, FollowerCount: 1000,
	}
	mockVideos := []domain.Video{
		{
			Id:            1,
			Author:        mockUser,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			IsFavorite:    false,
		},
		{
			Id:            2,
			Author:        mockUser,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			IsFavorite:    false,
		},
	}

	// ================ MOCK FUNCTION ARGS ===============
	mockId := mockUser.Id
	mockVid := []int64{mockVideos[0].Id, mockVideos[1].Id}
	mockErrorId := int64(999)
	// ================ MOCK FUNCTION ARGS END ===============

	// ================ MOCK FUNCTION RETURN ===============
	mockFavoriteListModels := []domain.FavoriteListModel{
		{
			UserID:  mockUser.Id,
			VideoID: mockVideos[0].Id,
			Status:  1,
		},
		{
			UserID:  mockUser.Id,
			VideoID: mockVideos[1].Id,
			Status:  1,
		},
	}
	mockVideoModles := []domain.VideoModel{
		{
			Id:            mockVideos[0].Id,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			UpdatedTime:   time.Unix(100000001, 0),
		},
		{
			Id:            mockVideos[1].Id,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			UpdatedTime:   time.Unix(100000000, 0),
		},
	}
	// ================ MOCK FUNCTION RETURN END ===============

	gomock.InOrder(
		mockDoerInteractRepo.EXPECT().GetFavoriteListByUserId(mockId).Return(mockFavoriteListModels, nil),
		mockDoerInteractRepo.EXPECT().GetVideoModelsById(mockVid).Return(mockVideoModles, nil),
		mockDoerBasicRepo.EXPECT().GetUserById(mockId).Return(mockUserModel).MaxTimes(2),
	)

	u := ucase.NewInteractUsecase(mockDoerBasicRepo, mockDoerInteractRepo)

	// ================ TEST CASES ===============
	test_Arg := mockUser.Id
	expected_Videos := mockVideos
	// ================ TEST CASES END ===============

	actual_Videos, err := u.GetFavoriteListByUserId(test_Arg)
	assert.NoError(t, err)
	assert.Equal(t, expected_Videos, actual_Videos)

	gomock.InOrder(
		mockDoerInteractRepo.EXPECT().GetFavoriteListByUserId(mockErrorId).Return(nil, errors.New("test")),
	)

	// ================ TEST CASES ===============
	test_Arg = mockErrorId
	expected_Videos = mockVideos
	// ================ TEST CASES END ===============

	actual_Videos, err = u.GetFavoriteListByUserId(test_Arg)

	assert.Error(t, err)
	assert.Nil(t, actual_Videos)

}

func Test_interactUsecase_FavoriteActionByUserId(t *testing.T) {

}
