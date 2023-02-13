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

func Test_basicUsecase_GetUserInfo(t *testing.T) {

}

func Test_basicUsecase_GetVideoByTime(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDoer := mocks.NewMockBasicRepository(mockCtrl)

	// ================ TEST CASES ===============
	mockUser := domain.User{
		Id: 123, Name: "xiaomin", FollowCount: 100, FollowerCount: 1000, IsFollow: false,
	}
	mockUserModel := domain.UserModel{
		Id: 123, Name: "xiaomin", Password: "password", FollowCount: 100, FollowerCount: 1000,
	}
	test_Arg1 := time.Unix(99999999, 0)
	test_Arg2 := int64(123)
	expected_Videos := []domain.Video{
		{
			Id:            1,
			Author:        mockUser,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			IsFavorite:    false,
			Title:         "vd1",
		},
		{
			Id:            2,
			Author:        mockUser,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			IsFavorite:    false,
			Title:         "vd2",
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
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			Title:         "vd1",
			UpdatedTime:   time.Unix(1000000, 0),
		},
		{
			Id:            2,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			Title:         "vd2",
			UpdatedTime:   time.Unix(100000000, 0),
		},
	}
	// ================ MOCK FUNCTION RETURN END ===============

	gomock.InOrder(
		mockDoer.EXPECT().GetVideoByTime(gomock.Any()).Return(mockVideoModel),
		mockDoer.EXPECT().GetUserById(mockUser.Id).Return(mockUserModel),
		mockDoer.EXPECT().IsFollow(gomock.Any(), gomock.Any()).Return(false),
		mockDoer.EXPECT().IsFavorite(mockUser.Id, gomock.Any()).Return(false),
		mockDoer.EXPECT().GetUserById(mockUser.Id).Return(mockUserModel),
		mockDoer.EXPECT().IsFollow(gomock.Any(), gomock.Any()).Return(false),
		mockDoer.EXPECT().IsFavorite(mockUser.Id, gomock.Any()).Return(false),
	)

	u := ucase.NewBasicUsecase(mockDoer)
	actual_Videos, actual_Time := u.GetVideoByTime(test_Arg1, test_Arg2)

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
		},
		{
			Id:            0,
			Name:          "xiaohan",
			Password:      "password",
			FollowCount:   0,
			FollowerCount: 0,
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
