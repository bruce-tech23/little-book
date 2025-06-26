package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now()) // 注册时间这个结构体类型，为了后面的 sess.Set(updateTimeKey, now) 可以存储这个结构体的字节切片
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 这两个接口要么没注册要么没登录，所以不需要校验是否登录
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// cookie 过期时间刷新放到这里，因为基本所有的接口都会走这个中间件
		/*
			过期时间设置在 web/user.go Login 接口的 session 设置
			sess.Options(sessions.Options{
					MaxAge: 900, // 过期时间，单位为秒
			})
		*/
		// 我怎么知道，要刷新了呢？
		// 假如说，我们的策略是每分钟刷一次，我怎么知道，已经过了一分钟？
		now := time.Now()
		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute*1 {
			// val == nil 第一次进来
			// !ok sess 获取上次刷新时间失败
			// now.Sub(lastUpdateTime) > time.Minute * 1 超过了一分钟，需要刷新
			sess.Set(updateTimeKey, now)
			// 由于 gin 的 session 是覆盖机制，就是上面 sess.Set(updateTimeKey, now) 会把 session 中的其他值清空，所以需要重新设置
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
