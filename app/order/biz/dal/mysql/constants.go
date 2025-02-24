package mysql

// 订单状态常量
const (
	OrderStatusPending  = 1 // 待支付
	OrderStatusPaid     = 2 // 已支付
	OrderStatusCanceled = 3 // 已取消
	OrderStatusComplete = 4 // 已完成
)

// 支付方式常量
const (
	PaymentMethodCreditCard = 1 // 信用卡
	PaymentMethodAlipay     = 2 // 支付宝
	PaymentMethodWechat     = 3 // 微信支付
)
