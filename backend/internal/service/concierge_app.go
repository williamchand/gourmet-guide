package service

import (
	"context"

	"github.com/gourmet-guide/backend/internal/agent"
	"github.com/gourmet-guide/backend/internal/domain"
)

type StartSessionInput struct {
	RestaurantID   string
	HardAllergens  []domain.Allergen
	PreferenceTags []string
	MenuItems      []domain.MenuItem
}

type StartSessionOutput struct {
	Session            domain.ConciergeSession
	SuggestedMenuItems []domain.MenuItem
}

type ExtractMenuOutput struct {
	ImagePath string
	MenuItems []domain.MenuItem
}

type ConciergeApp struct {
	concierge *agent.ConciergeService
}

func NewConciergeApp(concierge *agent.ConciergeService) *ConciergeApp {
	return &ConciergeApp{concierge: concierge}
}

func (a *ConciergeApp) StartSession(ctx context.Context, input StartSessionInput) (StartSessionOutput, error) {
	enriched, err := a.concierge.SaveMenuItems(ctx, input.RestaurantID, input.MenuItems)
	if err != nil {
		return StartSessionOutput{}, err
	}
	session, err := a.concierge.StartSession(ctx, input.RestaurantID, input.HardAllergens, input.PreferenceTags)
	if err != nil {
		return StartSessionOutput{}, err
	}
	return StartSessionOutput{Session: session, SuggestedMenuItems: enriched}, nil
}

func (a *ConciergeApp) GetSession(ctx context.Context, sessionID string) (domain.ConciergeSession, error) {
	return a.concierge.GetSession(ctx, sessionID)
}

func (a *ConciergeApp) EndSession(ctx context.Context, sessionID string) error {
	return a.concierge.EndSession(ctx, sessionID)
}

func (a *ConciergeApp) InterruptSession(ctx context.Context, sessionID string) error {
	return a.concierge.InterruptSession(ctx, sessionID)
}

func (a *ConciergeApp) SendMessage(ctx context.Context, sessionID, prompt string) (string, error) {
	return a.concierge.SendMessage(ctx, sessionID, prompt)
}

func (a *ConciergeApp) TagMenuItems(ctx context.Context, restaurantID string, items []domain.MenuItem) ([]domain.MenuItem, error) {
	return a.concierge.SaveMenuItems(ctx, restaurantID, items)
}

func (a *ConciergeApp) ExtractMenuFromImage(ctx context.Context, restaurantID, fileName string, content []byte) (ExtractMenuOutput, error) {
	items, imagePath, err := a.concierge.AutoExtractMenuFromImage(ctx, restaurantID, fileName, content)
	if err != nil {
		return ExtractMenuOutput{}, err
	}
	return ExtractMenuOutput{ImagePath: imagePath, MenuItems: items}, nil
}
