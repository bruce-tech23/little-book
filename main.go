package main

import (
	"geektime-basic-learning2/little-book/internal/repository"
	"geektime-basic-learning2/little-book/internal/repository/dao"
	"geektime-basic-learning2/little-book/internal/service"
	"geektime-basic-learning2/little-book/internal/web"
	"geektime-basic-learning2/little-book/internal/web/middleware"
	"geektime-basic-learning2/little-book/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	sessredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

//func main() {
//	db := initDB()
//	server := initWebServer()
//	initUserHandler(db, server)
//	err := server.Run(":8080")
//	if err != nil {
//		panic("Server run failed.")
//	}
//}

func main() {
	// Kubernetes 练习。去除对 MySQL 和 Redis 依赖
	/*
		 * 部署三个实例
		也就是说，需要一个 Service, 一个 Deployment，这个 Deployment 管着三个 Pod，每一个 Pod 是一个实例。
	*/
	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})
	err := server.Run(":8080")
	if err != nil {
		panic("Server start failed.")
	}
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/littlebook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 跨域处理
	server.Use(cors.New(cors.Config{
		AllowCredentials: true, // 是否允许带认证信息(例如 cookie)过来
		// 旧版本AllowHeaders可以不写，默认的就行，新版本跨域时这个 Content-Type 必须要显示的写上
		AllowHeaders: []string{"Content-Type", "Origin", "Authorization"},
		// 这个是允许跨域时前端访问后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			// if strings.Contains(origin, "localhost") {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build()) // 1秒最多100个请求
	useJWT(server)
	//useSession(server)
	return server
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	// session 本身初始化
	// 存储数据的 cookie，也就是 userId 存放的地方
	//store := cookie.NewStore([]byte("secret"))

	// redis 存储
	store, err := sessredis.NewStore(16, "tcp", "localhost:16379", "", "",
		[]byte("TGvRiyJZTMNSRrZhZndRrEIOQ14cqF7E"), []byte("TGvRiyJZTMNSRrZhZndRrEIOQ14cqF7E")) // 16 最大空闲连接数
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
