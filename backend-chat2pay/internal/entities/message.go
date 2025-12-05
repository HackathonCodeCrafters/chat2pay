package entities

import "time"

type Message struct {
	ID             string       `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID string       `gorm:"not null;index:idx_messages_conversation_id" json:"conversation_id"`
	SenderType     string       `gorm:"type:enum('customer','merchant_user','system');not null;index:idx_messages_sender_type_sender_id,priority:1" json:"sender_type"`
	SenderID       *string      `gorm:"index:idx_messages_sender_type_sender_id,priority:2" json:"sender_id,omitempty"`
	MessageText    string       `gorm:"type:text;not null" json:"message_text"`
	CreatedAt      time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
