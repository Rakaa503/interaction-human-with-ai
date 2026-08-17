package interaction

import "time"

type Interaction struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	UserID         uint64 `gorm:"not null;index" json:"userId"`
	ConversationID uint64 `gorm:"not null;index" json:"conversationId"`

	Input string `gorm:"type:text;not null" json:"input"`

	Intent     *string  `gorm:"type:varchar(100)" json:"intent,omitempty"`
	Emotion    *string  `gorm:"type:varchar(100)" json:"emotion,omitempty"`
	Topic      *string  `gorm:"type:varchar(100)" json:"topic,omitempty"`
	Confidence *float64 `json:"confidence,omitempty"`

	Response *string `gorm:"type:text" json:"response,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
