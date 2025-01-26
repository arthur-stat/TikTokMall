package biz

import (
    "TikTokMall/app/cart/biz/dal/mysql"
    "TikTokMall/app/cart/biz/dal/redis"
)

// Init initializes all business layer components
func Init() error {
    // Initialize Redis first
    redis.Init()
    
    // Then initialize MySQL
    mysql.Init()
    
    return nil
} 
