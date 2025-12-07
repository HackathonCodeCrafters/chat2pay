package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

type ChatMessageRepository interface {
	Create(ctx context.Context, msg *entities.ChatMessage) error
	FindByCustomerID(ctx context.Context, customerID string, limit int) ([]entities.ChatMessage, error)
	DeleteByCustomerID(ctx context.Context, customerID string) error
}

type chatMessageRepository struct {
	DB *sqlx.DB
}

func NewChatMessageRepository(db *sqlx.DB) ChatMessageRepository {
	return &chatMessageRepository{DB: db}
}

func (r *chatMessageRepository) Create(ctx context.Context, msg *entities.ChatMessage) error {
	var productsJSON []byte
	var err error

	if msg.Products != nil {
		productsJSON = msg.Products
	} else {
		productsJSON, _ = json.Marshal(nil)
	}

	query := `
		INSERT INTO chat_messages (id, customer_id, role, content, products)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.DB.ExecContext(ctx, query, msg.ID, msg.CustomerID, msg.Role, msg.Content, productsJSON)
	return err
}

func (r *chatMessageRepository) FindByCustomerID(ctx context.Context, customerID string, limit int) ([]entities.ChatMessage, error) {
	var messages []entities.ChatMessage

	query := `
		SELECT id, customer_id, role, content, products, created_at 
		FROM chat_messages 
		WHERE customer_id = $1 
		ORDER BY created_at ASC
		LIMIT $2
	`
	err := r.DB.SelectContext(ctx, &messages, query, customerID, limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *chatMessageRepository) DeleteByCustomerID(ctx context.Context, customerID string) error {
	query := `DELETE FROM chat_messages WHERE customer_id = $1`
	_, err := r.DB.ExecContext(ctx, query, customerID)
	return err
}
