package mysql

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"douyin-service/domain"
)

func Test_mysqlInteractRepository_GetVideoModelsById(t *testing.T) {
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := &mysqlInteractRepository{db}

	// ================ TEST CASES ===============
	test_Args := []int64{123}
	mockUser := domain.UserModel{
		Id: 123, Name: "xiaomin", Password: "1234567", FollowCount: 100, FollowerCount: 1000,
	}
	expected_VideoModels := []domain.VideoModel{
		{
			Id:            1,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			UpdatedTime:   time.Unix(1000001000, 0),
		},
		{
			Id:            2,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			UpdatedTime:   time.Unix(1000002000, 0),
		},
	}
	// ================ TEST CASES END ===============

	sql := "SELECT * FROM `video_models` WHERE `video_models`.`id` = ? ORDER BY updated_time desc LIMIT 0"

	mockRows := sqlmock.NewRows([]string{"id", "uid", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "updated_time"})
	for _, expd_v := range expected_VideoModels {
		mockRows.AddRow(expd_v.Id, expd_v.Uid, expd_v.PlayUrl, expd_v.CoverUrl, expd_v.FavoriteCount, expd_v.CommentCount, expd_v.UpdatedTime)
	}
	mock.ExpectQuery(sql).WillReturnRows(mockRows)

	actual_VideoModels, err := repository.GetVideoModelsById(test_Args)

	assert.NoError(t, err)
	assert.Equal(t, expected_VideoModels, actual_VideoModels)
}

func Test_mysqlInteractRepository_GetFavoriteListByUserId(t *testing.T) {

}

func Test_mysqlInteractRepository_FavoriteActionByUserId(t *testing.T) {

}
