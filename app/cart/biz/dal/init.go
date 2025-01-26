package dal

import (
    "TikTokMall/app/cart/biz/dal/mysql"
    "TikTokMall/app/cart/biz/dal/redis"
)

// Init 初始化数据访问层
func Init() error {
    // 初始化 MySQL
    if err := mysql.Init(); err != nil {
        return err
    }
    // 初始化 Redis
    if err := redis.Init(); err != nil {
        return err
    }
    return nil
} 
