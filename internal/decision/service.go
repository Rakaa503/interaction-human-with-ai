package decision

import (
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Decide(
	analysisIntent string,
	emotion string,
	topic string,
	confidence float64,
	ctx *appcontext.MessageContext,
) *Decision {

	intent := strings.ToLower(strings.TrimSpace(analysisIntent))

	// Memory question harus diprioritaskan sebelum hasil ML.
	if isMemoryQuestion(ctx) {
		return &Decision{
			Action:     ActionAnswerQuestion,
			Reason:     "user is asking about information from conversation memory",
			Confidence: confidence,
		}
	}

	switch intent {

	case "greeting":
		return &Decision{
			Action:     ActionGreeting,
			Reason:     "user is greeting AVIGO",
			Confidence: confidence,
		}

	case "question":
		return &Decision{
			Action:     ActionAnswerQuestion,
			Reason:     "user is asking a question",
			Confidence: confidence,
		}

	case "problem_solving":
		return &Decision{
			Action:     ActionSolveProblem,
			Reason:     "user reported a problem",
			Confidence: confidence,
		}

	case "request":
		return &Decision{
			Action:     ActionExecuteRequest,
			Reason:     "user is requesting an action",
			Confidence: confidence,
		}

	case "general":
		return &Decision{
			Action:     ActionGeneralConversation,
			Reason:     "input is general conversation",
			Confidence: confidence,
		}

	default:
		return &Decision{
			Action:     ActionClarify,
			Reason:     "intent could not be determined",
			Confidence: confidence,
		}
	}
}

func isMemoryQuestion(
	ctx *appcontext.MessageContext,
) bool {

	if ctx == nil || len(ctx.RecentMessages) == 0 {
		return false
	}

	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {

		message := ctx.RecentMessages[i]

		if !strings.EqualFold(
			strings.TrimSpace(message.Role),
			"user",
		) {
			continue
		}

		text := strings.ToLower(
			strings.TrimSpace(message.Content),
		)

		if text == "" {
			continue
		}

		// Pertanyaan identitas user.
		if containsAny(text,
			"siapa nama saya",
			"siapa nama aku",
			"siapa nama gue",
			"siapa nama gua",
			"siapa nama gw",
			"apa nama saya",
			"apa nama aku",
			"apa nama gue",
			"apa nama gua",
			"apa nama gw",
			"nama saya siapa",
			"nama aku siapa",
			"nama gue siapa",
			"nama gua siapa",
			"nama gw siapa",
		) {
			return true
		}

		// Pertanyaan tentang memory conversation.
		if containsAny(text,
			"apa yang saya bilang",
			"apa yang aku bilang",
			"apa yang gue bilang",
			"apa yang gua bilang",
			"apa yang gw bilang",
			"apa yang saya suka",
			"apa yang aku suka",
			"apa yang gue suka",
			"apa yang gua suka",
			"apa yang gw suka",
			"kamu ingat saya",
			"kamu ingat aku",
			"apakah kamu ingat saya",
			"apakah kamu ingat aku",
		) {
			return true
		}

		// Cukup periksa pesan user terbaru.
		return false
	}

	return false
}

func containsAny(
	text string,
	keywords ...string,
) bool {

	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}
