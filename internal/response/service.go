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

	if action == "" {
		return nil, fmt.Errorf(
			"response action cannot be empty",
		)
	}

	content := s.generateContent(
		action,
		ctx,
	)

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

		if answer := extractNameFromContext(ctx); answer != "" {
			return fmt.Sprintf(
				"Nama kamu %s.",
				answer,
			)
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

func extractNameFromContext(
	ctx *appcontext.MessageContext,
) string {

	if ctx == nil {
		return ""
	}

	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {

		message := ctx.RecentMessages[i]

		if message.Role != "user" {
			continue
		}

		text := strings.TrimSpace(message.Content)
		lower := strings.ToLower(text)

		prefixes := []string{
			"nama saya ",
			"nama aku ",
			"saya bernama ",
			"aku bernama ",
		}

		for _, prefix := range prefixes {

			if strings.HasPrefix(lower, prefix) {

				name := strings.TrimSpace(
					text[len(prefix):],
				)

				if name != "" {
					return name
				}
			}
		}
	}

	return ""
}
