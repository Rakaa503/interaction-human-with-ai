package orchestrator

import (
	"fmt"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
	"github.com/Rakaa503/AviGo/internal/conversation"
	"github.com/Rakaa503/AviGo/internal/decision"
	"github.com/Rakaa503/AviGo/internal/interaction"
	"github.com/Rakaa503/AviGo/internal/response"
)

type Service struct {
	interactionService  *interaction.Service
	contextService      *appcontext.Service
	decisionService     *decision.Service
	responseService     *response.Service
	conversationService *conversation.Service
}

func NewService(
	interactionService *interaction.Service,
	contextService *appcontext.Service,
	decisionService *decision.Service,
	responseService *response.Service,
	conversationService *conversation.Service,
) *Service {
	return &Service{
		interactionService:  interactionService,
		contextService:      contextService,
		decisionService:     decisionService,
		responseService:     responseService,
		conversationService: conversationService,
	}
}

type Result struct {
	Input      string  `json:"input"`
	Intent     string  `json:"intent"`
	Emotion    string  `json:"emotion"`
	Topic      string  `json:"topic"`
	Confidence float64 `json:"confidence"`
	Action     string  `json:"action"`
	Response   string  `json:"response"`
}

func (s *Service) Process(
	userID uint64,
	conversationID uint64,
	input string,
) (*Result, error) {

	if input == "" {
		return nil, fmt.Errorf("input cannot be empty")
	}

	// =========================
	// 1. Save User Message
	// =========================

	if _, err := s.conversationService.AddMessage(
		conversationID,
		"user",
		input,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to save user message: %w",
			err,
		)
	}

	// =========================
	// 2. Interaction Analysis
	// =========================

	analysis, err := s.interactionService.Analyze(input)
	if err != nil {
		return nil, fmt.Errorf(
			"interaction analysis failed: %w",
			err,
		)
	}

	// =========================
	// 3. Build Context
	// =========================

	ctx, err := s.contextService.Build(
		conversationID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"context building failed: %w",
			err,
		)
	}

	// =========================
	// 4. Decision Engine
	// =========================

	decisionResult := s.decisionService.Decide(
		analysis.Intent,
		analysis.Emotion,
		analysis.Topic,
		analysis.Confidence,
		ctx,
	)

	// =========================
	// 5. Response Engine
	// =========================

	responseResult, err := s.responseService.Generate(
		string(decisionResult.Action),
		decisionResult.Confidence,
		ctx,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"response generation failed: %w",
			err,
		)
	}

	// =========================
	// 6. Save Assistant Message
	// =========================

	if _, err := s.conversationService.AddMessage(
		conversationID,
		"assistant",
		responseResult.Content,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to save assistant message: %w",
			err,
		)
	}

	// =========================
	// 7. Save Interaction
	// =========================

	if _, err := s.interactionService.Save(
		userID,
		conversationID,
		input,
		analysis,
		responseResult.Content,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to save interaction: %w",
			err,
		)
	}

	// =========================
	// 8. Final Result
	// =========================

	return &Result{
		Input:      input,
		Intent:     analysis.Intent,
		Emotion:    analysis.Emotion,
		Topic:      analysis.Topic,
		Confidence: analysis.Confidence,
		Action:     string(decisionResult.Action),
		Response:   responseResult.Content,
	}, nil
}
