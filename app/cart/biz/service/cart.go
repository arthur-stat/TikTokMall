package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"TikTokMall/app/cart/biz/dal/redis"
	"TikTokMall/app/cart/biz/model"
	"TikTokMall/app/cart/kitex_gen/cart"
)

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than 0")
	ErrUserNotFound    = errors.New("user not found")
)

// CartItem 购物车商品
type CartItem struct {
	UserID    int64  `json:"user_id"`
	ProductID int64  `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Image     string `json:"image"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// CartService 购物车服务接口
type CartService interface {
	AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error)
	GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error)
	EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error)
}

// cartService 购物车服务实现
type cartService struct {
	// 可以添加数据库访问层等依赖
}

// NewCartService 创建购物车服务
func NewCartService() CartService {
	return &cartService{}
}

// AddItem 添加商品到购物车
func (s *cartService) AddItem(ctx context.Context, req *cart.AddItemReq) (*cart.AddItemResp, error) {
	// 这里应该实现添加商品到购物车的逻辑
	// 例如，将商品信息存储到Redis或MySQL中
	now := time.Now().Unix()

	// 示例：使用Redis存储购物车数据
	key := fmt.Sprintf("cart:%d", req.UserId)
	field := fmt.Sprintf("%d", req.Item.ProductId)

	// 如果数量为0，则删除该商品
	if req.Item.Quantity == 0 {
		err := redis.RDB.HDel(ctx, key, field).Err()
		if err != nil {
			return nil, fmt.Errorf("从购物车删除商品失败: %w", err)
		}
		return &cart.AddItemResp{}, nil
	}

	// 这里只是示例，实际应该将item序列化后存储
	value := fmt.Sprintf("%d:%d:%d", req.Item.Quantity, req.Item.ProductId, now)
	err := redis.RDB.HSet(ctx, key, field, value).Err()
	if err != nil {
		return nil, fmt.Errorf("添加商品到Redis失败: %w", err)
	}

	return &cart.AddItemResp{}, nil
}

// GetCart 获取购物车
func (s *cartService) GetCart(ctx context.Context, req *cart.GetCartReq) (*cart.GetCartResp, error) {
	// 这里应该实现获取购物车的逻辑
	// 例如，从Redis或MySQL中获取购物车商品信息
	key := fmt.Sprintf("cart:%d", req.UserId)

	// 从Redis获取购物车数据
	result, err := redis.RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("从Redis获取购物车失败: %w", err)
	}

	items := make([]*cart.CartItem, 0, len(result))
	for productIDStr, valueStr := range result {
		// 解析productID
		var productID int64
		fmt.Sscanf(productIDStr, "%d", &productID)

		// 解析存储的值
		var quantity int32
		var timestamp int64
		fmt.Sscanf(valueStr, "%d:%d:%d", &quantity, &productID, &timestamp)

		item := &cart.CartItem{
			ProductId: uint32(productID),
			Quantity:  quantity,
		}

		items = append(items, item)
	}

	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.UserId,
			Items:  items,
		},
	}, nil
}

// EmptyCart 清空购物车
func (s *cartService) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	// 这里应该实现清空购物车的逻辑
	key := fmt.Sprintf("cart:%d", req.UserId)

	// 从Redis删除购物车数据
	err := redis.RDB.Del(ctx, key).Err()
	if err != nil {
		return nil, fmt.Errorf("从Redis删除购物车失败: %w", err)
	}

	return &cart.EmptyCartResp{}, nil
}

// Helper function to convert database items to proto message
func convertToProtoCart(userID uint32, items []*model.CartItem) *cart.GetCartResp {
	protoItems := make([]*cart.CartItem, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, &cart.CartItem{
			ProductId: uint32(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}
	return &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: userID,
			Items:  protoItems,
		},
	}
}
