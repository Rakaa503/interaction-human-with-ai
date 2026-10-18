package lazarus

import "testing"

func BenchmarkLimiterAcquireRelease(b *testing.B) {
	limiter, err := NewLimiter(100)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		release, ok := limiter.TryAcquire()
		if !ok {
			b.Fatal("unexpected rejection")
		}
		release()
	}
}
