package seed

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gourmet-guide/backend/internal/domain"
)

const (
	defaultTextModel  = "gemini-2.0-flash"
	defaultImageModel = "gemini-2.0-flash-preview-image-generation"
)

// GeminiProvider calls Gemini API for menu + image generation.
type GeminiProvider struct {
	APIKey     string
	TextModel  string
	ImageModel string
	BaseURL    string
	Client     *http.Client
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		APIKey:     apiKey,
		TextModel:  defaultTextModel,
		ImageModel: defaultImageModel,
		BaseURL:    "https://generativelanguage.googleapis.com/v1beta/models",
		Client:     &http.Client{Timeout: 60 * time.Second},
	}
}

func (g *GeminiProvider) GenerateMenuConcepts(ctx context.Context, count int) ([]MenuConcept, error) {
	prompt := fmt.Sprintf("Generate %d unique restaurant menu items for a modern global bistro. Return strict JSON array where each object has keys: name, description, tags (array of short strings), allergens (array choosing only from dairy, egg, fish, peanut, shellfish, soy, tree_nut, wheat). No markdown.", count)
	text, err := g.generateText(ctx, prompt)
	if err != nil {
		return nil, err
	}

	clean := strings.TrimSpace(text)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")

	var raw []struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		Allergens   []string `json:"allergens"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(clean)), &raw); err != nil {
		return nil, fmt.Errorf("parse menu JSON from gemini: %w", err)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("gemini returned zero menu concepts")
	}

	concepts := make([]MenuConcept, 0, len(raw))
	for _, item := range raw {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		concepts = append(concepts, MenuConcept{
			Name:        item.Name,
			Description: item.Description,
			Tags:        item.Tags,
			Allergens:   parseAllergens(item.Allergens),
		})
	}
	if len(concepts) == 0 {
		return nil, fmt.Errorf("gemini response had no valid menu names")
	}
	return concepts, nil
}

func (g *GeminiProvider) GenerateFoodImage(ctx context.Context, menuName, menuDescription string) ([]byte, error) {
	prompt := fmt.Sprintf("Create a realistic food photography image of this dish: %s. Description: %s", menuName, menuDescription)
	payload := map[string]any{
		"contents": []map[string]any{{
			"parts": []map[string]any{{"text": prompt}},
		}},
	}

	responseBody, err := g.callGenerateContent(ctx, g.model(g.ImageModel, defaultImageModel), payload)
	if err != nil {
		return nil, err
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					InlineData struct {
						Data string `json:"data"`
					} `json:"inlineData"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("decode image response: %w", err)
	}

	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.InlineData.Data) == "" {
				continue
			}
			decoded, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
			if err != nil {
				return nil, fmt.Errorf("decode image base64: %w", err)
			}
			return decoded, nil
		}
	}

	return nil, fmt.Errorf("gemini did not return inline image data")
}

func (g *GeminiProvider) generateText(ctx context.Context, prompt string) (string, error) {
	payload := map[string]any{
		"contents": []map[string]any{{
			"parts": []map[string]any{{"text": prompt}},
		}},
	}
	responseBody, err := g.callGenerateContent(ctx, g.model(g.TextModel, defaultTextModel), payload)
	if err != nil {
		return "", err
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("decode text response: %w", err)
	}
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.Text) != "" {
				return part.Text, nil
			}
		}
	}
	return "", fmt.Errorf("gemini returned empty text response")
}

func (g *GeminiProvider) callGenerateContent(ctx context.Context, model string, payload map[string]any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s:generateContent?key=%s", strings.TrimRight(g.BaseURL, "/"), model, g.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini api status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return respBody, nil
}

func (g *GeminiProvider) model(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func parseAllergens(values []string) []domain.Allergen {
	allergens := make([]domain.Allergen, 0, len(values))
	for _, v := range values {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case string(domain.AllergenDairy):
			allergens = append(allergens, domain.AllergenDairy)
		case string(domain.AllergenEgg):
			allergens = append(allergens, domain.AllergenEgg)
		case string(domain.AllergenFish):
			allergens = append(allergens, domain.AllergenFish)
		case string(domain.AllergenPeanut):
			allergens = append(allergens, domain.AllergenPeanut)
		case string(domain.AllergenShellfish):
			allergens = append(allergens, domain.AllergenShellfish)
		case string(domain.AllergenSoy):
			allergens = append(allergens, domain.AllergenSoy)
		case string(domain.AllergenTreeNut):
			allergens = append(allergens, domain.AllergenTreeNut)
		case string(domain.AllergenWheat):
			allergens = append(allergens, domain.AllergenWheat)
		}
	}
	return allergens
}
