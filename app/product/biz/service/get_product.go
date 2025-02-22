package service

import (
	product "TikTokMall/app/product/kitex_gen/product"
	"TikTokMall/app/product/model"
	"context"

	"github.com/pkg/errors"
)

type productGetter interface {
	GetById(uint32) (model.Product, error)
}

type GetProductService struct {
	ctx          context.Context
	productQuery productGetter
} // NewGetProductService new GetProductService
func NewGetProductService(ctx context.Context, productQuery productGetter) *GetProductService {
	return &GetProductService{ctx: ctx, productQuery: productQuery}
}

// Run create note info
func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	// 1. params validation
	if req.GetId() == 0 {
		return nil, errors.New("product id is invalid.")
	}

	// 2. biz
	p, err := s.productQuery.GetById(req.GetId())
	if err != nil {
		return nil, errors.WithMessagef(err, "product id: %d", req.GetId())
	}

	// 3. respoonse
	var categories []string
	for _, c := range p.Categories {
		categories = append(categories, c.Name)
	}
	return &product.GetProductResp{
		Product: &product.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       float32(p.Price.InexactFloat64()),
			Stock:       p.Stock,
			Categories:  categories,
		},
	}, nil
}
