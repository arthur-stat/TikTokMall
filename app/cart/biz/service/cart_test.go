package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockRepo "TikTokMall/app/cart/biz/dal/mysql/mock"
	"TikTokMall/app/cart/biz/model"
	"TikTokMall/app/cart/kitex_gen/cart"
)

func TestCartService_AddItem(t *testing.T) {
	repo := new(mockRepo.MockCartRepository)
	svc := NewCartServiceWithRepo(repo)

	tests := []struct {
		name    string
		req     *cart.AddItemReq
		mockFn  func()
		wantErr bool
	}{
		{
			name: "valid item",
			req: &cart.AddItemReq{
				UserId: 1,
				Item: &cart.CartItem{
					ProductId: 101,
					Quantity:  2,
				},
			},
			mockFn: func() {
				repo.On("AddItem", mock.Anything, uint32(1), mock.AnythingOfType("*model.CartItem")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "zero quantity",
			req: &cart.AddItemReq{
				UserId: 1,
				Item: &cart.CartItem{
					ProductId: 101,
					Quantity:  0,
				},
			},
			mockFn:  func() {},
			wantErr: true,
		},
		{
			name: "repository error",
			req: &cart.AddItemReq{
				UserId: 1,
				Item: &cart.CartItem{
					ProductId: 101,
					Quantity:  2,
				},
			},
			mockFn: func() {
				repo.On("AddItem", mock.Anything, uint32(1), mock.AnythingOfType("*model.CartItem")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			repo.ExpectedCalls = nil

			tt.mockFn()

			resp, err := svc.AddItem(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			repo.AssertExpectations(t)
		})
	}
}

func TestCartService_GetCart(t *testing.T) {
	repo := new(mockRepo.MockCartRepository)
	svc := NewCartServiceWithRepo(repo)

	tests := []struct {
		name     string
		req      *cart.GetCartReq
		mockFn   func()
		wantResp *cart.GetCartResp
		wantErr  bool
	}{
		{
			name: "with items",
			req: &cart.GetCartReq{
				UserId: 1,
			},
			mockFn: func() {
				repo.On("GetItems", mock.Anything, uint32(1)).Return([]*model.CartItem{
					{
						ProductID: 101,
						Quantity:  2,
						Selected:  true,
					},
				}, nil)
			},
			wantResp: &cart.GetCartResp{
				Cart: &cart.Cart{
					UserId: 1,
					Items: []*cart.CartItem{
						{
							ProductId: 101,
							Quantity:  2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty cart",
			req: &cart.GetCartReq{
				UserId: 1,
			},
			mockFn: func() {
				repo.On("GetItems", mock.Anything, uint32(1)).Return([]*model.CartItem{}, nil)
			},
			wantResp: &cart.GetCartResp{
				Cart: &cart.Cart{
					UserId: 1,
					Items:  []*cart.CartItem{},
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: &cart.GetCartReq{
				UserId: 1,
			},
			mockFn: func() {
				repo.On("GetItems", mock.Anything, uint32(1)).Return(nil, errors.New("db error"))
			},
			wantResp: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			repo.ExpectedCalls = nil

			tt.mockFn()

			resp, err := svc.GetCart(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResp, resp)
			repo.AssertExpectations(t)
		})
	}
}

func TestCartService_EmptyCart(t *testing.T) {
	repo := new(mockRepo.MockCartRepository)
	svc := NewCartServiceWithRepo(repo)

	tests := []struct {
		name    string
		req     *cart.EmptyCartReq
		mockFn  func()
		wantErr bool
	}{
		{
			name: "success",
			req: &cart.EmptyCartReq{
				UserId: 1,
			},
			mockFn: func() {
				repo.On("EmptyCart", mock.Anything, uint32(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: &cart.EmptyCartReq{
				UserId: 1,
			},
			mockFn: func() {
				repo.On("EmptyCart", mock.Anything, uint32(1)).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			repo.ExpectedCalls = nil

			tt.mockFn()

			resp, err := svc.EmptyCart(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			repo.AssertExpectations(t)
		})
	}
}
