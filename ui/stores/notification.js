import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    loading: false,
    error: null,
    abortControllers: [],
    pagination: {
      page: 1,
      total_count: 0,
      per_page: 10,
      total_pages: 0
    }
  }),

  actions: {
    cancelRequests() {
      this.abortControllers.forEach((controller) => controller.abort())
      this.abortControllers = []
    },

    _createSignal() {
      const controller = new AbortController()
      this.abortControllers.push(controller)
      const signal = controller.signal
      signal.addEventListener('abort', () => {
        const index = this.abortControllers.indexOf(controller)
        if (index !== -1) {
          this.abortControllers.splice(index, 1)
        }
      })
      return signal
    },

    _showError(error) {
      this.error = error.message
      const { showError } = useErrorModal()
      showError(error.message)
      throw error
    },

    /**
     * Add a new notification from WebSocket push.
     * Inserts at the beginning (newest first) and increments unread count.
     */
    addNotification(notification) {
      this.notifications.unshift({
        id: notification.id || Date.now(),
        notification_method: notification.notification_method,
        notification_type: notification.notification_type,
        notification_type_id: notification.notification_type_id,
        subject: notification.subject,
        message: notification.message,
        sent_at: Date.now(),
        is_read: false
      })
      this.unreadCount++
    },

    /**
     * Fetch all notifications from the API.
     */
    async fetchPushNotifications(params = {}) {
      try {
        this.loading = true
        const { page, limit } = params

        const query = new URLSearchParams()
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/notifications/push?${query.toString()}`,
          {
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.notifications = result.data || []
        if (result.pagination) {
          this.pagination = result.pagination
        }

        // Calculate unread count from fetched data
        this.unreadCount = this.notifications.filter((n) => !n.is_read).length

        return result
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this._showError(error)
      } finally {
        this.loading = false
      }
    },

    /**
     * Mark a single notification as read.
     */
    async markAsRead(notificationId) {
      try {
        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/notifications/${notificationId}/read`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        // Update local state
        const notification = this.notifications.find((n) => n.id === notificationId)
        if (notification && !notification.is_read) {
          notification.is_read = true
          this.unreadCount = Math.max(0, this.unreadCount - 1)
        }

        return result
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this._showError(error)
      }
    },

    async markAsSent(notificationId) {
      try {
        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/notifications/${notificationId}/sent`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        return result
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this._showError(error)
      }
    },

    /**
     * Mark all notifications as read (local-only, batch).
     */
    markAllAsReadLocal() {
      this.notifications = this.notifications.map((n) => ({
        ...n,
        is_read: true
      }))
      this.unreadCount = 0
    },

    /**
     * Clear all notifications from state (e.g. on logout).
     */
    clearNotifications() {
      this.notifications = []
      this.unreadCount = 0
    }
  },

  getters: {
    unreadNotifications: (state) => state.notifications.filter((n) => !n.is_read),
    hasUnread: (state) => state.unreadCount > 0
  }
})
