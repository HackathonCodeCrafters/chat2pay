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
	mistral, err := mistral.New(
		mistral.WithAPIKey(apiKey),
		//mistral.WithModel("ministral-14b-latest"),
	)
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
	Your ONLY task is to classify the message. DO NOT answer the user.
	
	IMPORTANT: When user asks for specifications, details, or features of a product 
	(using words like "spesifikasi", "fitur", "detail", "apa saja"), ALWAYS classify as "follow_up".
	
	Classify the user's message into EXACTLY one of these:
	
	1. chit_chat
	- Small talk unrelated to products.
	Examples: "hi", "apa kabar", "lagi apa?"
	
	2. general_product_request
	- User mentions a product category WITHOUT ANY requirement.
	Examples:
	"I want a new watch"
	"Looking for a laptop"
	"Ada mouse bagus?"
	
	3. specific_product_search
	- User includes ANY requirement such as:
	- budget (murah, di bawah 1 juta)
	- brand
	- feature (waterproof, wireless, bluetooth)
	- use case (gaming, traveling)
	- specifications (RAM, storage, camera, etc.)
	Examples:
	"Jam tangan brand G-SHOCK"
	"mouse wireless LOGITECH"
	"Laptop untuk desain"
	"iPhone dengan RAM besar"
	
	4. follow_up
	- User refers to previous suggestions or asks for details about shown products.
	- Asking for specifications: "spesifikasi nya?", "fiturnya apa?", "detailnya?"
	- Asking for alternatives: "yang lebih murah", "ada warna lain?"
	- Requesting more information: "bisa lihat lebih detail?", "apa keunggulannya?"
	Examples:
	"yang lebih murah"
	"ada warna lain?"
	"show more options"
	"spesifikasinya apa?"
	"fiturnya apa saja?"
	"detail produknya?"
	"bisa jelaskan lebih lanjut?"
	"keunggulannya apa?"
	
	RULES:
	- Output MUST be ONLY one of the following:
	chit_chat
	general_product_request
	specific_product_search
	follow_up
	
	- No explanation, no formatting, no JSON. Just the label.`

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
