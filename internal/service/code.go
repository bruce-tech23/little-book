package service

import (
	"context"
	"errors"
	"fmt"
	"geektime-basic-learning2/little-book/internal/repository"
	"geektime-basic-learning2/little-book/internal/service/sms"
	"math/rand"
)

const codeSMSTplId = "12"

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送验证码
	return svc.sms.Send(ctx, codeSMSTplId, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 相当于我们对外面屏蔽了验证次数过多的错误，就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
