package http

import (
	"context"
	"douyin-service/basic/delivery/http/middleware"
	"douyin-service/domain"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"log"
	"net/http"
)
import "github.com/cloudwego/hertz/pkg/app"

type InteractHandler struct {
	IUsecase domain.InteractUsecase
}

func NewInteractHandler(h *server.Hertz, IUsecase domain.InteractUsecase, mid *middleware.DouyinMiddleware) {
	handler := InteractHandler{IUsecase}
	h.GET("/douyin/favorite/list", handler.GetFavoriteListByUserId)
	h.POST("/douyin/favorite/action", handler.FavoriteAction)
	authGroup := h.Group("/douyin")
	authGroup.Use(mid.TokenAuth())
}

func (t *InteractHandler) GetFavoriteListByUserId(ctx context.Context, c *app.RequestContext) {
	var r domain.FavoriteListRequest
	err := c.Bind(&r)
	fmt.Println("r: ", r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	fmt.Println(r.UserId)
	videos, err := t.IUsecase.GetFavoriteListByUserId(r.UserId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取视频失败",
		})
		return
	}
	fmt.Println(videos)
	c.JSON(http.StatusOK, domain.FavoriteListResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		VideoList: videos,
	})
}

func (t *InteractHandler) FavoriteAction(ctx context.Context, c *app.RequestContext) {
	var r domain.FavoriteActionRequest
	err := c.Bind(&r)
	fmt.Println("r: ", r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	fmt.Println(r.VideoId)
	validateToken := func(token string) (*domain.TokenClaims, error) {
		key := viper.GetString("jwt_key")
		var c domain.TokenClaims
		_, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		fmt.Println(c)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	userClaim, _ := validateToken(r.Token)
	videoBool, err := t.IUsecase.FavoriteActionByUserId(userClaim.Id, r.VideoId, r.ActionType)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，操作错误",
		})
		return
	}
	fmt.Println(videoBool)
	c.JSON(http.StatusOK, domain.FavoriteActionResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}
