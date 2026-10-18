package lazarus

import (
	"runtime"
	"testing"
)

func BenchmarkLimiterAcquireReleaseParallel(b *testing.B) {
	limiter, err := NewLimiter(100)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for {
				release, ok := limiter.TryAcquire()
				if ok {
					release()
					break
				}

				runtime.Gosched()
			}
		}
	})
}
