package bootstrap

import (
	"chat2pay/config/yaml"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sarulabs/di/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewAdapter() *[]di.Def {
	return &[]di.Def{
		{
			Name:  DatabaseAdapter,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				config := ctn.Get(ConfigDefName).(*yaml.Config)
				// Generate DSN string from config
				var generateConnectionString = func() string {
					return fmt.Sprintf(
						"host=%s port=%s dbname=%s user=%s password=%s sslmode=%v application_name=%s",
						config.DB.Host,
						config.DB.Port,
						config.DB.DbName,
						config.DB.Username,
						config.DB.Password,
						config.DB.SSLMode,
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
		//{
		//	Name:  SocketAdapter,
		//	Scope: di.App,
		//	Build: func(ctn di.Container) (interface{}, error) {
		//		config := ctn.Get(ConfigDefName).(*yaml.Config)
		//		// Define the WebSocket server URL
		//		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
		//
		//		// Establish the WebSocket connection
		//		ws, err := websocket.Dial(u.String(), "", u.String())
		//		if err != nil {
		//			log.Fatalf("Failed to dial WebSocket: %v", err)
		//		}
		//		defer ws.Close()
		//
		//		fmt.Println("Connected to WebSocket server.")
		//
		//		// Send a message
		//		message := []byte("Hello from Go client!")
		//		_, err = ws.Write(message)
		//		if err != nil {
		//			log.Fatalf("Failed to write to WebSocket: %v", err)
		//		}
		//		fmt.Printf("Sent: %s\n", message)
		//
		//		// Receive a message
		//		var msg = make([]byte, 512)
		//		n, err := ws.Read(msg)
		//		if err != nil {
		//			log.Fatalf("Failed to read from WebSocket: %v", err)
		//		}
		//		fmt.Printf("Received: %s\n", msg[:n])
		//
		//		// Keep the connection open for a bit
		//		time.Sleep(2 * time.Second)
		//		return db, err
		//	},
		//	Close: func(obj interface{}) error {
		//		return obj.(*sqlx.DB).Close()
		//	},
		//},
	}
}
