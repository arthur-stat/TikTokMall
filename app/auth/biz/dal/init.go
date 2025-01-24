package dal

import (
	"TikTokMall/app/auth/biz/dal/mysql"
	"TikTokMall/app/auth/biz/dal/redis"
)

// Init 初始化数据访问层
func Init() error {
	// 初始化Redis连接
	if err := redis.Init("localhost:6379", "", 0); err != nil {
		return err
	}

	// 初始化MySQL连接
	dsn := "root:root@tcp(localhost:3306)/tiktok_mall_test?charset=utf8mb4&parseTime=True&loc=Local"
	if err := mysql.Init(dsn); err != nil {
		return err
	}

	return nil
}
