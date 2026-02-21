package seed

import (
	"context"
	"fmt"

	"github.com/gourmet-guide/backend/internal/domain"
)

// FoodImage represents a generated visual for a menu item.
type FoodImage struct {
	FileName string
	Data     []byte
}

// MenuConcept contains AI-generated menu metadata.
type MenuConcept struct {
	Name        string
	Description string
	Tags        []string
	Allergens   []domain.Allergen
}

// AIProvider defines the behaviors needed from an LLM/image generator.
type AIProvider interface {
	GenerateMenuConcepts(ctx context.Context, count int) ([]MenuConcept, error)
	GenerateFoodImage(ctx context.Context, menuName, menuDescription string) ([]byte, error)
}

// GenerateRestaurantsFromAI builds demo restaurants and images via an AI provider.
func GenerateRestaurantsFromAI(ctx context.Context, provider AIProvider, menuCount int) ([]domain.Restaurant, map[string]FoodImage, error) {
	concepts, err := provider.GenerateMenuConcepts(ctx, menuCount)
	if err != nil {
		return nil, nil, err
	}

	restaurantID, err := newUUID()
	if err != nil {
		return nil, nil, err
	}
	restaurant := domain.Restaurant{
		ID:        restaurantID,
		Name:      "Gemini Demo Bistro",
		MenuItems: make([]domain.MenuItem, 0, len(concepts)),
		Combos:    []domain.Combo{},
	}

	images := make(map[string]FoodImage, len(concepts))
	for _, concept := range concepts {
		itemID, err := newUUID()
		if err != nil {
			return nil, nil, err
		}
		imageData, err := provider.GenerateFoodImage(ctx, concept.Name, concept.Description)
		if err != nil {
			return nil, nil, err
		}
		images[itemID] = FoodImage{FileName: fmt.Sprintf("%s.png", itemID), Data: imageData}
		restaurant.MenuItems = append(restaurant.MenuItems, domain.MenuItem{
			ID:          itemID,
			Name:        concept.Name,
			Description: concept.Description,
			Allergens:   concept.Allergens,
			Tags:        concept.Tags,
			ImageURL:    fmt.Sprintf("images/%s.png", itemID),
		})
	}

	for i := 0; i+1 < len(restaurant.MenuItems); i += 2 {
		comboID, err := newUUID()
		if err != nil {
			return nil, nil, err
		}
		restaurant.Combos = append(restaurant.Combos, domain.Combo{
			ID:          comboID,
			Name:        fmt.Sprintf("Chef Pairing %d", (i/2)+1),
			ItemIDs:     []string{restaurant.MenuItems[i].ID, restaurant.MenuItems[i+1].ID},
			Description: "Auto-paired by Gemini for demo recommendations.",
		})
	}

	return []domain.Restaurant{restaurant}, images, nil
}
