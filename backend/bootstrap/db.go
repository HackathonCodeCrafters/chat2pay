package bootstrap

import (
	"chat2pay/config/yaml"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// simple db connection
func DatabaseConnection(config *yaml.Config) (*gorm.DB, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf(`%s:%s@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local`,
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
