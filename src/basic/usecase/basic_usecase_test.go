// gotests -all -w .\basic_usecase.go
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

func Test_basicUsecase_GetVideoByTime(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDoer := mocks.NewMockBasicRepository(mockCtrl)

	// ================ TEST CASES ===============
	mockAuthor := domain.UserModel{
		Id: 123, Name: "xiaomin", Password: "1234567", FollowCount: 100, FollowerCount: 1000, IsFollow: true,
	}
	test_Args := time.Unix(99999999, 0)
	expected_Videos := []domain.Video{
		{
			Id:            1,
			Author:        mockAuthor,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			IsFavorite:    false,
		},
		{
			Id:            2,
			Author:        mockAuthor,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			IsFavorite:    false,
		},
	}
	expected_Time := time.Unix(100000000, 0)

	// ================ TEST CASE END ===============

	// ================ MOCK FUNCTION ARGS ===============
	// ================ MOCK FUNCTION ARGS END ===============

	// ================ MOCK FUNCTION RETURN ===============
	mockVideoModel := []domain.VideoModel{
		{
			Id:            1,
			Uid:           mockAuthor.Id,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			IsFavorite:    false,
			UpdatedTime:   time.Unix(1000000, 0),
		},
		{
			Id:            2,
			Uid:           mockAuthor.Id,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			IsFavorite:    false,
			UpdatedTime:   time.Unix(100000000, 0),
		},
	}
	// ================ MOCK FUNCTION RETURN END ===============

	gomock.InOrder(
		mockDoer.EXPECT().GetVideoByTime(test_Args).Return(mockVideoModel),
		mockDoer.EXPECT().GetUserById(mockAuthor.Id).Return(mockAuthor).AnyTimes(),
	)

	u := ucase.NewBasicUsecase(mockDoer)
	actual_Videos, actual_Time := u.GetVideoByTime(test_Args)

	assert.Equal(t, expected_Time, actual_Time)
	assert.Equal(t, expected_Videos, actual_Videos)
}

func Test_basicUsecase_UserRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDoer := mocks.NewMockBasicRepository(mockCtrl)

	// ================ MOCK FUNCTION ARGS ===============
	mockUserModels := []domain.UserModel{
		{
			Id:            0,
			Name:          "xiaoming",
			Password:      "password",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		},
		{
			Id:            0,
			Name:          "xiaohan",
			Password:      "password",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		},
	}
	// ================ MOCK FUNCTION ARGS END ===============

	// ================ MOCK FUNCTION RETURN ===============
	// ================ MOCK FUNCTION RETURN END ===============

	gomock.InOrder(
		mockDoer.EXPECT().CreateUser(mockUserModels[0]).Return(
			int64(1), nil,
		),
		mockDoer.EXPECT().CreateUser(mockUserModels[1]).Return(
			int64(0), errors.New("test"),
		),
	)

	u := ucase.NewBasicUsecase(mockDoer)

	// ================ TEST CASES ===============
	test_Args := domain.UserRegisterRequest{
		Username: "xiaoming",
		Password: "password",
	}
	expected_Uid := int64(1)
	// ================ TEST CASES END ===============

	autual_Uid, err := u.UserRegister(test_Args)
	assert.NoError(t, err)
	assert.Equal(t, expected_Uid, autual_Uid)

	// ================ TEST CASES ===============
	test_Args = domain.UserRegisterRequest{
		Username: "xiaohan",
		Password: "password",
	}
	expected_Uid = int64(0)
	// ================ TEST CASES END ===============

	autual_Uid, err = u.UserRegister(test_Args)
	assert.Error(t, err)
	assert.Equal(t, expected_Uid, autual_Uid)
}

func Test_basicUsecase_UserRequest(t *testing.T) {

}

func Test_basicUsecase_UserLogin(t *testing.T) {

}
