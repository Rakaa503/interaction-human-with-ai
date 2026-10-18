package lazarus

import (
	"sync"
	"testing"
)

func TestLimiterRejectsWhenFull(t *testing.T) {
	limiter, err := NewLimiter(2)
	if err != nil {
		t.Fatal(err)
	}

	release1, ok := limiter.TryAcquire()
	if !ok {
		t.Fatal("first acquire should succeed")
	}
	defer release1()

	release2, ok := limiter.TryAcquire()
	if !ok {
		t.Fatal("second acquire should succeed")
	}
	defer release2()

	if _, ok := limiter.TryAcquire(); ok {
		t.Fatal("third acquire should be rejected")
	}

	stats := limiter.Stats()
	if stats.Active != 2 || stats.Limit != 2 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestLimiterReleaseIsIdempotent(t *testing.T) {
	limiter, err := NewLimiter(1)
	if err != nil {
		t.Fatal(err)
	}

	release, ok := limiter.TryAcquire()
	if !ok {
		t.Fatal("acquire should succeed")
	}

	release()
	release()

	if stats := limiter.Stats(); stats.Active != 0 {
		t.Fatalf("active = %d, want 0", stats.Active)
	}

	if _, ok := limiter.TryAcquire(); !ok {
		t.Fatal("acquire should succeed after release")
	}
}

func TestLimiterDynamicLimit(t *testing.T) {
	limiter, err := NewLimiter(3)
	if err != nil {
		t.Fatal(err)
	}

	release1, _ := limiter.TryAcquire()
	release2, _ := limiter.TryAcquire()
	defer release1()
	defer release2()

	if err := limiter.SetLimit(1); err != nil {
		t.Fatal(err)
	}

	if _, ok := limiter.TryAcquire(); ok {
		t.Fatal("acquire should fail while active exceeds new limit")
	}

	release1()

	if _, ok := limiter.TryAcquire(); ok {
		t.Fatal("acquire should fail while active equals new limit")
	}

	release2()

	release3, ok := limiter.TryAcquire()
	if !ok {
		t.Fatal("acquire should succeed after active reaches zero")
	}
	release3()
}

func TestLimiterConcurrentSafety(t *testing.T) {
	const (
		limit   = 8
		workers = 100
	)

	limiter, err := NewLimiter(limit)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			release, ok := limiter.TryAcquire()
			if ok {
				defer release()
			}
		}()
	}
	wg.Wait()

	if stats := limiter.Stats(); stats.Active != 0 {
		t.Fatalf("active = %d after workers finish, want 0", stats.Active)
	}
}

func TestLimiterRejectsInvalidLimits(t *testing.T) {
	if _, err := NewLimiter(0); err == nil {
		t.Fatal("expected error for zero limit")
	}

	limiter, err := NewLimiter(1)
	if err != nil {
		t.Fatal(err)
	}

	if err := limiter.SetLimit(0); err == nil {
		t.Fatal("expected error for zero dynamic limit")
	}
}
