package service

import (
	"context"
	"testing"

	"TikTokMall/app/checkout/kitex_gen/checkout"
	"TikTokMall/app/checkout/kitex_gen/payment"
)

func BenchmarkCheckoutService_Run(b *testing.B) {
	svc := NewCheckoutService()
	req := &checkout.CheckoutReq{
		UserId:    1,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@example.com",
		Address: &checkout.Address{
			StreetAddress: "123 Main St",
			City:          "City",
			State:         "State",
			Country:       "Country",
			ZipCode:       "12345",
		},
		CreditCard: &payment.CreditCardInfo{
			CardNumber:  "4111111111111111",
			Cvv:         "123",
			ExpiryYear:  2025,
			ExpiryMonth: 12,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Run(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
