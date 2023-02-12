package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	htp "douyin-service/basic/delivery/http"
	"douyin-service/domain"
	mocks "douyin-service/domain/mocks"

	"github.com/golang/mock/gomock"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestInteractHandler_GetFavoriteListByUserId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockIUsecase := mocks.NewMockInteractUsecase(mockCtrl)

	// ================ MOCK FUNCTION RETURN ===============
	mockUser := domain.User{
		Id: 123, Name: "xiaomin", FollowCount: 100, FollowerCount: 1000, IsFollow: true,
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
	// ================ MOCK FUNCTION RETURN END ===============

	// ================ TEST CASES ===============
	test_Args := domain.FavoriteListRequest{
		UserId: 123,
		Token:  "",
	}
	test_returnStatus := http.StatusOK
	test_returnBody := domain.FavoriteListResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		VideoList: mockVideos,
	}
	// ================ TEST CASE END ===============

	mockIUsecase.EXPECT().GetFavoriteListByUserId(gomock.Any()).Return(mockVideos, nil)

	h := server.Default()
	hp := &htp.InteractHandler{
		IUsecase: mockIUsecase,
	}
	h.GET("/favorite/list/", hp.GetFavoriteListByUserId)

	w := ut.PerformRequest(
		h.Engine, "GET", "/favorite/list/?user_id="+strconv.FormatInt(test_Args.UserId, 10),
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

// func TestInteractHandler_FavoriteAction(t *testing.T) {
// }
