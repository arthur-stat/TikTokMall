package service

import (
	"TikTokMall/app/product/kitex_gen/product"
	"TikTokMall/app/product/model"
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryQuery struct {
	mock.Mock
}

func (m *MockCategoryQuery) GetProductsByCategoryName(name string) ([]model.Product, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func TestListProducts_Run(t *testing.T) {
	tests := []struct {
		name      string
		req       *product.ListProductsReq
		mockSetup func(*MockCategoryQuery)
		wantResp  *product.ListProductsResp
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			req: &product.ListProductsReq{
				CategoryName: "category1",
			},
			mockSetup: func(m *MockCategoryQuery) {
				price1, _ := decimal.NewFromString("99.99")
				price2, _ := decimal.NewFromString("99999.99")
				m.On("GetProductsByCategoryName", "category1").Return([]model.Product{
					{
						Base:        model.Base{ID: 1},
						Name:        "product1",
						Description: "description 1",
						Picture:     "http://example.com/image1.jpg",
						Price:       price1,
						Stock:       100,
						Categories:  []model.Category{{Name: "category1"}, {Name: "category2"}},
					},
					{
						Base:        model.Base{ID: 333},
						Name:        "product333",
						Description: "description 333",
						Picture:     "http://example.com/image333.jpg",
						Price:       price2,
						Stock:       1000,
						Categories:  []model.Category{{Name: "category1"}, {Name: "category10"}},
					},
				}, nil)
			},
			wantResp: &product.ListProductsResp{
				Products: []*product.Product{
					{
						Id:          1,
						Name:        "product1",
						Description: "description 1",
						Picture:     "http://example.com/image1.jpg",
						Price:       99.99,
						Stock:       100,
						Categories:  []string{"category1", "category2"},
					},
					{
						Id:          333,
						Name:        "product333",
						Description: "description 333",
						Picture:     "http://example.com/image333.jpg",
						Price:       99999.99,
						Stock:       1000,
						Categories:  []string{"category1", "category10"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "error_empty_category_name",
			req:       &product.ListProductsReq{},
			mockSetup: func(m *MockCategoryQuery) {},
			wantResp:  nil,
			wantErr:   true,
			errMsg:    "category name is empty.",
		},
		{
			name: "error_category_not_found",
			req: &product.ListProductsReq{
				CategoryName: "nonexistent",
			},
			mockSetup: func(m *MockCategoryQuery) {
				m.On("GetProductsByCategoryName", "nonexistent").Return(nil, errors.New("category nonexistent does not exist."))
			},
			wantResp: nil,
			wantErr:  true,
			errMsg:   "categoryName: \"nonexistent\": category nonexistent does not exist.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQuery := new(MockCategoryQuery)
			tt.mockSetup(mockQuery)

			// server mockSetup
			ctx := context.Background()
			s := NewListProductsService(ctx, mockQuery)

			// run tests
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err, "should return error")
				assert.Equal(t, tt.errMsg, err.Error(), "error message should match")
				assert.Nil(t, resp, "response should be nil")
			} else {
				assert.Nil(t, err, "should not return error")
				assert.Equal(t, tt.wantResp, resp, "response should match")
				assert.NotNil(t, resp, "response should not be nil")
			}

			// verify mock
			mockQuery.AssertExpectations(t)
		})
	}
}
