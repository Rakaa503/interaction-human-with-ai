package memory

import (
	"testing"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func TestExtractActivity(t *testing.T) {
	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:    "user",
				Content: "Saya sedang belajar machine learning",
			},
		},
	}

	result := ExtractActivity(ctx)

	if result == nil {
		t.Fatal("expected activity memory, got nil")
	}

	if result.Type != MemoryActivity {
		t.Fatalf("expected type %q, got %q", MemoryActivity, result.Type)
	}

	if result.Value != "machine learning" {
		t.Fatalf("expected value %q, got %q", "machine learning", result.Value)
	}

	if result.Confidence != 0.95 {
		t.Fatalf("expected confidence 0.95, got %f", result.Confidence)
	}
}

func TestExtractActivityUsesLatestUserMessage(t *testing.T) {
	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:    "user",
				Content: "Saya sedang belajar Python",
			},
			{
				Role:    "assistant",
				Content: "Baik.",
			},
			{
				Role:    "user",
				Content: "Saya sedang belajar machine learning",
			},
		},
	}

	result := ExtractActivity(ctx)

	if result == nil {
		t.Fatal("expected activity memory, got nil")
	}

	if result.Value != "machine learning" {
		t.Fatalf(
			"expected latest activity %q, got %q",
			"machine learning",
			result.Value,
		)
	}
}

func TestExtractActivityIgnoresAssistantMessage(t *testing.T) {
	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:    "assistant",
				Content: "Saya sedang belajar machine learning",
			},
		},
	}

	result := ExtractActivity(ctx)

	if result != nil {
		t.Fatalf("expected nil, got %+v", result)
	}
}

func TestExtractActivityReturnsNilWhenNoActivity(t *testing.T) {
	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:    "user",
				Content: "Halo AVIGO",
			},
		},
	}

	result := ExtractActivity(ctx)

	if result != nil {
		t.Fatalf("expected nil, got %+v", result)
	}
}
