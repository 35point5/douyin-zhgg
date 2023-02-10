package http

import (
	"context"
	"douyin-service/basic/delivery/http/middleware"
	"douyin-service/domain"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"net/http"
)

type SocialHandler struct {
	SUsecase domain.SocialUsecase
}

func NewSocialHandler(h *server.Hertz, SUsecase domain.SocialUsecase, mid *middleware.DouyinMiddleware) {
	handler := SocialHandler{SUsecase}
	authGroup := h.Group("/douyin/")
	authGroup.Use(mid.TokenAuth())
	authGroup.POST("/relation/action/", handler.RelationAction)
	authGroup.GET("/relation/follow/list/", handler.GetFollowListByUserId)
	authGroup.GET("/relation/follower/list", handler.GetFollowerListByUserId)
	authGroup.GET("/relation/friend/list", handler.GetFriendListByUserId)
}

func (t *SocialHandler) RelationAction(ctx context.Context, c *app.RequestContext) {
	var r domain.RelationActionRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	//fmt.Println(r.ToUserId)
	uid, _ := c.Get("uid")
	followBool, err := t.SUsecase.FollowActionByUserId(uid.(int64), r.ToUserId, r.ActionType)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "关注失败",
		})
		return
	}
	fmt.Println(followBool)
	c.JSON(http.StatusOK, domain.RelationActionResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}

func (t *SocialHandler) GetFollowListByUserId(ctx context.Context, c *app.RequestContext) {
	var r domain.FollowListRequest
	err := c.Bind(&r)
	if err != nil {
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	users, err := t.SUsecase.GetFollowListByUserId(r.UserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "获取关注列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, domain.FollowListResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: users,
	})
}

func (t *SocialHandler) GetFollowerListByUserId(ctx context.Context, c *app.RequestContext) {
	var r domain.FollowerListRequest
	err := c.Bind(&r)
	if err != nil {
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	users, err := t.SUsecase.GetFollowerListByUserId(r.UserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "获取粉丝列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, domain.FollowerListResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: users,
	})
}

func (t *SocialHandler) GetFriendListByUserId(ctx context.Context, c *app.RequestContext) {
	var r domain.FriendListRequest
	err := c.Bind(&r)
	if err != nil {
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	users, err := t.SUsecase.GetFriendListByUserId(r.UserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, domain.Response{
			StatusCode: 1,
			StatusMsg:  "获取好友列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, domain.FriendListResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: users,
	})
}
