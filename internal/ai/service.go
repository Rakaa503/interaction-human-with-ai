package ai

import "context"

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) Generate(
	ctx context.Context,
	input string,
) (*Response, error) {
	if input == "" {
		return nil, nil
	}

	return s.provider.Generate(ctx, Request{
		Input: input,
	})
}
