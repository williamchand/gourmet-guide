package agent

import (
	"strings"

	"github.com/gourmet-guide/backend/internal/domain"
)

var tagRules = []struct {
	tag      string
	keywords []string
}{
	{tag: "halal", keywords: []string{"halal"}},
	{tag: "no-pork", keywords: []string{"no pork", "without pork", "pork-free"}},
	{tag: "no-beef", keywords: []string{"no beef", "without beef", "beef-free"}},
	{tag: "no-lard", keywords: []string{"no lard", "without lard", "lard-free"}},
	{tag: "vegetarian", keywords: []string{"vegetarian"}},
	{tag: "vegan", keywords: []string{"vegan", "plant-based"}},
	{tag: "gluten-free", keywords: []string{"gluten free", "gluten-free"}},
	{tag: "dairy-free", keywords: []string{"dairy free", "dairy-free"}},
	{tag: "nut-free", keywords: []string{"nut free", "nut-free", "peanut-free"}},
}

func SuggestTags(item domain.MenuItem) []string {
	searchText := strings.ToLower(strings.TrimSpace(item.Name + " " + item.Description + " " + strings.Join(item.Tags, " ")))
	candidates := make(map[string]struct{}, len(item.Tags)+4)
	for _, existing := range item.Tags {
		if normalized := strings.ToLower(strings.TrimSpace(existing)); normalized != "" {
			candidates[normalized] = struct{}{}
		}
	}
	for _, rule := range tagRules {
		for _, keyword := range rule.keywords {
			if strings.Contains(searchText, keyword) {
				candidates[rule.tag] = struct{}{}
				break
			}
		}
	}
	for _, allergen := range item.Allergens {
		switch allergen {
		case domain.AllergenPeanut, domain.AllergenTreeNut:
			delete(candidates, "nut-free")
		case domain.AllergenDairy:
			delete(candidates, "dairy-free")
		case domain.AllergenWheat:
			delete(candidates, "gluten-free")
		}
	}

	result := make([]string, 0, len(candidates))
	for tag := range candidates {
		result = append(result, tag)
	}
	return result
}

func EnrichMenuItemsWithSuggestedTags(items []domain.MenuItem) []domain.MenuItem {
	enriched := make([]domain.MenuItem, len(items))
	for i, item := range items {
		item.Tags = SuggestTags(item)
		enriched[i] = item
	}
	return enriched
}
