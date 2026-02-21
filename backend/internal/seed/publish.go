package seed

import (
	"context"

	"github.com/gourmet-guide/backend/internal/domain"
)

// ImageUploader writes generated menu images to a remote image store.
type ImageUploader interface {
	UploadMenuImage(ctx context.Context, restaurantID, menuItemID, fileName string, data []byte) (string, error)
	Close() error
}

// RestaurantStore persists generated restaurants.
type RestaurantStore interface {
	SaveRestaurant(ctx context.Context, restaurant domain.Restaurant) error
	Close() error
}
