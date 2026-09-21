package memory

import (
	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Retrieve(
	ctx *appcontext.MessageContext,
	memoryType MemoryType,
) *Memory {
	switch memoryType {
	case MemoryName:
		return ExtractName(ctx)

	case MemoryActivity:
		return ExtractActivity(ctx)

	default:
		return nil
	}
}
