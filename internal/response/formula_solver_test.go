package response

import (
	"strings"
	"testing"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func TestSolveCosineSimilarityFromContext(t *testing.T) {
	ctx := &appcontext.MessageContext{
		RecentMessages: []appcontext.MessageSnapshot{
			{
				Role:    "user",
				Content: "Hitung cosine similarity antara A=(1,2) dan B=(2,1). Jelaskan langkah perhitungannya.",
			},
		},
	}

	answer, handled := solveCosineSimilarityFromContext(ctx)

	if !handled {
		t.Fatal("solver tidak menangani permintaan cosine similarity")
	}

	if !strings.Contains(answer, "0.800000") {
		t.Fatalf("hasil cosine similarity tidak ditemukan: %s", answer)
	}

	t.Log(answer)
}
