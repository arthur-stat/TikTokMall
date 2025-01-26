package mysql

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

func TestCartItem_CRUD(t *testing.T) {
    ctx := context.Background()
    
    // 测试创建购物车商品
    item := &model.CartItem{
        UserID:    1,
        ProductID: 1,
        Quantity:  2,
        Selected:  true,
    }
    err := CreateCartItem(ctx, item)
    assert.NoError(t, err)
    assert.NotZero(t, item.ID)

    // 测试获取购物车商品
    items, err := GetCartItems(ctx, 1)
    assert.NoError(t, err)
    assert.Len(t, items, 1)
    assert.Equal(t, uint32(1), items[0].UserID)
    assert.Equal(t, uint32(1), items[0].ProductID)
    assert.Equal(t, uint32(2), items[0].Quantity)
    assert.True(t, items[0].Selected)

    // 测试更新商品
    item.Quantity = 5
    err = UpdateCartItem(ctx, item)
    assert.NoError(t, err)

    items, err = GetCartItems(ctx, 1)
    assert.NoError(t, err)
    assert.Equal(t, uint32(5), items[0].Quantity)

    // 测试删除商品
    err = DeleteCartItem(ctx, 1, item.ID)
    assert.NoError(t, err)

    items, err = GetCartItems(ctx, 1)
    assert.NoError(t, err)
    assert.Len(t, items, 0)
}

func TestCartItem_BatchOperations(t *testing.T) {
    ctx := context.Background()
    
    // 准备测试数据
    items := []*model.CartItem{
        {UserID: 2, ProductID: 1, Quantity: 1, Selected: true},
        {UserID: 2, ProductID: 2, Quantity: 2, Selected: true},
        {UserID: 2, ProductID: 3, Quantity: 3, Selected: true},
    }

    err := BatchCreateCartItems(ctx, items)
    assert.NoError(t, err)

    // 测试获取购物车商品
    cartItems, err := GetCartItems(ctx, 2)
    assert.NoError(t, err)
    assert.Equal(t, 3, len(cartItems))

    // 测试清空购物车
    err = EmptyCart(ctx, 2)
    assert.NoError(t, err)

    cartItems, err = GetCartItems(ctx, 2)
    assert.NoError(t, err)
    assert.Equal(t, 0, len(cartItems))
}

func TestCartItem_Timestamps(t *testing.T) {
    ctx := context.Background()
    
    // 测试创建时间和更新时间
    item := &model.CartItem{
        UserID:    3,
        ProductID: 1,
        Quantity:  1,
        Selected:  true,
    }
    err := CreateCartItem(ctx, item)
    assert.NoError(t, err)

    // 验证时间戳
    assert.NotZero(t, item.CreatedAt)
    assert.NotZero(t, item.UpdatedAt)
    assert.True(t, item.UpdatedAt.Sub(item.CreatedAt) <= time.Second)

    // 更新商品
    time.Sleep(time.Second) // 确保时间戳有变化
    item.Quantity = 2
    err = UpdateCartItem(ctx, item)
    assert.NoError(t, err)

    items, err := GetCartItems(ctx, 3)
    assert.NoError(t, err)
    assert.True(t, items[0].UpdatedAt.After(items[0].CreatedAt))
}
