package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Timer 计时器
type Timer struct {
	observer prometheus.Observer
	start    float64
}

// NewTimer 创建新的计时器
func NewTimer(o prometheus.Observer) *Timer {
	return &Timer{
		observer: o,
		start:    float64(time.Now().UnixNano()) / 1e9,
	}
}

// ObserveDuration 观察持续时间
func (t *Timer) ObserveDuration() {
	if t.observer == nil {
		return
	}
	t.observer.Observe(float64(time.Now().UnixNano())/1e9 - t.start)
}
