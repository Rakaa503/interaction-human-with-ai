package orchestrator

import (
	"testing"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
	"github.com/Rakaa503/AviGo/internal/conversation"
	"github.com/Rakaa503/AviGo/internal/decision"
	"github.com/Rakaa503/AviGo/internal/interaction"
	"github.com/Rakaa503/AviGo/internal/response"
)

type mockAnalyzer struct{}

func (m *mockAnalyzer) Analyze(
	input string,
) (*interaction.AnalysisResult, error) {
	return &interaction.AnalysisResult{
		Intent:     "greeting",
		Emotion:    "neutral",
		Topic:      "general",
		Confidence: 0.90,
	}, nil
}

type mockConversationRepository struct{}

func (m *mockConversationRepository) GetMessages(
	conversationID uint64,
) ([]conversation.Message, error) {
	return []conversation.Message{}, nil
}

type mockInteractionRepository struct{}

func (m *mockInteractionRepository) GetByConversationID(
	conversationID uint64,
) ([]interaction.Interaction, error) {
	return []interaction.Interaction{}, nil
}

func TestOrchestrator(t *testing.T) {
	contextService := appcontext.NewService(
		&mockConversationRepository{},
		&mockInteractionRepository{},
	)

	decisionService := decision.NewService()
	responseService := response.NewService()
	analyzer := &mockAnalyzer{}

	service := NewService(
		analyzer,
		contextService,
		decisionService,
		responseService,
	)

	result, err := service.Process(
		3,
		5,
		"Hai AVIGO",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Intent != "greeting" {
		t.Fatalf(
			"expected greeting intent, got %q",
			result.Intent,
		)
	}

	if result.Action != "respond_greeting" {
		t.Fatalf(
			"expected respond_greeting action, got %q",
			result.Action,
		)
	}

	if result.Response == "" {
		t.Fatal("expected response content")
	}
}

func TestOrchestratorEmptyInput(t *testing.T) {
	contextService := appcontext.NewService(
		&mockConversationRepository{},
		&mockInteractionRepository{},
	)

	service := NewService(
		&mockAnalyzer{},
		contextService,
		decision.NewService(),
		response.NewService(),
	)

	_, err := service.Process(
		3,
		5,
		"",
	)

	if err == nil {
		t.Fatal("expected error for empty input")
	}
}
