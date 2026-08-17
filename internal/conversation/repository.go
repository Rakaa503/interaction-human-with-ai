package conversation

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateConversation(
	conversation *Conversation,
) error {
	return r.db.Create(conversation).Error
}

func (r *Repository) GetConversation(
	id uint64,
) (*Conversation, error) {
	var conversation Conversation

	err := r.db.
		Preload("Messages").
		First(&conversation, id).
		Error

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *Repository) AddMessage(
	message *Message,
) error {
	return r.db.Create(message).Error
}

func (r *Repository) GetMessages(
	conversationID uint64,
) ([]Message, error) {
	var messages []Message

	err := r.db.
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).
		Error

	return messages, err
}
