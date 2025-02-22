package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	UserDB *gorm.DB
)

// Init 初始化 MySQL 连接
func Init(dsn string) error {
	var err error
	UserDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info), // 启用详细日志
	})
	if err != nil {
		return fmt.Errorf("connect mysql failed: %w", err)
	}

	sqlDB, err := UserDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置连接池中的最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 设置数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接的最大生命周期

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping mysql failed: %w", err)
	}

	return nil
}

// Close 关闭 MySQL 连接
func Close() error {
	sqlDB, err := UserDB.DB()
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
