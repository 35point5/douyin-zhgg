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
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type User struct {
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
	Id    int64  `query:"user_id" json:"user_id"`
	Token string `query:"token" json:"token"`
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
	User User `json:"user"`
}

type BasicRepository interface {
	GetVideoByTime(t time.Time) []VideoModel
	//SetToken(uid int64, token string) error
	GetUserById(id int64) UserModel
	GetUserByName(name string) UserModel
	CreateUser(user UserModel) (int64, error)
	//UserRegister(user UserRegisterRequest) (UserModel, string)
}

type BasicUsecase interface {
	GetVideoByTime(t time.Time) ([]Video, time.Time)
	UserRegister(user UserRegisterRequest) (int64, error)
	UserRequest(userauth UserAuth) UserModel
	UserLogin(user UserRegisterRequest) UserModel
}

type FavoriteListRequest struct {
	UserId int64  `query:"user_id"` // 用户id
	Token  string // 用户鉴权token
}

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

type FavoriteListModel struct {
	UserID  int64  `gorm:"primaryKey;autoIncrement:false"`
	VideoID int64  `gorm:"primaryKey;autoIncrement:false"`
	Status  uint32 `json:"status" gorm:"default:1"` //记录是否有效
	//CreatedAt time.Time
}

type InteractRepository interface {
	GetVideoModelsById(id []int64) ([]VideoModel, error)
	GetFavoriteListByUserId(id int64) ([]FavoriteListModel, error)
	FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error)
}

type InteractUsecase interface {
	GetFavoriteListByUserId(id int64) ([]Video, error)
	FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error)
}

type FavoriteActionRequest struct {
	Token      string `query:"token"`       // 用户鉴权token
	VideoId    int64  `query:"video_id"`    // 视频id
	ActionType int32  `query:"action_type"` // 1-点赞，2-取消点赞
}

type FavoriteActionResponse struct {
	Response
}

// 社交功能
type FollowListModel struct {
	UserID   int64  `gorm:"primaryKey;autoIncrement:false"`
	ToUserID int64  `gorm:"primaryKey;autoIncrement:false"`
	Status   uint32 `json:"status" gorm:"default:0"`
	// 0 代表 UserID 关注了 ToUserID
	// 1 代表 UserID 取关了 ToUserID
	// 2 代表 UserID 和 ToUserID 互相关注
}

type SocialRepository interface {
	GetFollowListByUserId(id int64) ([]FollowListModel, error)
	GetFollowerListByUserId(id int64) ([]FollowListModel, error)
	GetFriendListByUserId(id int64) ([]FollowListModel, error)
	FollowActionByUserId(user_id int64, to_user_id int64, action_type int32) (bool, error)
}
type SocialUsecase interface {
	GetFollowListByUserId(id int64) ([]User, error)
	GetFollowerListByUserId(id int64) ([]User, error)
	GetFriendListByUserId(id int64) ([]User, error)
	FollowActionByUserId(user_id int64, to_user_id int64, action_type int32) (bool, error)
}
type RelationActionRequest struct {
	Token      string `query:"token"`
	ToUserId   int64  `query:"to_user_id"`
	ActionType int32  `query:"action_type"`
}

type RelationActionResponse struct {
	Response
}

type FollowListRequest struct {
	UserId int64 `query:"user_id"`
	Token  string
}

type FollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type FollowerListRequest struct {
	UserId int64 `query:"user_id"`
	Token  string
}

type FollowerListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type FriendListRequest struct {
	UserId int64 `query:"user_id"`
	Token  string
}

type FriendListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
