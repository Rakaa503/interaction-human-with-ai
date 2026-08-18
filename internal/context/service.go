package context

import (
	"time"

	"github.com/Rakaa503/AviGo/internal/conversation"
	"github.com/Rakaa503/AviGo/internal/interaction"
)

type ConversationRepository interface {
	GetMessages(conversationID uint64) ([]conversation.Message, error)
}

type InteractionRepository interface {
	GetByConversationID(
		conversationID uint64,
	) ([]interaction.Interaction, error)
}

type Service struct {
	conversationRepository ConversationRepository
	interactionRepository  InteractionRepository
}

func NewService(
	conversationRepository ConversationRepository,
	interactionRepository InteractionRepository,
) *Service {
	return &Service{
		conversationRepository: conversationRepository,
		interactionRepository:  interactionRepository,
	}
}

func (s *Service) Build(
	conversationID uint64,
	userID uint64,
) (*MessageContext, error) {

	messages, err := s.conversationRepository.GetMessages(
		conversationID,
	)

	if err != nil {
		return nil, err
	}

	interactions, err := s.interactionRepository.GetByConversationID(
		conversationID,
	)

	if err != nil {
		return nil, err
	}

	ctx := &MessageContext{
		ConversationID: conversationID,
		UserID:         userID,
		RecentMessages: make([]MessageSnapshot, 0),
		UpdatedAt:      time.Now(),
	}

	// Ambil maksimal 10 pesan terakhir.
	start := 0

	if len(messages) > 10 {
		start = len(messages) - 10
	}

	for _, message := range messages[start:] {
		ctx.RecentMessages = append(
			ctx.RecentMessages,
			MessageSnapshot{
				Role:      message.Role,
				Content:   message.Content,
				CreatedAt: message.CreatedAt,
			},
		)
	}

	// Ambil interaction terbaru.
	if len(interactions) > 0 {
		last := interactions[len(interactions)-1]

		if last.Intent != nil {
			ctx.LastIntent = *last.Intent
		}

		if last.Emotion != nil {
			ctx.LastEmotion = *last.Emotion
		}

		if last.Topic != nil {
			ctx.LastTopic = *last.Topic
		}

		if last.Confidence != nil {
			ctx.LastConfidence = *last.Confidence
		}
	}

	return ctx, nil
}
