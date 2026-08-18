package context

import (
	"testing"
	"time"

	"github.com/Rakaa503/AviGo/internal/conversation"
	"github.com/Rakaa503/AviGo/internal/interaction"
)

type mockConversationRepository struct {
	messages []conversation.Message
}

func (m *mockConversationRepository) GetMessages(
	conversationID uint64,
) ([]conversation.Message, error) {
	return m.messages, nil
}

type mockInteractionRepository struct {
	interactions []interaction.Interaction
}

func (m *mockInteractionRepository) GetByConversationID(
	conversationID uint64,
) ([]interaction.Interaction, error) {
	return m.interactions, nil
}

func TestBuildContext(t *testing.T) {
	now := time.Now()

	messages := make([]conversation.Message, 0, 12)

	for i := 1; i <= 12; i++ {
		messages = append(messages, conversation.Message{
			ID:             uint64(i),
			ConversationID: 5,
			Role:           "user",
			Content:        "message",
			CreatedAt:      now.Add(time.Duration(i) * time.Minute),
		})
	}

	intent := "problem_solving"
	emotion := "frustrated"
	topic := "technical"
	confidence := 0.92

	interactions := []interaction.Interaction{
		{
			ID:             1,
			UserID:         3,
			ConversationID: 5,
			Input:          "Database saya error",
			Intent:         &intent,
			Emotion:        &emotion,
			Topic:          &topic,
			Confidence:     &confidence,
			CreatedAt:      now,
		},
	}

	service := NewService(
		&mockConversationRepository{
			messages: messages,
		},
		&mockInteractionRepository{
			interactions: interactions,
		},
	)

	ctx, err := service.Build(5, 3)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.ConversationID != 5 {
		t.Fatalf(
			"expected conversation ID 5, got %d",
			ctx.ConversationID,
		)
	}

	if ctx.UserID != 3 {
		t.Fatalf(
			"expected user ID 3, got %d",
			ctx.UserID,
		)
	}

	if len(ctx.RecentMessages) != 10 {
		t.Fatalf(
			"expected 10 recent messages, got %d",
			len(ctx.RecentMessages),
		)
	}

	if ctx.RecentMessages[0].Content != "message" {
		t.Fatal("expected first recent message to exist")
	}

	if ctx.LastIntent != "problem_solving" {
		t.Fatalf(
			"expected intent problem_solving, got %s",
			ctx.LastIntent,
		)
	}

	if ctx.LastEmotion != "frustrated" {
		t.Fatalf(
			"expected emotion frustrated, got %s",
			ctx.LastEmotion,
		)
	}

	if ctx.LastTopic != "technical" {
		t.Fatalf(
			"expected topic technical, got %s",
			ctx.LastTopic,
		)
	}

	if ctx.LastConfidence != 0.92 {
		t.Fatalf(
			"expected confidence 0.92, got %.2f",
			ctx.LastConfidence,
		)
	}

	t.Log("Context Engine successfully built conversation context")
}
