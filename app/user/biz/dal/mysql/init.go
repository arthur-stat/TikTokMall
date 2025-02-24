package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 全局DB实例（与auth服务隔离）
var DB *gorm.DB

// Init 初始化MySQL连接（接收DSN参数）
func Init(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 与auth服务表结构一致
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("mysql connection failed: %w", err)
	}

	// 配置连接池
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("mysql ping failed: %w", err)
	}

	return nil
}

// Close 关闭连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// // 旧mysql初始化代码
//package mysql
//
//import (
//	"TikTokMall/app/user/conf"
//
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//var (
//	DB  *gorm.DB
//	err error
//)
//
//func Init() {
//	DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
//		&gorm.Config{
//			PrepareStmt:            true,
//			SkipDefaultTransaction: true,
//		},
//	)
//	if err != nil {
//		panic(err)
//	}
//}
