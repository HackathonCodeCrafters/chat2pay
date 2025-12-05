package kolosal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"net/http"
)

// KolosalLLM implements the LLM interface for your hosted model
type KolosalLLM struct {
	Url        string
	APIKey     string
	ModelName  string
	HTTPClient *http.Client
}

// NewKolosalLLM creates a new instance of your custom LLM.
func NewKolosalLLM(url, apiKey, modelName string) *KolosalLLM {
	return &KolosalLLM{
		Url:       url,
		APIKey:    apiKey,
		ModelName: modelName,
	}
}

// Call implements the [llms.Model] interface.
func (c *KolosalLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, c, prompt, options...)
}

func (c *KolosalLLM) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {

	// Convert messages → plain text for model prompt
	var finalPrompt string
	for _, msg := range messages {
		for _, part := range msg.Parts {
			if text, ok := part.(llms.TextContent); ok {
				finalPrompt += text.Text + "\n"
			}
		}
	}

	body := map[string]any{
		"model":      c.ModelName,
		"max_tokens": 500,
		"messages": []map[string]string{
			{
				"role":    "user", // simplified assumption
				"content": finalPrompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", c.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, errEmptyResponseFromModel
	}

	output := result.Choices[0].Message.Content

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: output,
			},
		},
	}, nil
}
func (c *KolosalLLM) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	// Implement your custom LLM's logic for generating embeddings.
	return nil, fmt.Errorf("embedding not implemented for custom LLM")
}

func (c *KolosalLLM) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	// Implement your custom LLM's logic for generating query embeddings.
	return nil, fmt.Errorf("query embedding not implemented for custom LLM")
}
