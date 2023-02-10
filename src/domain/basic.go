package domain

import (
	"mime/multipart"
	"time"
)

type VideoModel struct {
	Id            int64  `json:"id,omitempty" gorm:"primarykey;AUTO_INCREMENT"`
	Uid           int64  `json:"author" gorm:"embedded;embeddedPrefix:author_"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	// TODO IsFavorite 这个字段是没用的，并且是不正确的，应该通过 FavoriteListModel 去确认是否点赞
	IsFavorite  bool   `json:"is_favorite,omitempty"`
	Title       string `json:"title,omitempty"`
	UpdatedTime time.Time
}

type UserModel struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty" gorm:"unique"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	// TODO IsFollow 这个字段同样是没用且不正确的，应该通过 UserFollowModel 去确认是否关注
	IsFollow bool `json:"is_follow,omitempty"`
}

type FavoriteListModel struct {
	UserID  int64  `gorm:"primaryKey;autoIncrement:false"`
	VideoID int64  `gorm:"primaryKey;autoIncrement:false"`
	Status  uint32 `json:"status" gorm:"default:1"` //记录是否有效
	//CreatedAt time.Time
}

type UserFollowModel struct {
	// UserId follow TargetUserId
	UserId       int64 `gorm:"primaryKey"`
	TargetUserId int64 `gorm:"primaryKey"`
}

type Video struct {
	Id            int64  `json:"id,omitempty" gorm:"primarykey;AUTO_INCREMENT"`
	Author        User   `json:"author" gorm:"embedded;embeddedPrefix:author_"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty" gorm:"unique"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type FeedRequest struct {
	LatestTime int64  `query:"latest_time"`
	Token      string `query:"token"`
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
	Uid int64 `json:"uid"`
	//jwt.StandardClaims
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
	IsFollow(id int64, fid int64) bool
	//GetFollowCnt(id int64) int64
	//GetFollowerCnt(id int64) int64
	IsFavorite(uid int64, vid int64) bool
	//GetFavoriteCnt(vid int64) int64
}

type BasicUsecase interface {
	GetVideoByTime(t time.Time, uid int64) ([]Video, time.Time)
	UserRegister(user UserRegisterRequest) (int64, error)
	UserRequest(userauth UserAuth) User
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

// publish相关接口
type PublishActionRequest struct {
	Token string               `json:"token" form:"token" query:"token"` // 用户鉴权token
	Data  multipart.FileHeader `json:"data" form:"data" query:"data"`    // 视频数据
	Title string               `json:"title" form:"title" query:"title"` // 视频标题
}

type PublishActionResponse struct {
	StatusCode int32  `json:"status_code" form:"status_code" query:"status_code"`        // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg,omitempty" form:"status_msg" query:"status_msg"` // 返回状态描述
}
type PublishListRequest struct {
	UserId int64  `json:"user_id" form:"user_id" query:"user_id"` // 用户id
	Token  string `json:"token" form:"token" query:"token"`       // 用户鉴权token
}
type PublishListResponse struct {
	StatusCode int32   `json:"status_code" form:"status_code" query:"status_code"`        // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg,omitempty" form:"status_msg" query:"status_msg"` // 返回状态描述
	VideoList  []Video `json:"video_list" form:"video_list" query:"video_list"`           // 用户发布的视频列表
}

type PublishRepository interface {
	AddVideo(v *VideoModel) error
	LisVideoByUserId(userId int64) (videoList []VideoModel, err error)
}
