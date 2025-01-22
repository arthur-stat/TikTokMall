package dal

import (
	"TikTokMall/app/user/biz/dal/mysql"
	"TikTokMall/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
