package failover

import (
	"errors"
	"geektime-basic-learning2/little-book/internal/service/sms"
	smsmocks "geektime-basic-learning2/little-book/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"testing"
)

func TestFOSMSService_Send(t *testing.T) {
	testCases := []struct {
		name  string
		mocks func(ctrl *gomock.Controller) []sms.Service

		expectErr error
	}{
		{
			name: "一次发送成功",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
		},

		{
			name: "第二次发送成功",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed"))

				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
		},

		{
			name: "全部失败",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed0"))

				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed1"))
				return []sms.Service{svc0, svc1}
			},
			expectErr: errors.New("failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFOSMSService(tc.mocks(ctrl))
			err := svc.Send(context.Background(), "test", []string{"1234"}, "13312341234")
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
