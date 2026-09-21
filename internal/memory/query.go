package memory

import (
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func RequestedType(ctx *appcontext.MessageContext) (MemoryType, bool) {
	if ctx == nil || len(ctx.RecentMessages) == 0 {
		return "", false
	}

	latest := latestUserMessage(ctx)
	if latest == "" {
		return "", false
	}

	// Activity memory
	if containsAny(
		latest,
		"apa yang sedang saya pelajari",
		"apa yang sedang aku pelajari",
		"apa yang sedang gue pelajari",
		"apa yang sedang gua pelajari",
		"apa yang sedang gw pelajari",
		"apa yang saya pelajari",
		"apa yang aku pelajari",
		"apa yang gue pelajari",
		"apa yang gua pelajari",
		"apa yang gw pelajari",
		"saya sedang belajar apa",
		"aku sedang belajar apa",
		"gue sedang belajar apa",
		"gua sedang belajar apa",
		"gw sedang belajar apa",
		"saya lagi belajar apa",
		"aku lagi belajar apa",
		"gue lagi belajar apa",
		"gua lagi belajar apa",
		"gw lagi belajar apa",
	) {
		return MemoryActivity, true
	}

	// Name memory
	if containsAny(
		latest,
		"siapa nama saya",
		"siapa nama aku",
		"siapa nama gue",
		"siapa nama gua",
		"siapa nama gw",
		"apa nama saya",
		"apa nama aku",
		"apa nama gue",
		"apa nama gua",
		"apa nama gw",
		"nama saya siapa",
		"nama aku siapa",
		"nama gue siapa",
		"nama gua siapa",
		"nama gw siapa",
	) {
		return MemoryName, true
	}

	return "", false
}

func latestUserMessage(ctx *appcontext.MessageContext) string {
	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {
		message := ctx.RecentMessages[i]

		if !strings.EqualFold(strings.TrimSpace(message.Role), "user") {
			continue
		}

		text := strings.TrimSpace(message.Content)
		if text == "" {
			continue
		}

		return strings.ToLower(text)
	}

	return ""
}

func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}
