package main

import (
	_basicDelivery "douyin-service/basic/delivery/http"
	"douyin-service/basic/delivery/http/middleware"
	_basicRepo "douyin-service/basic/repository/mysql"
	_basicUC "douyin-service/basic/usecase"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var debug bool

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?&parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	//redisHost := viper.GetString("redis.host")
	//redisPort := viper.GetString("redis.port")
	//redisPass := viper.GetString("redis.pass")
	//redisDb := viper.GetInt("redis.db")
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     redisHost + ":" + redisPort,
	//	Password: redisPass,
	//	DB:       redisDb,
	//})
	//_, err = rdb.Ping().Result()
	//if err != nil {
	//	panic(err)
	//}
	debug = viper.GetBool("debug")
	basicRepo := _basicRepo.NewMysqlBasicRepository(db, debug)
	basicUC := _basicUC.NewBasicUsecase(basicRepo)
	domainName := viper.GetString("listen")
	h := server.Default(server.WithHostPorts(domainName))
	middle := middleware.DouyinMiddleware{}
	_basicDelivery.NewBasicHandler(h, basicUC, &middle)
	log.Fatal(h.Run())
}
