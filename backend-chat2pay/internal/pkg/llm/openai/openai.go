package openai

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// OpenAI implements the LLM interface for your hosted model
type OpenAI struct {
	openai *openai.LLM
}

// NewOpenAI creates a new instance of your custom LLM.
func NewOpenAI(apiKey string) *OpenAI {
	openai, err := openai.New(openai.WithToken(apiKey))
	if err != nil {
		panic(err)
	}
	return &OpenAI{
		openai: openai,
	}
}

// Call implements the [llms.Model] interface.
func (c *OpenAI) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, c, prompt, options...)
}

func (c *OpenAI) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return c.openai.GenerateContent(ctx, messages)
}

func (c *OpenAI) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	embed, err := c.openai.CreateEmbedding(ctx, texts)
	if err != nil {
		// Todo add log here
		return nil, err
	}

	return embed, nil
}

func (c *OpenAI) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embed, err := c.openai.CreateEmbedding(ctx, []string{text})
	if err != nil {
		// Todo add log here
		return nil, err
	}

	return embed[0], nil
}
