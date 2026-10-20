package lazarus

import (
	"context"
	"errors"
	"math"
	"time"
)

var (
	ErrInvalidAdaptiveConfig = errors.New("invalid adaptive controller configuration")
	ErrInvalidMetricsSource  = errors.New("metrics source cannot be nil")
)

type AdaptiveConfig struct {
	// Alpha mengatur pengaruh metrik terbaru pada exponential moving average.
	// Nilai harus > 0 dan <= 1.
	Alpha float64

	// Hysteresis mencegah perpindahan zona akibat fluktuasi kecil.
	Hysteresis float64

	// RecoveryStep adalah persentase MaxConcurrency yang boleh
	// ditambahkan pada limit dalam satu siklus recovery.
	RecoveryStep float64
}

func DefaultAdaptiveConfig() AdaptiveConfig {
	return AdaptiveConfig{
		Alpha:        0.35,
		Hysteresis:   5,
		RecoveryStep: 0.10,
	}
}

type AdaptiveController struct {
	controller *Controller
	limiter    *Limiter
	config     AdaptiveConfig

	initialized bool
	smoothed    float64
	zone        Zone
}

func NewAdaptiveController(
	limiter *Limiter,
	config AdaptiveConfig,
) (*AdaptiveController, error) {
	if limiter == nil {
		return nil, ErrNilLimiter
	}
	if math.IsNaN(config.Alpha) ||
		math.IsInf(config.Alpha, 0) ||
		config.Alpha <= 0 || config.Alpha > 1 ||
		math.IsNaN(config.Hysteresis) ||
		math.IsInf(config.Hysteresis, 0) ||
		config.Hysteresis < 0 ||
		math.IsNaN(config.RecoveryStep) ||
		math.IsInf(config.RecoveryStep, 0) ||
		config.RecoveryStep <= 0 || config.RecoveryStep > 1 {
		return nil, ErrInvalidAdaptiveConfig
	}

	controller, err := NewController(limiter)
	if err != nil {
		return nil, err
	}

	return &AdaptiveController{
		controller: controller,
		limiter:    limiter,
		config:     config,
	}, nil
}

// Update mengevaluasi metrik terbaru, menghaluskan score,
// menerapkan hysteresis, dan mengubah concurrency limit secara langsung.
// Penurunan limit dilakukan segera; kenaikannya dilakukan bertahap.
func (a *AdaptiveController) Update(metrics Metrics) (Result, error) {
	raw, err := Evaluate(metrics)
	if err != nil {
		return Result{}, err
	}

	if !a.initialized {
		a.smoothed = raw.Score
		a.zone = raw.Zone
		a.initialized = true
	} else {
		a.smoothed = a.config.Alpha*raw.Score +
			(1-a.config.Alpha)*a.smoothed
		a.zone = zoneWithHysteresis(
			a.smoothed,
			a.zone,
			a.config.Hysteresis,
		)
	}

	// Hard limit dan antrean penuh tidak boleh ditunda oleh smoothing.
	if raw.RejectNew {
		a.zone = ZoneProtection
	}

	result := raw
	result.Score = a.smoothed
	result.Zone = a.zone
	result.SuggestedConcurrency = concurrencyForZone(
		metrics.MaxConcurrency,
		a.zone,
	)

	if raw.RejectNew {
		result.Zone = ZoneProtection
		result.SuggestedConcurrency = concurrencyForZone(
			metrics.MaxConcurrency,
			ZoneProtection,
		)
	}

	current := a.limiter.Stats().Limit
	target := result.SuggestedConcurrency

	if raw.RejectNew || target < current {
		// Degradasi kapasitas harus langsung diterapkan.
		current = target
	} else if target > current {
		// Recovery bertahap, tidak langsung melonjak ke kapasitas penuh.
		step := int(math.Ceil(
			float64(metrics.MaxConcurrency) * a.config.RecoveryStep,
		))
		if step < 1 {
			step = 1
		}

		current += step
		if current > target {
			current = target
		}
	}

	if err := a.limiter.ApplyPolicy(current, raw.RejectNew); err != nil {
		return Result{}, err
	}

	return result, nil
}

func zoneWithHysteresis(score float64, current Zone, margin float64) Zone {
	rank := zoneRank(current)
	if rank < 0 {
		return zoneForScore(score)
	}

	thresholds := []float64{
		WarningThreshold,
		CriticalThreshold,
		ProtectThreshold,
	}

	// Naik zona hanya ketika melewati threshold + margin.
	for rank < len(thresholds) &&
		score >= thresholds[rank]+margin {
		rank++
	}

	// Turun zona hanya ketika melewati threshold - margin.
	for rank > 0 &&
		score < thresholds[rank-1]-margin {
		rank--
	}

	switch rank {
	case 0:
		return ZoneNormal
	case 1:
		return ZoneWarning
	case 2:
		return ZoneCritical
	default:
		return ZoneProtection
	}
}

func zoneRank(zone Zone) int {
	switch zone {
	case ZoneNormal:
		return 0
	case ZoneWarning:
		return 1
	case ZoneCritical:
		return 2
	case ZoneProtection:
		return 3
	default:
		return -1
	}
}

// Run memanggil sumber metrik secara berkala sampai context dibatalkan.
// collect wajib mengambil metrik runtime yang nyata dari aplikasi.
func (a *AdaptiveController) Run(
	ctx context.Context,
	interval time.Duration,
	collect func() (Metrics, error),
) error {
	if interval <= 0 {
		return errors.New("monitoring interval must be greater than zero")
	}
	if collect == nil {
		return ErrInvalidMetricsSource
	}

	update := func() error {
		metrics, err := collect()
		if err != nil {
			return err
		}
		_, err = a.Update(metrics)
		return err
	}

	if err := update(); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := update(); err != nil {
				return err
			}
		}
	}
}
