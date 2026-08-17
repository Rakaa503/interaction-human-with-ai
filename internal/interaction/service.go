package interaction

import "fmt"

type RepositoryInterface interface {
	Create(interaction *Interaction) error
	GetByID(id uint64) (*Interaction, error)
	GetByConversationID(conversationID uint64) ([]Interaction, error)
}

type Service struct {
	repository RepositoryInterface
	analyzer   Analyzer
}

func NewService(
	repository RepositoryInterface,
	analyzer Analyzer,
) *Service {
	return &Service{
		repository: repository,
		analyzer:   analyzer,
	}
}

func (s *Service) Process(
	userID uint64,
	conversationID uint64,
	input string,
) (*Interaction, error) {

	if input == "" {
		return nil, fmt.Errorf("interaction input cannot be empty")
	}

	result, err := s.analyzer.Analyze(input)
	if err != nil {
		return nil, err
	}

	interaction := &Interaction{
		UserID:         userID,
		ConversationID: conversationID,
		Input:          input,
		Intent:         &result.Intent,
		Emotion:        &result.Emotion,
		Topic:          &result.Topic,
		Confidence:     &result.Confidence,
	}

	if err := s.repository.Create(interaction); err != nil {
		return nil, err
	}

	return interaction, nil
}

func (s *Service) GetByID(
	id uint64,
) (*Interaction, error) {
	return s.repository.GetByID(id)
}

func (s *Service) GetByConversationID(
	conversationID uint64,
) ([]Interaction, error) {
	return s.repository.GetByConversationID(conversationID)
}
