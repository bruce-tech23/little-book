package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		expectErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)

				mockRes := sqlmock.NewResult(10, 1)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(mockRes)

				return db
			},
			ctx: context.Background(),
			user: User{ // 由于所有字段都是 mock 的，所以这些字段不赋值都可以
				Nickname: "A",
			},
		},

		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)

				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnError(&mysqlDriver.MySQLError{Number: 1062})

				return db
			},
			ctx: context.Background(),
			user: User{ // 由于所有字段都是 mock 的，所以这些字段不赋值都可以
				Nickname: "A",
			},
			expectErr: ErrDuplicatedEmail,
		},

		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)

				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("db err"))

				return db
			},
			ctx: context.Background(),
			user: User{ // 由于所有字段都是 mock 的，所以这些字段不赋值都可以
				Nickname: "A",
			},
			expectErr: errors.New("db err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true, // 跳过查询数据库服务的版本
			}), &gorm.Config{
				DisableAutomaticPing:   true, // 不让 gorm 自动发 ping
				SkipDefaultTransaction: true, // 不自动添加事务提交
			})
			assert.NoError(t, err)
			dao := NewUserDao(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
