
package lazarus

import (
	"math"
	"testing"
	"time"
)

func validMetrics() Metrics {
	return Metrics{
		CPUUsage:       70,
		MemoryUsage:    60,
		QueueDepth:     80,
		QueueLimit:     100,
		P95Latency:     900 * time.Millisecond,
		LatencyTarget:  time.Second,
		MaxConcurrency: 20,
	}
}

func TestEvaluateScore(t *testing.T) {
	result, err := Evaluate(validMetrics())
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(result.Score-76.5) > 1e-9 {
		t.Fatalf("score = %f, want 76.5", result.Score)
	}

	if result.Zone != ZoneCritical {
		t.Fatalf("zone = %q, want %q", result.Zone, ZoneCritical)
	}

	if result.SuggestedConcurrency != 10 {
		t.Fatalf(
			"concurrency = %d, want 10",
			result.SuggestedConcurrency,
		)
	}

	if result.RejectNew {
		t.Fatal("unexpected rejection without hard limit")
	}
}

func TestEvaluateHardLimits(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Metrics)
	}{
		{
			name: "CPU",
			change: func(m *Metrics) {
				m.CPUUsage = 95
			},
		},
		{
			name: "memory",
			change: func(m *Metrics) {
				m.MemoryUsage = 90
			},
		},
		{
			name: "full queue",
			change: func(m *Metrics) {
				m.QueueDepth = m.QueueLimit
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validMetrics()
			tt.change(&m)

			result, err := Evaluate(m)
			if err != nil {
				t.Fatal(err)
			}

			if !result.RejectNew {
				t.Fatal("expected new requests to be rejected")
			}

			if result.Zone != ZoneProtection {
				t.Fatalf("zone = %q, want protection", result.Zone)
			}
		})
	}
}

func TestEvaluateInvalidMetrics(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Metrics)
	}{
		{
			name: "CPU above 100",
			change: func(m *Metrics) {
				m.CPUUsage = 101
			},
		},
		{
			name: "negative memory",
			change: func(m *Metrics) {
				m.MemoryUsage = -1
			},
		},
		{
			name: "zero queue limit",
			change: func(m *Metrics) {
				m.QueueLimit = 0
			},
		},
		{
			name: "zero latency target",
			change: func(m *Metrics) {
				m.LatencyTarget = 0
			},
		},
		{
			name: "zero concurrency",
			change: func(m *Metrics) {
				m.MaxConcurrency = 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validMetrics()
			tt.change(&m)

			if _, err := Evaluate(m); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestEvaluateNormalZone(t *testing.T) {
	m := validMetrics()
	m.CPUUsage = 10
	m.MemoryUsage = 10
	m.QueueDepth = 0
	m.P95Latency = 100 * time.Millisecond

	result, err := Evaluate(m)
	if err != nil {
		t.Fatal(err)
	}

	if result.Zone != ZoneNormal {
		t.Fatalf("zone = %q, want normal", result.Zone)
	}

	if result.SuggestedConcurrency != m.MaxConcurrency {
		t.Fatalf(
			"concurrency = %d, want %d",
			result.SuggestedConcurrency,
			m.MaxConcurrency,
		)
	}
}
