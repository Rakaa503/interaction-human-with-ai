package lazarus

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLimiterStressConcurrency(t *testing.T) {
	const (
		limit      = int64(8)
		workers    = 500
		iterations = 100
	)

	limiter, err := NewLimiter(int(limit))
	if err != nil {
		t.Fatal(err)
	}

	var active atomic.Int64
	var maxActive atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				var release func()

				for {
					var ok bool
					release, ok = limiter.TryAcquire()
					if ok {
						break
					}
					runtime.Gosched()
				}

				current := active.Add(1)

				for {
					maximum := maxActive.Load()
					if current <= maximum ||
						maxActive.CompareAndSwap(maximum, current) {
						break
					}
				}

				runtime.Gosched()

				active.Add(-1)
				release()
			}
		}()
	}

	wg.Wait()

	if got := maxActive.Load(); got > limit {
		t.Fatalf("maximum active = %d, limit = %d", got, limit)
	}

	if got := active.Load(); got != 0 {
		t.Fatalf("active after completion = %d, want 0", got)
	}

	if got := limiter.Stats().Active; got != 0 {
		t.Fatalf("limiter active after completion = %d, want 0", got)
	}
}
