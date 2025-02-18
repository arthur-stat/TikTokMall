package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	AuthTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_total",
			Help: "认证请求总数",
		},
		[]string{"type", "status"},
	)

	AuthDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_duration_seconds",
			Help:    "认证处理时间",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	TokenTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_total",
			Help: "令牌操作总数",
		},
		[]string{"operation", "status"},
	)
)
