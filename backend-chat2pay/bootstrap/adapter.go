package bootstrap

import (
	"chat2pay/config/yaml"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

func loadAdapter(builder *di.Builder, config *yaml.Config) {
	builder.Add([]di.Def{
		{
			Name:  DatabaseAdapter,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				// Generate DSN string from config
				var generateConnectionString = func() string {
					return fmt.Sprintf(
						"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s application_name=%s",
						config.DB.Host,
						config.DB.Port,
						config.DB.DbName,
						config.DB.Username,
						config.DB.Password,
						true,
						config.App.Name,
					)
				}

				db, err := sqlx.Connect("postgres", generateConnectionString())
				if err != nil {
					log.Printf("Error while initialize db provider. Detail: %s", err.Error())
					return nil, err
				}
				db.SetMaxOpenConns(50)
				db.SetConnMaxLifetime(time.Minute * 15)
				db.SetMaxIdleConns(10)
				return db, err
			},
			Close: func(obj interface{}) error {
				return obj.(*sqlx.DB).Close()
			},
		},
	}...)
}
