package model

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ProductModelTestSuite struct {
	suite.Suite
	db           *gorm.DB
	ctx          context.Context
	query        IProductQuery
	testProducts []Product
}

func (s *ProductModelTestSuite) prepareTestRecords() {
	s.testProducts = []Product{
		{
			Base:        Base{ID: 1},
			Name:        "iPhone 15 Pro",
			Description: "Latest Apple iPhone",
			Picture:     "http://example.com/iphone15.jpg",
			Price:       decimal.NewFromFloat(6999.00),
			Stock:       1000,
		},
		{
			Base:        Base{ID: 2},
			Name:        "Xiaomi 14",
			Description: "Latest Xiaomi phone Flagship Series",
			Picture:     "http://example.com/mi14.jpg",
			Price:       decimal.NewFromFloat(4999.00),
			Stock:       2000,
		},
		{
			Base:        Base{ID: 3},
			Name:        "Foxconn Tablet",
			Description: "Standard Tablet",
			Picture:     "http://example.com/huawei-pad.jpg",
			Price:       decimal.NewFromFloat(3999.00),
			Stock:       1500,
		},
	}
}

func (s *ProductModelTestSuite) SetupSuite() {
	// setup test requirements
	dsn := "root:root@tcp(127.0.0.1:3306)/productdb_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	s.Require().NoError(err)
	s.db = db

	s.db.AutoMigrate(&Product{}, &Category{})
	s.prepareTestRecords()

	s.ctx = context.Background()
	s.query = NewProductQuery(s.ctx, s.db)
}

func (s *ProductModelTestSuite) SetupTest() {
	s.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	s.db.Exec("TRUNCATE TABLE product")
	s.db.Exec("TRUNCATE TABLE category")
	s.db.Exec("TRUNCATE TABLE product_category")
	s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func (s *ProductModelTestSuite) TestGetById() {
	// prepare records
	product := s.testProducts[0]
	s.db.Create(&product)

	// exec test
	result, err := s.query.GetById(1)

	// veryfy return
	s.Require().NoError(err)
	s.Equal(product.ID, result.ID)
	s.Equal(product.Name, result.Name)
	s.Equal(product.Description, result.Description)
	s.Equal(product.Picture, result.Picture)
	s.Equal(product.Price.String(), result.Price.String())
	s.Equal(product.Stock, result.Stock)
}

func (s *ProductModelTestSuite) TestGetById_NotFound() {
	_, err := s.query.GetById(999)
	s.Require().Error(err)
	s.Contains(err.Error(), "product id does not exist")
}

func (s *ProductModelTestSuite) TestSearchProducts() {
	s.db.Create(&s.testProducts)

	// exec test
	s.Run("sucess_search_by_name", func() {
		results, err := s.query.SearchProducts("iphone")
		s.Require().NoError(err)
		s.Len(results, 1)

		s.Equal("iPhone 15 Pro", results[0].Name)
	})

	s.Run("sucess_search_by_description", func() {
		results, err := s.query.SearchProducts("Latest")
		s.Require().NoError(err)
		s.Len(results, 2)

		names := make([]string, len(results))
		for i, result := range results {
			names[i] = result.Name
		}
		s.Contains(names, "iPhone 15 Pro")
		s.Contains(names, "Xiaomi 14")
	})

	s.Run("sucess_search_in_both_fields", func() {
		results, err := s.query.SearchProducts("phone")
		s.Require().NoError(err)
		s.Len(results, 2)

		names := make([]string, len(results))
		for i, result := range results {
			names[i] = result.Name
		}
		s.Contains(names, "iPhone 15 Pro")
		s.Contains(names, "Xiaomi 14")
	})
}

func (s *ProductModelTestSuite) TestSearchProducts_NoResults() {
	s.db.Create(&s.testProducts)

	s.Run("error_no_match_found", func() {
		results, err := s.query.SearchProducts("nonexistent product")
		s.Require().Error(err)
		s.Contains(err.Error(), "no products found")
		s.Nil(results)
	})

	s.Run("error_search_with_spaces", func() {
		results, err := s.query.SearchProducts("   ")
		s.Require().Error(err)
		s.Contains(err.Error(), "no products found")
		s.Nil(results)
	})
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductModelTestSuite))
}
