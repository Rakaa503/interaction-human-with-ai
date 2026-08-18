package decision

type Action string

const (
	ActionGreeting            Action = "respond_greeting"
	ActionAnswerQuestion      Action = "answer_question"
	ActionSolveProblem        Action = "solve_problem"
	ActionExecuteRequest      Action = "execute_request"
	ActionGeneralConversation Action = "general_conversation"
	ActionClarify             Action = "ask_clarification"
)

type Decision struct {
	Action     Action
	Reason     string
	Confidence float64
}
