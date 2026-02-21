package agent

import (
	"bytes"
	"context"
	"regexp"
	"strings"

	"github.com/gourmet-guide/backend/internal/domain"
)

var menuLinePattern = regexp.MustCompile(`(?i)[a-z][a-z0-9\s,&'/-]{2,}`)

// MenuExtractor provides menu extraction from uploaded menu images.
type MenuExtractor interface {
	ExtractMenuItems(ctx context.Context, content []byte) ([]domain.MenuItem, error)
}

// HeuristicMenuExtractor provides a lightweight local fallback for menu extraction.
type HeuristicMenuExtractor struct{}

func (h *HeuristicMenuExtractor) ExtractMenuItems(_ context.Context, content []byte) ([]domain.MenuItem, error) {
	lines := bytes.Split(content, []byte("\n"))
	items := make([]domain.MenuItem, 0, len(lines))
	seen := map[string]struct{}{}
	for _, line := range lines {
		candidate := strings.TrimSpace(string(line))
		if !menuLinePattern.MatchString(candidate) {
			continue
		}
		candidate = strings.Join(strings.Fields(candidate), " ")
		if len(candidate) > 80 {
			candidate = candidate[:80]
		}
		if _, ok := seen[strings.ToLower(candidate)]; ok {
			continue
		}
		seen[strings.ToLower(candidate)] = struct{}{}
		items = append(items, domain.MenuItem{
			ID:          slugify(candidate),
			Name:        candidate,
			Description: "Auto-extracted from uploaded menu image. Review before publishing.",
		})
	}
	if len(items) > 12 {
		items = items[:12]
	}
	return items, nil
}

func slugify(v string) string {
	lower := strings.ToLower(strings.TrimSpace(v))
	lower = strings.ReplaceAll(lower, " ", "-")
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, lower)
}
