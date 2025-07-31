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

func TestTimeoutFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []sms.Service
		threshold int32
		idx       int32
		cnt       int32

		expectErr error
		expectCnt int32
		expectIdx int32
	}{
		{
			name: "没有触发切换",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			idx:       0,
			cnt:       12,
			threshold: 15,
			expectErr: nil,
			expectIdx: 0,
			expectCnt: 0,
		},

		{
			name: "触发切换，成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			expectErr: nil,
			expectCnt: 0,
			expectIdx: 1,
		},

		{
			name: "触发切换，未知错误失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("field"))
				return []sms.Service{svc0, svc1}
			},
			idx:       1,
			cnt:       15,
			threshold: 15,
			expectErr: errors.New("field"),
			expectCnt: 0,
			expectIdx: 0,
		},

		{
			name: "触发切换，超时失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			idx:       1,
			cnt:       15,
			threshold: 15,
			expectErr: context.DeadlineExceeded,
			expectCnt: 1,
			expectIdx: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailOverSMSService(tc.mock(ctrl), tc.threshold)
			svc.cnt = tc.cnt
			svc.idx = tc.idx
			err := svc.Send(context.Background(), "test", []string{"1234"}, "13312341234")
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, tc.expectCnt, svc.cnt)
			assert.Equal(t, tc.expectIdx, svc.idx)
		})
	}
}
