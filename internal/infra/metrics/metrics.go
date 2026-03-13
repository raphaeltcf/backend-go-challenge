package metrics

import (
	"fmt"
	"sync/atomic"
)

type Metrics struct {
	TotalProcessedOrders int64
	TotalFailedOrders    int64
	TotalInvalidOrders   int64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncProcessedOrders() {
	atomic.AddInt64(&m.TotalProcessedOrders, 1)
}
func (m *Metrics) IncFailedOrders() {
	atomic.AddInt64(&m.TotalFailedOrders, 1)
}
func (m *Metrics) IncInvalidOrders() {
	atomic.AddInt64(&m.TotalInvalidOrders, 1)
}

func (m *Metrics) FailureRate() string {

	total := m.TotalProcessedOrders + m.TotalFailedOrders + m.TotalInvalidOrders
	if total == 0 {
		return "0.00%"
	}
	failures := m.TotalFailedOrders + m.TotalInvalidOrders
	rate := float64(failures) / float64(total) * 100
	return fmt.Sprintf("%.2f%%", rate)
}
