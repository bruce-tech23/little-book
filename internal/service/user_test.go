package service

import (
	"context"
	"errors"
	"fmt"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/repository"
	repomocks "geektime-basic-learning2/little-book/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 预期输入
		ctx      context.Context
		email    string
		password string

		// 预期输出
		expectUser domain.User
		expectErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				retUser := domain.User{Id: 1, Email: "123@qq.com", Password: "$2a$10$K24lVYxr2xeK0OPIx3sDMOuZoWCaWw3Ez5NRqH8OtHv8.UITIbHIy"}
				repo.EXPECT().FindByEmail(gomock.Any(), retUser.Email).Return(retUser, nil)
				return repo
			},
			email:      "123@qq.com",
			password:   "123",
			expectUser: domain.User{Id: 1, Email: "123@qq.com", Password: "$2a$10$K24lVYxr2xeK0OPIx3sDMOuZoWCaWw3Ez5NRqH8OtHv8.UITIbHIy"},
		},
		{
			name: "没有这个用户",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				retUser := domain.User{Email: "123@qq.com"}
				repo.EXPECT().FindByEmail(gomock.Any(), retUser.Email).Return(retUser, repository.ErrUserNotFound)
				return repo
			},
			email:      "123@qq.com",
			password:   "123",
			expectErr:  ErrInvalidUserOrPassword,
			expectUser: domain.User{},
		},
		{
			name: "数据库出错",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				retUser := domain.User{Email: "123@qq.com"}
				repo.EXPECT().FindByEmail(gomock.Any(), retUser.Email).Return(retUser, errors.New("db error"))
				return repo
			},
			email:      "123@qq.com",
			password:   "123",
			expectErr:  errors.New("db error"),
			expectUser: domain.User{},
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				retUser := domain.User{Id: 1, Email: "123@qq.com", Password: "$2a$10$"}
				repo.EXPECT().FindByEmail(gomock.Any(), retUser.Email).Return(retUser, nil)
				return repo
			},
			email:      "123@qq.com",
			password:   "123",
			expectUser: domain.User{},
			expectErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.expectUser, user)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123") // 注意：bcrypt 加密字符串长度不超过72字节
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	fmt.Println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("1234#test"))
	assert.NoError(t, err)
}
