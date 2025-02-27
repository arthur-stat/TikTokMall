package service

import (
	"reflect"
	"testing"

	"TikTokMall/app/checkout/kitex_gen/payment"
)

// TestInspectTypes 帮助我们了解生成的类型
func TestInspectTypes(t *testing.T) {
	// 创建信用卡信息
	card := &payment.CreditCardInfo{}

	// 使用反射获取类型信息
	typ := reflect.TypeOf(card).Elem()

	// 输出类型信息
	t.Logf("CreditCardInfo 类型: %v", typ)

	// 输出所有字段信息
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		t.Logf("字段 #%d: 名称=%s, 类型=%v, 标签=%v",
			i, field.Name, field.Type, field.Tag)
	}
}
