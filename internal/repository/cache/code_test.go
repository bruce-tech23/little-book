package cache

import (
	"context"
	"errors"
	"fmt"
	"geektime-basic-learning2/little-book/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	biz := "login-test"
	phone := "13312341234"
	code := "123456"
	// 因为这里调不了 key 方法，所以复制过来
	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		ctx   context.Context
		biz   string
		phone string
		code  string

		expectErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdAble := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				cmdAble.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{keyFunc(biz, phone)},
					code,
				).Return(cmd)
				return cmdAble
			},
			ctx:       context.Background(),
			biz:       biz,
			phone:     phone,
			code:      code,
			expectErr: nil,
		},

		{
			name: "设置失败，Redis 出错",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdAble := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis err"))
				cmd.SetVal(int64(0))
				cmdAble.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{keyFunc(biz, phone)},
					code,
				).Return(cmd)
				return cmdAble
			},
			ctx:       context.Background(),
			biz:       biz,
			phone:     phone,
			code:      code,
			expectErr: errors.New("redis err"),
		},

		{
			name: "设置失败，Redis key 出错",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdAble := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-1))
				cmdAble.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{keyFunc(biz, phone)},
					code,
				).Return(cmd)
				return cmdAble
			},
			ctx:       context.Background(),
			biz:       biz,
			phone:     phone,
			code:      code,
			expectErr: errors.New("system error"),
		},

		{
			name: "设置失败，发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdAble := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-2))
				cmdAble.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{keyFunc(biz, phone)},
					code,
				).Return(cmd)
				return cmdAble
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
