package factory

import (
	"TikTokMall/app/product/biz/dal/mysql"
	"TikTokMall/app/product/model"
	"context"
)

type ProductFactory struct{}

func GetProductFactory() *ProductFactory {
	return &ProductFactory{}
}

func (f *ProductFactory) NewProductQuery(ctx context.Context) model.IProductQuery {
	return model.NewProductQuery(ctx, mysql.DB)
}

func (f *ProductFactory) NewCategoryQuery(ctx context.Context) model.ICategoryQuery {
    return model.NewCategoryQuery(ctx, mysql.DB)
}
