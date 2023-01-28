# 项目架构

参考go-clean-arch(https://github.com/bxcodec/go-clean-arch)，以及第四课（https://juejin.cn/course/bytetech/7140987981803814919/section/7141273296397402148）

![clean-arch-1](.\pic\clean-arch-1.png)

![image-20230128155852372](.\pic\image-20230128155852372.png)

```
./src
|-- app //程序入口
|   `-- main.go
|-- basic //对应大项目基础接口（视频、用户）
|   |-- delivery //视图层
|   |   `-- http //web服务
|   |       |-- basic_handler.go //请求处理
|   |       `-- middleware //中间件，例如校验token
|   |           `-- middleware.go
|   |-- repository //数据层
|   |   `-- mysql //与mysql交互
|   |       `-- mysql_basic.go
|   `-- usecase //逻辑层
|       `-- basic_usecase.go
|-- config.json //配置文件
|-- domain //各种结构体、接口的定义
|   `-- basic.go
|-- go.mod
|-- go.sum
`-- public //存放视频等静态文件
```

# 接口编写方法

以feed接口为例，返回视频列表。

## 定义结构体、接口等

```
type VideoModel struct { //数据层使用，与数据库中的表结构相同，用于与数据库交互
	Id            int64  `gorm:"primarykey;AUTO_INCREMENT"` //主键
	Uid           int64  `gorm:"embedded;embeddedPrefix:author_"`
	PlayUrl       string 
	CoverUrl      string 
	FavoriteCount int64  
	CommentCount  int64  
	IsFavorite    bool   
	UpdatedTime   time.Time
}
type Video struct { //逻辑层、视图层使用，用于构建响应体
	Id            int64     `json:"id,omitempty"`
	Author        UserModel `json:"author"`
	PlayUrl       string    `json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
}
type FeedRequest struct { //请求体，用于绑定参数
	LatestTime int64 `query:"latest_time"`
	Token      string
}
type FeedResponse struct { //响应体
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"0`
}
type BasicRepository interface { //定义数据层要实现的方法，提供给逻辑层使用
	GetVideoByTime(t time.Time) []VideoModel
	GetUserById(id int64) UserModel
	CreateUser(user UserModel) (int64, error)
}
type BasicUsecase interface { //定义逻辑层要实现的方法，提供给视图层使用
	GetVideoByTime(t time.Time) ([]Video, time.Time)
	UserRegister(user UserRegisterRequest) (int64, error)
}
```

## 数据层

使用GORM与数据库交互

首先建表，用刚才定义好的数据模型自动迁移即可

```go
func NewMysqlBasicRepository(conn *gorm.DB, debug bool) domain.BasicRepository {
	if err := conn.AutoMigrate(&domain.VideoModel{}, &domain.UserModel{}); err != nil {
		log.Fatal(err)
	}
	return &mysqlBasicRepository{conn}
}
```

然后实现`BasicRepository`接口中的方法，例如根据时间查询视频、根据ID查找用户等

```go
func (m *mysqlBasicRepository) GetVideoByTime(t time.Time) []domain.VideoModel {
	var res []domain.VideoModel
	videoCount := viper.GetInt("video_limit")
	m.Mysql.Where("updated_time < ?", t).Order("updated_time desc").Limit(videoCount).Find(&res)
	return res
}

func (m *mysqlBasicRepository) GetUserById(id int64) domain.UserModel {
	var res domain.UserModel
	res.Id = id
	m.Mysql.First(&res)
	return res
}
```

## 逻辑层

实现`BasicUsecase`接口中的方法，把原始数据包装一下交给视图层

```go
func (u *basicUsecase) GetVideoByTime(t time.Time) ([]domain.Video, time.Time) {
	vms := u.basicRepo.GetVideoByTime(t)
	res := make([]domain.Video, 0, len(vms))
	for _, vm := range vms {
		user := u.basicRepo.GetUserById(vm.Uid)
		res = append(res, domain.Video{
			Id:            vm.Id,
			Author:        user,
			PlayUrl:       vm.PlayUrl,
			CoverUrl:      vm.CoverUrl,
			FavoriteCount: vm.FavoriteCount,
			CommentCount:  vm.CommentCount,
			IsFavorite:    vm.IsFavorite,
		})
	}
	if len(vms) == 0 {
		return res, time.Unix(0, 0)
	}
	return res, vms[len(vms)-1].UpdatedTime
}
```

## 视图层

使用hertz框架，接受web请求，从中获取请求参数，利用逻辑层提供的方法，产生响应体返回

编写处理函数

```go
func (t *BasicHandler) GetVideoByTime(ctx context.Context, c *app.RequestContext) {
	var r domain.FeedRequest
	err := c.Bind(&r) //绑定参数
	fmt.Println("r: ", r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "系统错误，获取参数失败",
		})
		return
	}
	videos, lastTime := t.BUsecase.GetVideoByTime(time.Unix(r.LatestTime/1000, 0)) //调用逻辑层方法
	fmt.Println(videos)
	c.JSON(http.StatusOK, domain.FeedResponse{ //响应
		Response: domain.Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		VideoList: videos,
		NextTime:  lastTime.Unix() * 1000,
	})
}
```

绑定路由

```go
func NewBasicHandler(h *server.Hertz, BUsecase domain.BasicUsecase, mid *middleware.DouyinMiddleware) {
	handler := BasicHandler{BUsecase}
	h.GET("/douyin/feed", handler.GetVideoByTime)
	h.POST(...)
}
```

# 分工

- 基础接口
  - user
  - publish
- 互动接口
  - favourite
  - comment
- 社交接口
  - relation
- 测试、文档等