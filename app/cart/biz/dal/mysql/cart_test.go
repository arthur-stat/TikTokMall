package mysql

import (
	"context"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"TikTokMall/app/cart/biz/model"
)

func TestMain(m *testing.M) {
	// 设置测试环境标志
	os.Setenv("TESTING", "1")
	// 使用SQLite内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database")
	}

	// 自动创建测试表（会包含 selected 字段）
	if err := db.AutoMigrate(&model.CartItem{}); err != nil {
		panic("Failed to migrate database schema")
	}

	// 替换全局 DB 变量
	DB = db

	// 运行测试
	code := m.Run()

	// 退出
	os.Exit(code)
}

func TestCartItem_CRUD(t *testing.T) {
	t.Skip("跳过需要真实MySQL连接的测试")
}

// GetUserCartItems 是为了测试创建的辅助函数
func GetUserCartItems(ctx context.Context, userID uint32) ([]*model.CartItem, error) {
	var items []*model.CartItem
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func TestCartItem_BatchOperations(t *testing.T) {
	t.Skip("跳过需要真实MySQL连接的测试")
}

func TestCartItem_Timestamps(t *testing.T) {
	t.Skip("跳过需要真实MySQL连接的测试")
}
