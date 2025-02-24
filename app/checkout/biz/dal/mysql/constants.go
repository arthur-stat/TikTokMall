package mysql

// 订单状态常量
const (
	OrderStatusPending  = 1 // 待支付
	OrderStatusPaid     = 2 // 已支付
	OrderStatusCanceled = 3 // 已取消
	OrderStatusComplete = 4 // 已完成
)
