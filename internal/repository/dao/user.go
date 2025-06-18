package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	return dao.db.WithContext(ctx).Create(&u).Error
}

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Email    string
	Password string
	Ctime    int64 // 创建时间
	Utime    int64 // 更新时间
}
