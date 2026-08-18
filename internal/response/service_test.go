package response

import "testing"

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Generate(
				tt.action,
				0.95,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
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
		})
	}
}

func TestResponseServiceEmptyAction(t *testing.T) {
	service := NewService()

	_, err := service.Generate("", 0.9)

	if err == nil {
		t.Fatal("expected error for empty action")
	}
}
