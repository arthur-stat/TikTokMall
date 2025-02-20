package mysql

import (
	"TikTokMall/app/product/conf"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DB"))
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			TranslateError: true,
		},
	)
	if err != nil {
		panic(err)
	}
}
