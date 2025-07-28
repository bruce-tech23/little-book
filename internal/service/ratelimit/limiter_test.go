package ratelimit

import (
	"errors"
	"geektime-basic-learning2/little-book/internal/service/sms"
	smsmocks "geektime-basic-learning2/little-book/internal/service/sms/mocks"
	"geektime-basic-learning2/little-book/pkg/limiter"
	limitermocks "geektime-basic-learning2/little-book/pkg/limiter/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"testing"
)

func TestRLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		// 一个输入都没有

		// 预期输出
		expectErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), "sms-limiter").Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
		},

		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return svc, l
			},
			expectErr: errLimited,
		},

		{
			name: "limiter 异常",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("unknown err"))
				return svc, l
			},
			expectErr: errors.New("unknown err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)
			svc := NewLimitSMSService(smsSvc, l)
			err := svc.Send(context.Background(), "test", []string{"1234"}, "13312341234")
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
