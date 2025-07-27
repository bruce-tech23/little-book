package repository

import (
	"context"
	"database/sql"
	"errors"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/repository/cache"
	cachemocks "geektime-basic-learning2/little-book/internal/repository/cache/mocks"
	"geektime-basic-learning2/little-book/internal/repository/dao"
	daomocks "geektime-basic-learning2/little-book/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao dao.UserDao, c cache.UserCache) // 这里的返回参数要参考 NewCachedUserRepository 的参数

		// 这里是参考 FindById 方法的参数
		ctx context.Context
		uid int64
		//  这里是参考 FindById 方法的返回值
		expectUser domain.User
		expectErr  error
	}{
		{
			name: "查询成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				uid := int64(123)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(
					dao.User{
						Id:       uid,
						Email:    sql.NullString{String: "123@qq.com", Valid: true},
						Password: "123",
						Birthday: 100,
						AboutMe:  "",
						Phone:    sql.NullString{String: "13312341234", Valid: true},
						Ctime:    101,
						Utime:    102,
					},
					nil,
				)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123",
					Birthday: time.UnixMilli(100),
					AboutMe:  "",
					Phone:    "13312341234",
					Ctime:    time.UnixMilli(101),
				}).Return(nil)
				return d, c
			},
			uid: 123,
			ctx: context.Background(),
			expectUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123",
				Birthday: time.UnixMilli(100),
				AboutMe:  "",
				Phone:    "13312341234",
				Ctime:    time.UnixMilli(101),
			},
			expectErr: nil,
		},
		{
			name: "查询成功，缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				uid := int64(123)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123",
					Birthday: time.UnixMilli(100),
					AboutMe:  "",
					Phone:    "13312341234",
					Ctime:    time.UnixMilli(100),
				}, nil)
				return d, c
			},
			uid: 123,
			ctx: context.Background(),
			expectUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123",
				Birthday: time.UnixMilli(100),
				AboutMe:  "",
				Phone:    "13312341234",
				Ctime:    time.UnixMilli(100),
			},
			expectErr: nil,
		},
		{
			name: "查询失败，未找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				uid := int64(123)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(
					dao.User{},
					errors.New("db err"),
				)
				return d, c
			},
			uid:        123,
			ctx:        context.Background(),
			expectUser: domain.User{},
			expectErr:  errors.New("db err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			svc := NewCachedUserRepository(ud, uc)
			user, err := svc.FindById(tc.ctx, tc.uid)
			assert.Equal(t, tc.expectUser, user)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
