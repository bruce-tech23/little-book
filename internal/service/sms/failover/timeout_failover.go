package failover

import (
	"errors"
	"geektime-basic-learning2/little-book/internal/service/sms"
	"golang.org/x/net/context"
	"sync/atomic"
)

type TimeoutFailOverSMSService struct {
	svcs []sms.Service

	idx int32 // 当前正在使用节点
	cnt int32 // 连续几个超时了

	threshold int32 // 切换的阈值，只读。也就是每个节点超时的最大次数
}

func NewTimeoutFailOverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{svcs: svcs, threshold: threshold}
}

func (t *TimeoutFailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold { // 超过阈值，执行切换
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 重置 cnt 计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}

	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch {
	case err == nil:
		// 连续超时，所以不超时的时候要重置到 0
		atomic.StoreInt32(&t.cnt, 0)
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&t.cnt, 1)
		//default:
		//	// 遇到了不是超时的错误，这个时候，要考虑怎么设计
		//	// 可以增加超时次数，也可以不增加
		//	// 如果强调一定是超时，那么就不要增加
		//	// 如果是 EOF 之类的错误，还可以考虑直接切换
	}
	return err
}
