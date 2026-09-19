package ai

import "context"

// Provider is the abstraction for an AI reasoning provider.
//
// Implementations may use:
//   - OpenAI
//   - Gemini
//   - Ollama
//   - another compatible AI service
//   - a local model
type Provider interface {
	Generate(ctx context.Context, request Request) (*Response, error)
}
