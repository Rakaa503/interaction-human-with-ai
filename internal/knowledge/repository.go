package knowledge

import "gorm.io/gorm"

type RepositoryInterface interface {
	Create(document *Document) error
	FindByHash(hash string) (*Document, error)
	FindAll() ([]Document, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(document *Document) error {
	return r.db.Create(document).Error
}

func (r *Repository) FindByHash(hash string) (*Document, error) {
	var document Document

	err := r.db.
		Where("content_hash = ?", hash).
		First(&document).
		Error

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *Repository) FindAll() ([]Document, error) {
	var documents []Document

	err := r.db.
		Order("created_at DESC").
		Find(&documents).
		Error

	return documents, err
}
