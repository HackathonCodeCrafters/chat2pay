package entities

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	ID         string          `json:"id" db:"id"`
	CustomerID string          `json:"customer_id" db:"customer_id"`
	Role       string          `json:"role" db:"role"` // "user" or "assistant"
	Content    string          `json:"content" db:"content"`
	Products   json.RawMessage `json:"products,omitempty" db:"products"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}
