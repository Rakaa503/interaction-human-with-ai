package conversation

import "time"

type Conversation struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"userId"`
	Title     *string   `gorm:"type:varchar(255)" json:"title,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Messages []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

type Message struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	ConversationID uint64    `gorm:"not null;index" json:"conversationId"`
	Role           string    `gorm:"type:varchar(20);not null" json:"role"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}
