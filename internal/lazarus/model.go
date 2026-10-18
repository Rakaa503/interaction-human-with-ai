
package lazarus

import "time"

type Zone string

const (
	ZoneNormal     Zone = "normal"
	ZoneWarning    Zone = "warning"
	ZoneCritical   Zone = "critical"
	ZoneProtection Zone = "protection"
)

type Metrics struct {
	CPUUsage      float64
	MemoryUsage   float64
	QueueDepth    int
	QueueLimit    int
	P95Latency    time.Duration
	LatencyTarget time.Duration
	MaxConcurrency int
}

type Result struct {
	Score               float64
	Zone                Zone
	SuggestedConcurrency int
	RejectNew           bool
	Reasons             []string
}
