package mistral

import (
	"context"
	"encoding/json"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
)

type ChatClassify struct {
	Intent string `json:"intent"`
}

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

func (c *MistralLLM) ClassifyIntent(ctx context.Context, userMessage string) (*ChatClassify, error) {
	systemPrompt := `
		You are an Intent Classification AI for a shopping assistant.
		
		Your job is to classify the user's latest message into ONE of these intent categories:
		
		1. chit_chat
		   - Greeting or casual conversation not related to shopping.
		   Example: "hi", "how are you", "tell me a joke"
		
		2. general_product_request
		   - The user expresses interest in a product but is not specific enough.
		   Example: "I need a mouse", "I'm looking for headphones", 
					"recommend a laptop", "I want a new phone"
		
		3. specific_product_search
		   - The user provides enough details such as budget, specification, brand,
			 use-case, feature request, or preference.
		   Example:
		   - "Best wireless mouse under $50"
		   - "Gaming keyboard with RGB"
		   - "Headphones for gym"
		   - "Ergonomic desk chair for back pain"
		
		4. follow_up
		   - The user message references a previous answer or continues the context.
		   Example:
		   - "yes wireless"
		   - "cheaper one"
		   - "any other option?"
		   - "show me more models"
		
		---
		
		RESPONSE FORMAT (MUST BE JSON, NO EXTRA TEXT):
		
		{
		  "intent": "<one of: chit_chat, general_product_request, specific_product_search, follow_up>"
		}
		
		Rules:
		- Do NOT explain the classification.
		- Do NOT include confidence scores.
		- Do NOT respond as a chatbot.
		- ONLY return JSON with the detected intent.
		`

	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextPart(systemPrompt)},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(userMessage)},
		},
	}

	resp, err := c.GenerateContent(ctx, messages)
	if err != nil {
		return nil, err
	}

	// Extract text from response safely
	var (
		result   string
		classify ChatClassify
	)

	for _, gen := range resp.Choices {
		result = gen.Content

	}

	if err = json.Unmarshal([]byte(result), &classify); err != nil {
		return nil, err
	}

	// Basic JSON extraction (you could decode properly)
	return &classify, nil
}

func (c *MistralLLM) Chat(ctx context.Context, userMessage string) (string, error) {
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(userMessage)},
		},
	}

	resp, err := c.GenerateContent(ctx, messages)
	if err != nil {
		return "", err
	}

	// Extract text from response safely
	var result string
	for _, gen := range resp.Choices {
		result = gen.Content

	}

	// Basic JSON extraction (you could decode properly)
	return result, nil
}
