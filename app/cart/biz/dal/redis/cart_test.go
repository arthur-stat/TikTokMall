package redis

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "TikTokMall/app/cart/biz/model"
)

func TestMain(m *testing.M) {
    if err := Init(); err != nil {
        panic(err)
    }
    m.Run()
}

func TestCartCache_Basic(t *testing.T) {
    ctx := context.Background()
    userID := uint32(1)

    // 测试缓存购物车商品
    items := []*model.CartItem{
        {
            UserID:    userID,
            ProductID: 1,
            Quantity:  2,
            Selected:  true,
        },
    }
    err := CacheCart(ctx, userID, items)
    assert.NoError(t, err)

    // 测试获取缓存的商品
    cachedItems, err := GetCachedCart(ctx, userID)
    assert.NoError(t, err)
    assert.Len(t, cachedItems, 1)
    assert.Equal(t, uint32(2), cachedItems[0].Quantity)
    assert.True(t, cachedItems[0].Selected)

    // 测试更新缓存的商品
    items[0].Quantity = 5
    err = UpdateCartCache(ctx, userID, items)
    assert.NoError(t, err)

    cachedItems, err = GetCachedCart(ctx, userID)
    assert.NoError(t, err)
    assert.Len(t, cachedItems, 1)
    assert.Equal(t, uint32(5), cachedItems[0].Quantity)

    // 测试删除缓存
    err = InvalidateCartCache(ctx, userID)
    assert.NoError(t, err)

    cachedItems, err = GetCachedCart(ctx, userID)
    assert.Error(t, err) // 应该返回 redis.Nil 错误
}

func TestCartCache_BatchOperations(t *testing.T) {
    ctx := context.Background()
    userID := uint32(2)

    // 准备测试数据
    items := []*model.CartItem{
        {UserID: userID, ProductID: 1, Quantity: 1, Selected: true},
        {UserID: userID, ProductID: 2, Quantity: 2, Selected: true},
        {UserID: userID, ProductID: 3, Quantity: 3, Selected: false},
    }

    // 缓存多个商品
    err := CacheCart(ctx, userID, items)
    assert.NoError(t, err)

    // 测试获取用户的所有缓存商品
    cachedItems, err := GetCachedCart(ctx, userID)
    assert.NoError(t, err)
    assert.Len(t, cachedItems, len(items))

    // 测试清空用户的购物车缓存
    err = InvalidateCartCache(ctx, userID)
    assert.NoError(t, err)

    cachedItems, err = GetCachedCart(ctx, userID)
    assert.Error(t, err) // 应该返回 redis.Nil 错误
}

func TestCartCache_Expiration(t *testing.T) {
    ctx := context.Background()
    userID := uint32(3)

    // 准备测试数据
    items := []*model.CartItem{
        {UserID: userID, ProductID: 1, Quantity: 1, Selected: true},
    }

    // 缓存商品
    err := CacheCart(ctx, userID, items)
    assert.NoError(t, err)

    // 验证商品已被缓存
    cachedItems, err := GetCachedCart(ctx, userID)
    assert.NoError(t, err)
    assert.Len(t, cachedItems, 1)

    // 等待缓存过期（注意：这需要修改 cartExpiration 常量为较小的值进行测试）
    time.Sleep(cartExpiration + time.Second)

    // 验证缓存已过期
    cachedItems, err = GetCachedCart(ctx, userID)
    assert.Error(t, err) // 应该返回 redis.Nil 错误
}
