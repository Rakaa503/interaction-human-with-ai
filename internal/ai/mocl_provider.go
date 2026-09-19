package ai

import "context"

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Generate(
	ctx context.Context,
	request Request,
) (*Response, error) {
	return &Response{
		Content: "AVIGO AI: Saya menerima pertanyaan kamu dan siap membantu menganalisisnya.",
	}, nil
}
