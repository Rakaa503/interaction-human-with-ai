package knowledge

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrInvalidDocument   = errors.New("invalid knowledge document")
	ErrDuplicateDocument = errors.New("knowledge document already exists")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddDocument(
	title string,
	content string,
	url *string,
	source *string,
	category *string,
) (*Document, error) {

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return nil, ErrInvalidDocument
	}

	contentHash := GenerateContentHash(content)

	_, err := s.repository.FindByHash(contentHash)

	if err == nil {
		return nil, ErrDuplicateDocument
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	document := &Document{
		Title:       title,
		Content:     content,
		URL:         url,
		Source:      source,
		Category:    category,
		ContentHash: contentHash,
	}

	if err := s.repository.Create(document); err != nil {
		return nil, err
	}

	return document, nil
}
