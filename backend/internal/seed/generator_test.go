package seed

import (
	"context"
	"regexp"
	"testing"

	"github.com/gourmet-guide/backend/internal/domain"
)

var uuidPattern = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

type fakeProvider struct{}

func (f *fakeProvider) GenerateMenuConcepts(_ context.Context, count int) ([]MenuConcept, error) {
	concepts := make([]MenuConcept, 0, count)
	for i := 0; i < count; i++ {
		concepts = append(concepts, MenuConcept{
			Name:        "Dish",
			Description: "Demo description",
			Tags:        []string{"demo"},
			Allergens:   []domain.Allergen{domain.AllergenSoy},
		})
	}
	return concepts, nil
}

func (f *fakeProvider) GenerateFoodImage(_ context.Context, _, _ string) ([]byte, error) {
	return []byte("png-bytes"), nil
}

func TestGenerateRestaurantsFromAIIncludesImagesCombosAndUUIDs(t *testing.T) {
	restaurants, images, err := GenerateRestaurantsFromAI(context.Background(), &fakeProvider{}, 6)
	if err != nil {
		t.Fatalf("GenerateRestaurantsFromAI returned error: %v", err)
	}
	if len(restaurants) != 1 {
		t.Fatalf("expected 1 restaurant, got %d", len(restaurants))
	}
	restaurant := restaurants[0]
	if !uuidPattern.MatchString(restaurant.ID) {
		t.Fatalf("restaurant id should be uuid, got %q", restaurant.ID)
	}
	if len(restaurant.MenuItems) != 6 {
		t.Fatalf("expected 6 menu items, got %d", len(restaurant.MenuItems))
	}
	if len(images) != 6 {
		t.Fatalf("expected 6 images, got %d", len(images))
	}
	if len(restaurant.Combos) != 3 {
		t.Fatalf("expected 3 combos, got %d", len(restaurant.Combos))
	}

	for _, item := range restaurant.MenuItems {
		if !uuidPattern.MatchString(item.ID) {
			t.Fatalf("menu item id should be uuid, got %q", item.ID)
		}
		if _, ok := images[item.ID]; !ok {
			t.Fatalf("expected image for menu item %q", item.ID)
		}
	}
	for _, combo := range restaurant.Combos {
		if !uuidPattern.MatchString(combo.ID) {
			t.Fatalf("combo id should be uuid, got %q", combo.ID)
		}
	}
}

func TestParseAllergensFiltersUnknown(t *testing.T) {
	got := parseAllergens([]string{"dairy", "soy", "invalid", "tree_nut"})
	if len(got) != 3 {
		t.Fatalf("expected 3 allergens, got %d", len(got))
	}
}
