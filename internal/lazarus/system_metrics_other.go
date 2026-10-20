//go:build !windows

package lazarus

import "errors"

var ErrUnsupportedSystemMetrics = errors.New("system metrics collector is only implemented for Windows")

type SystemMetricsCollector struct{}

func NewSystemMetricsCollector() *SystemMetricsCollector {
	return &SystemMetricsCollector{}
}

func (s *SystemMetricsCollector) Collect() (float64, float64, error) {
	return 0, 0, ErrUnsupportedSystemMetrics
}
