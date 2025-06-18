package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	// 严格地来说这不是好的实践
	// 1. 还是应该走审批流程去建表
	// 2. 这个方法放到 dao 里不伦不类，意味着 dao 和 DB 是强耦合的。如果没有强耦合，dao可以在不同的 DB 之间切换
	// 3. 有其他的 表 也需要再在这里添加
	return db.AutoMigrate(&User{})
}
