package dto

import (
	"chat2pay/internal/entities"
	"encoding/json"
	"time"
)

type ChatMessageRequest struct {
	Role     string          `json:"role" validate:"required"`
	Content  string          `json:"content" validate:"required"`
	Products json.RawMessage `json:"products,omitempty"`
}

type ChatMessageResponse struct {
	ID        string          `json:"id"`
	Role      string          `json:"role"`
	Content   string          `json:"content"`
	Products  json.RawMessage `json:"products,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type ChatHistoryResponse struct {
	Messages []ChatMessageResponse `json:"messages"`
}

func ToChatMessageResponse(msg *entities.ChatMessage) ChatMessageResponse {
	return ChatMessageResponse{
		ID:        msg.ID,
		Role:      msg.Role,
		Content:   msg.Content,
		Products:  msg.Products,
		CreatedAt: msg.CreatedAt,
	}
}

func ToChatHistoryResponse(messages []entities.ChatMessage) ChatHistoryResponse {
	resp := ChatHistoryResponse{
		Messages: make([]ChatMessageResponse, len(messages)),
	}
	for i, msg := range messages {
		resp.Messages[i] = ToChatMessageResponse(&msg)
	}
	return resp
}
