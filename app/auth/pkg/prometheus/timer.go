// TODO: 等待 Prometheus 完全兼容 Go 1.23.4
package prometheus

/*
import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Timer struct {
	observer prometheus.Observer
	start    float64
}

func NewTimer(o prometheus.Observer) *Timer {
	return &Timer{
		observer: o,
		start:    float64(time.Now().UnixNano()) / 1e9,
	}
}

func (t *Timer) ObserveDuration() {
	if t.observer == nil {
		return
	}
	t.observer.Observe(float64(time.Now().UnixNano())/1e9 - t.start)
}
*/
