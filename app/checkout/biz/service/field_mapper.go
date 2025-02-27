package service

import (
	"reflect"

	"TikTokMall/app/checkout/kitex_gen/payment"
)

// 获取 CreditCardInfo 结构的字段信息
func inspectCreditCardInfo() map[string]reflect.Kind {
	// 用反射来检查实际字段
	card := &payment.CreditCardInfo{}
	cardType := reflect.TypeOf(card).Elem()

	fields := make(map[string]reflect.Kind)
	for i := 0; i < cardType.NumField(); i++ {
		field := cardType.Field(i)
		fields[field.Name] = field.Type.Kind()
	}

	return fields
}

// 将测试使用的通用字段名映射到实际字段名
// 这个函数可以在测试中调用，获取实际的字段名
func getCreditCardFieldNames() (cardNumber, cvv, expiryYear, expiryMonth string) {
	// 默认字段名
	cardNumber = "CardNumber"
	cvv = "Cvv"
	expiryYear = "ExpiryYear"
	expiryMonth = "ExpiryMonth"

	// 检查实际字段，进行可能的映射
	fields := inspectCreditCardInfo()

	// 这里可以添加字段映射逻辑
	// 例如，如果实际字段是 Number 而不是 CardNumber:
	if _, ok := fields["Number"]; ok {
		cardNumber = "Number"
	}

	return
}
