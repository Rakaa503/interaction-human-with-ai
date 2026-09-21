package response

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rakaa503/AviGo/internal/ai"
	appcontext "github.com/Rakaa503/AviGo/internal/context"
	"github.com/Rakaa503/AviGo/internal/memory"
)

type Service struct {
	aiService *ai.Service
}

func NewService(aiServices ...*ai.Service) *Service {
	var aiService *ai.Service

	if len(aiServices) > 0 {
		aiService = aiServices[0]
	}

	return &Service{
		aiService: aiService,
	}
}

func (s *Service) Generate(
	action string,
	confidence float64,
	ctx *appcontext.MessageContext,
) (*Response, error) {

	if strings.TrimSpace(action) == "" {
		return nil, fmt.Errorf("response action cannot be empty")
	}

	content, err := s.generateContent(action, ctx)
	if err != nil {
		return nil, err
	}

	return &Response{
		Action:     action,
		Content:    content,
		Confidence: confidence,
	}, nil
}

func (s *Service) generateContent(
	action string,
	ctx *appcontext.MessageContext,
) (string, error) {

	// =========================================================
	// MEMORY HAS PRIORITY
	// =========================================================

	if action == ActionAnswerQuestion {

		// NAME MEMORY
		if name := extractNameFromContext(ctx); name != "" {
			return fmt.Sprintf("Nama kamu %s.", name), nil
		}

		// ACTIVITY MEMORY
		if activity := memory.ExtractActivity(ctx); activity != nil {
			return fmt.Sprintf(
				"Kamu sedang belajar %s.",
				activity.Value,
			), nil
		}
	}

	// =========================================================
	// AI REASONING
	// =========================================================

	if s.aiService != nil {
		input := buildAIInput(action, ctx)

		if strings.TrimSpace(input) != "" {
			aiResponse, err := s.aiService.Generate(
				context.Background(),
				input,
			)

			if err == nil &&
				aiResponse != nil &&
				strings.TrimSpace(aiResponse.Content) != "" {

				return aiResponse.Content, nil
			}
		}
	}

	// =========================================================
	// FALLBACK
	// =========================================================

	switch action {
	case ActionGreeting:
		return "Halo! 👋 Ada yang bisa AVIGO bantu?", nil

	case ActionAnswerQuestion:
		return "Tentu, saya akan membantu menjawab pertanyaan kamu.", nil

	case ActionSolveProblem:
		return "Baik, saya akan membantu menganalisis dan menyelesaikan masalah kamu.", nil

	case ActionExecuteRequest:
		return "Siap, saya akan membantu mengerjakan permintaan kamu.", nil

	case ActionGeneralConversation:
		return "Baik, mari kita lanjutkan percakapannya.", nil

	case ActionClarify:
		return "Boleh jelaskan lebih detail supaya saya bisa membantu dengan tepat?", nil

	default:
		return "Saya belum memahami tindakan yang harus dilakukan.", nil
	}
}

func buildAIInput(
	action string,
	ctx *appcontext.MessageContext,
) string {

	if ctx == nil || len(ctx.RecentMessages) == 0 {
		return ""
	}

	var builder strings.Builder

	builder.WriteString("You are AVIGO, an AI interaction assistant.\n")
	builder.WriteString("Action: ")
	builder.WriteString(action)
	builder.WriteString("\n\nConversation:\n")

	for _, message := range ctx.RecentMessages {
		role := strings.TrimSpace(message.Role)
		content := strings.TrimSpace(message.Content)

		if role == "" || content == "" {
			continue
		}

		builder.WriteString(role)
		builder.WriteString(": ")
		builder.WriteString(content)
		builder.WriteString("\n")
	}

	builder.WriteString(
		"\nGenerate a helpful response to the user's latest message.",
	)

	return builder.String()
}

func extractNameFromContext(
	ctx *appcontext.MessageContext,
) string {

	if ctx == nil || len(ctx.RecentMessages) == 0 {
		return ""
	}

	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {
		message := ctx.RecentMessages[i]

		if !strings.EqualFold(
			strings.TrimSpace(message.Role),
			"user",
		) {
			continue
		}

		text := strings.TrimSpace(message.Content)

		if text == "" {
			continue
		}

		if name := extractName(text); name != "" {
			return name
		}
	}

	return ""
}

func extractName(text string) string {

	original := strings.TrimSpace(text)

	if original == "" {
		return ""
	}

	lower := strings.ToLower(original)

	prefixes := []string{
		"nama saya ",
		"nama aku ",
		"nama gue ",
		"nama gua ",
		"nama gw ",
		"saya bernama ",
		"aku bernama ",
		"gue bernama ",
		"gua bernama ",
		"gw bernama ",
	}

	for _, prefix := range prefixes {
		if !strings.HasPrefix(lower, prefix) {
			continue
		}

		name := strings.TrimSpace(
			original[len(prefix):],
		)

		return extractNamePart(name)
	}

	shortPrefixes := []string{
		"saya ",
		"aku ",
		"gue ",
		"gua ",
		"gw ",
	}

	for _, prefix := range shortPrefixes {
		if !strings.HasPrefix(lower, prefix) {
			continue
		}

		name := strings.TrimSpace(
			original[len(prefix):],
		)

		if name == "" {
			return ""
		}

		firstWord := strings.Fields(name)[0]

		if isBlockedWord(firstWord) {
			return ""
		}

		return extractNamePart(name)
	}

	return ""
}

func extractNamePart(text string) string {

	text = strings.TrimSpace(text)

	if text == "" {
		return ""
	}

	lower := strings.ToLower(text)

	stopWords := []string{
		" saya ",
		" aku ",
		" gue ",
		" gua ",
		" gw ",
		" dan ",
		" karena ",
		" yang ",
		" sedang ",
		" lagi ",
		" adalah ",
		" merupakan ",
	}

	cutIndex := len(text)

	for _, word := range stopWords {
		if index := strings.Index(lower, word); index >= 0 {
			if index < cutIndex {
				cutIndex = index
			}
		}
	}

	text = strings.TrimSpace(
		text[:cutIndex],
	)

	return cleanName(text)
}

func cleanName(name string) string {

	name = strings.TrimSpace(name)

	if name == "" {
		return ""
	}

	name = strings.TrimRight(
		name,
		".,!?;:",
	)

	name = strings.TrimSpace(name)

	if name == "" {
		return ""
	}

	words := strings.Fields(name)

	if len(words) > 4 {
		return ""
	}

	return name
}

func isBlockedWord(word string) bool {

	word = strings.ToLower(
		strings.TrimSpace(word),
	)

	switch word {
	case "sedang",
		"lagi",
		"membuat",
		"ingin",
		"mau",
		"adalah",
		"merupakan",
		"tidak",
		"bukan",
		"suka",
		"punya",
		"memiliki":

		return true
	}

	return false
}
