package bootstrap

import (
	"chat2pay/config/yaml"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// simple db connection
func DatabaseConnection(config *yaml.Config) (*gorm.DB, error) {
	// refer https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL for details
	dsn := fmt.Sprintf(`host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta`,
		config.DB.Host,
		config.DB.Username,
		config.DB.Password,
		config.DB.DbName,
		config.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
