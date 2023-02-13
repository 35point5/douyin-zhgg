package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"douyin-service/domain"
)

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func Test_mysqlBasicRepository_GetVideoByTime(t *testing.T) {
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := &mysqlBasicRepository{db}

	// ================ TEST CASES ===============
	test_Args := time.Now()
	mockUser := domain.UserModel{
		Id: 123, Name: "xiaomin", Password: "1234567", FollowCount: 100, FollowerCount: 1000, IsFollow: true,
	}
	expected_VideoModels := []domain.VideoModel{
		{
			Id:            1,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl1",
			CoverUrl:      "CoverUrl1",
			FavoriteCount: 1,
			CommentCount:  2,
			IsFavorite:    false,
			UpdatedTime:   time.Unix(1000001000, 0),
		},
		{
			Id:            2,
			Uid:           mockUser.Id,
			PlayUrl:       "PlayUrl2",
			CoverUrl:      "CoverUrl2",
			FavoriteCount: 12,
			CommentCount:  26,
			IsFavorite:    false,
			UpdatedTime:   time.Unix(1000002000, 0),
		},
	}
	// ================ TEST CASES END ===============

	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE updated_time <= ? ORDER BY updated_time desc LIMIT 0", "video_models")

	mockRows := sqlmock.NewRows([]string{"id", "uid", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "updated_time"})
	for _, expd_v := range expected_VideoModels {
		mockRows.AddRow(expd_v.Id, expd_v.Uid, expd_v.PlayUrl, expd_v.CoverUrl, expd_v.FavoriteCount, expd_v.CommentCount, expd_v.IsFavorite, expd_v.UpdatedTime)
	}
	mock.ExpectQuery(sql).WillReturnRows(mockRows)

	actual_VideoModels := repository.GetVideoByTime(test_Args)

	assert.Equal(t, expected_VideoModels, actual_VideoModels)
}

func Test_mysqlBasicRepository_GetUserById(t *testing.T) {
	db, mock, err := getDBMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := &mysqlBasicRepository{db}

	// ================ TEST CASES ===============
	test_Args := int64(123)
	expected_User := domain.UserModel{
		Id: 123, Name: "xiaomin", Password: "1234567", FollowCount: 100, FollowerCount: 1000, IsFollow: true,
	}
	// ================ TEST CASES END ===============

	sql := "SELECT * FROM `user_models` WHERE `user_models`.`id` = ? ORDER BY `user_models`.`id` LIMIT 1"

	mockRows := sqlmock.
		NewRows([]string{"id", "name", "password", "follow_count", "follower_count", "is_follow"}).
		AddRow(expected_User.Id, expected_User.Name, expected_User.Password, expected_User.FollowCount, expected_User.FollowerCount, expected_User.IsFollow)

	mock.ExpectQuery(sql).WillReturnRows(mockRows)

	actual_User := repository.GetUserById(test_Args)

	assert.Equal(t, expected_User, actual_User)
}

// func Test_mysqlBasicRepository_GetUserByName(t *testing.T) {

// }

// func Test_mysqlBasicRepository_CreateUser(t *testing.T) {

// }
