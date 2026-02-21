package agent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gourmet-guide/backend/internal/domain"
	"github.com/gourmet-guide/backend/internal/gcp"
)

const highRiskDisclaimer = "I cannot confidently guarantee safety for that request. Please confirm ingredients and cross-contamination policy with restaurant staff before ordering."

type ConciergeService struct {
	store         gcp.SessionStore
	imageStore    gcp.ImageStore
	menuExtractor MenuExtractor
	runtime       *Runtime

	mu      sync.Mutex
	ongoing map[string]context.CancelFunc
}

func NewConciergeService(store gcp.SessionStore, imageStore gcp.ImageStore, runtime *Runtime) *ConciergeService {
	return &ConciergeService{
		store:         store,
		imageStore:    imageStore,
		menuExtractor: &HeuristicMenuExtractor{},
		runtime:       runtime,
		ongoing:       map[string]context.CancelFunc{},
	}
}

func (s *ConciergeService) SaveMenuItems(ctx context.Context, restaurantID string, items []domain.MenuItem) ([]domain.MenuItem, error) {
	enriched := EnrichMenuItemsWithSuggestedTags(items)
	if err := s.store.SaveMenuSafetyMetadata(ctx, restaurantID, enriched); err != nil {
		return nil, err
	}
	return enriched, nil
}

func (s *ConciergeService) StartSession(ctx context.Context, restaurantID string, hardAllergens []domain.Allergen, preferenceTags []string) (domain.ConciergeSession, error) {
	now := time.Now().UTC()
	session := domain.ConciergeSession{
		ID:             newSessionID(),
		RestaurantID:   restaurantID,
		HardAllergens:  hardAllergens,
		PreferenceTags: preferenceTags,
		Status:         domain.SessionStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return session, s.store.SaveSession(ctx, session)
}

func (s *ConciergeService) SendMessage(ctx context.Context, sessionID, prompt string) (string, error) {
	session, err := s.store.LoadSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if session.ID == "" {
		return "", errors.New("session not found")
	}

	items, err := s.store.LoadMenuSafetyMetadata(ctx, session.RestaurantID)
	if err != nil {
		return "", err
	}

	safeItems, warning := applySafetyPolicies(items, session.HardAllergens, session.PreferenceTags)
	if len(safeItems) == 0 {
		return highRiskDisclaimer, nil
	}

	menuNames := make([]string, 0, len(safeItems))
	for _, item := range safeItems {
		menuNames = append(menuNames, item.Name)
	}

	turnCtx, cancel := context.WithCancel(ctx)
	s.setOngoingCancel(sessionID, cancel)
	defer s.clearOngoingCancel(sessionID)

	reply, err := s.runtime.Respond(turnCtx, sessionID, prompt, menuNames)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "response interrupted, ready for your next request", nil
		}
		return "", err
	}
	if warning != "" {
		reply = fmt.Sprintf("%s\n\nSafety note: %s", reply, warning)
	}

	session.Status = domain.SessionStatusActive
	session.LastAssistantMsg = reply
	session.UpdatedAt = time.Now().UTC()
	if err := s.store.SaveSession(ctx, session); err != nil {
		return "", err
	}
	return reply, nil
}

func (s *ConciergeService) InterruptSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	cancel := s.ongoing[sessionID]
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	session, err := s.store.LoadSession(ctx, sessionID)
	if err != nil {
		return err
	}
	session.Status = domain.SessionStatusInterrupted
	session.UpdatedAt = time.Now().UTC()
	return s.store.SaveSession(ctx, session)
}

func (s *ConciergeService) EndSession(ctx context.Context, sessionID string) error {
	session, err := s.store.LoadSession(ctx, sessionID)
	if err != nil {
		return err
	}
	session.Status = domain.SessionStatusCompleted
	session.UpdatedAt = time.Now().UTC()
	return s.store.SaveSession(ctx, session)
}

func (s *ConciergeService) AutoExtractMenuFromImage(ctx context.Context, restaurantID, fileName string, content []byte) ([]domain.MenuItem, string, error) {
	imagePath, err := s.imageStore.SaveSessionImage(ctx, restaurantID, fileName, content)
	if err != nil {
		return nil, "", err
	}
	items, err := s.menuExtractor.ExtractMenuItems(ctx, content)
	if err != nil {
		return nil, "", err
	}
	enriched, err := s.SaveMenuItems(ctx, restaurantID, items)
	if err != nil {
		return nil, "", err
	}
	return enriched, imagePath, nil
}
func (s *ConciergeService) GetSession(ctx context.Context, sessionID string) (domain.ConciergeSession, error) {
	return s.store.LoadSession(ctx, sessionID)
}

func (s *ConciergeService) setOngoingCancel(sessionID string, cancel context.CancelFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ongoing[sessionID] = cancel
}

func (s *ConciergeService) clearOngoingCancel(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ongoing, sessionID)
}

func applySafetyPolicies(items []domain.MenuItem, hardAllergens []domain.Allergen, preferenceTags []string) ([]domain.MenuItem, string) {
	allergenSet := map[domain.Allergen]struct{}{}
	for _, allergen := range hardAllergens {
		allergenSet[allergen] = struct{}{}
	}
	filtered := make([]domain.MenuItem, 0, len(items))
	crossContaminationWarning := false
	for _, item := range items {
		if containsAnyAllergen(item.Allergens, allergenSet) {
			continue
		}
		if containsAnyAllergen(item.CrossContaminationRisk, allergenSet) {
			crossContaminationWarning = true
			continue
		}
		filtered = append(filtered, item)
	}

	dietaryFiltered := false
	if len(preferenceTags) > 0 {
		strictFiltered := make([]domain.MenuItem, 0, len(filtered))
		// Dietary constraints are treated as hard requirements in-session for safety.
		for _, item := range filtered {
			if hasAllRequiredTags(item, preferenceTags) {
				strictFiltered = append(strictFiltered, item)
			} else {
				dietaryFiltered = true
			}
		}
		filtered = strictFiltered
		sort.SliceStable(filtered, func(i, j int) bool {
			return preferenceScore(filtered[i], preferenceTags) > preferenceScore(filtered[j], preferenceTags)
		})
	}

	if crossContaminationWarning {
		return filtered, "Some items were excluded due to cross-contamination risk."
	}
	if dietaryFiltered {
		return filtered, "Some menu items were excluded because they did not satisfy required dietary tags."
	}
	if len(filtered) < len(items) {
		return filtered, "Some menu items were removed by hard allergen filters."
	}
	return filtered, ""
}

func containsAnyAllergen(itemAllergens []domain.Allergen, restricted map[domain.Allergen]struct{}) bool {
	for _, allergen := range itemAllergens {
		if _, blocked := restricted[allergen]; blocked {
			return true
		}
	}
	return false
}

func hasAllRequiredTags(item domain.MenuItem, requiredTags []string) bool {
	tagSet := map[string]struct{}{}
	for _, tag := range item.Tags {
		tagSet[strings.ToLower(strings.TrimSpace(tag))] = struct{}{}
	}
	for _, required := range requiredTags {
		normalized := strings.ToLower(strings.TrimSpace(required))
		if normalized == "" {
			continue
		}
		if _, ok := tagSet[normalized]; !ok {
			return false
		}
	}
	return true
}

func preferenceScore(item domain.MenuItem, preferenceTags []string) int {
	score := 0
	for _, preference := range preferenceTags {
		for _, tag := range item.Tags {
			if strings.EqualFold(preference, tag) {
				score++
			}
		}
	}
	return score
}

func newSessionID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
