package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestRedisUserCache_Set_e2e(t *testing.T) {
	biz := "redis-code-test"
	phone := "13312341234"
	code := "123456"
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:16379"})
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		ctx   context.Context
		biz   string
		phone string
		code  string

		expectErr error
	}{
		{
			name: "设置成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := fmt.Sprintf("phone_code:%s:%s", biz, phone)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*4+time.Second*50)

				rdCode, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, code, rdCode)

			},
			ctx:       context.Background(),
			biz:       biz,
			phone:     phone,
			code:      code,
			expectErr: nil,
		},

		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := fmt.Sprintf("phone_code:%s:%s", biz, phone)
				err := rdb.Set(ctx, key, code, time.Minute*4+time.Second*55).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := fmt.Sprintf("phone_code:%s:%s", biz, phone)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*4+time.Second*50)

				rdCode, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, code, rdCode)

			},
			ctx:       context.Background(),
			biz:       biz,
			phone:     phone,
			code:      code,
			expectErr: ErrCodeSendTooMany,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			c := NewCodeCache(rdb)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
