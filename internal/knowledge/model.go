package knowledge

import "time"

type Document struct {
	ID          uint64  `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Content     string  `gorm:"not null"`
	URL         *string `gorm:"type:text"`
	Source      *string `gorm:"type:text"`
	Category    *string `gorm:"type:varchar(100)"`
	ContentHash string  `gorm:"type:char(64);uniqueIndex;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Document) TableName() string {
	return "knowledge_documents"
}
