package service

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/pkg/ai"
	"chat2pay/internal/pkg/logger"
	"context"
	"errors"
	"fmt"
)

type AiService interface {
	Prompt(ctx context.Context, aiAsk *entities.AiAsk) *presenter.Response
}

type aiService struct {
	cfg   *yaml.Config
	model ai.AIModel
}

func NewAiService(cfg *yaml.Config, model ai.AIModel) AiService {
	return &aiService{
		cfg:   cfg,
		model: model,
	}
}

func (s *aiService) Prompt(ctx context.Context, aiAsk *entities.AiAsk) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("login_service", s.cfg.Logger.Enable)
	)

	result, err := s.model.Prompt("")

	if err != nil {
		log.Error(fmt.Sprintf(`error generating jwt token got %s`, err))
		return response.WithCode(500).WithError(errors.New("something went wrong!"))
	}

	return response.WithCode(200).WithData(result)
}
