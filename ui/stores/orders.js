import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useOrderStore = defineStore('orders', {
  state: () => ({
    orders: [],
    loading: false,
    error: null,
    order: {},
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

    async fetchOrders(params = {}) {
      try {
        this.loading = true
        const { search, page, limit, status } = params

        const query = new URLSearchParams()
        if (search) query.append('search', search)
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)
        if (status) query.append('status', status)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/orders?${query.toString()}`,
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

        this.orders = result.data
        this.pagination = result.pagination

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

    async fetchOrdersAdmin(params = {}) {
      try {
        this.loading = true
        const { search, page, limit, status } = params

        const query = new URLSearchParams()
        if (search) query.append('search', search)
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)
        if (status) query.append('status', status)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/orders/admin?${query.toString()}`,
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

        this.orders = result.data
        this.pagination = result.pagination

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

    async createOrders(orderData, lat, lng, idempotencyKey) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/orders?lat=${lat}&lng=${lng}`,
          {
            method: 'POST',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json',
              'X-Idempotency-Key': idempotencyKey
            },
            body: JSON.stringify(orderData),
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.order = result.data

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

    async getDetailOrders(orderId, isAdmin = false) {
      try {
        this.loading = true

        let url = `${useRuntimeConfig().public.apiGatewayBaseUrl}/orders/${orderId}`
        if (isAdmin) {
          url = `${useRuntimeConfig().public.apiGatewayBaseUrl}/orders/${orderId}/admin`
        }

        const response = await fetch(url, {
          method: 'GET',
          headers: {
            Authorization: `Bearer ${useCookie('token').value}`,
            'Content-Type': 'application/json'
          },
          signal: this._createSignal()
        })

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.order = result.data

        return result
      } catch (err) {
        if (err.name === 'AbortError') {
          throw err
        }
        this._showError(err)
      } finally {
        this.loading = false
      }
    },

    async updateStatusOrder(orderId, dataUpdate) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/orders/${orderId}/status`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(dataUpdate),
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        return result
      } catch (err) {
        if (err.name === 'AbortError') {
          throw err
        }
        this._showError(err)
      } finally {
        this.loading = false
      }
    }
  }
})
