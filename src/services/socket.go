package services

import (
	"dockerwizard-api/src/utils"
	"dockerwizard-api/store"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketMessage represents the JSON message
type WebSocketMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// WebSocketUpgrader handles WebSocket upgrade and events
func WebSocketUpgrader() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userID := c.Query("userId")
		if userID == "" {
			log.Println("Missing userId in query params")
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"error":"userId is required"}`))
			_ = c.Close()
			return
		}

		uid := utils.ConvertStringToUint(userID)
		store.SetUserSocket(uid, c)
		log.Printf("Client connected: %s (userId: %d)", c.RemoteAddr(), uid)

		defer func() {
			_ = c.Close()
			store.RemoveUserSocket(uid)
			log.Printf("Client disconnected: %s (userId: %d)", c.RemoteAddr(), uid)
		}()

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				break
			}

			var wsMessage WebSocketMessage
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				log.Printf("JSON unmarshal error: %v", err)
				continue
			}

			handleMessage(c, messageType, wsMessage)
		}
	})
}

// handleMessage processes incoming WebSocket messages
func handleMessage(conn *websocket.Conn, messageType int, wsMessage WebSocketMessage) {
	var response WebSocketMessage

	switch wsMessage.Type {
	case "notice":
		response = WebSocketMessage{Type: "notice", Message: "Received: " + wsMessage.Message}
	case "bye":
		response = WebSocketMessage{Type: "bye", Message: "Goodbye: " + wsMessage.Message}
	case "ping":
		response = WebSocketMessage{Type: "ping", Message: "Pong: " + wsMessage.Message}
	default:
		response = WebSocketMessage{Type: "echo", Message: "Echo: " + wsMessage.Message}
	}

	if err := sendJSONMessage(conn, messageType, response); err != nil {
		log.Printf("Send error: %v", err)
	}
}

// sendJSONMessage sends a JSON-encoded message to the client
func sendJSONMessage(conn *websocket.Conn, messageType int, response WebSocketMessage) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return conn.WriteMessage(messageType, data)
}

// SendMessageToUser sends a message to a specific user via their WebSocket connection
func SendMessageToUser(userID uint, message WebSocketMessage) error {
	conn, exists := store.GetUserSocket(userID)
	if !exists {
		return fmt.Errorf("user %d not connected", userID)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}
