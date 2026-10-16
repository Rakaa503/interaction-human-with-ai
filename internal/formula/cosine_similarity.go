package formula

import (
	"fmt"
	"math"
)

type CosineSimilarityResult struct {
	VectorA    []float64 `json:"vector_a"`
	VectorB    []float64 `json:"vector_b"`
	DotProduct float64   `json:"dot_product"`
	MagnitudeA float64   `json:"magnitude_a"`
	MagnitudeB float64   `json:"magnitude_b"`
	Similarity float64   `json:"similarity"`
	Steps      []string  `json:"steps"`
}

func SolveCosineSimilarity(
	vectorA []float64,
	vectorB []float64,
) (*CosineSimilarityResult, error) {
	if len(vectorA) == 0 || len(vectorB) == 0 {
		return nil, fmt.Errorf("vektor tidak boleh kosong")
	}

	if len(vectorA) != len(vectorB) {
		return nil, fmt.Errorf(
			"dimensi vektor harus sama: A=%d, B=%d",
			len(vectorA),
			len(vectorB),
		)
	}

	var dotProduct, squaredA, squaredB float64

	for i := range vectorA {
		a, b := vectorA[i], vectorB[i]

		if math.IsNaN(a) || math.IsInf(a, 0) ||
			math.IsNaN(b) || math.IsInf(b, 0) {
			return nil, fmt.Errorf(
				"elemen vektor harus berupa angka finite",
			)
		}

		dotProduct += a * b
		squaredA += a * a
		squaredB += b * b
	}

	magnitudeA := math.Sqrt(squaredA)
	magnitudeB := math.Sqrt(squaredB)

	if magnitudeA == 0 || magnitudeB == 0 {
		return nil, fmt.Errorf(
			"cosine similarity tidak terdefinisi untuk vektor nol",
		)
	}

	similarity := dotProduct / (magnitudeA * magnitudeB)

	steps := []string{
		fmt.Sprintf("Diketahui A = %v dan B = %v.", vectorA, vectorB),
		fmt.Sprintf(
			"Dot product = %g.",
			dotProduct,
		),
		fmt.Sprintf(
			"Magnitudo A = sqrt(sum(A_i^2)) = %.6f.",
			magnitudeA,
		),
		fmt.Sprintf(
			"Magnitudo B = sqrt(sum(B_i^2)) = %.6f.",
			magnitudeB,
		),
		fmt.Sprintf(
			"Cosine similarity = %.6f / (%.6f * %.6f) = %.6f.",
			dotProduct,
			magnitudeA,
			magnitudeB,
			similarity,
		),
	}

	return &CosineSimilarityResult{
		VectorA:    append([]float64(nil), vectorA...),
		VectorB:    append([]float64(nil), vectorB...),
		DotProduct: dotProduct,
		MagnitudeA: magnitudeA,
		MagnitudeB: magnitudeB,
		Similarity: similarity,
		Steps:      steps,
	}, nil
}
