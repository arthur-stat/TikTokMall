package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCartItem(t *testing.T) {
	// 测试模型基本属性
	item := &CartItem{
		ID:        1,
		UserID:    100,
		ProductID: 200,
		Quantity:  2,
		Selected:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 测试字段值是否正确
	assert.Equal(t, uint32(1), item.ID)
	assert.Equal(t, uint32(100), item.UserID)
	assert.Equal(t, uint32(200), item.ProductID)
	assert.Equal(t, uint32(2), item.Quantity)
	assert.True(t, item.Selected)
	assert.NotZero(t, item.CreatedAt)
	assert.NotZero(t, item.UpdatedAt)

	// 测试表名方法
	assert.Equal(t, "cart_items", item.TableName())
}
