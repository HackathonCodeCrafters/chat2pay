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
	//&entities.Order{},
	//&entities.OrderItem{},
	//&entities.Payment{},
	//&entities.Shipment{},
	//&entities.Conversation{},
	//&entities.Message{},
}

func AutoMigration(db *gorm.DB) {
	//db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models...)
	err := db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}
}
