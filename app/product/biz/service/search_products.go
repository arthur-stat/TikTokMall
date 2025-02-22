package service

import (
	product "TikTokMall/app/product/kitex_gen/product"
	"TikTokMall/app/product/model"
	"context"

	"github.com/pkg/errors"
)

type ProductSearcher interface {
    SearchProducts(string) ([]model.Product, error)
}

type SearchProductsService struct {
	ctx context.Context
    productQuery ProductSearcher
} // NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context, productQuery ProductSearcher) *SearchProductsService {
	return &SearchProductsService{ctx: ctx, productQuery: productQuery}
}

// Run create note info
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
    // 1. param validation
	query := req.GetQuery()
	if query == "" {
		return nil, errors.New("product query is empty.")
	}

    // 2. biz
	products, err := s.productQuery.SearchProducts(query)
    if err != nil {
        return nil, errors.WithMessagef(err, "query: %q", query)
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
    return &product.SearchProductsResp{
        Results: results,
    }, nil
}
