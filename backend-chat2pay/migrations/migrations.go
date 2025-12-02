package migrations

import (
	"chat2pay/internal/entities"
	"gorm.io/gorm"
)

var models = []interface{}{
	&entities.User{},
	&entities.Merchant{},
	&entities.MerchantUser{},
	&entities.Outlet{},
	&entities.ProductCategory{},
	&entities.Customer{},
	&entities.Product{},
	&entities.ProductImage{},
	&entities.Order{},
	&entities.OrderItem{},
	&entities.Payment{},
	&entities.Shipment{},
	&entities.Conversation{},
	&entities.Message{},
}

func AutoMigration(db *gorm.DB) {
	// Disabled AutoMigrate because we use SQL migrations (001_create_merchants_table.sql)
	// AutoMigrate can cause issues with existing constraints and enums in PostgreSQL
	// If you need to add new tables, update the SQL migration file instead

	// db.AutoMigrate(models...)
}
