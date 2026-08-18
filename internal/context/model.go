package context

import "time"

type MessageContext struct {
	ConversationID uint64
	UserID         uint64

	RecentMessages []MessageSnapshot

	LastIntent     string
	LastEmotion    string
	LastTopic      string
	LastConfidence float64

	UpdatedAt time.Time
}

type MessageSnapshot struct {
	Role      string
	Content   string
	CreatedAt time.Time
}
