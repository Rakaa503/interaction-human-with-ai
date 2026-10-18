package lazarus

import (
	"errors"
	"sync"
)

var ErrInvalidLimit = errors.New("concurrency limit must be greater than zero")

type Limiter struct {
	mu        sync.Mutex
	active    int
	limit     int
	rejectNew bool
}

type LimiterStats struct {
	Active    int
	Limit     int
	RejectNew bool
}

func NewLimiter(limit int) (*Limiter, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	return &Limiter{limit: limit}, nil
}

// TryAcquire menolak request jika proteksi aktif atau kapasitas penuh.
func (l *Limiter) TryAcquire() (release func(), ok bool) {
	l.mu.Lock()

	if l.rejectNew || l.active >= l.limit {
		l.mu.Unlock()
		return nil, false
	}

	l.active++
	l.mu.Unlock()

	var once sync.Once
	release = func() {
		once.Do(func() {
			l.mu.Lock()
			l.active--
			l.mu.Unlock()
		})
	}

	return release, true
}

func (l *Limiter) SetLimit(limit int) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}

	l.mu.Lock()
	l.limit = limit
	l.mu.Unlock()

	return nil
}

// ApplyPolicy memperbarui limit dan status proteksi secara atomik.
func (l *Limiter) ApplyPolicy(limit int, rejectNew bool) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}

	l.mu.Lock()
	l.limit = limit
	l.rejectNew = rejectNew
	l.mu.Unlock()

	return nil
}

func (l *Limiter) Stats() LimiterStats {
	l.mu.Lock()
	defer l.mu.Unlock()

	return LimiterStats{
		Active:    l.active,
		Limit:     l.limit,
		RejectNew: l.rejectNew,
	}
}
