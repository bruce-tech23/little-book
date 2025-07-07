package service

import (
	"context"
	"errors"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedEmail       = repository.ErrDuplicatedUser
	ErrInvalidUserOrPassword = errors.New("用户名或密码不正确")
	ErrCodeSendTooMany       = repository.ErrCodeVerifyTooMany
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encryptPassword)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	// UpdateNicknameAndXXAnd
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *UserService) FindById(ctx context.Context,
	uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先找一下，大部分用户是已经存在的
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		// 两种情况
		// err == nil u是可用的
		// err != nil，系统错误
		return u, nil
	}
	// 用户没找到，注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 两种可能
	// 一种是 err 恰好是唯一索引冲突(phone)
	// 一种是 err != nil，系统错误
	if err != nil && !errors.Is(err, repository.ErrDuplicatedUser) {
		return domain.User{}, err
	}

	// 要么 err 是 nil，要么是 ErrDuplicatedUser，都代表用户存在的
	// 主从延迟，理论上来讲，这里需要强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
