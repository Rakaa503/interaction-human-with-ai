package orchestrator

import (
	"fmt"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
	"github.com/Rakaa503/AviGo/internal/decision"
	"github.com/Rakaa503/AviGo/internal/interaction"
	"github.com/Rakaa503/AviGo/internal/response"
)

type Service struct {
	interactionAnalyzer interaction.Analyzer
	contextService      *appcontext.Service
	decisionService     *decision.Service
	responseService     *response.Service
}

func NewService(
	interactionAnalyzer interaction.Analyzer,
	contextService *appcontext.Service,
	decisionService *decision.Service,
	responseService *response.Service,
) *Service {
	return &Service{
		interactionAnalyzer: interactionAnalyzer,
		contextService:      contextService,
		decisionService:     decisionService,
		responseService:     responseService,
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
	// 1. Interaction Analysis
	// =========================

	analysis, err := s.interactionAnalyzer.Analyze(input)
	if err != nil {
		return nil, fmt.Errorf(
			"interaction analysis failed: %w",
			err,
		)
	}

	// =========================
	// 2. Build Context
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
	// 3. Decision Engine
	// =========================

	decisionResult := s.decisionService.Decide(
		analysis.Intent,
		analysis.Emotion,
		analysis.Topic,
		analysis.Confidence,
		ctx,
	)

	// =========================
	// 4. Response Engine
	// =========================

	responseResult, err := s.responseService.Generate(
		string(decisionResult.Action),
		decisionResult.Confidence,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"response generation failed: %w",
			err,
		)
	}

	// =========================
	// 5. Final Result
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
