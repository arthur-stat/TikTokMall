package model

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ICategoryQuery is a category query interface
type ICategoryQuery interface {
	GetProductsByCategoryName(string) ([]Product, error)
}

type Category struct {
	Base

	Name     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Products []Product `gorm:"many2many:product_category"`
}

func (c Category) TableName() string {
	return "category"
}

// CategoryQuery implements ICategoryQuery
type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewCategoryQuery(ctx context.Context, db *gorm.DB) ICategoryQuery {
	return &CategoryQuery{ctx: ctx, db: db}
}

func (c CategoryQuery) GetProductsByCategoryName(name string) (products []Product, err error) {
	var category Category
	err = c.db.WithContext(c.ctx).Preload("Products").Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Errorf("category does not exist")
		}
		return nil, errors.Wrap(err, "query db by category failed")
	}
	return category.Products, nil
}
