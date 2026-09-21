package memory

import (
	"strings"

	appcontext "github.com/Rakaa503/AviGo/internal/context"
)

func ExtractName(ctx *appcontext.MessageContext) *Memory {
	if ctx == nil {
		return nil
	}

	for i := len(ctx.RecentMessages) - 1; i >= 0; i-- {
		message := ctx.RecentMessages[i]

		if !strings.EqualFold(strings.TrimSpace(message.Role), "user") {
			continue
		}

		name := extractName(message.Content)

		if name == "" {
			continue
		}

		return &Memory{
			Type:       MemoryName,
			Value:      name,
			Confidence: 0.95,
		}
	}

	return nil
}

func extractName(input string) string {
	text := strings.TrimSpace(input)

	if text == "" {
		return ""
	}

	lower := strings.ToLower(text)

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
		"gue ",
		"gua ",
		"gw ",
	}

	for _, prefix := range prefixes {
		if !strings.HasPrefix(lower, prefix) {
			continue
		}

		name := strings.TrimSpace(text[len(prefix):])

		name = extractNamePart(name)

		if cleanName(name) == "" {
			continue
		}

		return cleanName(name)
	}

	return ""
}

func extractNamePart(input string) string {
	text := strings.TrimSpace(input)
	lower := strings.ToLower(text)

	stopWords := []string{
		" saya ",
		" aku ",
		" gue ",
		" gua ",
		" gw ",
		" dan ",
		" karena ",
		" yang ",
		" sedang ",
		" lagi ",
		" adalah ",
		" merupakan ",
	}

	for _, stopWord := range stopWords {
		if index := strings.Index(lower, stopWord); index >= 0 {
			text = text[:index]
			break
		}
	}

	return strings.TrimSpace(text)
}

func cleanName(input string) string {
	name := strings.TrimSpace(input)
	name = strings.Trim(name, ".,!?;:\"'")

	if name == "" {
		return ""
	}

	words := strings.Fields(name)

	if len(words) > 4 {
		return ""
	}

	if isBlockedWord(strings.ToLower(name)) {
		return ""
	}

	return name
}

func isBlockedWord(value string) bool {
	blocked := map[string]struct{}{
		"sedang":    {},
		"lagi":      {},
		"membuat":   {},
		"ingin":     {},
		"mau":       {},
		"adalah":    {},
		"merupakan": {},
		"tidak":     {},
		"bukan":     {},
		"suka":      {},
		"punya":     {},
		"memiliki":  {},
	}

	_, exists := blocked[value]

	return exists
}
