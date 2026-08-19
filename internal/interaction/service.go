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

// Analyze menjalankan analyzer tanpa menyimpan interaction.
func (s *Service) Analyze(
	input string,
) (*AnalysisResult, error) {

	if input == "" {
		return nil, fmt.Errorf(
			"interaction input cannot be empty",
		)
	}

	return s.analyzer.Analyze(input)
}

// Save menyimpan hasil analysis menjadi interaction.
func (s *Service) Save(
	userID uint64,
	conversationID uint64,
	input string,
	result *AnalysisResult,
	response string,
) (*Interaction, error) {

	if input == "" {
		return nil, fmt.Errorf(
			"interaction input cannot be empty",
		)
	}

	if result == nil {
		return nil, fmt.Errorf(
			"analysis result cannot be nil",
		)
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

	if response != "" {
		interaction.Response = &response
	}

	if err := s.repository.Create(interaction); err != nil {
		return nil, err
	}

	return interaction, nil
}

// Process menjalankan analysis sekaligus menyimpan interaction.
func (s *Service) Process(
	userID uint64,
	conversationID uint64,
	input string,
) (*Interaction, error) {

	result, err := s.Analyze(input)
	if err != nil {
		return nil, err
	}

	return s.Save(
		userID,
		conversationID,
		input,
		result,
		"",
	)
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
