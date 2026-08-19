package response

import (
	"testing"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func TestResponseService(t *testing.T) {
	service := NewService()

	tests := []struct {
		name   string
		action string
		want   string
	}{
		{
			name:   "greeting",
			action: ActionGreeting,
			want:   "Halo! 👋 Ada yang bisa AVIGO bantu?",
		},
		{
			name:   "question",
			action: ActionAnswerQuestion,
			want:   "Tentu, saya akan membantu menjawab pertanyaan kamu.",
		},
		{
			name:   "problem solving",
			action: ActionSolveProblem,
			want:   "Baik, saya akan membantu menganalisis dan menyelesaikan masalah kamu.",
		},
		{
			name:   "request",
			action: ActionExecuteRequest,
			want:   "Siap, saya akan membantu mengerjakan permintaan kamu.",
		},
		{
			name:   "general",
			action: ActionGeneralConversation,
			want:   "Baik, mari kita lanjutkan percakapannya.",
		},
		{
			name:   "clarify",
			action: ActionClarify,
			want:   "Boleh jelaskan lebih detail supaya saya bisa membantu dengan tepat?",
		},
	}

	ctx := &appcontext.MessageContext{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Generate(
				tt.action,
				0.90,
				ctx,
			)

			if err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}

			if result.Action != tt.action {
				t.Fatalf(
					"expected action %q, got %q",
					tt.action,
					result.Action,
				)
			}

			if result.Content != tt.want {
				t.Fatalf(
					"expected content %q, got %q",
					tt.want,
					result.Content,
				)
			}

			if result.Confidence != 0.90 {
				t.Fatalf(
					"expected confidence 0.90, got %f",
					result.Confidence,
				)
			}
		})
	}
}

func TestResponseServiceEmptyAction(t *testing.T) {
	service := NewService()

	ctx := &appcontext.MessageContext{}

	_, err := service.Generate(
		"",
		0.90,
		ctx,
	)

	if err == nil {
		t.Fatal("expected error for empty action")
	}
}

func TestResponseServiceMemory(t *testing.T) {
	service := NewService()

	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:    "user",
				Content: "Nama saya Rakha",
			},
		},
	}

	result, err := service.Generate(
		ActionAnswerQuestion,
		0.90,
		ctx,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	expected := "Nama kamu Nama saya Rakha."

	if result.Content == expected {
		t.Fatalf(
			"memory extraction returned incorrect duplicated name: %q",
			result.Content,
		)
	}

	if result.Content != "Nama kamu Rakha." {
		t.Fatalf(
			"expected memory response %q, got %q",
			"Nama kamu Rakha.",
			result.Content,
		)
	}
}
