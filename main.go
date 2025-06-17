package main

import (
	"geektime-basic-learning2/little-book/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	server := gin.Default()
	// 跨域处理
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		// 旧版本AllowHeaders可以不写，默认的就行，新版本跨域时这个 Origin 必须要显示的写上
		AllowHeaders: []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	hdl := web.NewUserHandler()
	hdl.RegisterRoutes(server)
	err := server.Run(":8080")
	if err != nil {
		panic("Server run failed.")
	}
}
