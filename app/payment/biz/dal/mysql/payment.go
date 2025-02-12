package mysql

import (
	"TikTokMall/app/payment/biz/model"
	"context"
	"fmt"
	"gorm.io/gorm"
)

// CreatePayment 插入一条支付记录到数据库
func CreatePayment(db *gorm.DB, ctx context.Context, payment *model.Payments) error {
	err := db.WithContext(ctx).Model(&model.Payments{}).Create(payment).Error
	if err != nil {
		if err.Error() == "Error 1062 (23000): Duplicate entry" {
			return fmt.Errorf("payment already exists for order ID %d", payment.OrderID)
		}
		return err
	}
	return nil
}

// GetPaymentByOrderID 根据订单 ID 获取支付记录
func GetPaymentByOrderID(db *gorm.DB, ctx context.Context, orderID int64) (*model.Payments, error) {
	var payment model.Payments
	err := db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
