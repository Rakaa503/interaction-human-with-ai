package interaction

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMLAnalyzer(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/predict" {
				t.Fatalf(
					"expected /predict, got %s",
					r.URL.Path,
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte(`{
				"intent": "problem_solving",
				"confidence": 0.5078
			}`))
		}),
	)

	defer server.Close()

	analyzer := NewMLAnalyzer(
		server.URL,
	)

	result, err := analyzer.Analyze(
		"Aplikasi saya mengalami bug saat login",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.Intent != "problem_solving" {
		t.Fatalf(
			"expected problem_solving, got %s",
			result.Intent,
		)
	}

	if result.Confidence != 0.5078 {
		t.Fatalf(
			"expected confidence 0.5078, got %f",
			result.Confidence,
		)
	}

	if result.Topic != "technical" {
		t.Fatalf(
			"expected technical, got %s",
			result.Topic,
		)
	}
}
