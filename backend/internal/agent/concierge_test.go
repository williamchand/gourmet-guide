package agent

import (
	"context"
	"testing"

	"github.com/gourmet-guide/backend/internal/domain"
	"github.com/gourmet-guide/backend/internal/gcp"
)

func TestApplySafetyPoliciesFiltersHardAllergensAndCrossContamination(t *testing.T) {
	t.Parallel()
	items := []domain.MenuItem{
		{Name: "Peanut Curry", Allergens: []domain.Allergen{domain.AllergenPeanut}, Tags: []string{"spicy"}},
		{Name: "House Salad", Tags: []string{"vegan", "light"}},
		{Name: "Fries", CrossContaminationRisk: []domain.Allergen{domain.AllergenPeanut}, Tags: []string{"vegan"}},
	}

	safe, warning := applySafetyPolicies(items, []domain.Allergen{domain.AllergenPeanut}, []string{"vegan"})
	if len(safe) != 1 {
		t.Fatalf("expected 1 safe item, got %d", len(safe))
	}
	if safe[0].Name != "House Salad" {
		t.Fatalf("expected House Salad, got %q", safe[0].Name)
	}
	if warning == "" {
		t.Fatal("expected warning to mention filtered items")
	}
}

func TestInterruptSessionUpdatesStatus(t *testing.T) {
	t.Parallel()
	store := gcp.NewMemoryStore()
	runtime := NewRuntime("gemini", store)
	service := NewConciergeService(store, gcp.NewMemoryImageStore(), runtime)

	session, err := service.StartSession(context.Background(), "rest-1", nil, nil)
	if err != nil {
		t.Fatalf("start session: %v", err)
	}
	if err := service.InterruptSession(context.Background(), session.ID); err != nil {
		t.Fatalf("interrupt: %v", err)
	}
	updated, err := service.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	if updated.Status != domain.SessionStatusInterrupted {
		t.Fatalf("expected interrupted status, got %s", updated.Status)
	}
}

func TestApplySafetyPoliciesTreatsDietaryTagsAsHardRequirement(t *testing.T) {
	t.Parallel()
	items := []domain.MenuItem{
		{Name: "Halal Salad", Tags: []string{"halal", "no-pork", "no-lard"}},
		{Name: "Pork Ramen", Tags: []string{"spicy"}},
	}

	safe, warning := applySafetyPolicies(items, nil, []string{"halal", "no-pork"})
	if len(safe) != 1 {
		t.Fatalf("expected 1 dietary-safe item, got %d", len(safe))
	}
	if safe[0].Name != "Halal Salad" {
		t.Fatalf("expected Halal Salad, got %q", safe[0].Name)
	}
	if warning == "" {
		t.Fatal("expected warning for excluded dietary mismatches")
	}
}

func TestHasAllRequiredTags(t *testing.T) {
	t.Parallel()
	item := domain.MenuItem{Tags: []string{"Halal", "No-Pork", "no-lard"}}
	if !hasAllRequiredTags(item, []string{"halal", "no-pork"}) {
		t.Fatal("expected tags to satisfy required dietary tags")
	}
	if hasAllRequiredTags(item, []string{"halal", "vegan"}) {
		t.Fatal("expected false when missing required dietary tag")
	}
}
