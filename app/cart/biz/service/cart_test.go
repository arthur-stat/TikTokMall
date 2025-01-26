package service

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"

    "TikTokMall/app/cart/biz/dal"
    "TikTokMall/app/cart/kitex_gen/cart"
)

func TestMain(m *testing.M) {
    if err := dal.Init(); err != nil {
        panic(err)
    }
    m.Run()
}

func TestCartService_AddItem(t *testing.T) {
    svc := NewCartService()
    ctx := context.Background()

    req := &cart.AddItemReq{
        UserId: 1,
        Item: &cart.CartItem{
            ProductId: 1,
            Quantity: 2,
        },
    }
    err := svc.AddItem(ctx, req)
    assert.NoError(t, err)

    result, err := svc.GetCart(ctx, 1)
    assert.NoError(t, err)
    assert.Equal(t, uint32(1), result.UserId)
    assert.Len(t, result.Items, 1)
    assert.Equal(t, uint32(1), result.Items[0].ProductId)
    assert.Equal(t, int32(2), result.Items[0].Quantity)
}

func TestCartService_GetCart(t *testing.T) {
    svc := NewCartService()
    ctx := context.Background()

    req1 := &cart.AddItemReq{
        UserId: 2,
        Item: &cart.CartItem{
            ProductId: 1,
            Quantity: 1,
        },
    }
    err := svc.AddItem(ctx, req1)
    assert.NoError(t, err)

    req2 := &cart.AddItemReq{
        UserId: 2,
        Item: &cart.CartItem{
            ProductId: 2,
            Quantity: 2,
        },
    }
    err = svc.AddItem(ctx, req2)
    assert.NoError(t, err)

    result, err := svc.GetCart(ctx, 2)
    assert.NoError(t, err)
    assert.Equal(t, uint32(2), result.UserId)
    assert.Len(t, result.Items, 2)
}

func TestCartService_EmptyCart(t *testing.T) {
    svc := NewCartService()
    ctx := context.Background()

    req := &cart.AddItemReq{
        UserId: 5,
        Item: &cart.CartItem{
            ProductId: 1,
            Quantity: 1,
        },
    }
    err := svc.AddItem(ctx, req)
    assert.NoError(t, err)

    err = svc.EmptyCart(ctx, 5)
    assert.NoError(t, err)

    result, err := svc.GetCart(ctx, 5)
    assert.NoError(t, err)
    assert.Equal(t, uint32(5), result.UserId)
    assert.Len(t, result.Items, 0)
}
