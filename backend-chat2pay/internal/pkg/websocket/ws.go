package websocket

import (
	"chat2pay/config/yaml"
	"context"
	"golang.org/x/net/websocket"
	"net/url"
)

type Websocket interface {
	Write(ctx context.Context, message string) error
}

type socket struct {
	ws *websocket.Conn
}

func NewWebsocket(config yaml.Config) Websocket {
	u := url.URL{Scheme: "ws", Host: config.Websocket.Host, Path: config.Websocket.Path}

	// Establish the WebSocket connection
	ws, err := websocket.Dial(u.String(), "", u.String())
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	return &socket{
		ws: ws,
	}
}

func (s *socket) Write(ctx context.Context, message string) error {
	_, err := s.ws.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}
