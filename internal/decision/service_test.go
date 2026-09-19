package decision

import (
	"testing"

	"github.com/Rakaa503/AviGo/internal/context"
)

func TestDecisionEngine(t *testing.T) {
	service := NewService()

	tests := []struct {
		name   string
		intent string
		want   Action
	}{
		{
			name:   "greeting",
			intent: "greeting",
			want:   ActionGreeting,
		},
		{
			name:   "question",
			intent: "question",
			want:   ActionAnswerQuestion,
		},
		{
			name:   "problem solving",
			intent: "problem_solving",
			want:   ActionSolveProblem,
		},
		{
			name:   "request",
			intent: "request",
			want:   ActionExecuteRequest,
		},
		{
			name:   "general",
			intent: "general",
			want:   ActionGeneralConversation,
		},
		{
			name:   "unknown",
			intent: "unknown",
			want:   ActionClarify,
		},
	}

	ctx := &context.MessageContext{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.Decide(
				tt.intent,
				"neutral",
				"general",
				0.9,
				ctx,
			)

			if result.Action != tt.want {
				t.Fatalf(
					"expected action %q, got %q",
					tt.want,
					result.Action,
				)
			}
		})
	}
}

func TestDecisionEngineActivityMemoryQuestion(t *testing.T) {
	service := NewService()

	ctx := &context.MessageContext{
		RecentMessages: []context.MessageSnapshot{
			{
				Role:    "user",
				Content: "Saya sedang belajar machine learning",
			},
			{
				Role:    "assistant",
				Content: "Baik.",
			},
			{
				Role:    "user",
				Content: "Apa yang sedang saya pelajari?",
			},
		},
	}

	result := service.Decide(
		"general",
		"neutral",
		"general",
		0.30,
		ctx,
	)

	if result == nil {
		t.Fatal("expected decision, got nil")
	}

	if result.Action != ActionAnswerQuestion {
		t.Fatalf(
			"expected action %q, got %q",
			ActionAnswerQuestion,
			result.Action,
		)
	}
}
