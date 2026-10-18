package lazarus

import (
	"fmt"
	"math"
)

const (
	CPUHardLimit    = 95.0
	MemoryHardLimit = 90.0

	WarningThreshold  = 50.0
	CriticalThreshold = 70.0
	ProtectThreshold  = 85.0
)

func Evaluate(m Metrics) (Result, error) {
	if !validPercentage(m.CPUUsage) {
		return Result{}, fmt.Errorf("CPU usage must be between 0 and 100")
	}
	if !validPercentage(m.MemoryUsage) {
		return Result{}, fmt.Errorf("memory usage must be between 0 and 100")
	}
	if m.QueueDepth < 0 || m.QueueLimit <= 0 {
		return Result{}, fmt.Errorf("invalid queue depth or queue limit")
	}
	if m.P95Latency < 0 || m.LatencyTarget <= 0 {
		return Result{}, fmt.Errorf("invalid latency or latency target")
	}
	if m.MaxConcurrency <= 0 {
		return Result{}, fmt.Errorf("max concurrency must be greater than zero")
	}

	cpu := m.CPUUsage / 100
	memory := m.MemoryUsage / 100
	queue := math.Min(float64(m.QueueDepth)/float64(m.QueueLimit), 1)
	latency := math.Min(float64(m.P95Latency)/float64(m.LatencyTarget), 1)

	score := 100 * (0.25*cpu + 0.20*memory + 0.25*queue + 0.30*latency)

	zone := zoneForScore(score)
	result := Result{
		Score:                score,
		Zone:                 zone,
		SuggestedConcurrency: concurrencyForZone(m.MaxConcurrency, zone),
	}

	if m.CPUUsage >= CPUHardLimit {
		result.RejectNew = true
		result.Reasons = append(result.Reasons, "CPU reached the hard limit")
	}
	if m.MemoryUsage >= MemoryHardLimit {
		result.RejectNew = true
		result.Reasons = append(result.Reasons, "memory reached the hard limit")
	}
	if m.QueueDepth >= m.QueueLimit {
		result.RejectNew = true
		result.Reasons = append(result.Reasons, "request queue is full")
	}

	if result.RejectNew {
		result.Zone = ZoneProtection
		result.SuggestedConcurrency = concurrencyForZone(m.MaxConcurrency, ZoneProtection)
	}

	return result, nil
}

func validPercentage(value float64) bool {
	return !math.IsNaN(value) &&
		!math.IsInf(value, 0) &&
		value >= 0 &&
		value <= 100
}

func zoneForScore(score float64) Zone {
	switch {
	case score >= ProtectThreshold:
		return ZoneProtection
	case score >= CriticalThreshold:
		return ZoneCritical
	case score >= WarningThreshold:
		return ZoneWarning
	default:
		return ZoneNormal
	}
}

func concurrencyForZone(max int, zone Zone) int {
	ratio := 1.0

	switch zone {
	case ZoneWarning:
		ratio = 0.80
	case ZoneCritical:
		ratio = 0.50
	case ZoneProtection:
		ratio = 0.25
	}

	limit := int(math.Ceil(float64(max) * ratio))
	if limit < 1 {
		return 1
	}
	if limit > max {
		return max
	}
	return limit
}
