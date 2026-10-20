package lazarus

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var ErrInvalidLatencyWindow = errors.New("latency window must be greater than zero")

// Telemetry menyimpan sampel latency request yang selesai diproses.
// Kapasitas tetap mencegah pertumbuhan penggunaan memori tanpa batas.
type Telemetry struct {
	mu        sync.Mutex
	latencies []time.Duration
	capacity  int
	next      int
	count     int
	active    int
	total     uint64
	rejected  uint64
}

func NewTelemetry(capacity int) (*Telemetry, error) {
	if capacity <= 0 {
		return nil, ErrInvalidLatencyWindow
	}

	return &Telemetry{
		latencies: make([]time.Duration, capacity),
		capacity:  capacity,
	}, nil
}

// StartRequest mencatat request yang mulai diproses.
// Fungsi kembalian harus dipanggil setelah request selesai.
func (t *Telemetry) StartRequest() func() {
	t.mu.Lock()
	t.active++
	t.mu.Unlock()

	start := time.Now()
	var once sync.Once

	return func() {
		once.Do(func() {
			elapsed := time.Since(start)

			t.mu.Lock()
			defer t.mu.Unlock()

			t.active--
			t.total++

			t.latencies[t.next] = elapsed
			t.next = (t.next + 1) % t.capacity

			if t.count < t.capacity {
				t.count++
			}
		})
	}
}

func (t *Telemetry) RecordRejected() {
	t.mu.Lock()
	t.rejected++
	t.mu.Unlock()
}

type TelemetrySnapshot struct {
	ActiveRequests int
	Completed      uint64
	Rejected       uint64
	P95Latency     time.Duration
}

func (t *Telemetry) Snapshot() TelemetrySnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()

	result := TelemetrySnapshot{
		ActiveRequests: t.active,
		Completed:      t.total,
		Rejected:       t.rejected,
	}

	if t.count == 0 {
		return result
	}

	samples := make([]time.Duration, t.count)
	copy(samples, t.latencies[:t.count])

	// Setelah buffer penuh, urutan fisik tidak lagi kronologis.
	// Untuk P95, yang dibutuhkan adalah distribusi sampel.
	sort.Slice(samples, func(i, j int) bool {
		return samples[i] < samples[j]
	})

	index := int(0.95*float64(len(samples)-1) + 0.999999)
	if index >= len(samples) {
		index = len(samples) - 1
	}
	result.P95Latency = samples[index]

	return result
}
