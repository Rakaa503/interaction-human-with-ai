package interaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MLAnalyzer struct {
	baseURL string
	client  *http.Client
}

type mlPredictionRequest struct {
	Input string `json:"input"`
}

type mlPredictionResponse struct {
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`
}

func NewMLAnalyzer(baseURL string) *MLAnalyzer {
	return &MLAnalyzer{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (a *MLAnalyzer) Analyze(
	input string,
) (*AnalysisResult, error) {

	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("input tidak boleh kosong")
	}

	requestBody := mlPredictionRequest{
		Input: input,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to encode ML request: %w",
			err,
		)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		a.baseURL+"/predict",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create ML request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to ML service: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		return nil, fmt.Errorf(
			"ML service returned status %d",
			resp.StatusCode,
		)
	}

	var result mlPredictionResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&result); err != nil {

		return nil, fmt.Errorf(
			"failed to decode ML response: %w",
			err,
		)
	}

	return &AnalysisResult{
		Intent:     result.Intent,
		Emotion:    "neutral",
		Topic:      inferTopic(result.Intent),
		Confidence: result.Confidence,
	}, nil
}

func inferTopic(intent string) string {
	switch intent {
	case "problem_solving":
		return "technical"

	case "question":
		return "general"

	case "request":
		return "general"

	default:
		return "general"
	}
}
