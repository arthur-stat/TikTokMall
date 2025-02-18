package main

import (
	"TikTokMall/app/cart/kitex_gen/cart/cartservice"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

// NewCartServiceImpl creates a new CartServiceImpl.
func NewCartServiceImpl() cartservice.CartService {
	return &CartServiceImpl{}
}
