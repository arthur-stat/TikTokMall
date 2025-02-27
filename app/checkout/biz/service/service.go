package service

import (
	"fmt"
)

// 错误定义
var (
	ErrInvalidInput      = fmt.Errorf("无效的输入参数")
	ErrPaymentFailed     = fmt.Errorf("支付处理失败")
	ErrAddressInvalid    = fmt.Errorf("地址信息无效")
	ErrCreditCardInvalid = fmt.Errorf("信用卡信息无效")
	ErrOrderCreateFailed = fmt.Errorf("订单创建失败")
)
