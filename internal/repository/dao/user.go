package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"` // 这些字段要去官网复制，以免写错难以排查问题 https://gorm.io/zh_CN/docs/models.html
	Email    string `gorm:"unique"`
	Password string
	// 不要使用任何和时区有关的数据，最好的事使用 UTC 0 的毫秒数，要处理时区要去 domain 的对象上定义，然后在给前端的逻辑中处理
	Ctime int64 // 创建时间
	Utime int64 // 更新时间
}
