package utils

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	wsClients     = make(map[int64]*websocket.Conn)
	wsClientMutex = sync.RWMutex{}
)

func AddWebSocketConn(userId int64, conn *websocket.Conn) {
	wsClientMutex.Lock()
	defer wsClientMutex.Unlock()
	wsClients[userId] = conn
}

func GetWebSocketConn(userId int64) *websocket.Conn {
	wsClientMutex.RLock()
	defer wsClientMutex.RUnlock()
	return wsClients[userId]
}

func RemoveWebSocketConn(userId int64) {
	wsClientMutex.Lock()
	defer wsClientMutex.Unlock()
	delete(wsClients, userId)
}
