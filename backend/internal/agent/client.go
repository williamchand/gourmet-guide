package agent

import "context"

// Client describes LLM generation behavior used by runtime.
type Client interface {
	Generate(ctx context.Context, modelName, prompt string) (string, error)
}

type echoClient struct{}

func (e *echoClient) Generate(_ context.Context, _, prompt string) (string, error) {
	return "received and processed with Gemini: " + prompt, nil
}

func newDefaultClient() Client {
	return &echoClient{}
}
