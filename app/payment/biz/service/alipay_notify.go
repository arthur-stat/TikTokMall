package service

import (
	"TikTokMall/app/payment/biz/dal/mysql"
	"TikTokMall/app/payment/biz/model"
	"TikTokMall/app/payment/conf"
	payment "TikTokMall/app/payment/kitex_gen/payment"
	"context"
	"fmt"
	kkerrors "github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/smartwalle/alipay/v3"
	"time"
)

// AlipayNotifyService 处理支付宝的通知
type AlipayNotifyService struct {
	ctx context.Context
}

// NewAlipayNotifyService 创建一个新的 AlipayNotifyService 实例
func NewAlipayNotifyService(ctx context.Context) *AlipayNotifyService {
	return &AlipayNotifyService{ctx: ctx}
}

// Run 处理支付宝通知并更新支付状态
func (s *AlipayNotifyService) Run(req *payment.AlipayNotifyReq) (*payment.AlipayNotifyResp, error) {
	// 获取支付宝的配置
	appID := conf.GetConf().Alipay.AppID               // appID
	privateKey := conf.GetConf().Alipay.PrivateKey     // 私钥
	aliPublicKey := conf.GetConf().Alipay.AliPublicKey // 支付宝的公钥

	// 创建支付宝客户端
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4004009, fmt.Sprintf("创建支付宝客户端失败: %v", err))
	}

	// 加载支付宝公钥
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		return nil, kkerrors.NewBizStatusError(4004010, fmt.Sprintf("加载支付宝公钥失败: %v", err))
	}

	// 如果通知验证成功，检查交易状态
	if req.TradeStatus != "TRADE_SUCCESS" && req.TradeStatus != "TRADE_FINISHED" {
		klog.Error("订单 %s 未成功，交易状态: %s", req.OutTradeNo, req.TradeStatus)
		return nil, kkerrors.NewBizStatusError(4005006, fmt.Sprintf("订单未成功，交易状态: %s", req.TradeStatus))
	}

	// 查询支付状态
	existingPayment, err := mysql.GetPaymentByOrderID(mysql.DB, s.ctx, req.OutTradeNo)
	if err != nil {
		klog.Error("查询订单 %s 的支付记录时发生错误: %s", req.OutTradeNo, err.Error())
		return nil, kkerrors.NewBizStatusError(4005007, fmt.Sprintf("查询订单支付记录失败: %s", err.Error()))
	}

	// 果支付记录不存在，则创建支付记录
	if existingPayment == nil {
		thePayment := &model.Payments{
			OrderID:       req.OutTradeNo,
			TransactionID: req.TradeNo,
			Status:        1, // 支付成功
			Amount:        req.TotalAmount,
			PaymentMethod: "alipay",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err = mysql.CreatePayment(mysql.DB, s.ctx, thePayment)
		if err != nil {
			klog.Error("创建支付记录时发生错误，订单 %s: %s", req.OutTradeNo, err.Error())
			return nil, kkerrors.NewBizStatusError(4005011, fmt.Sprintf("创建支付记录失败: %s", err.Error()))
		}
		klog.Info("成功创建支付记录，订单 %s", req.OutTradeNo)
	} else {
		// 如果支付记录已存在，且支付状态不为已支付，则更新支付状态
		if existingPayment.Status != 1 {
			err = mysql.StatusToChargeStatus(mysql.DB, s.ctx, req.OutTradeNo)
			if err != nil {
				klog.Error("更新订单 %s 支付状态时发生错误: %s", req.OutTradeNo, err.Error())
				return nil, kkerrors.NewBizStatusError(4005012, fmt.Sprintf("更新支付状态失败: %s", err.Error()))
			}
		}
		klog.Info("订单 %s 已处理", req.OutTradeNo)
	}

	// 向支付宝发送确认收到通知的响应
	client.ACKNotification(nil)

	// 返回成功的响应
	return &payment.AlipayNotifyResp{Status: "success"}, nil
}
