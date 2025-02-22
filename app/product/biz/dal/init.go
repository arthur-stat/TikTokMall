package dal

import (
	"TikTokMall/app/product/biz/dal/mysql"
	// "TikTokMall/app/product/biz/dal/redis"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
