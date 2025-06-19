package service

import (
	"context"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicatedEmail = repository.ErrDuplicatedEmail

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
