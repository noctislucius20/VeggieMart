package proxy

import (
	"api-gateway/response"
	"api-gateway/utils/helper"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ForwardWebSocket proxies a WebSocket connection from the client to a backend service.
// It upgrades the incoming HTTP request to WebSocket on both the client side and the
// backend side, then bidirectionally pipes messages between them.
func ForwardWebSocket(c echo.Context, targetUrl string) error {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return c.JSON(http.StatusBadGateway, response.ResponseFailed("invalid target url"))
	}

	// Convert http(s) scheme to ws(s) for the backend dial
	switch target.Scheme {
	case "http":
		target.Scheme = "ws"
	case "https":
		target.Scheme = "wss"
	}

	// Prepare headers to forward to the backend WebSocket server
	reqHeader := http.Header{}

	helper.SetGatewayHeaders(c, &reqHeader)
	target.RawQuery = c.Request().URL.RawQuery

	// Dial the backend WebSocket server
	backendConn, resp, err := websocket.DefaultDialer.Dial(target.String(), reqHeader)
	if err != nil {
		if resp != nil {
			return c.JSON(resp.StatusCode, response.ResponseFailed("failed to connect to backend WebSocket"))
		}
		return c.JSON(http.StatusBadGateway, response.ResponseFailed("failed to connect to backend WebSocket"))
	}
	defer backendConn.Close()

	// Upgrade the client connection to WebSocket
	clientConn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer clientConn.Close()

	errChan := make(chan error, 2)

	// Client → Backend
	go func() {
		for {
			msgType, msg, err := clientConn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if err := backendConn.WriteMessage(msgType, msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Backend → Client
	go func() {
		for {
			msgType, msg, err := backendConn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if err := clientConn.WriteMessage(msgType, msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Wait until one side closes, then let deferred Close() clean up both
	<-errChan

	// Attempt graceful close on both sides (best-effort, ignore errors)
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	_ = clientConn.WriteMessage(websocket.CloseMessage, closeMsg)
	_ = backendConn.WriteMessage(websocket.CloseMessage, closeMsg)

	// Drain the response body if present
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	return nil
}
