package llm

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/pkg/llm/kolosal"
	"chat2pay/internal/pkg/llm/mistral"
	"chat2pay/internal/pkg/redis"
	"context"
	"github.com/tmc/langchaingo/llms"
)

type (
	ChatClassify struct {
		Intent string `json:"intent"`
	}
)

type LLM interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
	GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error)

	//Custom behavior
	ChatWithHistory(ctx context.Context, sessionId string, userMessage string) (string, error)
	ClassifyIntent(ctx context.Context, userMessage string) (string, error)
	NewConnection(ctx context.Context, sessionId string) error
	GetLastMessageContext(ctx context.Context, sessionID string) (string, error)
}

type llm struct {
	llm LLM
}

func NewLLM(cfg *yaml.Config, redisClient redis.RedisClient) LLM {
	switch cfg.LLM.Provider {
	case "kolosal":
		llmProvider := kolosal.NewKolosalLLM(cfg.LLM.Kolosal.URL, cfg.LLM.Kolosal.APIKey, cfg.LLM.Kolosal.ModelName)
		return &llm{
			llm: llmProvider,
		}

	case "mistral":
		llmProvider := mistral.NewMistralLLM(cfg.LLM.Mistral.APIKey, redisClient)

		return &llm{
			llm: llmProvider,
		}
	default:
		llmProvider := mistral.NewMistralLLM(cfg.LLM.Mistral.APIKey, redisClient)

		return &llm{
			llm: llmProvider,
		}
	}
}

func (l *llm) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return l.llm.GenerateContent(ctx, messages)
}

func (l *llm) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return l.llm.Call(ctx, prompt, options...)
}

func (l *llm) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	return l.llm.EmbedDocuments(ctx, texts)
}

func (l *llm) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return l.llm.EmbedQuery(ctx, text)
}

func (c *llm) ChatWithHistory(ctx context.Context, sessionId string, userMessage string) (string, error) {
	return c.llm.ChatWithHistory(ctx, sessionId, userMessage)
}

func (c *llm) ClassifyIntent(ctx context.Context, userMessage string) (string, error) {
	return c.llm.ClassifyIntent(ctx, userMessage)
}

func (c *llm) NewConnection(ctx context.Context, sessionId string) error {
	return c.llm.NewConnection(ctx, sessionId)
}

func (c *llm) GetLastMessageContext(ctx context.Context, sessionID string) (string, error) {
	return c.llm.GetLastMessageContext(ctx, sessionID)
}
