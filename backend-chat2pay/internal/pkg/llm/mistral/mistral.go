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
		mistral.WithModel("ministral-14b-latest"),
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
	
	Classify the user's message into EXACTLY one of these:
	
	1. chit_chat
	- Small talk unrelated to products.
	Examples: "hi", "apa kabar", "lagi apa?"
	
	2. general_product_request
	- User mentions a product category WITHOUT specific requirements.
	- First time asking about a product without details.
	Examples:
	"I want a new watch"
	"Looking for a laptop"
	"Ada mouse bagus?"
	"Cari laptop dong"
	
	3. specific_product_search
	- User EXPLICITLY mentions a specific product with clear requirements in ONE message.
	- Must include BOTH product type AND at least one specific requirement.
	Examples:
	"Laptop gaming budget 15 juta"
	"Mouse wireless LOGITECH"
	"Laptop untuk desain grafis"
	"iPhone dengan RAM besar"
	"Cari laptop harga 15 jutaan ke atas"
	
	4. product_clarification
	- User is ANSWERING a previous question to provide more details for product search.
	- User provides additional preferences/budget/use-case AFTER being asked.
	Examples:
	"buat sehari-hari"
	"maksimal 17 juta"
	"yang penting speknya oke"
	"untuk gaming sih"
	"budget sekitar 10 juta"
	"gak ada sih, yang penting bagus"
	"paling buat game aja sih"
	
	5. product_question
	- User is asking questions ABOUT a product that was already shown/recommended.
	- User wants explanation, specs, or reasoning about shown products.
	- User asks WHY a product was recommended.
	Examples:
	"kenapa kamu menyarankan ini?"
	"speknya apa?"
	"apa kelebihannya?"
	"kenapa ini cocok?"
	"jelaskan lebih detail"
	"fiturnya apa saja?"
	"bisa jelasin gak?"
	"apa bedanya dengan yang lain?"
	"review nya gimana?"
	
	6. follow_up
	- User asks for alternatives or modifications to shown products.
	- User wants to see more options or different products.
	Examples:
	"yang lebih murah"
	"ada warna lain?"
	"ada yang lain?"
	"show more options"
	"yang lebih bagus?"
	"ada alternatif?"
	
	IMPORTANT RULES:
	- If user asks "kenapa", "mengapa", "apa speknya", "jelaskan" about shown products → "product_question"
	- If user provides preferences/budget as answer to clarifying question → "product_clarification"  
	- If user asks for alternatives/more options → "follow_up"
	- Output MUST be ONLY one of: chit_chat, general_product_request, specific_product_search, product_clarification, product_question, follow_up
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

// func (c *MistralLLM) ClassifyIntent(ctx context.Context, userMessage string) (*ChatClassify, error) {
// 	systemPrompt := `
// 		You are an Intent Classification AI for a shopping assistant in Indonesia.
// 		You MUST understand Indonesian language (Bahasa Indonesia).
		
// 		Your job is to classify the user's latest message into ONE of these intent categories:
		
// 		1. chit_chat
// 		   - Greeting or casual conversation not related to shopping.
// 		   Example: "hi", "halo", "apa kabar", "tell me a joke"
		
// 		2. general_product_request
// 		   - The user expresses interest in a product but is not specific enough.
// 		   Example: "I need a mouse", "cari laptop", "mau beli hp", "ada makanan gak"
		
// 		3. specific_product_search
// 		   - The user mentions ANY product name, category, type, or keyword.
// 		   - Even simple product searches should be classified here.
// 		   Example:
// 		   - "basreng" (snack name)
// 		   - "makanan ringan" (snacks)
// 		   - "laptop gaming"
// 		   - "cariin produk basreng"
// 		   - "cari makanan"
// 		   - "headphone murah"
// 		   - "tas ransel"
		
// 		4. follow_up
// 		   - The user message references a previous answer or continues the context.
// 		   Example:
// 		   - "yang lebih murah"
// 		   - "ada warna lain?"
// 		   - "show me more"
		
// 		---
		
// 		IMPORTANT: If the user mentions ANY product keyword in Indonesian or English, 
// 		classify as "specific_product_search". Be generous with this classification.
		
// 		RESPONSE FORMAT (MUST BE JSON, NO EXTRA TEXT):
		
// 		{
// 		  "intent": "<one of: chit_chat, general_product_request, specific_product_search, follow_up>"
// 		}
		
// 		Rules:
// 		- Do NOT explain the classification.
// 		- Do NOT include confidence scores.
// 		- Do NOT respond as a chatbot.
// 		- ONLY return JSON with the detected intent.
// 		`

// 	messages := []llms.MessageContent{
// 		{
// 			Role:  llms.ChatMessageTypeSystem,
// 			Parts: []llms.ContentPart{llms.TextPart(systemPrompt)},
// 		},
// 		{
// 			Role:  llms.ChatMessageTypeHuman,
// 			Parts: []llms.ContentPart{llms.TextPart(userMessage)},
// 		},
// 	}

// 	resp, err := c.GenerateContent(ctx, messages)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Extract text from response safely
// 	var (
// 		result   string
// 		classify ChatClassify
// 	)

// 	for _, gen := range resp.Choices {
// 		result = gen.Content

// 	}

// 	if err = json.Unmarshal([]byte(result), &classify); err != nil {
// 		return nil, err
// 	}

// 	// Basic JSON extraction (you could decode properly)
// 	return &classify, nil
// }

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

	sessionId, ok := ctx.Value("session_id").(string)
	if !ok || sessionId == "" {
		return "", fmt.Errorf("missing session_id in context")
	}

	key := fmt.Sprintf("history_context:%s", sessionId)

	// Fetch Redis history
	raw, err := c.redisClient.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if raw == nil {
		return "", nil
	}

	var history []llms.MessageContent

	if raw != nil {
		if err := json.Unmarshal([]byte(*raw), &history); err != nil {
			// corrupt JSON → reset history
			history = []llms.MessageContent{}
		}
	}

	// ---- Prevent panic: check length ----
	if len(history) == 0 {
		return "", fmt.Errorf("no history available")
	}

	// If only 1 message exists, return that safely
	if len(history) == 1 {
		return extractMessage(history[0]), nil
	}

	// If 2 or more messages → return last AI message
	last := history[len(history)-1]

	return extractMessage(last), nil
}

func extractMessage(msg llms.MessageContent) string {
	if len(msg.Parts) == 0 {
		return ""
	}
	if text, ok := msg.Parts[0].(llms.TextContent); ok {
		return text.Text
	}
	return ""
}

func contentPartToString(part llms.ContentPart) string {
	switch v := part.(type) {
	case llms.TextContent: // normal text
		return v.Text

	default:
		return ""
	}
}
