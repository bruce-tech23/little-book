package ratelimit

import (
	"errors"
	"geektime-basic-learning2/little-book/internal/service/sms"
	"geektime-basic-learning2/little-book/pkg/limiter"
	"golang.org/x/net/context"
)

var errLimited = errors.New("up to limit")

type RLimitSMSService struct {
	svc     sms.Service // **重要** 被装饰的
	limiter limiter.Limiter
	key     string
}

func (r *RLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewLimitSMSService(svc sms.Service, l limiter.Limiter) *RLimitSMSService {
	return &RLimitSMSService{svc: svc, limiter: l, key: "sms-limiter"}
}
