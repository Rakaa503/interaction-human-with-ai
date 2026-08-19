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

func (m *mockConversationRepository) CreateConversation(
	conversation *conversation.Conversation,
) error {
	return nil
}

func (m *mockConversationRepository) GetConversation(
	id uint64,
) (*conversation.Conversation, error) {
	return &conversation.Conversation{
		ID: id,
	}, nil
}

func (m *mockConversationRepository) AddMessage(
	message *conversation.Message,
) error {
	return nil
}

type mockInteractionRepository struct{}

func (m *mockInteractionRepository) Create(
	interaction *interaction.Interaction,
) error {
	return nil
}

func (m *mockInteractionRepository) GetByID(
	id uint64,
) (*interaction.Interaction, error) {
	return &interaction.Interaction{
		ID: id,
	}, nil
}

func (m *mockInteractionRepository) GetByConversationID(
	conversationID uint64,
) ([]interaction.Interaction, error) {
	return []interaction.Interaction{}, nil
}

func newTestOrchestrator() *Service {
	conversationRepository := &mockConversationRepository{}
	interactionRepository := &mockInteractionRepository{}

	conversationService := conversation.NewService(
		conversationRepository,
	)

	interactionService := interaction.NewService(
		interactionRepository,
		&mockAnalyzer{},
	)

	contextService := appcontext.NewService(
		conversationRepository,
		interactionRepository,
	)

	decisionService := decision.NewService()
	responseService := response.NewService()

	return NewService(
		interactionService,
		contextService,
		decisionService,
		responseService,
		conversationService,
	)
}

func TestOrchestrator(t *testing.T) {
	service := newTestOrchestrator()

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
	service := newTestOrchestrator()

	_, err := service.Process(
		3,
		5,
		"",
	)

	if err == nil {
		t.Fatal("expected error for empty input")
	}
}
