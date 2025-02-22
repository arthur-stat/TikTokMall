package model

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// IProductQuery defines the interface for product query
type IProductQuery interface {
	GetById(uint32) (Product, error)
	SearchProducts(string) ([]Product, error)
}

type Product struct {
	Base

	Name        string          `gorm:"type:varchar(255);not null;index"`
	Description string          `gorm:"type:text"`
	Picture     string          `gorm:"type:varchar(1024)"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Stock       uint32          `gorm:"not null;default:0"`

	Categories []Category `gorm:"many2many:product_category"`
}

func (p Product) TableName() string {
	return "product"
}

// ProductQuery is an implementation of IProductQuery
type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewProductQuery(ctx context.Context, db *gorm.DB) IProductQuery {
	return &ProductQuery{ctx: ctx, db: db}
}

func (p ProductQuery) GetById(id uint32) (product Product, err error) {
	err = p.db.WithContext(p.ctx).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return product, errors.New("product id does not exist")
		}
		return product, errors.Wrap(err, "query db by id failed")
	}
	return
}

func (p ProductQuery) SearchProducts(q string) (products []Product, err error) {
	err = p.db.WithContext(p.ctx).Where("name LIKE ? OR description LIKE ?", "%"+q+"%", "%"+q+"%").Find(&products).Error
	if err != nil {
		return nil, errors.Wrap(err, "search products failed")
	}

	if len(products) == 0 {
		return nil, errors.New("no products found")
	}
	return
}
