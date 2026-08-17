package interaction

import "testing"

func TestRuleBasedAnalyzer(t *testing.T) {
	analyzer := NewRuleBasedAnalyzer()

	tests := []struct {
		name           string
		input          string
		expectedIntent string
		expectedTopic  string
	}{
		{
			name:           "technical bug",
			input:          "Saya menemukan bug di aplikasi",
			expectedIntent: "problem_solving",
			expectedTopic:  "technical",
		},
		{
			name:           "technical error",
			input:          "Ada error saat login",
			expectedIntent: "problem_solving",
			expectedTopic:  "technical",
		},
		{
			name:           "greeting halo",
			input:          "Halo AVIGO",
			expectedIntent: "greeting",
			expectedTopic:  "general",
		},
		{
			name:           "greeting hai",
			input:          "Hai AVIGO",
			expectedIntent: "greeting",
			expectedTopic:  "general",
		},
		{
			name:           "general input",
			input:          "Saya ingin berbicara dengan AVIGO",
			expectedIntent: "general",
			expectedTopic:  "general",
		},
		{
			name:           "bug with greeting",
			input:          "Halo AVIGO, ada bug di aplikasi saya",
			expectedIntent: "problem_solving",
			expectedTopic:  "technical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzer.Analyze(tt.input)

			if err != nil {
				t.Fatalf(
					"Analyze() returned error: %v",
					err,
				)
			}

			if result.Intent != tt.expectedIntent {
				t.Errorf(
					"expected intent %q, got %q",
					tt.expectedIntent,
					result.Intent,
				)
			}

			if result.Topic != tt.expectedTopic {
				t.Errorf(
					"expected topic %q, got %q",
					tt.expectedTopic,
					result.Topic,
				)
			}

			if result.Confidence <= 0 {
				t.Errorf(
					"expected confidence > 0, got %f",
					result.Confidence,
				)
			}
		})
	}
}
