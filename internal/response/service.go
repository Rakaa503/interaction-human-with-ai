package response

import "fmt"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(
	action string,
	confidence float64,
) (*Response, error) {

	if action == "" {
		return nil, fmt.Errorf("response action cannot be empty")
	}

	var content string

	switch action {
	case ActionGreeting:
		content = "Halo! 👋 Ada yang bisa AVIGO bantu?"

	case ActionAnswerQuestion:
		content = "Tentu, saya akan membantu menjawab pertanyaan kamu."

	case ActionSolveProblem:
		content = "Baik, saya akan membantu menganalisis dan menyelesaikan masalah kamu."

	case ActionExecuteRequest:
		content = "Siap, saya akan membantu mengerjakan permintaan kamu."

	case ActionGeneralConversation:
		content = "Baik, mari kita lanjutkan percakapannya."

	case ActionClarify:
		content = "Boleh jelaskan lebih detail supaya saya bisa membantu dengan tepat?"

	default:
		content = "Saya belum memahami tindakan yang harus dilakukan."
	}

	return &Response{
		Action:     action,
		Content:    content,
		Confidence: confidence,
	}, nil
}
