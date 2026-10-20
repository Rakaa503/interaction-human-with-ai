package lazarus

import (
	"testing"
	"time"
)

func TestAdaptiveControllerReducesLimitAndRecovers(t *testing.T) {
	limiter, err := NewLimiter(100)
	if err != nil {
		t.Fatal(err)
	}

	adaptive, err := NewAdaptiveController(
		limiter,
		DefaultAdaptiveConfig(),
	)
	if err != nil {
		t.Fatal(err)
	}

	healthy := Metrics{
		CPUUsage:       10,
		MemoryUsage:    20,
		QueueDepth:     0,
		QueueLimit:     100,
		P95Latency:     10 * time.Millisecond,
		LatencyTarget:  100 * time.Millisecond,
		MaxConcurrency: 100,
	}

	overloaded := healthy
	overloaded.CPUUsage = 96

	// Kondisi sehat: kapasitas penuh.
	result, err := adaptive.Update(healthy)
	if err != nil {
		t.Fatal(err)
	}
	if got := limiter.Stats().Limit; got != 100 {
		t.Fatalf("healthy limit = %d, want 100", got)
	}
	if result.Zone != ZoneNormal {
		t.Fatalf("healthy zone = %s, want normal", result.Zone)
	}

	// CPU hard limit: kapasitas turun dan request baru ditolak.
	result, err = adaptive.Update(overloaded)
	if err != nil {
		t.Fatal(err)
	}

	stats := limiter.Stats()
	if stats.Limit >= 100 {
		t.Fatalf("overload limit = %d, want less than 100", stats.Limit)
	}
	if !stats.RejectNew {
		t.Fatal("expected new requests to be rejected at CPU hard limit")
	}
	if result.Zone != ZoneProtection {
		t.Fatalf("overload zone = %s, want protection", result.Zone)
	}

	// CPU kembali sehat: proteksi dilepas dan kapasitas pulih bertahap.
	result, err = adaptive.Update(healthy)
	if err != nil {
		t.Fatal(err)
	}
	if limiter.Stats().RejectNew {
		t.Fatal("expected rejection policy to be disabled after recovery")
	}

	previousLimit := limiter.Stats().Limit

	for i := 0; i < 60; i++ {
		if _, err := adaptive.Update(healthy); err != nil {
			t.Fatal(err)
		}
	}

	recoveredLimit := limiter.Stats().Limit
	if recoveredLimit <= previousLimit {
		t.Fatalf(
			"limit did not recover: before=%d after=%d",
			previousLimit,
			recoveredLimit,
		)
	}
	if recoveredLimit != 100 {
		t.Fatalf("recovered limit = %d, want 100", recoveredLimit)
	}
}

func TestAdaptiveControllerRejectsInvalidConfig(t *testing.T) {
	limiter, err := NewLimiter(10)
	if err != nil {
		t.Fatal(err)
	}

	config := DefaultAdaptiveConfig()
	config.Alpha = 0

	if _, err := NewAdaptiveController(limiter, config); err == nil {
		t.Fatal("expected invalid config error")
	}
}

func TestAdaptiveControllerHysteresis(t *testing.T) {
	if got := zoneWithHysteresis(52, ZoneNormal, 5); got != ZoneNormal {
		t.Fatalf("zone at score 52 = %s, want normal", got)
	}

	if got := zoneWithHysteresis(56, ZoneNormal, 5); got != ZoneWarning {
		t.Fatalf("zone at score 56 = %s, want warning", got)
	}

	if got := zoneWithHysteresis(48, ZoneWarning, 5); got != ZoneWarning {
		t.Fatalf("zone at score 48 = %s, want warning", got)
	}

	if got := zoneWithHysteresis(44, ZoneWarning, 5); got != ZoneNormal {
		t.Fatalf("zone at score 44 = %s, want normal", got)
	}
}
