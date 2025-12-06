package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/entities"
	"chat2pay/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

func NewSocketEvent(router fiber.Router, ctn di.Container) {

	// WS Upgrade check
	router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// ---- SOCKET.IO EVENTS ----
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Println("🔌 User connected:", ep.Kws.GetStringAttribute("user_id"))
	})

	// Handle chat messages
	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		fmt.Println("📨 Message:", string(ep.Data))

		productService := ctn.Get("product.service").(service.ProductService)

		response := productService.AskProduct(
			context.Background(),
			&dto.AskProduct{
				SessionId: ep.Kws.GetStringAttribute("user_id"),
				Prompt:    string(ep.Data)},
		)

		if response.Code == 200 {

			if data, ok := response.Data.(dto.LLMResponse); ok {

				if len(data.Products) != 0 {
					productsBytes, _ := json.Marshal(data.Products)
					ep.Kws.Emit(productsBytes, socketio.TextMessage)
				} else {
					ep.Kws.Emit([]byte(data.Message), socketio.TextMessage)
				}
			} else {
				fmt.Println("error converting struct")
			}

		} else {
			ep.Kws.Emit([]byte("Something went wrong, try again later!"), socketio.TextMessage)
		}

	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		delete(entities.SocketClients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Println("❌ User disconnected")
	})

	// ---- SOCKET.IO ENTRYPOINT ----
	router.Get("/ws/chat/:id", socketio.New(func(kws *socketio.Websocket) {
		userID := kws.Params("id")
		entities.SocketClients[userID] = kws.UUID
		kws.SetAttribute("user_id", userID)

		//llmPackage := ctn.Get("llm.package").(llm.LLM)

		//llmPackage.NewConnection(userID)

		kws.Emit([]byte(fmt.Sprintf("👋 Welcome %s!", userID)), socketio.TextMessage)
		//kws.Broadcast([]byte(fmt.Sprintf("📢 %s joined the chat", userID)), true, socketio.TextMessage)
	}))
}
