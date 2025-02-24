package service

import (
	"TikTokMall/app/payment/biz/dal/mysql"
	"TikTokMall/app/payment/biz/dal/redis"
	"TikTokMall/app/payment/biz/model"
	"TikTokMall/app/payment/conf"
	"TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"fmt"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/smartwalle/alipay/v3"
	"math"
	"strconv"
	"time"
)

type AlipayChargeService struct {
	ctx context.Context
}

func NewAlipayChargeService(ctx context.Context) *AlipayChargeService {
	return &AlipayChargeService{ctx: ctx}
}

func (s *AlipayChargeService) Run(req *payment.AlipayChargeReq) (*payment.AlipayChargeResp, error) {
	appID := conf.GetConf().Alipay.AppID               // appID
	privateKey := conf.GetConf().Alipay.PrivateKey     // 私钥
	aliPublicKey := conf.GetConf().Alipay.AliPublicKey // 支付宝的公钥

	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4004009, fmt.Sprintf("failed to create alipay client: %v", err))
	}

	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4004010, fmt.Sprintf("failed to load alipay public key: %v", err))
	}

	p := alipay.TradePagePay{}
	p.NotifyURL = conf.GetConf().Alipay.NotifyUrl                               // 支付宝回调
	p.ReturnURL = req.ReturnUrl                                                 // 支付后调转页面
	p.Subject = "tik_tok_mall"                                                  // 标题
	p.OutTradeNo = strconv.FormatInt(req.OrderId, 10)                           // 传递一个唯一单号
	p.TotalAmount = fmt.Sprintf("%.2f", math.Ceil(float64(req.Amount)*100)/100) // 金额
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"                                    // 网页支付
	p.TimeoutExpress = "10m"                                                    // 定时取消订单的超时时间

	// 检查 Redis 是否已经有该订单的 payUrl
	existingPayUrl, err := redis.GetPayUrlFromCache(s.ctx, req.OrderId)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005006, fmt.Sprintf("failed to get payUrl from cache: %v", err))
	}

	if existingPayUrl != "" {
		// 如果 Redis 中已经有该 payUrl，直接返回
		return &payment.AlipayChargeResp{PayUrl: existingPayUrl}, nil
	}

	// 如果 Redis 中没有 payUrl，则生成支付链接
	payURL, err := client.TradePagePay(p)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4005007, fmt.Sprintf("failed to create trade page pay: %v", err))
	}

	// 将生成的 payUrl 存储到 Redis 10分钟
	err = redis.SetPayUrlToCache(s.ctx, req.OrderId, payURL.String())
	if err != nil {
		// 如果 Redis 缓存失败，可以记录日志，但这里返回错误表示缓存失败
		return nil, kkerrors.NewBizStatusError(4005008, fmt.Sprintf("failed to set payUrl to cache: %v", err))
	}

	// 检查是否已存在相同的订单
	existingPayment, err := mysql.GetPaymentByOrderID(mysql.DB, s.ctx, req.OrderId)

	// 如果数据库查询返回 record not found，表示订单不存在
	if err != nil && err.Error() == "record not found" {
		// 创建支付记录
		transactionId := p.OutTradeNo
		thePayment := &model.Payments{
			UserID:        req.UserId,
			OrderID:       req.OrderId,
			TransactionID: transactionId,
			Status:        0, // 初始状态
			Amount:        req.Amount,
			PaymentMethod: "alipay",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// 插入支付记录到数据库
		err = mysql.CreatePayment(mysql.DB, s.ctx, thePayment)
		if err != nil {
			return nil, kkerrors.NewBizStatusError(4005010, fmt.Sprintf("failed to create payment record: %v", err))
		}
	} else if existingPayment != nil {
		// 如果存在支付记录，不做重复插入
		return &payment.AlipayChargeResp{PayUrl: payURL.String()}, nil
	}

	// 返回生成的支付页面链接
	return &payment.AlipayChargeResp{PayUrl: payURL.String()}, nil
}
