package service

import (
	product "TikTokMall/app/product/kitex_gen/product"
	"TikTokMall/app/product/model"
	"context"

	"github.com/pkg/errors"
)

type ListProductsService struct {
	ctx context.Context
    categoryQuery model.ICategoryQuery
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context, categoryQuery model.ICategoryQuery) *ListProductsService {
	return &ListProductsService{ctx: ctx, categoryQuery: categoryQuery}
}

// Run create note info
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
    categoryName := req.GetCategoryName()
    // 1. params check
    if categoryName == "" {
        return nil, errors.New("category name is empty.")
    }

    // 2. query products
    products, err := s.categoryQuery.GetProductsByCategoryName(categoryName)
    if err != nil {
        return nil, errors.WithMessagef(err, "categoryName: %q", categoryName)
    }

    // 3. response
    var results []*product.Product
    for _, p := range products {
        var categories []string
        for _, c := range p.Categories {
            categories = append(categories, c.Name)
        }
        results = append(results, &product.Product{
            Id: p.ID,
            Name: p.Name,
            Description: p.Description,
            Picture: p.Picture,
            Price: float32(p.Price.InexactFloat64()),
            Stock: p.Stock,
            Categories: categories,
        })
    }
    return &product.ListProductsResp{
        Products: results,
    }, nil
}
