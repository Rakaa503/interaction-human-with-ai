package response

import (
	"fmt"
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(
	action string,
	confidence float64,
	ctx *appcontext.MessageContext,
) (*Response, error) {

	if action == "" {
		return nil, fmt.Errorf(
			"response action cannot be empty",
		)
	}

	content := s.generateContent(
		action,
		ctx,
	)

	return &Response{
		Action:     action,
		Content:    content,
		Confidence: confidence,
	}, nil
}

func (s *Service) generateContent(
	action string,
	ctx *appcontext.MessageContext,
) string {

	switch action {

	case ActionGreeting:
		return "Halo! 👋 Ada yang bisa AVIGO bantu?"

	case ActionAnswerQuestion:

		if name := extractNameFromContext(ctx); name != "" {
			return fmt.Sprintf(
				"Nama kamu %s.",
				name,
			)
		}

		return "Tentu, saya akan membantu menjawab pertanyaan kamu."

	case ActionSolveProblem:
		return "Baik, saya akan membantu menganalisis dan menyelesaikan masalah kamu."

	case ActionExecuteRequest:
		return "Siap, saya akan membantu mengerjakan permintaan kamu."

	case ActionGeneralConversation:
		return "Baik, mari kita lanjutkan percakapannya."

	case ActionClarify:
		return "Boleh jelaskan lebih detail supaya saya bisa membantu dengan tepat?"

	default:
		return "Saya belum memahami tindakan yang harus dilakukan."
	}
}

// =========================
// Memory Extraction
// =========================

func extractNameFromContext(
	ctx *appcontext.MessageContext,
) string {

	if ctx == nil {
		return ""
	}

	// Cari dari pesan user paling baru menuju pesan lama.
	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {

		message := ctx.RecentMessages[i]

		// Role dibuat case-insensitive supaya
		// "user", "User", dan "USER" tetap terbaca.
		if !strings.EqualFold(
			strings.TrimSpace(message.Role),
			"user",
		) {
			continue
		}

		text := strings.TrimSpace(message.Content)

		if text == "" {
			continue
		}

		name := extractName(text)

		if name != "" {
			return name
		}
	}

	return ""
}

// extractName mencoba beberapa pola kalimat
// yang umum digunakan ketika user memperkenalkan
// namanya.

func extractName(text string) string {

	original := strings.TrimSpace(text)
	lower := strings.ToLower(original)

	prefixes := []string{
		"nama saya ",
		"nama aku ",
		"nama gue ",
		"nama gua ",
		"nama gw ",
		"saya bernama ",
		"aku bernama ",
		"gue bernama ",
		"gua bernama ",
		"gw bernama ",
		"saya ",
		"aku ",
	}

	for _, prefix := range prefixes {

		if !strings.HasPrefix(lower, prefix) {
			continue
		}

		name := strings.TrimSpace(
			original[len(prefix):],
		)

		name = cleanName(name)

		if name != "" {
			return name
		}
	}

	return ""
}

// cleanName membersihkan bagian tambahan
// dari kalimat setelah nama.

func cleanName(name string) string {

	name = strings.TrimSpace(name)

	if name == "" {
		return ""
	}

	// Buang tanda baca di akhir.
	name = strings.TrimRight(
		name,
		".,!?;:",
	)

	name = strings.TrimSpace(name)

	if name == "" {
		return ""
	}

	// Jangan menganggap kalimat panjang sebagai nama.
	//
	// Contoh:
	// "saya sedang membuat AI"
	//
	// Tidak seharusnya seluruh kalimat dianggap
	// sebagai nama.

	words := strings.Fields(name)

	if len(words) > 4 {
		return ""
	}

	return name
}
