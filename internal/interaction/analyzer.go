package interaction

import "strings"

type AnalysisResult struct {
	Intent     string
	Emotion    string
	Topic      string
	Confidence float64
}

type Analyzer interface {
	Analyze(input string) (*AnalysisResult, error)
}

type RuleBasedAnalyzer struct{}

func NewRuleBasedAnalyzer() *RuleBasedAnalyzer {
	return &RuleBasedAnalyzer{}
}

func (a *RuleBasedAnalyzer) Analyze(
	input string,
) (*AnalysisResult, error) {

	text := strings.ToLower(input)

	result := &AnalysisResult{
		Intent:     "general",
		Emotion:    "neutral",
		Topic:      "general",
		Confidence: 0.50,
	}

	// Technical problem has higher priority
	// than a general greeting.
	if strings.Contains(text, "error") ||
		strings.Contains(text, "bug") {

		result.Intent = "problem_solving"
		result.Topic = "technical"
		result.Confidence = 0.80

		return result, nil
	}

	// Greeting
	if strings.Contains(text, "halo") ||
		strings.Contains(text, "hai") {

		result.Intent = "greeting"
		result.Confidence = 0.90

		return result, nil
	}

	return result, nil
}
