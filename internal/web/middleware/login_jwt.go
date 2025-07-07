package middleware

import (
	"encoding/gob"
	"geektime-basic-learning2/little-book/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now()) // 注册时间这个结构体类型，为了后面的 sess.Set(updateTimeKey, now) 可以存储这个结构体的字节切片
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" {
			// 这两个接口要么没注册要么没登录，所以不需要校验是否登录
			return
		}
		// 根据约定，token 在 Authorization 头部
		authCode := ctx.GetHeader("Authorization") // 得到的 authCode 格式是 Bearer XXX
		segs := strings.Split(authCode, " ")
		// if len(segs) != 2 包含一种情况就是 authCode 是空字符串的情况，就是没有 Authorization，也就是没有登录。
		if len(segs) != 2 {
			// Authorization 值是乱传的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			// token 不对，token 是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			// token 解析出来了，但是 token 可能是非法的，或者过期了的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			// 监控告警这里需要埋点
			// 能够进来这个分支的，大概率是攻击者
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		// 这个不判断都可以。上面的 token.Valid 可能处理了，只是文档没说明
		//if expireTime.Before(time.Now()) {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		// 剩余过期时间 < 50s 就刷新
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				// 这边不要中断，因为仅仅是过期时间没刷新成功，但是用户是登录了的
				log.Println(err)
			}
		}

		// 将 uc 缓存到 ctx 中，方便后续业务获取 Uid
		ctx.Set("user", uc)
	}
}
