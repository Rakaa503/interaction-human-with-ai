package response

type Response struct {
	Action     string  `json:"action"`
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"`
}

const (
	ActionGreeting            = "respond_greeting"
	ActionAnswerQuestion      = "answer_question"
	ActionSolveProblem        = "solve_problem"
	ActionExecuteRequest      = "execute_request"
	ActionGeneralConversation = "general_conversation"
	ActionClarify             = "ask_clarification"
)
