/**
 * Composable for managing a WebSocket connection to the notification service.
 * Handles connection lifecycle, auto-reconnect with exponential backoff,
 * and message dispatching.
 */
export function useWebSocket() {
  const authStore = useAuthStore()
  let ws = null
  let reconnectTimer = null
  let reconnectAttempts = 0
  const maxReconnectAttempts = 10
  const baseReconnectDelay = 1000 // 1 second

  const isConnected = ref(false)

  /**
   * Get the JWT token from the auth store or cookie.
   */
  function getToken() {
    return authStore?.token || useCookie('token').value || ''
  }

  /**
   * Build the WebSocket URL from the gateway base URL.
   * Converts http(s) to ws(s) and appends the notification ws path.
   */
  function buildWsUrl() {
    const gatewayBaseUrl = useRuntimeConfig().public.apiGatewayBaseUrl
    const wsUrl = gatewayBaseUrl.replace(/^https:/, 'wss:').replace(/^http:/, 'ws:')
    const token = getToken()

    if (authStore?.userRole?.toLowerCase() === 'customer') {
      // Pass JWT token as query param for customer WebSocket connections
      return `${wsUrl}/notifications/ws?token=${encodeURIComponent(token)}`
    } else {
      // Pass JWT token as query param for admin WebSocket connections
      return `${wsUrl}/notifications/admin/ws?token=${encodeURIComponent(token)}`
    }
  }

  /**
   * Connect to the WebSocket server.
   * @param {(data: object) => void} onMessage - Callback for incoming messages
   * @param {() => void} onReconnect - Callback when successfully reconnected after a disconnect
   */
  function connect(onMessage, onReconnect) {
    if (!getToken()) {
      console.warn('[WebSocket] No token available, skipping connection')
      return
    }

    // Prevent duplicate connections
    if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
      return
    }

    const url = buildWsUrl()

    try {
      ws = new WebSocket(url)
    } catch (err) {
      console.error('[WebSocket] Failed to create WebSocket:', err)
      scheduleReconnect(onMessage, onReconnect)
      return
    }

    ws.onopen = () => {
      console.log('[WebSocket] Connected')
      isConnected.value = true

      // If this is a reconnection (attempts > 0), trigger the onReconnect callback
      if (reconnectAttempts > 0 && onReconnect) {
        onReconnect()
      }

      reconnectAttempts = 0
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (onMessage) {
          onMessage(data.Payload)
        }
      } catch (err) {
        console.error('[WebSocket] Failed to parse message:', err)
      }
    }

    ws.onclose = (event) => {
      console.log('[WebSocket] Disconnected:', event.code, event.reason)
      isConnected.value = false
      ws = null

      // Only auto-reconnect if not intentionally closed
      if (event.code !== 1000) {
        scheduleReconnect(onMessage, onReconnect)
      }
    }

    ws.onerror = (err) => {
      console.error('[WebSocket] Error:', err)
    }
  }

  /**
   * Schedule a reconnection attempt with exponential backoff.
   */
  function scheduleReconnect(onMessage, onReconnect) {
    if (reconnectAttempts >= maxReconnectAttempts) {
      console.warn('[WebSocket] Max reconnect attempts reached. Giving up.')
      return
    }

    const delay = baseReconnectDelay * Math.pow(2, reconnectAttempts)
    reconnectAttempts++

    console.log(
      `[WebSocket] Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})...`
    )

    reconnectTimer = setTimeout(() => {
      connect(onMessage, onReconnect)
    }, delay)
  }

  /**
   * Gracefully close the WebSocket connection.
   */
  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }

    reconnectAttempts = maxReconnectAttempts // Prevent auto-reconnect

    if (ws) {
      ws.close(1000, 'Client disconnect')
      ws = null
    }

    isConnected.value = false
  }

  return {
    isConnected,
    connect,
    disconnect
  }
}
