package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	wsClients      = make(map[int64]*websocket.Conn)
	wsClientMutex  = sync.RWMutex{}
	adminWsClients = make(map[*websocket.Conn]bool)
	adminWsMutex   = sync.RWMutex{}
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

func AddAdminWebSocketConn(conn *websocket.Conn) {
	adminWsMutex.Lock()
	defer adminWsMutex.Unlock()
	adminWsClients[conn] = true
}

func RemoveAdminWebSocketConn(conn *websocket.Conn) {
	adminWsMutex.Lock()
	defer adminWsMutex.Unlock()
	delete(adminWsClients, conn)
}

func BroadcastAdminMessage(msg interface{}) {
	adminWsMutex.RLock()
	defer adminWsMutex.RUnlock()

	for conn := range adminWsClients {
		if err := conn.WriteJSON(msg); err != nil {
			conn.Close()
			delete(adminWsClients, conn)
		}
	}
}
