package lazarus

import "testing"

func TestControllerAppliesFormulaResult(t *testing.T) {
	limiter, err := NewLimiter(100)
	if err != nil {
		t.Fatal(err)
	}

	controller, err := NewController(limiter)
	if err != nil {
		t.Fatal(err)
	}

	metrics := Metrics{
		CPUUsage:       10,
		MemoryUsage:    20,
		QueueDepth:     1,
		QueueLimit:     100,
		P95Latency:     20,
		LatencyTarget:  100,
		MaxConcurrency: 100,
	}

	result, err := controller.Update(metrics)
	if err != nil {
		t.Fatalf("controller update failed: %v", err)
	}

	stats := limiter.Stats()

	if stats.Limit != result.SuggestedConcurrency {
		t.Fatalf(
			"limit = %d, want %d",
			stats.Limit,
			result.SuggestedConcurrency,
		)
	}

	if stats.RejectNew != result.RejectNew {
		t.Fatalf(
			"RejectNew = %v, want %v",
			stats.RejectNew,
			result.RejectNew,
		)
	}
}

func TestControllerRejectsNewRequestsInProtection(t *testing.T) {
	limiter, err := NewLimiter(100)
	if err != nil {
		t.Fatal(err)
	}

	controller, err := NewController(limiter)
	if err != nil {
		t.Fatal(err)
	}

	metrics := Metrics{
		CPUUsage:       96,
		MemoryUsage:    50,
		QueueDepth:     1,
		QueueLimit:     100,
		P95Latency:     20,
		LatencyTarget:  100,
		MaxConcurrency: 100,
	}

	result, err := controller.Update(metrics)
	if err != nil {
		t.Fatalf("controller update failed: %v", err)
	}

	if !result.RejectNew {
		t.Fatal("expected formula to enable request rejection")
	}

	_, ok := limiter.TryAcquire()
	if ok {
		release, _ := limiter.TryAcquire()
		if release != nil {
			release()
		}
		t.Fatal("expected new request to be rejected")
	}

	if stats := limiter.Stats(); !stats.RejectNew {
		t.Fatal("expected RejectNew to be enabled")
	}
}
