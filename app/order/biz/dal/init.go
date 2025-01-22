package dal

import (
	"TikTokMall/app/order/biz/dal/mysql"
	"TikTokMall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
