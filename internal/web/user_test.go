package web

import (
	"bytes"
	"context"
	"errors"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/service"
	svcmocks "geektime-basic-learning2/little-book/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		// mock
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 构造请求，预期中的输入
		reqBuilder func(t *testing.T) *http.Request
		// 预期中的输出
		expectCode int
		expectBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "123#test",
				}).Return(nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "123#test",
"confirmPassword": "123#test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "注册成功",
		},
		{
			// 请求体的 Body 不是 JSON
			name: "Bind 出错",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "123#test",
"confirmPass
`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123.com",
"password": "123#test",
"confirmPassword": "123#test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "邮箱格式不正确",
		},
		{
			name: "两次密码不匹配",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "123#test",
"confirmPassword": "123test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "两次密码不匹配",
		},
		// 密码系统错误测试不了的原因是已经提前编译过了，但业务代码里不要去删掉 err 判断，因为可能被人改了不知道哪里出错了
		{
			name: "密码格式不正确",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "1",
"confirmPassword": "1"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "密码格式不正确",
		},
		{
			name: "邮箱已注册",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "123#test",
				}).Return(service.ErrDuplicatedEmail)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "123#test",
"confirmPassword": "123#test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "邮箱已注册",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "123#test",
				}).Return(errors.New("db error"))
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "123#test",
"confirmPassword": "123#test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { // 注意这里的两个 t 是不一样的。
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造 handler
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			// 准备服务器，注册路由
			server := gin.Default()
			hdl.RegisterRoutes(server)

			// 准备 Req 和记录的 recorder
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.expectCode, recorder.Code)          // 比较响应码
			assert.Equal(t, tc.expectBody, recorder.Body.String()) // 比较 body,注意 recorder.Body 是 buffer(byte)
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 现在的版本(高版本)只要是上面 NewController 传入的是指针的 t 可以不调用这里的 Finish.
	// mock 实现，模拟实现
	userSvc := svcmocks.NewMockUserService(ctrl)
	// 注意调用 userSvc 中的任何业务方法必须先调用 EXPECT
	// 设置了一个模拟场景
	userSvc.EXPECT().Signup(gomock.Any(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	}).Return(nil) // 这里的 Return 要看原来的 user.go 方法具体的返回值，这里返回 nil，就是说预期这里传入的 domain.User 数据是通过的
	err := userSvc.Signup(context.Background(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	})
	t.Log(err)

	userSvc.EXPECT().Signup(gomock.Any(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	}).Return(errors.New("db error"))
	err1 := userSvc.Signup(context.Background(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	})
	t.Log(err1) // db error
}
