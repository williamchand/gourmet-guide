package agent

import (
	"fmt"
	"strings"
)

func validatePrompt(prompt string) (string, error) {
	trimmed := strings.TrimSpace(prompt)
	if trimmed == "" {
		return "", fmt.Errorf("prompt cannot be empty")
	}
	return trimmed, nil
}
