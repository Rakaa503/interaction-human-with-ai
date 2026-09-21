package memory

import (
	"testing"
	"time"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func TestRequestedType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected MemoryType
		found    bool
	}{
		{
			name:     "activity question",
			input:    "Apa yang sedang saya pelajari?",
			expected: MemoryActivity,
			found:    true,
		},
		{
			name:     "name question",
			input:    "Siapa nama saya?",
			expected: MemoryName,
			found:    true,
		},
		{
			name:     "normal conversation",
			input:    "Halo AVIGO",
			expected: "",
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &appcontext.MessageContext{
				RecentMessages: []appcontext.MessageSnapshot{
					{
						Role:      "user",
						Content:   tt.input,
						CreatedAt: time.Now(),
					},
				},
			}

			got, found := RequestedType(ctx)

			if found != tt.found {
				t.Fatalf("expected found=%v, got %v", tt.found, found)
			}

			if got != tt.expected {
				t.Fatalf("expected type=%q, got %q", tt.expected, got)
			}
		})
	}
}