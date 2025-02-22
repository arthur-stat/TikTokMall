package main

import (
	"TikTokMall/app/product/biz/service"
	"TikTokMall/app/product/factory"
	product "TikTokMall/app/product/kitex_gen/product"
	"context"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct {
	factory *factory.ProductFactory
}

// NewProductCatalogServiceImpl creates a new ProductCatalogServiceImpl instance
func NewProductCatalogServiceImpl() *ProductCatalogServiceImpl {
	return &ProductCatalogServiceImpl{
		factory: factory.GetProductFactory(),
	}
}

// ListProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
    categoryQuery := s.factory.NewCategoryQuery(ctx)
	resp, err = service.NewListProductsService(ctx, categoryQuery).Run(req)
	return resp, err
}

// GetProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	productQuery := s.factory.NewProductQuery(ctx)
	resp, err = service.NewGetProductService(ctx, productQuery).Run(req)
	return resp, err
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	productQuery := s.factory.NewProductQuery(ctx)
	resp, err = service.NewSearchProductsService(ctx, productQuery).Run(req)
	return resp, err
}
