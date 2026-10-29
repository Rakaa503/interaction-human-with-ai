package lazarus

import (
	"testing"
	"time"
)

func TestLimiterRejectsWhenCapacityIsFull(t *testing.T) {
	limiter, err := NewLimiter(1)
	if err != nil {
		t.Fatalf("NewLimiter() error = %v", err)
	}

	release, ok := limiter.TryAcquire()
	if !ok {
		t.Fatal("first request should acquire the slot")
	}

	if _, ok := limiter.TryAcquire(); ok {
		t.Fatal("second request should be rejected at capacity")
	}

	stats := limiter.Stats()
	if stats.Active != 1 {
		t.Fatalf("Active = %d, want 1", stats.Active)
	}

	release()

	stats = limiter.Stats()
	if stats.Active != 0 {
		t.Fatalf("Active after release = %d, want 0", stats.Active)
	}

	release() // Memastikan release aman dipanggil dua kali.

	stats = limiter.Stats()
	if stats.Active != 0 {
		t.Fatalf("Active after second release = %d, want 0", stats.Active)
	}

	if _, ok := limiter.TryAcquire(); !ok {
		t.Fatal("slot should be available after release")
	}
}

func TestLimiterApplyPolicyRejectsNewRequests(t *testing.T) {
	limiter, err := NewLimiter(5)
	if err != nil {
		t.Fatal(err)
	}

	if err := limiter.ApplyPolicy(2, true); err != nil {
		t.Fatal(err)
	}

	stats := limiter.Stats()
	if stats.Limit != 2 || !stats.RejectNew {
		t.Fatalf("unexpected policy state: %+v", stats)
	}

	if _, ok := limiter.TryAcquire(); ok {
		t.Fatal("request should be rejected while protection is active")
	}

	if err := limiter.ApplyPolicy(2, false); err != nil {
		t.Fatal(err)
	}

	release, ok := limiter.TryAcquire()
	if !ok {
		t.Fatal("request should be accepted after protection is disabled")
	}
	release()
}

func TestTelemetryTracksCompletedAndRejectedRequests(t *testing.T) {
	telemetry, err := NewTelemetry(8)
	if err != nil {
		t.Fatal(err)
	}

	finish := telemetry.StartRequest()
	time.Sleep(2 * time.Millisecond)
	finish()
	finish() // Satu request hanya dihitung selesai sekali.

	telemetry.RecordRejected()

	snapshot := telemetry.Snapshot()

	if snapshot.ActiveRequests != 0 {
		t.Errorf("ActiveRequests = %d, want 0", snapshot.ActiveRequests)
	}
	if snapshot.Completed != 1 {
		t.Errorf("Completed = %d, want 1", snapshot.Completed)
	}
	if snapshot.Rejected != 1 {
		t.Errorf("Rejected = %d, want 1", snapshot.Rejected)
	}
	if snapshot.P95Latency <= 0 {
		t.Errorf("P95Latency = %v, want greater than zero", snapshot.P95Latency)
	}
}
