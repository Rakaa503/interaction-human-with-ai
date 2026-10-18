package lazarus

import "errors"

var ErrNilLimiter = errors.New("lazarus limiter cannot be nil")

// Controller menghubungkan hasil formula dengan kebijakan limiter.
type Controller struct {
	limiter *Limiter
}

func NewController(limiter *Limiter) (*Controller, error) {
	if limiter == nil {
		return nil, ErrNilLimiter
	}

	return &Controller{
		limiter: limiter,
	}, nil
}

// Update mengevaluasi metrik dan menerapkan kebijakan hasil evaluasi.
func (c *Controller) Update(metrics Metrics) (Result, error) {
	result, err := Evaluate(metrics)
	if err != nil {
		return Result{}, err
	}

	if err := c.limiter.ApplyPolicy(
		result.SuggestedConcurrency,
		result.RejectNew,
	); err != nil {
		return Result{}, err
	}

	return result, nil
}
