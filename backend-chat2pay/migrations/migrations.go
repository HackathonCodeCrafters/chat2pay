package migrations

import (
	"chat2pay/internal/entities"
	"gorm.io/gorm"
)

var models = []interface{}{
	&entities.User{},
}

func AutoMigration(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models...)
}
