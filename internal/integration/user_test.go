package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"geektime-basic-learning2/little-book/internal/integration/startup"
	"geektime-basic-learning2/little-book/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// init 方法设置了 gin 为 Release 模式，可以减少日志输出。如果找日志能力很强，去掉就行。
func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestUserHandler_SendSMSCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()
	testCases := []struct {
		name string

		// before
		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		expectCode int
		expectBody web.Result
	}{
		{
			name: "验证码发送成功",
			before: func(t *testing.T) {

			},

			after: func(t *testing.T) {
				// 验证 Redis 中有值且有正确的过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := fmt.Sprintf("phone_code:%s:%s", "login", "19912341234")
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*4+time.Second*50)
				t.Log("dur is ", dur)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},

			phone:      "19912341234",
			expectCode: http.StatusOK,
			expectBody: web.Result{
				Msg: "succeed",
			},
		},

		{
			name: "手机号码是空",
			before: func(t *testing.T) {

			},

			after: func(t *testing.T) {
			},

			phone:      "",
			expectCode: http.StatusOK,
			expectBody: web.Result{
				Code: 4,
				Msg:  "Please input your phone number",
			},
		},

		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := fmt.Sprintf("phone_code:%s:%s", "login", "19912341234")
				err := rdb.Set(ctx, key, "123456", time.Minute*4+time.Second*10).Err()
				assert.NoError(t, err)
			},

			after: func(t *testing.T) {
				// 验证 Redis 中有值且有正确的过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := fmt.Sprintf("phone_code:%s:%s", "login", "19912341234")
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},

			phone:      "19912341234",
			expectCode: http.StatusOK,
			expectBody: web.Result{
				Code: 4,
				Msg:  "send too many",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			// 准备 Req 和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.expectCode, recorder.Code) // 比较响应码
			if tc.expectCode != http.StatusOK {
				return
			}
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectBody, res)
		})
	}
}
