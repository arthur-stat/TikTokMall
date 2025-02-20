package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CheckoutTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "checkout_total",
			Help: "结账请求总数",
		},
		[]string{"status"},
	)

	CheckoutDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "checkout_duration_seconds",
			Help:    "结账处理时间",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	PaymentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_total",
			Help: "支付处理总数",
		},
		[]string{"status"},
	)
)
