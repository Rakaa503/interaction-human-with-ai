package response

import (
	"fmt"
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(
	action string,
	confidence float64,
	ctx *appcontext.MessageContext,
) (*Response, error) {

	if strings.TrimSpace(action) == "" {
		return nil, fmt.Errorf("response action cannot be empty")
	}

	content := s.generateContent(action, ctx)

	return &Response{
		Action:     action,
		Content:    content,
		Confidence: confidence,
	}, nil
}

func (s *Service) generateContent(
	action string,
	ctx *appcontext.MessageContext,
) string {

	switch action {

	case ActionGreeting:
		return "Halo! 👋 Ada yang bisa AVIGO bantu?"

	case ActionAnswerQuestion:
		if name := extractNameFromContext(ctx); name != "" {
			return fmt.Sprintf("Nama kamu %s.", name)
		}

		return "Tentu, saya akan membantu menjawab pertanyaan kamu."

	case ActionSolveProblem:
		return "Baik, saya akan membantu menganalisis dan menyelesaikan masalah kamu."

	case ActionExecuteRequest:
		return "Siap, saya akan membantu mengerjakan permintaan kamu."

	case ActionGeneralConversation:
		return "Baik, mari kita lanjutkan percakapannya."

	case ActionClarify:
		return "Boleh jelaskan lebih detail supaya saya bisa membantu dengan tepat?"

	default:
		return "Saya belum memahami tindakan yang harus dilakukan."
	}
}

// =========================================================
// MEMORY EXTRACTION
// =========================================================

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

// =========================================================
// NAME EXTRACTION
// =========================================================

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

// =========================================================
// NAME PART
// =========================================================

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

// =========================================================
// CLEAN NAME
// =========================================================

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

// =========================================================
// BLOCKED WORD
// =========================================================

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
