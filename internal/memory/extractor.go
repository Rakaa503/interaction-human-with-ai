package memory

import (
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func ExtractActivity(ctx *appcontext.MessageContext) *Memory {
	if ctx == nil {
		return nil
	}

	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {
		message := ctx.RecentMessages[i]

		if !strings.EqualFold(strings.TrimSpace(message.Role), "user") {
			continue
		}

		activity := extractActivity(message.Content)

		if activity == "" {
			continue
		}

		return &Memory{
			Type:       MemoryActivity,
			Value:      activity,
			Confidence: 0.95,
		}
	}

	return nil
}

func extractActivity(input string) string {
	text := strings.TrimSpace(input)
	lower := strings.ToLower(text)

	prefixes := []string{
		"saya sedang belajar ",
		"saya lagi belajar ",
		"saya sedang mempelajari ",
		"saya lagi mempelajari ",
		"aku sedang belajar ",
		"aku lagi belajar ",
		"aku sedang mempelajari ",
		"aku lagi mempelajari ",
		"gue sedang belajar ",
		"gue lagi belajar ",
		"gue sedang mempelajari ",
		"gue lagi mempelajari ",
		"gua sedang belajar ",
		"gua lagi belajar ",
		"gua sedang mempelajari ",
		"gua lagi mempelajari ",
		"gw sedang belajar ",
		"gw lagi belajar ",
		"gw sedang mempelajari ",
		"gw lagi mempelajari ",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			value := strings.TrimSpace(text[len(prefix):])

			if value == "" {
				return ""
			}

			return strings.TrimRight(value, ".,!?")
		}
	}

	return ""
}
