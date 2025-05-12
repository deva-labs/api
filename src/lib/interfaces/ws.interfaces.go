package interfaces

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type SafeConn struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

type WorkflowStep struct {
	Name     string
	Command  string
	Action   string
	EnvVars  map[string]string
	Required bool
}

func (sc *SafeConn) SafeWrite(msgType int, data []byte) error {
	sc.Mu.Lock()
	defer sc.Mu.Unlock()
	return sc.Conn.WriteMessage(msgType, data)
}
