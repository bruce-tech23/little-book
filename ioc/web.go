package ioc

import (
	"geektime-basic-learning2/little-book/internal/web"
	"geektime-basic-learning2/little-book/internal/web/middleware"
	"geektime-basic-learning2/little-book/pkg/ginx/middleware/ratelimit"
	"geektime-basic-learning2/little-book/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
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
		}),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)).Build(),
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}
