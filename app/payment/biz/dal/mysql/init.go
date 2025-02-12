package mysql

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"TikTokMall/app/payment/biz/model"
)

var DB *gorm.DB

// Init 初始化 MySQL 连接
func Init() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&model.Payments{})
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %v", err)
	}

	return nil
}
