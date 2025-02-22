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

type MockProductSearcher struct {
	mock.Mock
}

func (m *MockProductSearcher) SearchProducts(q string) ([]model.Product, error) {
	args := m.Called(q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)

}

func TestSearchProducts_Run(t *testing.T) {
	tests := []struct {
		name      string
		req       *product.SearchProductsReq
		mockSetup func(*MockProductSearcher)
		wantResp  *product.SearchProductsResp
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			req: &product.SearchProductsReq{
				Query: "test",
			},
			mockSetup: func(m *MockProductSearcher) {
				price, _ := decimal.NewFromString("99.99")
				m.On("SearchProducts", "test").Return([]model.Product{
					{
						Base:        model.Base{ID: 1},
						Name:        "product 1",
						Description: "description 1",
						Picture:     "http://example.com/image1.jpg",
						Price:       price,
						Stock:       100,
						Categories:  []model.Category{{Name: "category1"}, {Name: "category2"}},
					},
				}, nil)
			},
			wantResp: &product.SearchProductsResp{
				Results: []*product.Product{
					{
						Id:          1,
						Name:        "product 1",
						Description: "description 1",
						Picture:     "http://example.com/image1.jpg",
						Price:       99.99,
						Stock:       100,
						Categories:  []string{"category1", "category2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error_empty_query",
			req: &product.SearchProductsReq{
				Query: "",
			},
			mockSetup: func(m *MockProductSearcher) {},
			wantResp:  nil,
			wantErr:   true,
			errMsg:    "product query is empty.",
		},
		{
			name: "error_no_products_found",
			req: &product.SearchProductsReq{
				Query: "none",
			},
			mockSetup: func(m *MockProductSearcher) {
				m.On("SearchProducts", "none").Return(nil, errors.New("no products found"))
			},
			wantResp: nil,
			wantErr:  true,
			errMsg:   "query: \"none\": no products found",
		},
		{
			name: "success_multiple_results",
			req: &product.SearchProductsReq{
				Query: "multiple",
			},
			mockSetup: func(m *MockProductSearcher) {
				price1, _ := decimal.NewFromString("99.99")
				price2, _ := decimal.NewFromString("199.99")
				m.On("SearchProducts", "multiple").Return([]model.Product{
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
						Base:        model.Base{ID: 2},
						Name:        "product2",
						Description: "description 2",
						Picture:     "http://example.com/image2.jpg",
						Price:       price2,
						Stock:       200,
						Categories:  []model.Category{{Name: "category2"}},
					},
				}, nil)
			},
			wantResp: &product.SearchProductsResp{
				Results: []*product.Product{
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
						Id:          2,
						Name:        "product2",
						Description: "description 2",
						Picture:     "http://example.com/image2.jpg",
						Price:       199.99,
						Stock:       200,
						Categories:  []string{"category2"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create mock
			mockSearcher := new(MockProductSearcher)
			tt.mockSetup(mockSearcher)

			// create service
			ctx := context.Background()
			s := NewSearchProductsService(ctx, mockSearcher)

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
			mockSearcher.AssertExpectations(t)
		})
	}
}
