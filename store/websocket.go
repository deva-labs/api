package store

import (
	"github.com/google/uuid"
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	userSocketMap = make(map[uuid.UUID]*websocket.Conn)
	mapMutex      sync.RWMutex // Protects concurrent access
)

// SetUserSocket stores the WebSocket connection for a users
func SetUserSocket(userID uuid.UUID, conn *websocket.Conn) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	userSocketMap[userID] = conn
}

// GetUserSocket retrieves the WebSocket connection for a users
func GetUserSocket(userID uuid.UUID) (*websocket.Conn, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	conn, exists := userSocketMap[userID]
	return conn, exists
}

// RemoveUserSocket deletes the WebSocket connection for a users
func RemoveUserSocket(userID uuid.UUID) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	delete(userSocketMap, userID)
}
