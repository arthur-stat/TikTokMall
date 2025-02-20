package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CartTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_total",
			Help: "购物车操作总数",
		},
		[]string{"operation", "status"},
	)

	CartDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cart_duration_seconds",
			Help:    "购物车操作处理时间",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	CartItemTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cart_item_total",
			Help: "购物车商品数量",
		},
		[]string{"user_id"},
	)
)
