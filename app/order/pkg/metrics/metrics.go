package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrderTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_total",
			Help: "订单操作总数",
		},
		[]string{"operation", "status"},
	)

	OrderDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_duration_seconds",
			Help:    "订单操作处理时间",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	OrderAmount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_amount",
			Help:    "订单金额分布",
			Buckets: []float64{10, 50, 100, 500, 1000, 5000},
		},
		[]string{"currency"},
	)
)
