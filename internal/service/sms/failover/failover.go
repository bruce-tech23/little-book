package failover

import (
	"errors"
	"geektime-basic-learning2/little-book/internal/service/sms"
	"golang.org/x/net/context"
	"log"
	"sync/atomic"
)

type FOSMSService struct {
	svcs []sms.Service // 候选的服务商

	// V1 的字段
	idx uint64 // 当前服务商下标
}

func NewFOSMSService(svcs []sms.Service) *FOSMSService {
	return &FOSMSService{svcs: svcs}
}

func (f *FOSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}

	// 轮询了所有的服务商，但是发送都失败了
	return errors.New("failed")
}

func (f *FOSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))

	// 迭代 length 次
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			// 前者是被调用者主动取消，后者是超时
			return err
		}
		log.Println(err)
	}

	return errors.New("final fail")
}
