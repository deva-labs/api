package store

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	userSocketMap = make(map[uint]*websocket.Conn)
	mapMutex      sync.RWMutex // Protects concurrent access
)

// SetUserSocket stores the WebSocket connection for a user
func SetUserSocket(userID uint, conn *websocket.Conn) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	userSocketMap[userID] = conn
}

// GetUserSocket retrieves the WebSocket connection for a user
func GetUserSocket(userID uint) (*websocket.Conn, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	conn, exists := userSocketMap[userID]
	return conn, exists
}

// RemoveUserSocket deletes the WebSocket connection for a user
func RemoveUserSocket(userID uint) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	delete(userSocketMap, userID)
}
