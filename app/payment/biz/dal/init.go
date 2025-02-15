package dal

import (
	"TikTokMall/app/payment/biz/dal/mysql"
	"TikTokMall/app/payment/biz/dal/redis"
)

func Init() error {
	if err := mysql.Init(); err != nil {
		return err
	}
	if err := redis.Init(); err != nil {
		return err
	}
	return nil
}
