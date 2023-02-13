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
	Title         string `json:"title,omitempty"`
	UpdatedTime   time.Time
}

type UserModel struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty" gorm:"unique"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
}

type FavoriteListModel struct {
	UserID  int64  `gorm:"primaryKey;autoIncrement:false"`
	VideoID int64  `gorm:"primaryKey;autoIncrement:false"`
	Status  uint32 `json:"status" gorm:"default:1"` //记录是否有效
	//CreatedAt time.Time
}

type CommentModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement:true"`
	UserID      int64  `gorm:"autoIncrement:false"`
	VideoID     int64  `gorm:"autoIncrement:false"`
	Status      uint32 `json:"status" gorm:"default:1"`
	CommentText string `json:"comment_text"`
	CreateDate  string
}

// type UserFollowModel struct {
// 	// UserId follow TargetUserId
// 	UserId       int64 `gorm:"primaryKey"`
// 	TargetUserId int64 `gorm:"primaryKey"`
// }

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

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
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

type CommentActionRequest struct {
	Token       string `query:"token"`
	VideoId     int64  `query:"video_id"`
	ActionType  int    `query:"action_type"`
	CommentText string `query:"comment_text"`
	CommentId   int64  `query:"comment_id"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type CommentListRequest struct {
	Token   string `query:"token"`
	VideoId int64  `query:"video_id"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}

type InteractRepository interface {
	GetVideoModelsById(id []int64) ([]VideoModel, error)
	GetFavoriteListByUserId(id int64) ([]FavoriteListModel, error)
	FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error)
	GetCommentListByVideoId(video_id int64) ([]CommentModel, error)
	GetUser(user_id int64) (UserModel, error)
	AddCommentByUserId(user_id int64, video_id int64, content string) (CommentModel, error)
	DeleteCommentById(user_id int64, comment_id int64) (CommentModel, error)
}

type InteractUsecase interface {
	GetFavoriteListByUserId(id int64) ([]Video, error)
	FavoriteActionByUserId(user_id int64, video_id int64, action_type int32) (bool, error)
	GetCommentListByVideoId(video_id int64) ([]Comment, error)
	CommentAction(user_id int64, video_id int64, content string, comment_id int64, action_type int) (Comment, error)
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
