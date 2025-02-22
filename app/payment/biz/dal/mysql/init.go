package mysql

import (
	"fmt"
	"log"
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
		getEnvWithFallback("MYSQL_USER", "gorm"),
		getEnvWithFallback("MYSQL_PASSWORD", "gorm"),
		getEnvWithFallback("MYSQL_HOST", "127.0.0.1"),
		getEnvWithFallback("MYSQL_PORT", "3307"),
		getEnvWithFallback("MYSQL_DATABASE", "gorm"),
	)
	log.Printf("Connecting to MySQL with DSN: %s", dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	if os.Getenv("GO_ENV") != "online" {
		err := DB.AutoMigrate(
			&model.Payments{},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func getEnvWithFallback(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
