package http

import (
	"bytes"
	"context"
	"douyin-service/basic/delivery/http/middleware"
	"douyin-service/domain"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"time"
)

type PublishHandler struct {
	r          domain.PublishRepository
	basicR     domain.BasicRepository
	domainName string
	staticURL  string
	staticPath string
}

func NewPublishHandler(h *server.Hertz, r domain.PublishRepository, basicR domain.BasicRepository, mid *middleware.DouyinMiddleware) {
	handler := PublishHandler{r, basicR, viper.GetString("domain"), viper.GetString("static_url"), viper.GetString("static_path") + viper.GetString("static_url")}
	g := h.Group("/douyin/publish")
	g.POST("/action/", mid.TokenAuth(), handler.Publish)
	g.GET("/list/", mid.TokenAuth(), handler.List)
}

func WriteVideo(DirPath, videoName string, data []byte) error {
	videoFullPath := DirPath + "/" + videoName
	file, err := os.Create(videoFullPath)
	if err != nil {
		log.Fatal("创建视频文件失败：", err)
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		log.Fatal("写视频文件失败：", err)
		return err
	}
	err = file.Close()
	if err != nil {
		log.Fatal("关闭视频文件失败：", err)
		return err
	}
	return nil
}

func GetFirstFrame(DirPath, videoName string, frameNum int) (imageName string, err error) {
	videoFullPath := DirPath + "/" + videoName
	imageName = videoName + ".png"
	imageFullPath := DirPath + "/" + imageName
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoFullPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, imageFullPath)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}
	return imageName, nil
}

// Publish .
// @router /douyin/publish/action [POST]
func (h *PublishHandler) Publish(ctx context.Context, c *app.RequestContext) {
	var err error
	var req domain.PublishActionRequest
	resp := new(domain.PublishActionResponse)
	err = c.Bind(&req)
	if err != nil {
		log.Println(err)
		resp.StatusCode = 1
		resp.StatusMsg = "参数错误"
		c.JSON(consts.StatusOK, resp)
		return
	}
	f, _ := req.Data.Open()
	data := make([]byte, req.Data.Size)
	f.Read(data)
	uid, _ := c.Get("uid")
	userId := uid.(int64)
	videoName := uuid.New().String() + ".mp4"
	if err = WriteVideo(h.staticPath, videoName, data); err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "保存Video失败"
		c.JSON(consts.StatusOK, resp)
		return
	}
	imageName, err := GetFirstFrame(h.staticPath, videoName, 1)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "截取首帧失败"
		c.JSON(consts.StatusOK, resp)
		return
	}

	vm := domain.VideoModel{
		Uid:           userId,
		PlayUrl:       h.domainName + h.staticURL + videoName,
		CoverUrl:      h.domainName + h.staticURL + imageName,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         req.Title,
		UpdatedTime:   time.Now(),
	}
	if err = h.r.AddVideo(&vm); err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "插入数据库失败"
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp.StatusCode = 0
	c.JSON(consts.StatusOK, resp)
}

// List .
// @router /douyin/publish/list [GET]
func (h *PublishHandler) List(ctx context.Context, c *app.RequestContext) {
	var err error
	var req domain.PublishListRequest
	resp := new(domain.PublishListResponse)
	err = c.Bind(&req)
	if err != nil {
		resp.StatusCode = 1
		c.JSON(consts.StatusOK, resp)
		return
	}

	uid, _ := c.Get("uid")
	userId := uid.(int64)
	vmList, err := h.r.LisVideoByUserId(userId)
	if err != nil {
		resp.StatusCode = 1
		c.JSON(consts.StatusOK, resp)
		return
	}
	log.Println(vmList)
	videoList := make([]domain.Video, 0)
	for _, vm := range vmList {
		user := h.basicR.GetUserById(vm.Uid)
		video := domain.Video{
			Id: vm.Id,
			Author: domain.User{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false, // 自己不会关注自己？
			},
			PlayUrl:       vm.PlayUrl,
			CoverUrl:      vm.CoverUrl,
			FavoriteCount: vm.FavoriteCount,
			CommentCount:  vm.CommentCount,
		}
		videoList = append(videoList, video)
	}

	resp.StatusCode = 0
	resp.VideoList = videoList
	c.JSON(consts.StatusOK, resp)
}
