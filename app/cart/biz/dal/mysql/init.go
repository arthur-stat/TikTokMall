package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"TikTokMall/app/cart/biz/model"
	"TikTokMall/app/cart/conf"
)

var DB *gorm.DB

// Init 初始化 MySQL 连接
func Init() error {
	config := conf.GetConf()
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&model.CartItem{})
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %v", err)
	}

	return nil
}
