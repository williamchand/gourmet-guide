package domain

import "time"

// Allergen captures allergens that can trigger severe reactions.
type Allergen string

const (
	AllergenDairy     Allergen = "dairy"
	AllergenEgg       Allergen = "egg"
	AllergenFish      Allergen = "fish"
	AllergenPeanut    Allergen = "peanut"
	AllergenShellfish Allergen = "shellfish"
	AllergenSoy       Allergen = "soy"
	AllergenTreeNut   Allergen = "tree_nut"
	AllergenWheat     Allergen = "wheat"
)

// MenuItem represents a single dish.
type MenuItem struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	Description            string     `json:"description"`
	Allergens              []Allergen `json:"allergens"`
	CrossContaminationRisk []Allergen `json:"crossContaminationRisk,omitempty"`
	Tags                   []string   `json:"tags,omitempty"`
	ImageURL               string     `json:"imageUrl,omitempty"`
}

// Combo defines a curated pairing of menu items.
type Combo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ItemIDs     []string `json:"itemIds"`
	Description string   `json:"description,omitempty"`
}

// Restaurant collects a restaurant menu and combo metadata.
type Restaurant struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	MenuItems []MenuItem `json:"menuItems"`
	Combos    []Combo    `json:"combos"`
}

// SessionStatus indicates current lifecycle state.
type SessionStatus string

const (
	SessionStatusActive      SessionStatus = "active"
	SessionStatusInterrupted SessionStatus = "interrupted"
	SessionStatusCompleted   SessionStatus = "completed"
)

// ConciergeSession is the long-lived conversation session.
type ConciergeSession struct {
	ID               string        `json:"id"`
	RestaurantID     string        `json:"restaurantId"`
	HardAllergens    []Allergen    `json:"hardAllergens"`
	PreferenceTags   []string      `json:"preferenceTags"`
	Status           SessionStatus `json:"status"`
	LastAssistantMsg string        `json:"lastAssistantMessage,omitempty"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
}
