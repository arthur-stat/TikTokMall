package model

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestCartItem_TableName(t *testing.T) {
    item := CartItem{}
    assert.Equal(t, "cart_items", item.TableName())
}

func TestCartItem_Fields(t *testing.T) {
    now := time.Now()
    item := CartItem{
        ID:        1,
        UserID:    100,
        ProductID: 200,
        Quantity:  2,
        Selected:  true,
        CreatedAt: now,
        UpdatedAt: now,
    }

    assert.Equal(t, uint32(1), item.ID)
    assert.Equal(t, uint32(100), item.UserID)
    assert.Equal(t, uint32(200), item.ProductID)
    assert.Equal(t, uint32(2), item.Quantity)
    assert.True(t, item.Selected)
    assert.Equal(t, now, item.CreatedAt)
    assert.Equal(t, now, item.UpdatedAt)
} 
