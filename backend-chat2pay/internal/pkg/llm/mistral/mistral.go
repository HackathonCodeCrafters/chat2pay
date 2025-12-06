package mistral

import (
	"chat2pay/internal/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
)

// MistralLLM implements the LLM interface for your hosted model
type MistralLLM struct {
	mistral     *mistral.Model
	redisClient redis.RedisClient
}

// NewMistralLLM creates a new instance of your custom LLM.
func NewMistralLLM(apiKey string, client redis.RedisClient) *MistralLLM {
	mistral, err := mistral.New(mistral.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}
	return &MistralLLM{
		mistral:     mistral,
		redisClient: client,
	}
}

func (c *MistralLLM) ChatWithHistory(ctx context.Context, sessionId string, userMessage string) (string, error) {
	val, err := c.redisClient.Get(ctx, fmt.Sprintf(`history_context:%s`, sessionId))
	if err != nil {
		return "", err
	}

	history := []llms.MessageContent{}
	if val != nil {
		err = json.Unmarshal([]byte(*val), &history)
		if err != nil {
			return "", nil
		}
	}

	history = append(history, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextPart(userMessage)},
	})

	resp, err := c.GenerateContent(ctx, history)
	if err != nil {
		return "", err
	}

	var result string
	for _, gen := range resp.Choices {
		result = gen.Content
	}

	history = append(history, llms.MessageContent{
		Role:  llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{llms.TextPart(result)},
	})

	b, _ := json.Marshal(history)

	_, err = c.redisClient.Set(ctx, fmt.Sprintf(`history_context:%s`, sessionId), string(b))
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *MistralLLM) ClassifyIntent(ctx context.Context, userMessage string) (string, error) {
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
		
		RESPONSE FORMAT (MUST BE STRING, NO EXTRA TEXT):
		
		<one of: chit_chat, general_product_request, specific_product_search, follow_up>
        example: general_product_request
		
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
		return "", err
	}

	// Extract text from response safely
	var (
		result string
	)

	for _, gen := range resp.Choices {
		result = gen.Content
	}

	// Basic JSON extraction (you could decode properly)
	return result, nil
}

func (c *MistralLLM) NewConnection(ctx context.Context, sessionId string) error {

	b, _ := json.Marshal([]llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextPart("You are a product shopping assistant. Help the user with product recommendations."),
			},
		},
	})

	_, err := c.redisClient.Set(ctx, fmt.Sprintf(`history_context:%s`, sessionId), string(b))

	return err
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

func (c *MistralLLM) GetLastMessageContext(ctx context.Context, sessionID string) (string, error) {

	val, err := c.redisClient.Get(ctx, fmt.Sprintf(`history_context:%s`, sessionID))
	if err != nil {
		return "", err
	}

	history := []llms.MessageContent{}
	err = json.Unmarshal([]byte(*val), &history)
	if err != nil {
		return "", nil
	}

	part := history[len(history)-1].Parts[len(history[len(history)-1].Parts)]

	result := contentPartToString(part)

	fmt.Println("result --> ", result)
	return result, nil
}

func contentPartToString(part llms.ContentPart) string {
	switch v := part.(type) {
	case llms.TextContent: // normal text
		return v.Text

	default:
		return ""
	}
}
