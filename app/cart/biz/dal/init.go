package dal

import (
	"TikTokMall/app/cart/biz/dal/mysql"
	"TikTokMall/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
