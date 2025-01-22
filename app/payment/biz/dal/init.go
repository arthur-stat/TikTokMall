package dal

import (
	"TikTokMall/app/payment/biz/dal/mysql"
	"TikTokMall/app/payment/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
