package model

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CategoryModelTestSuite struct {
	suite.Suite
	db             *gorm.DB
	ctx            context.Context
	query          ICategoryQuery
	testCategories []Category
	testProducts   []Product
}

func (s *CategoryModelTestSuite) prepareTestRecords() {
	s.testCategories = []Category{
		{
			Base: Base{ID: 1},
			Name: "Electronics",
		},
		{
			Base: Base{ID: 2},
			Name: "Phones",
		},
	}
	s.testProducts = []Product{
		{
			Base:        Base{ID: 1},
			Name:        "iPhone15",
			Description: "Apple",
			Picture:     "http://example.com/iphone15.jpg",
			Price:       decimal.NewFromFloat(6999.00),
			Stock:       1000,
			Categories:  s.testCategories,
		},
		{
			Base:        Base{ID: 2},
			Name:        "Xiaomi14",
			Description: "Xiaomi",
			Picture:     "http://example.com/mi14.jpg",
			Price:       decimal.NewFromFloat(4999.00),
			Stock:       2000,
			Categories:  []Category{s.testCategories[1]},
		},
	}
}

func (s *CategoryModelTestSuite) SetupSuite() {
	dsn := "root:root@tcp(127.0.0.1:3306)/productdb_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	s.Require().NoError(err)
	s.db = db

	s.db.AutoMigrate(&Category{}, &Product{})
	s.prepareTestRecords()

	s.ctx = context.Background()
	s.query = NewCategoryQuery(s.ctx, s.db)
}

func (s *CategoryModelTestSuite) SetupTest() {
	s.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	s.db.Exec("TRUNCATE TABLE category")
	s.db.Exec("TRUNCATE TABLE product")
	s.db.Exec("TRUNCATE TABLE product_category")
	s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func (s *CategoryModelTestSuite) TestGetProductsByCategoryName() {
	s.db.Create(&s.testCategories)
	s.db.Create(&s.testProducts)

	s.Run("sucess_get_products_by_category_name", func() {
		products, err := s.query.GetProductsByCategoryName("Electronics")
		s.Require().NoError(err)
		s.Len(products, 1)

		s.Contains(products[0].Name, "iPhone15")
	})
}

func (s *CategoryModelTestSuite) TestGetProductsByCategoryName_NotFound() {
	s.db.Create(&s.testCategories)
	s.db.Create(&s.testProducts)

	s.Run("error_category_not_found", func() {
		products, err := s.query.GetProductsByCategoryName("NonexistentCategory")
		s.Require().Error(err)
		s.Contains(err.Error(), "category does not exist")
		s.Nil(products)
	})
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategoryModelTestSuite))
}
