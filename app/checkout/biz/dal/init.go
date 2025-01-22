package dal

import (
	"TikTokMall/app/checkout/biz/dal/mysql"
	"TikTokMall/app/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
