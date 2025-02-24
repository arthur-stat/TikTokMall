package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"TikTokMall/app/payment/biz/dal/mysql"
	"TikTokMall/app/payment/biz/dal/redis"
	payment "TikTokMall/app/payment/kitex_gen/payment"
)

// TestMain 在所有测试运行前初始化 MySQL 和 Redis 连接
func TestMain(m *testing.M) {
	// 设置 MySQL 相关环境变量
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_PORT", "3307")
	os.Setenv("MYSQL_USER", "gorm")
	os.Setenv("MYSQL_PASSWORD", "gorm")
	os.Setenv("MYSQL_DATABASE", "gorm")

	// 设置 Redis 相关环境变量
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_PASSWORD", "")

	// 设置环境变量为 test
	_ = os.Setenv("GO_ENV", "test")
	// 初始化配置文件路径
	confPath := filepath.Join("conf", "test", "conf.yaml")
	_ = os.Setenv("CONF_PATH", confPath)

	// 初始化 Redis
	if err := redis.Init(); err != nil {
		fmt.Printf("redis.Init error: %v\n", err)
		os.Exit(1)
	}
	// 初始化 MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("mysql.Init error: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	// 测试后清理数据
	mysql.DB.Exec("DELETE FROM payments")
	redis.Client.FlushDB(context.Background())

	os.Exit(code)
}

// validChargeReq 返回一个合法的 ChargeReq 对象，用于成功流程测试
func validChargeReq() *payment.ChargeReq {
	return &payment.ChargeReq{
		Amount:        100.0,
		PaymentMethod: "credit_card",
		OrderId:       0, // 测试时会设置为唯一值
		UserId:        1,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          "4111111111111111",
			CreditCardCvv:             123,
			CreditCardExpirationMonth: 12,
			CreditCardExpirationYear:  2030,
		},
	}
}

// invalidCardChargeReq 返回一个信用卡校验失败的 ChargeReq
func invalidCardChargeReq() *payment.ChargeReq {
	return &payment.ChargeReq{
		Amount:        50.0,
		PaymentMethod: "credit_card",
		OrderId:       0,
		UserId:        2,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          "1234567890123456",
			CreditCardCvv:             123,
			CreditCardExpirationMonth: 12,
			CreditCardExpirationYear:  2030,
		},
	}
}

// TestChargeService_Run_Success 测试合法请求处理流程
func TestChargeService_Run_Success(t *testing.T) {
	ctx := context.Background()
	svc := NewChargeService(ctx)

	req := validChargeReq()
	req.OrderId = time.Now().UnixNano()

	resp, err := svc.Run(req)
	assert.NoError(t, err, "合法请求不应返回错误")
	assert.NotNil(t, resp, "返回结果不应为 nil")
	assert.NotEmpty(t, resp.TransactionId, "TransactionId 不应为空")
}

// TestChargeService_Run_InvalidCard 测试信用卡校验失败的情况
func TestChargeService_Run_InvalidCard(t *testing.T) {
	ctx := context.Background()
	svc := NewChargeService(ctx)

	req := invalidCardChargeReq()
	req.OrderId = time.Now().UnixNano()

	resp, err := svc.Run(req)
	assert.Error(t, err, "信用卡校验失败应返回错误")
	assert.Nil(t, resp, "错误情况下返回结果应为 nil")
	assert.Contains(t, err.Error(), "4004001", "错误码应为 4004001")
}
