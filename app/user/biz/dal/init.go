package dal

import (
	"TikTokMall/app/user/biz/dal/mysql"
	"TikTokMall/app/user/biz/dal/redis"
	"TikTokMall/app/user/conf" // 假设存在配置模块
)

// Init 初始化数据访问层
func Init() {
	// 初始化MySQL（使用user服务的独立配置）
	err := mysql.Init(conf.GetConf().MySQL.DSN)
	if err != nil {
		return
	}

	// 初始化Redis（使用user服务的独立配置）
	redisConf := conf.GetConf().Redis
	err = redis.Init(
		redisConf.Address,
		redisConf.Password,
		redisConf.DB,
	)
	if err != nil {
		return
	}
}
