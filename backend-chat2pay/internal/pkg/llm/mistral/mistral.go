package mistral

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
)

// MistralLLM implements the LLM interface for your hosted model
type MistralLLM struct {
	mistral *mistral.Model
}

// NewMistralLLM creates a new instance of your custom LLM.
func NewMistralLLM(apiKey string) *MistralLLM {
	mistral, err := mistral.New(mistral.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}
	return &MistralLLM{
		mistral: mistral,
	}
}

// Call implements the [llms.Model] interface.
func (c *MistralLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, c, prompt, options...)
}

func (c *MistralLLM) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return c.mistral.GenerateContent(ctx, messages)
}

func (c *MistralLLM) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	embed, err := c.mistral.CreateEmbedding(ctx, texts)
	if err != nil {
		// Todo add log here
		return nil, err
	}

	return embed, nil
}

func (c *MistralLLM) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embed, err := c.mistral.CreateEmbedding(ctx, []string{text})
	if err != nil {
		// Todo add log here
		return nil, err
	}

	return embed[0], nil
}
