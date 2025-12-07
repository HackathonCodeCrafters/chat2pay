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

func (c *MistralLLM) ChatWithHistory(ctx context.Context, userMessage string) (string, error) {
	sessionId := ctx.Value("session_id").(string)
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

		Classify the user's latest message into ONE of these categories:
		
		1. chit_chat
		- Casual conversation unrelated to shopping.
		Examples:
		"hi", "how are you", "where are you from"
		
		2. general_product_request
		- User asks about a product category WITHOUT specific requirements.
		- They mention only the PRODUCT TYPE or ask for advice.
		Examples:
		"I need a mouse", 
		"Looking for a keyboard", 
		"Any recommendation for a laptop?",
		"Saya ingin cari mouse, ada rekomendasi?"
		
		⚠ IMPORTANT:
		- If the user only says a product name + "cocok / bagus / rekomendasi", treat it as general request.
		
		3. specific_product_search
		- User provides ANY detail such as:
		  - Budget (expensive / murah / under 1M)
		  - Features (wireless, bluetooth, DPI adjustable, silent click)
		  - Purpose (gaming, design, work, traveling)
		  - Brand (Logitech, Razer)
		  - Specs (RGB, rechargeable, ergonomic)
		Examples:
		"mouse wireless murah",
		"best mouse for designer",
		"wireless keyboard under $50",
		"ergonomic mouse for wrist pain"
		
		4. follow_up
		- User message refers to previous discussion.
		Examples:
		"yes wireless",
		"any cheaper?",
		"show more options"
		
		----
		
		Rules:
		- MUST return plain text intent (not JSON).
		- No explanation or extra words.
		- Response must be exactly one of:
		  chit_chat
		  general_product_request
		  specific_product_search
		  follow_up
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

func (c *MistralLLM) NewConnection(ctx context.Context) error {
	sessionId := ctx.Value("session_id").(string)
	b, _ := json.Marshal([]llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextPart("" +
					"You are a product shopping assistant. " +
					"Help the user choose their product and give recommendations according to their needs!" +
					"speak ONLY native indonesian"),
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

func (c *MistralLLM) ClassifyIntent(ctx context.Context, userMessage string) (*ChatClassify, error) {
	systemPrompt := `
		You are an Intent Classification AI for a shopping assistant in Indonesia.
		You MUST understand Indonesian language (Bahasa Indonesia).
		
		Your job is to classify the user's latest message into ONE of these intent categories:
		
		1. chit_chat
		   - Greeting or casual conversation not related to shopping.
		   Example: "hi", "halo", "apa kabar", "tell me a joke"
		
		2. general_product_request
		   - The user expresses interest in a product but is not specific enough.
		   Example: "I need a mouse", "cari laptop", "mau beli hp", "ada makanan gak"
		
		3. specific_product_search
		   - The user mentions ANY product name, category, type, or keyword.
		   - Even simple product searches should be classified here.
		   Example:
		   - "basreng" (snack name)
		   - "makanan ringan" (snacks)
		   - "laptop gaming"
		   - "cariin produk basreng"
		   - "cari makanan"
		   - "headphone murah"
		   - "tas ransel"
		
		4. follow_up
		   - The user message references a previous answer or continues the context.
		   Example:
		   - "yang lebih murah"
		   - "ada warna lain?"
		   - "show me more"
		
		---
		
		IMPORTANT: If the user mentions ANY product keyword in Indonesian or English, 
		classify as "specific_product_search". Be generous with this classification.
		
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
	systemPrompt := `Kamu adalah asisten belanja online yang ramah dan membantu.
Kamu HARUS menjawab dalam Bahasa Indonesia.
Bantu pengguna mencari produk yang mereka butuhkan.
Jawab dengan singkat dan jelas.`

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
	var result string
	for _, gen := range resp.Choices {
		result = gen.Content

	}

	// Basic JSON extraction (you could decode properly)
	return result, nil
}

func (c *MistralLLM) GetLastMessageContext(ctx context.Context) (string, error) {
	sessionId := ctx.Value("session_id")
	val, err := c.redisClient.Get(ctx, fmt.Sprintf(`history_context:%s`, sessionId))
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
