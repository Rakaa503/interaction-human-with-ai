package decision

import (
	"strings"

	"github.com/Rakaa503/AviGo/internal/context"
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
	ctx *context.MessageContext,
) *Decision {

	intent := strings.ToLower(strings.TrimSpace(analysisIntent))
	emotion = strings.ToLower(strings.TrimSpace(emotion))
	topic = strings.ToLower(strings.TrimSpace(topic))

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
		if ctx != nil && len(ctx.RecentMessages) > 0 {
			return &Decision{
				Action:     ActionSolveProblem,
				Reason:     "user has a problem and conversation context is available",
				Confidence: confidence,
			}
		}

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
