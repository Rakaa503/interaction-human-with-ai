package interaction

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	interaction *Interaction,
) error {
	return r.db.Create(interaction).Error
}

func (r *Repository) GetByID(
	id uint64,
) (*Interaction, error) {

	var interaction Interaction

	err := r.db.
		First(&interaction, id).
		Error

	if err != nil {
		return nil, err
	}

	return &interaction, nil
}

func (r *Repository) GetByConversationID(
	conversationID uint64,
) ([]Interaction, error) {

	var interactions []Interaction

	err := r.db.
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&interactions).
		Error

	return interactions, err
}
