package http

import (
	"context"
	"douyin-service/basic/delivery/http/middleware"
	"douyin-service/domain"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BasicHandler struct {
	BUsecase domain.BasicUsecase
	Mid      *middleware.DouyinMiddleware
}

func NewBasicHandler(h *server.Hertz, BUsecase domain.BasicUsecase, mid *middleware.DouyinMiddleware) {
	handler := BasicHandler{BUsecase, mid}
	//staticURL := viper.GetString("static_url")
	// 这三个不用Token验证
	h.GET("/douyin/feed", handler.GetVideoByTime)
	h.POST("/douyin/user/register/", handler.UserRegister)
	h.POST("/douyin/user/login/", handler.UserLogin)
	//这需要Token验证
	authGroup := h.Group("/douyin/")
	authGroup.Use(mid.TokenAuth())
	authGroup.GET("/ping/", ping)
	authGroup.GET("/user/", handler.UserRequest)
	//h.GET("/douyin/user/", mid.TokenAuth(), handler.UserRequest)
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
	claims, err := t.Mid.JWTMid.GetClaimsFromJWT(ctx, c)
	uid := int64(0)
	if err == nil {
		uid, _ = strconv.ParseInt(claims["uid"].(string), 10, 64)
		log.Println("feed uid:", uid)
	}
	videos, lastTime := t.BUsecase.GetVideoByTime(time.Unix(r.LatestTime/1000, 0), uid)
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

func (t *BasicHandler) generateToken(claims *domain.TokenClaims) (string, error) {
	str, _, err := t.Mid.JWTMid.TokenGenerator(claims)
	return str, err
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//key := viper.GetString("jwt_key")
	//return token.SignedString([]byte(key))
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
	//domainName := viper.GetString("domain")
	token, err := t.generateToken(&domain.TokenClaims{
		Uid: uid,
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

func (t *BasicHandler) UserRequest(ctx context.Context, c *app.RequestContext) {
	var r domain.UserAuth
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	log.Println(r)
	user := t.BUsecase.UserRequest(r)
	fmt.Println("user_req", user)
	c.JSON(http.StatusOK, domain.UserRequesetResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		User: user,
	})
}

func (t *BasicHandler) UserLogin(ctx context.Context, c *app.RequestContext) {
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
	user := t.BUsecase.UserLogin(r)
	if r.Password != user.Password {
		log.Println("Password Wrong")
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "密码输入错误",
		})
		return
	}
	uid := user.Id
	//domainName := viper.GetString("domain")
	token, err := t.generateToken(&domain.TokenClaims{
		Uid: uid,
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
