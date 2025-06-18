package main

import (
	"geektime-basic-learning2/little-book/internal/repository"
	"geektime-basic-learning2/little-book/internal/repository/dao"
	"geektime-basic-learning2/little-book/internal/service"
	"geektime-basic-learning2/little-book/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	initUserHandler(db, server)
	err := server.Run(":8080")
	if err != nil {
		panic("Server run failed.")
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
		AllowHeaders: []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			// if strings.Contains(origin, "localhost") {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	return server
}
