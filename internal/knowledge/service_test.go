package knowledge

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

type fakeRepository struct {
	documents map[string]*Document
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		documents: make(map[string]*Document),
	}
}

func (f *fakeRepository) Create(document *Document) error {
	if _, exists := f.documents[document.ContentHash]; exists {
		return errors.New("duplicate content hash")
	}

	f.documents[document.ContentHash] = document

	return nil
}

func (f *fakeRepository) FindByHash(hash string) (*Document, error) {
	document, exists := f.documents[hash]

	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	return document, nil
}

func (f *fakeRepository) FindAll() ([]Document, error) {
	documents := make([]Document, 0, len(f.documents))

	for _, document := range f.documents {
		documents = append(documents, *document)
	}

	return documents, nil
}

func TestGenerateContentHash(t *testing.T) {
	content := "Machine learning adalah cabang dari artificial intelligence."

	hash := GenerateContentHash(content)

	if hash == "" {
		t.Fatal("expected hash, got empty string")
	}

	if len(hash) != 64 {
		t.Fatalf("expected hash length 64, got %d", len(hash))
	}

	hashAgain := GenerateContentHash(content)

	if hash != hashAgain {
		t.Fatal("expected same content to generate same hash")
	}
}

func TestAddDocument(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

	document, err := service.AddDocument(
		"Machine Learning",
		"Machine learning adalah cabang dari artificial intelligence.",
		nil,
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if document == nil {
		t.Fatal("expected document, got nil")
	}

	if document.Title != "Machine Learning" {
		t.Fatalf("expected title %q, got %q", "Machine Learning", document.Title)
	}

	if document.ContentHash == "" {
		t.Fatal("expected content hash, got empty string")
	}

	if len(repository.documents) != 1 {
		t.Fatalf(
			"expected 1 document in repository, got %d",
			len(repository.documents),
		)
	}
}

func TestAddDocumentDuplicate(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

	content := "Machine learning adalah cabang dari artificial intelligence."

	_, err := service.AddDocument(
		"Machine Learning",
		content,
		nil,
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf("expected first insert to succeed, got %v", err)
	}

	_, err = service.AddDocument(
		"Machine Learning Duplicate",
		content,
		nil,
		nil,
		nil,
	)

	if !errors.Is(err, ErrDuplicateDocument) {
		t.Fatalf(
			"expected ErrDuplicateDocument, got %v",
			err,
		)
	}

	if len(repository.documents) != 1 {
		t.Fatalf(
			"expected repository to contain 1 document, got %d",
			len(repository.documents),
		)
	}
}

func TestAddDocumentInvalid(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

	_, err := service.AddDocument(
		"",
		"",
		nil,
		nil,
		nil,
	)

	if !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf(
			"expected ErrInvalidDocument, got %v",
			err,
		)
	}

	if len(repository.documents) != 0 {
		t.Fatalf(
			"expected repository to remain empty, got %d documents",
			len(repository.documents),
		)
	}
}
