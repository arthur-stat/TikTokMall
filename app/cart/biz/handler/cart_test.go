package handler

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"

    "TikTokMall/app/cart/biz/dal/mysql"
    "TikTokMall/app/cart/biz/dal/redis"
    "TikTokMall/app/cart/kitex_gen/cart"
)

func init() {
    // 初始化数据库连接
    if err := mysql.Init(); err != nil {
        panic(err)
    }
    if err := redis.Init(); err != nil {
        panic(err)
    }
}

// TestNewCartServiceImpl tests the creation of a new cart service implementation
func TestNewCartServiceImpl(t *testing.T) {
    handler := NewCartServiceImpl()
    assert.NotNil(t, handler)
    assert.NotNil(t, handler.svc)
}

// TestCartServiceImpl_AddItem tests the AddItem method
func TestCartServiceImpl_AddItem(t *testing.T) {
    handler := NewCartServiceImpl()
    ctx := context.Background()
    req := &cart.AddItemReq{
        UserId: 101,
        Item: &cart.CartItem{
            ProductId: 1,
            Quantity: 1,
        },
    }

    // Test success case
    resp, err := handler.AddItem(ctx, req)
    assert.NoError(t, err)
    assert.NotNil(t, resp)

    // Test invalid quantity case
    req.Item.Quantity = 0
    resp, err = handler.AddItem(ctx, req)
    assert.Error(t, err)
    assert.NotNil(t, resp)
}

// TestCartServiceImpl_GetCart tests the GetCart method
func TestCartServiceImpl_GetCart(t *testing.T) {
    handler := NewCartServiceImpl()
    ctx := context.Background()

    // First add an item to the cart
    addReq := &cart.AddItemReq{
        UserId: 102,
        Item: &cart.CartItem{
            ProductId: 1,
            Quantity: 1,
        },
    }
    _, err := handler.AddItem(ctx, addReq)
    assert.NoError(t, err)

    // Test getting the cart
    getReq := &cart.GetCartReq{
        UserId: 102,
    }
    resp, err := handler.GetCart(ctx, getReq)
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.NotNil(t, resp.Cart)
    assert.NotEmpty(t, resp.Cart.Items)

    // Test getting non-existent cart
    getReq.UserId = 9999
    resp, err = handler.GetCart(ctx, getReq)
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.NotNil(t, resp.Cart)
    assert.Empty(t, resp.Cart.Items)
}

// TestCartServiceImpl_EmptyCart tests the EmptyCart method
func TestCartServiceImpl_EmptyCart(t *testing.T) {
    handler := NewCartServiceImpl()
    ctx := context.Background()

    // First add some items to the cart
    addReq := &cart.AddItemReq{
        UserId: 103,
        Item: &cart.CartItem{
            ProductId: 1,
            Quantity: 1,
        },
    }
    _, err := handler.AddItem(ctx, addReq)
    assert.NoError(t, err)

    // Test emptying the cart
    emptyReq := &cart.EmptyCartReq{
        UserId: 103,
    }
    resp, err := handler.EmptyCart(ctx, emptyReq)
    assert.NoError(t, err)
    assert.NotNil(t, resp)

    // Verify cart is empty
    getReq := &cart.GetCartReq{
        UserId: 103,
    }
    getResp, err := handler.GetCart(ctx, getReq)
    assert.NoError(t, err)
    assert.NotNil(t, getResp)
    assert.NotNil(t, getResp.Cart)
    assert.Empty(t, getResp.Cart.Items)
}
