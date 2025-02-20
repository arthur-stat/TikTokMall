package service

import (
	product "TikTokMall/app/product/kitex_gen/product"
	"TikTokMall/app/product/model"
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductGetter implements productGetter interface, for testing
type MockProductGetter struct {
	mock.Mock
}

func (m *MockProductGetter) GetById(id uint32) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func TestGetProduct_Run(t *testing.T) {
	tests := []struct {
		name      string
		req       *product.GetProductReq
		mockSetup func(*MockProductGetter)
		wantResp  *product.GetProductResp
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			req: &product.GetProductReq{
				Id: 1,
			},
			mockSetup: func(m *MockProductGetter) {
				price, _ := decimal.NewFromString("99.99")
				m.On("GetById", uint32(1)).Return(model.Product{
					Base: model.Base{
						ID: 1,
					},
					Name:        "test product 1",
					Description: "description 1",
					Picture:     "http://example.com/image1.jpg",
					Price:       price,
					Stock:       100,
					Categories:  []model.Category{{Name: "category1"}, {Name: "category2"}},
				}, nil)
			},
			wantResp: &product.GetProductResp{
				Product: &product.Product{
					Id:          1,
					Name:        "test product 1",
					Description: "description 1",
					Picture:     "http://example.com/image1.jpg",
					Price:       99.99,
					Stock:       100,
					Categories:  []string{"category1", "category2"},
				},
			},
			wantErr: false,
		},
		{
			name: "error_empty_id",
			req: &product.GetProductReq{
				Id: 0,
			},
			mockSetup: func(m *MockProductGetter) {},
			wantResp:  nil,
			wantErr:   true,
			errMsg:    "product id is invalid.",
		},
		{
			name: "error_id_not_found",
			req: &product.GetProductReq{
				Id: 333,
			},
			mockSetup: func(m *MockProductGetter) {
				m.On("GetById", uint32(333)).Return(model.Product{}, errors.New("product id does not exist"))
			},
			wantResp: nil,
			wantErr:  true,
			errMsg:   "product id: 333: product id does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create mock
			mockGetter := new(MockProductGetter)
			tt.mockSetup(mockGetter)

			// create service
			ctx := context.Background()
			s := NewGetProductService(ctx, mockGetter)

			// run tests
			resp, err := s.Run(tt.req)
			if tt.wantErr {
				assert.Error(t, err, "should return error")
				assert.Equal(t, tt.errMsg, err.Error(), "error message should match")
				assert.Nil(t, resp, "response should be nil")
			} else {
				assert.NoError(t, err, "should not return error")
				assert.Equal(t, tt.wantResp, resp, "response should match")
				assert.NotNil(t, resp, "response should not be nil")
			}

			// verify mock
			mockGetter.AssertExpectations(t)
		})
	}
}
