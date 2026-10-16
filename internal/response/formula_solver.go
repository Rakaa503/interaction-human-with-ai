package response

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
	"github.com/Rakaa503/AviGo/internal/formula"
)

var (
	cosineRequestPattern = regexp.MustCompile(
		`(?i)cosine\s*similarity|cos\s*similarity|kemiripan\s*cosine`,
	)

	vectorPattern = regexp.MustCompile(
		`[\[(]\s*([0-9eE+\-.,\s]+)\s*[\])]`,
	)
)

func solveCosineSimilarityFromContext(
	ctx *appcontext.MessageContext,
) (string, bool) {
	if ctx == nil {
		return "", false
	}

	// Periksa pesan pengguna terbaru.
	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {
		message := ctx.RecentMessages[i]

		if !strings.EqualFold(
			strings.TrimSpace(message.Role),
			"user",
		) {
			continue
		}

		input := strings.TrimSpace(message.Content)

		// Jangan mengambil pertanyaan lama jika pesan terbaru
		// bukan permintaan cosine similarity.
		if !cosineRequestPattern.MatchString(input) {
			return "", false
		}

		matches := vectorPattern.FindAllStringSubmatch(input, -1)
		if len(matches) != 2 {
			return "", false
		}

		vectorA, err := parseVector(matches[0][1])
		if err != nil {
			return "", false
		}

		vectorB, err := parseVector(matches[1][1])
		if err != nil {
			return "", false
		}

		result, err := formula.SolveCosineSimilarity(
			vectorA,
			vectorB,
		)
		if err != nil {
			return fmt.Sprintf(
				"AVIGO belum bisa menghitung cosine similarity: %s.",
				err.Error(),
			), true
		}

		return formatCosineSolution(result), true
	}

	return "", false
}

func parseVector(raw string) ([]float64, error) {
	parts := strings.Split(raw, ",")
	if len(parts) == 0 {
		return nil, fmt.Errorf("vektor kosong")
	}

	vector := make([]float64, 0, len(parts))

	for _, part := range parts {
		value, err := strconv.ParseFloat(
			strings.TrimSpace(part),
			64,
		)
		if err != nil {
			return nil, fmt.Errorf("elemen vektor tidak valid")
		}

		vector = append(vector, value)
	}

	return vector, nil
}

func formatCosineSolution(
	result *formula.CosineSimilarityResult,
) string {
	var builder strings.Builder

	builder.WriteString("Penyelesaian Cosine Similarity\n\n")

	builder.WriteString("1. Vektor yang diketahui\n")
	builder.WriteString(fmt.Sprintf(
		"A = %v dan B = %v.\n\n",
		result.VectorA,
		result.VectorB,
	))

	builder.WriteString("2. Rumus\n")
	builder.WriteString(
		"cos(A, B) = dot(A, B) / (norm(A) * norm(B))\n\n",
	)

	builder.WriteString("3. Hitung dot product\n")

	terms := make([]string, len(result.VectorA))
	for i := range result.VectorA {
		terms[i] = fmt.Sprintf(
			"(%g * %g)",
			result.VectorA[i],
			result.VectorB[i],
		)
	}

	builder.WriteString(strings.Join(terms, " + "))
	builder.WriteString(fmt.Sprintf(
		" = %g\n\n",
		result.DotProduct,
	))

	builder.WriteString("4. Hitung magnitudo vektor\n")
	builder.WriteString(fmt.Sprintf(
		"norm(A) = sqrt(sum(A_i^2)) = %.6f\n",
		result.MagnitudeA,
	))
	builder.WriteString(fmt.Sprintf(
		"norm(B) = sqrt(sum(B_i^2)) = %.6f\n\n",
		result.MagnitudeB,
	))

	builder.WriteString("5. Substitusi ke rumus\n")
	builder.WriteString(fmt.Sprintf(
		"cos(A, B) = %g / (%.6f * %.6f)\n",
		result.DotProduct,
		result.MagnitudeA,
		result.MagnitudeB,
	))
	builder.WriteString(fmt.Sprintf(
		"Hasil akhir: %.6f\n\n",
		result.Similarity,
	))

	builder.WriteString(
		"Interpretasi: nilai ini mengukur kemiripan arah kedua " +
			"vektor, bukan probabilitas atau persentase kebenaran.",
	)

	return builder.String()
}
