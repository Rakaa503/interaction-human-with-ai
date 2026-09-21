package memory

import (
	"testing"
	"time"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func TestMemoryServiceRetrieveName(t *testing.T) {
	service := NewService()

	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:      "user",
				Content:   "Nama saya Rakha",
				CreatedAt: time.Now(),
			},
		},
	}

	result := service.Retrieve(ctx, MemoryName)

	if result == nil {
		t.Fatal("expected name memory, got nil")
	}

	if result.Type != MemoryName {
		t.Fatalf(
			"expected memory type %q, got %q",
			MemoryName,
			result.Type,
		)
	}

	if result.Value != "Rakha" {
		t.Fatalf(
			"expected name %q, got %q",
			"Rakha",
			result.Value,
		)
	}

	if result.Confidence != 0.95 {
		t.Fatalf(
			"expected confidence %.2f, got %.2f",
			0.95,
			result.Confidence,
		)
	}
}

func TestMemoryServiceRetrieveActivity(t *testing.T) {
	service := NewService()

	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:      "user",
				Content:   "Saya sedang belajar machine learning",
				CreatedAt: time.Now(),
			},
		},
	}

	result := service.Retrieve(ctx, MemoryActivity)

	if result == nil {
		t.Fatal("expected activity memory, got nil")
	}

	if result.Type != MemoryActivity {
		t.Fatalf(
			"expected memory type %q, got %q",
			MemoryActivity,
			result.Type,
		)
	}

	if result.Value != "machine learning" {
		t.Fatalf(
			"expected activity %q, got %q",
			"machine learning",
			result.Value,
		)
	}

	if result.Confidence != 0.95 {
		t.Fatalf(
			"expected confidence %.2f, got %.2f",
			0.95,
			result.Confidence,
		)
	}
}

func TestMemoryServiceRetrieveUnknownType(t *testing.T) {
	service := NewService()

	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:      "user",
				Content:   "Nama saya Rakha",
				CreatedAt: time.Now(),
			},
		},
	}

	result := service.Retrieve(
		ctx,
		MemoryType("unknown"),
	)

	if result != nil {
		t.Fatalf(
			"expected nil for unknown memory type, got %+v",
			result,
		)
	}
}

func TestMemoryServiceRetrieveNilContext(t *testing.T) {
	service := NewService()

	result := service.Retrieve(nil, MemoryName)

	if result != nil {
		t.Fatalf(
			"expected nil for nil context, got %+v",
			result,
		)
	}
}
