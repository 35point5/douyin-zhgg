package http

import (
	"context"
	"douyin-service/basic/delivery/http/middleware"
	"douyin-service/domain"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)
import "github.com/cloudwego/hertz/pkg/app"

type BasicHandler struct {
	BUsecase domain.BasicUsecase
}

func NewBasicHandler(h *server.Hertz, BUsecase domain.BasicUsecase, mid *middleware.DouyinMiddleware) {
	handler := BasicHandler{BUsecase}
	//staticURL := viper.GetString("static_url")
	staticPath := viper.GetString("static_path")
	h.Static("/douyin/static/", staticPath)
	h.GET("/douyin/feed", handler.GetVideoByTime)
	h.POST("/douyin/user/register", handler.UserRegister)
	authGroup := h.Group("/douyin")
	authGroup.Use(mid.TokenAuth())
	authGroup.GET("/ping", ping)
}

func ping(ctx context.Context, c *app.RequestContext) {
	uid, _ := c.Get("uid")
	c.JSON(http.StatusOK, utils.H{"uid": uid})
}

func (t *BasicHandler) GetVideoByTime(ctx context.Context, c *app.RequestContext) {
	var r domain.FeedRequest
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
	videos, lastTime := t.BUsecase.GetVideoByTime(time.Unix(r.LatestTime/1000, 0))
	fmt.Println(videos)
	c.JSON(http.StatusOK, domain.FeedResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		VideoList: videos,
		NextTime:  lastTime.Unix() * 1000,
	})
}

func generateToken(claims domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := viper.GetString("jwt_key")
	return token.SignedString([]byte(key))
}

func (t *BasicHandler) UserRegister(ctx context.Context, c *app.RequestContext) {
	var r domain.UserRegisterRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	uid, err := t.BUsecase.UserRegister(r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，创建用户失败",
		})
		return
	}
	domainName := viper.GetString("domain")
	token, err := generateToken(domain.TokenClaims{
		Id: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    domainName,
		},
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，创建token失败",
		})
		return
	}
	c.JSON(http.StatusOK, domain.UserRegisterResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserAuth: domain.UserAuth{
			Id:    uid,
			Token: token,
		},
	})
}
