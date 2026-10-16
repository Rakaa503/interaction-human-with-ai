package formula

import (
	"math"
	"testing"
)

func TestSolveCosineSimilarity(t *testing.T) {
	result, err := SolveCosineSimilarity(
		[]float64{1, 2},
		[]float64{2, 1},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if math.Abs(result.DotProduct-4) > 1e-9 {
		t.Fatalf("dot product: got %v, want 4", result.DotProduct)
	}

	if math.Abs(result.Similarity-0.8) > 1e-9 {
		t.Fatalf("similarity: got %v, want 0.8", result.Similarity)
	}

	if len(result.Steps) == 0 {
		t.Fatal("expected calculation steps")
	}
}

func TestSolveCosineSimilarityRejectsInvalidVectors(t *testing.T) {
	tests := []struct {
		name string
		a    []float64
		b    []float64
	}{
		{"empty vectors", nil, nil},
		{"different dimensions", []float64{1, 2}, []float64{1}},
		{"zero vector", []float64{0, 0}, []float64{1, 2}},
		{"infinite value", []float64{math.Inf(1)}, []float64{1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := SolveCosineSimilarity(test.a, test.b)
			if err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
