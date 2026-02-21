//go:build gcp

package agent

import (
	"context"

	"github.com/google/generative-ai-go/genai"
)

type geminiClient struct {
	client *genai.Client
}

// newGeminiClient configures Google GenAI SDK to route traffic through Vertex AI.
func newGeminiClient(ctx context.Context) (Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{Backend: genai.BackendVertexAI})
	if err != nil {
		return nil, err
	}
	return &geminiClient{client: client}, nil
}

func (g *geminiClient) Generate(ctx context.Context, modelName, prompt string) (string, error) {
	config := &genai.GenerateContentConfig{MaxOutputTokens: 256}
	res, err := g.client.Models.GenerateContent(ctx, modelName, genai.Text(prompt), config)
	if err != nil {
		return "", err
	}
	return res.Text(), nil
}
