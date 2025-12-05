package gemini

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

// GeminiLLM implements the LLM interface for your hosted model
type GeminiLLM struct {
	gemini *googleai.GoogleAI
}

// NewGeminiLLM creates a new instance of your custom LLM.
func NewGeminiLLM(apiKey string) *GeminiLLM {
	gemini, err := googleai.New(context.Background(), googleai.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}
	return &GeminiLLM{
		gemini: gemini,
	}
}

// Call implements the [llms.Model] interface.
func (c *GeminiLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, c, prompt, options...)
}

func (c *GeminiLLM) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return c.gemini.GenerateContent(ctx, messages)
}

func (c *GeminiLLM) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	embed, err := c.gemini.CreateEmbedding(ctx, texts)
	if err != nil {
		// Todo add log here
		return nil, err
	}

	return embed, nil
}

func (c *GeminiLLM) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embed, err := c.gemini.CreateEmbedding(ctx, []string{text})
	if err != nil {
		// Todo add log here
		return nil, err
	}

	return embed[0], nil
}
