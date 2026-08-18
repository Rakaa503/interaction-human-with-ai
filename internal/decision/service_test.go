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
