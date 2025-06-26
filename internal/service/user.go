package service

import (
	"context"
	"errors"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedEmail       = repository.ErrDuplicatedEmail
	ErrInvalidUserOrPassword = errors.New("用户名或密码不正确")
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

//func (svc *UserService) FindById(ctx *gin.Context, uid int64) (interface{}, interface{}) {
//
//}
