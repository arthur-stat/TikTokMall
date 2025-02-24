package service

import (
	"TikTokMall/app/payment/biz/dal/mysql"
	"TikTokMall/app/payment/biz/dal/redis"
	"TikTokMall/app/payment/biz/model"
	"TikTokMall/app/payment/kitex_gen/payment"
	"context"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type ChargeService struct {
	ctx context.Context
}

// NewChargeService 构造 ChargeService 实例
func NewChargeService(ctx context.Context) *ChargeService {
	return &ChargeService{ctx: ctx}
}

func (s *ChargeService) Run(req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	// 构造信用卡信息对象并验证
	card := creditcard.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    strconv.Itoa(int(req.CreditCard.CreditCardCvv)),
		Month:  strconv.Itoa(int(req.CreditCard.CreditCardExpirationMonth)),
		Year:   strconv.Itoa(int(req.CreditCard.CreditCardExpirationYear)),
	}
	if err = card.Validate(true); err != nil {
		return nil, kkerrors.NewBizStatusError(4004001, err.Error())
	}

	// 生成交易 ID（暂未使用支付接口）
	transactionId, err := uuid.NewRandom()
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005001, err.Error())
	}

	// 尝试从 Redis 缓存中获取支付记录
	paymentFromCache, err := redis.GetPaymentCache(s.ctx, req.OrderId)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005003, err.Error())
	}
	if paymentFromCache != nil {
		// 若缓存中存在则直接返回缓存中的交易 ID
		return &payment.ChargeResp{TransactionId: paymentFromCache.TransactionID}, nil
	}

	// 检查是否已存在相同的订单
	existingPayment, err := mysql.GetPaymentByOrderID(mysql.DB, s.ctx, req.OrderId)
	if err == nil && existingPayment != nil {
		return nil, kkerrors.NewBizStatusError(4005002, "订单已存在，不能重复支付")
	}

	// 创建支付记录
	thePayment := &model.Payments{
		UserID:        req.UserId,
		OrderID:       req.OrderId,
		TransactionID: transactionId.String(),
		Status:        1,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 插入支付记录到数据库
	err = mysql.CreatePayment(mysql.DB, s.ctx, thePayment)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005002, err.Error())
	}

	// 写入成功后，将支付记录缓存到 Redis
	err = redis.SetPaymentCache(s.ctx, thePayment)
	if err != nil {
		// 如果 Redis 缓存失败，可以记录日志，但这里返回错误表示缓存失败
		return nil, kkerrors.NewBizStatusError(4005004, err.Error())
	}

	// 返回生成的交易 ID
	return &payment.ChargeResp{TransactionId: transactionId.String()}, nil
}
