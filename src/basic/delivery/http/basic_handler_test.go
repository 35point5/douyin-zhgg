package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	htp "douyin-service/basic/delivery/http"
	"douyin-service/domain"
	mocks "douyin-service/domain/mocks"

	"github.com/golang/mock/gomock"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestBasicHandler_GetVideoByTime(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBasicUsecase := mocks.NewMockBasicUsecase(mockCtrl)

	// ================ MOCK FUNCTION RETURN ===============
	mockUser := domain.User{
		Id: 123, Name: "xiaomin", FollowCount: 100, FollowerCount: 1000, IsFollow: false,
	}
	// mockUserModel := domain.UserModel{
	// 	Id: 123, Name: "xiaomin", Password: "password", FollowCount: 100, FollowerCount: 1000, IsFollow: false,
	// }
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
	// ================ MOCK FUNCTION RETURN END ===============

	// ================ TEST CASES ===============
	test_Args := domain.FeedRequest{
		LatestTime: int64(1675010951000),
		Token:      "",
	}
	test_returnStatus := http.StatusOK
	test_returnBody := domain.FeedResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		VideoList: mockVideos,
		NextTime:  test_Args.LatestTime,
	}
	// ================ TEST CASE END ===============

	mockBasicUsecase.EXPECT().GetVideoByTime(gomock.Any(), gomock.Any()).Return(mockVideos, time.Unix(test_Args.LatestTime/1000, 0))

	h := server.Default()
	hp := &htp.BasicHandler{
		BUsecase: mockBasicUsecase,
	}
	h.GET("/douyin/feed", hp.GetVideoByTime)

	w := ut.PerformRequest(
		h.Engine, "GET", "/douyin/feed?latest_time="+strconv.FormatInt(test_Args.LatestTime, 10),
		&ut.Body{
			Body: bytes.NewBufferString("1"), Len: 1,
		},
		ut.Header{Key: "Connection", Value: "close"},
	)
	resp := w.Result()

	assert.DeepEqual(t, test_returnStatus, resp.StatusCode())
	test_returnBody_String, err := json.Marshal(test_returnBody)
	if err != nil {
		t.Errorf("Marshal failed")
	} else {
		assert.DeepEqual(t, string(test_returnBody_String), string(resp.Body()))
	}
}

func TestBasicHandler_UserRegister(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBasicUsecase := mocks.NewMockBasicUsecase(mockCtrl)

	// ================ MOCK FUNCTION RETURN ===============

	// ================ MOCK FUNCTION RETURN END ===============

	// ================ TEST CASES ===============
	test_Args := domain.UserRegisterRequest{
		Username: "xiaoming",
		Password: "password",
	}
	test_returnStatus := http.StatusOK
	test_returnBody := domain.UserRegisterResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserAuth: domain.UserAuth{
			Id:    233,
			Token: "",
		},
	}
	// ================ TEST CASE END ===============

	mockBasicUsecase.EXPECT().UserRegister(gomock.Any()).Return(int64(233), nil)

	h := server.Default()
	hp := &htp.BasicHandler{
		BUsecase: mockBasicUsecase,
	}
	h.POST("/douyin/user/register/", hp.UserRegister)

	w := ut.PerformRequest(
		h.Engine, "POST",
		"/douyin/user/register/?username="+test_Args.Username+"password="+test_Args.Password,
		&ut.Body{
			Body: bytes.NewBufferString("1"), Len: 1,
		},
		ut.Header{Key: "Connection", Value: "close"},
	)
	resp := w.Result()

	assert.DeepEqual(t, test_returnStatus, resp.StatusCode())
	var actual_return *domain.UserRegisterResponse
	json.Unmarshal(resp.Body(), actual_return)
	assert.DeepEqual(t, test_returnBody.Id, actual_return.Id)
	// test_returnBody_String, err := json.Marshal(test_returnBody)
	// if err != nil {
	// 	t.Errorf("Marshal failed")
	// } else {
	// 	assert.DeepEqual(t, string(test_returnBody_String), string(resp.Body()))
	// }
}

func TestBasicHandler_UserRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBasicUsecase := mocks.NewMockBasicUsecase(mockCtrl)

	// ================ MOCK FUNCTION RETURN ===============
	mockUser := domain.User{
		Id: 123, Name: "xiaomin", FollowCount: 100, FollowerCount: 1000, IsFollow: false,
	}

	// ================ MOCK FUNCTION RETURN END ===============

	// ================ TEST CASES ===============
	test_Arg := domain.UserAuth{
		Id:    123,
		Token: "",
	}
	test_returnStatus := http.StatusOK
	test_returnBody := domain.UserRequesetResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		User: domain.User{
			Id:            mockUser.Id,
			Name:          mockUser.Name,
			FollowCount:   mockUser.FollowCount,
			FollowerCount: mockUser.FollowerCount,
			IsFollow:      mockUser.IsFollow,
		},
	}
	// ================ TEST CASE END ===============

	mockBasicUsecase.EXPECT().UserRequest(gomock.Any()).Return(mockUser)

	h := server.Default()
	hp := &htp.BasicHandler{
		BUsecase: mockBasicUsecase,
	}
	h.GET("/user/", hp.UserRequest)

	w := ut.PerformRequest(
		h.Engine, "GET",
		"/user?user_id="+strconv.FormatInt(test_Arg.Id, 10)+"token="+test_Arg.Token,
		&ut.Body{
			Body: bytes.NewBufferString("1"), Len: 1,
		},
		ut.Header{Key: "Connection", Value: "close"},
	)
	resp := w.Result()

	assert.DeepEqual(t, test_returnStatus, resp.StatusCode())
	test_returnBody_String, err := json.Marshal(test_returnBody)
	if err != nil {
		t.Errorf("Marshal failed")
	} else {
		assert.DeepEqual(t, string(test_returnBody_String), string(resp.Body()))
	}
}

// func TestBasicHandler_UserLogin(t *testing.T) {
// }
