//go:build wireinject

package startup

import (
	"geektime-basic-learning2/little-book/internal/repository"
	"geektime-basic-learning2/little-book/internal/repository/cache"
	"geektime-basic-learning2/little-book/internal/repository/dao"
	"geektime-basic-learning2/little-book/internal/service"
	"geektime-basic-learning2/little-book/internal/web"
	"geektime-basic-learning2/little-book/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		InitRedis, ioc.InitDB,

		//dao
		dao.NewUserDao,

		// cache
		cache.NewCodeCache, cache.NewUserCache,

		// repository
		repository.NewCodeRepository, repository.NewCachedUserRepository,

		// service
		ioc.InitSMSService, service.NewCodeService, service.NewUserService,

		// handler
		web.NewUserHandler,

		// middlewares
		ioc.InitMiddlewares,
		// 整个 webserver
		ioc.InitWebServer,
	)
	return gin.Default()
}
