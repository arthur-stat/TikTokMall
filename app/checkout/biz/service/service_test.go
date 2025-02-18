package service

import (
	"context"
	"testing"

	"TikTokMall/app/checkout/kitex_gen/checkout"
	"TikTokMall/app/checkout/kitex_gen/payment"
)

func TestCheckoutService_Run(t *testing.T) {
	svc := NewCheckoutService()

	validAddress := &checkout.Address{
		StreetAddress: "123 Main St",
		City:          "City",
		State:         "State",
		Country:       "Country",
		ZipCode:       "12345",
	}

	validCreditCard := &payment.CreditCardInfo{
		CardNumber:  "4111111111111111",
		Cvv:         "123",
		ExpiryYear:  2025,
		ExpiryMonth: 12,
	}

	tests := []struct {
		name    string
		req     *checkout.CheckoutReq
		wantErr bool
	}{
		{
			name: "valid request",
			req: &checkout.CheckoutReq{
				UserId:     1,
				Firstname:  "John",
				Lastname:   "Doe",
				Email:      "john@example.com",
				Address:    validAddress,
				CreditCard: validCreditCard,
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing user_id",
			req: &checkout.CheckoutReq{
				Firstname:  "John",
				Lastname:   "Doe",
				Email:      "john@example.com",
				Address:    validAddress,
				CreditCard: validCreditCard,
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing address",
			req: &checkout.CheckoutReq{
				UserId:     1,
				Firstname:  "John",
				Lastname:   "Doe",
				Email:      "john@example.com",
				CreditCard: validCreditCard,
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing credit card",
			req: &checkout.CheckoutReq{
				UserId:    1,
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john@example.com",
				Address:   validAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing firstname",
			req: &checkout.CheckoutReq{
				UserId:     1,
				Lastname:   "Doe",
				Email:      "john@example.com",
				Address:    validAddress,
				CreditCard: validCreditCard,
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing email",
			req: &checkout.CheckoutReq{
				UserId:     1,
				Firstname:  "John",
				Lastname:   "Doe",
				Address:    validAddress,
				CreditCard: validCreditCard,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Run(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Run() returned nil response for valid request")
				}
				if resp.OrderId == "" {
					t.Error("Run() returned empty order_id")
				}
				if resp.TransactionId == "" {
					t.Error("Run() returned empty transaction_id")
				}
			}
		})
	}
}
