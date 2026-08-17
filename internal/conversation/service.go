package conversation

import "fmt"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateConversation(
	userID uint64,
	title *string,
) (*Conversation, error) {
	conversation := &Conversation{
		UserID: userID,
		Title:  title,
	}

	if err := s.repository.CreateConversation(conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *Service) AddMessage(
	conversationID uint64,
	role string,
	content string,
) (*Message, error) {
	if role != "system" &&
		role != "user" &&
		role != "assistant" {
		return nil, fmt.Errorf("invalid message role")
	}

	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	message := &Message{
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
	}

	if err := s.repository.AddMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *Service) GetConversation(
	id uint64,
) (*Conversation, error) {
	return s.repository.GetConversation(id)
}
