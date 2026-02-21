package agent

import (
	"strings"
	"testing"

	"github.com/gourmet-guide/backend/internal/domain"
)

func TestSuggestTagsAddsDietaryAndPolicyTags(t *testing.T) {
	t.Parallel()
	item := domain.MenuItem{Name: "Halal Chicken Bowl", Description: "No pork, no lard recipe"}
	tags := SuggestTags(item)
	joined := strings.Join(tags, ",")
	for _, expected := range []string{"halal", "no-pork", "no-lard"} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("expected tag %q in %q", expected, joined)
		}
	}
}

func TestSuggestTagsRemovesContradictoryAllergenClaims(t *testing.T) {
	t.Parallel()
	item := domain.MenuItem{
		Name:      "Peanut Noodles",
		Tags:      []string{"nut-free"},
		Allergens: []domain.Allergen{domain.AllergenPeanut},
	}
	tags := SuggestTags(item)
	for _, tag := range tags {
		if tag == "nut-free" {
			t.Fatal("nut-free should be removed when peanut allergen is present")
		}
	}
}
