package repository

import (
	"context"
	"geektime-basic-learning2/little-book/internal/domain"
	"geektime-basic-learning2/little-book/internal/repository/dao"
)

var ErrDuplicatedEmail = dao.ErrDuplicatedEmail // 别名，这样方法里可以直接返回dao层定义的错误

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
