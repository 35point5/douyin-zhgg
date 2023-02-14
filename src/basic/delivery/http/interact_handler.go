package http

import (
	"context"
	"douyin-service/basic/delivery/http/middleware"
	"douyin-service/domain"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type InteractHandler struct {
	IUsecase domain.InteractUsecase
}

func NewInteractHandler(h *server.Hertz, IUsecase domain.InteractUsecase, mid *middleware.DouyinMiddleware) {
	handler := InteractHandler{IUsecase}
	authGroup := h.Group("/douyin")
	authGroup.Use(mid.TokenAuth())
	authGroup.POST("/favorite/action/", handler.FavoriteAction)
	authGroup.GET("/favorite/list/", handler.GetFavoriteListByUserId)
	authGroup.POST("/comment/action/", handler.CommentAction)
	authGroup.GET("/comment/list/", handler.GetCommentByVideoId)
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
	log.Println("fav videos:", videos)
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
	//validateToken := func(token string) (*domain.TokenClaims, error) {
	//	key := viper.GetString("jwt_key")
	//	var c domain.TokenClaims
	//	_, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
	//		return []byte(key), nil
	//	})
	//	fmt.Println(c)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &c, nil
	//}
	//
	//userClaim, _ := validateToken(r.Token)
	uid, _ := c.Get("uid")
	videoBool, err := t.IUsecase.FavoriteActionByUserId(uid.(int64), r.VideoId, r.ActionType)
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

func (t *InteractHandler) GetCommentByVideoId(ctx context.Context, c *app.RequestContext) {
	var request domain.CommentListRequest
	err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	uid, _ := c.Get("uid")
	comments, err := t.IUsecase.GetCommentListByVideoId(request.VideoId, uid.(int64))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取评论列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, domain.CommentListResponse{
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		CommentList: comments,
	})
}

func (t *InteractHandler) CommentAction(ctx context.Context, c *app.RequestContext) {
	var request domain.CommentActionRequest
	err := c.Bind(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.CommentActionResponse{
			Response: domain.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			Comment: domain.Comment{},
		})
		return
	}
	uid, _ := c.Get("uid")
	comment, err := t.IUsecase.CommentAction(uid.(int64), request.VideoId, request.CommentText, request.CommentId, request.ActionType)
	fmt.Println(comment)
	if err != nil {
		c.JSON(http.StatusOK, domain.CommentActionResponse{
			Response: domain.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			Comment: comment,
		})
	} else {
		c.JSON(http.StatusOK, domain.CommentActionResponse{
			Response: domain.Response{
				StatusCode: 0,
				StatusMsg:  "OK",
			},
			Comment: comment,
		})
	}
}
