package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Video struct {
	Id            int64     `json:"id,omitempty" gorm:"primarykey;AUTO_INCREMENT"`
	Author        UserModel `json:"author" gorm:"embedded;embeddedPrefix:author_"`
	PlayUrl       string    `json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
}

type VideoModel struct {
	Id            int64  `json:"id,omitempty" gorm:"primarykey;AUTO_INCREMENT"`
	Uid           int64  `json:"author" gorm:"embedded;embeddedPrefix:author_"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	UpdatedTime   time.Time
}

type UserModel struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty" gorm:"unique"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type FeedRequest struct {
	LatestTime int64 `query:"latest_time"`
	Token      string
}

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserRegisterRequest struct {
	Username string `query:"username"`
	Password string `query:"password"`
}

type UserAuth struct {
	Id    int64  `query:"userid"`
	Token string `query:"token"`
}

type TokenClaims struct {
	Id int64
	jwt.StandardClaims
}

type UserRegisterResponse struct {
	Response
	UserAuth
}

type UserRequesetResponse struct {
	Response
	UserModel
}
type BasicRepository interface {
	GetVideoByTime(t time.Time) []VideoModel
	//SetToken(uid int64, token string) error
	GetUserById(id int64) UserModel
	CreateUser(user UserModel) (int64, error)
	//UserRegister(user UserRegisterRequest) (UserModel, string)
}

type BasicUsecase interface {
	GetVideoByTime(t time.Time) ([]Video, time.Time)
	UserRegister(user UserRegisterRequest) (int64, error)
	UserRequest(userauth UserAuth) UserModel
}
