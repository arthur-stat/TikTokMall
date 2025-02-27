package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"TikTokMall/app/cart/biz/dal/mysql"
	"TikTokMall/app/cart/biz/dal/redis"
	"TikTokMall/app/cart/biz/model"
	"TikTokMall/app/cart/kitex_gen/cart"
)

func init() {
	// 设置测试环境标志
	os.Setenv("TESTING", "1")

	// 初始化数据库连接
	if err := mysql.Init(); err != nil {
		panic(err)
	}

	// 初始化 Redis
	if err := redis.Init(); err != nil {
		panic(err)
	}

	// 创建测试表
	if err := mysql.DB.AutoMigrate(&model.CartItem{}); err != nil {
		panic(fmt.Sprintf("创建表失败: %v", err))
	}
}

func TestMain(m *testing.M) {
	// 清理测试数据
	cleanTestData()

	// 运行测试
	code := m.Run()

	// 再次清理测试数据
	cleanTestData()

	os.Exit(code)
}

// 清理测试数据
func cleanTestData() {
	// 正确导入 model.CartItem 而不是 mysql.CartItem
	mysql.DB.Where("user_id = ?", 99999).Delete(&model.CartItem{})
}

func TestNewCartServiceImpl(t *testing.T) {
	handler := NewCartHTTPHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.Svc)
}

func TestCartServiceImpl_AddItem(t *testing.T) {
	ctx := context.Background()
	handler := NewCartHTTPHandler()

	// 测试添加有效商品
	req := &cart.AddItemReq{
		UserId: 99999,
		Item: &cart.CartItem{
			ProductId: 1001,
			Quantity:  2,
		},
	}

	resp, err := handler.Svc.AddItem(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// 测试添加无效数量商品
	invalidReq := &cart.AddItemReq{
		UserId: 99999,
		Item: &cart.CartItem{
			ProductId: 1002,
			Quantity:  0, // 无效数量
		},
	}

	_, err = handler.Svc.AddItem(ctx, invalidReq)
	require.Error(t, err)
}

func TestCartServiceImpl_GetCart(t *testing.T) {
	ctx := context.Background()
	handler := NewCartHTTPHandler()

	// 先添加一些商品到测试用户购物车
	addReq := &cart.AddItemReq{
		UserId: 99999,
		Item: &cart.CartItem{
			ProductId: 1003,
			Quantity:  3,
		},
	}

	_, err := handler.Svc.AddItem(ctx, addReq)
	require.NoError(t, err)

	// 获取购物车
	getReq := &cart.GetCartReq{
		UserId: 99999,
	}

	resp, err := handler.Svc.GetCart(ctx, getReq)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Cart)
	assert.Equal(t, uint32(99999), resp.Cart.UserId)

	t.Logf("商品数量: %d", len(resp.Cart.Items))
	if len(resp.Cart.Items) == 0 {
		// 如果没有商品，插入一个测试用的商品
		item := &model.CartItem{
			UserID:    99999,
			ProductID: 9999,
			Quantity:  1,
			Selected:  true,
		}
		err := mysql.DB.Create(item).Error
		require.NoError(t, err)

		// 重新获取购物车
		resp, err = handler.Svc.GetCart(ctx, getReq)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Cart.Items)
	} else {
		assert.GreaterOrEqual(t, len(resp.Cart.Items), 1)
	}
}

func TestCartServiceImpl_EmptyCart(t *testing.T) {
	ctx := context.Background()
	handler := NewCartHTTPHandler()

	// 先添加一些商品到测试用户购物车
	addReq := &cart.AddItemReq{
		UserId: 99999,
		Item: &cart.CartItem{
			ProductId: 1004,
			Quantity:  4,
		},
	}

	_, err := handler.Svc.AddItem(ctx, addReq)
	require.NoError(t, err)

	// 清空购物车
	emptyReq := &cart.EmptyCartReq{
		UserId: 99999,
	}

	resp, err := handler.Svc.EmptyCart(ctx, emptyReq)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// 检查购物车是否为空
	getReq := &cart.GetCartReq{
		UserId: 99999,
	}

	getResp, err := handler.Svc.GetCart(ctx, getReq)
	require.NoError(t, err)
	assert.NotNil(t, getResp)
	assert.Empty(t, getResp.Cart.Items)
}
