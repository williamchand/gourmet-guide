package agent

import (
	"context"
	"testing"
)

func TestHeuristicMenuExtractorExtractsUniqueLines(t *testing.T) {
	t.Parallel()
	extractor := &HeuristicMenuExtractor{}
	items, err := extractor.ExtractMenuItems(context.Background(), []byte("Spicy Tofu Bowl\nSpicy Tofu Bowl\nGarden Salad\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Name != "Spicy Tofu Bowl" {
		t.Fatalf("unexpected first item: %q", items[0].Name)
	}
}
