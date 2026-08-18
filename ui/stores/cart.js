import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useCartStore = defineStore('carts', {
  state: () => ({
    carts: [],
    loading: false,
    error: null,
    abortControllers: []
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

    clearCarts() {
      this.carts = []
    },

    async fetchCarts() {
      try {
        this.loading = true
        this.error = null
        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/carts`,
          {
            method: 'GET',
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

        this.carts = Array.isArray(result.data) ? result.data : []

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

    async addToCart(productId, quantity) {
      try {
        this.loading = true
        this.error = null
        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/carts`,
          {
            method: 'POST',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              product_id: productId,
              quantity: quantity
            }),
            signal: this._createSignal()
          }
        )
        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        await this.fetchCarts()

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

    async deleteCart(productId) {
      try {
        this.loading = true
        this.error = null
        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/carts?product_id=${productId}`,
          {
            method: 'DELETE',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        if (!response.ok) {
          const result = await response.json().catch(() => ({}))
          throw new Error(result?.message || 'Something went wrong!')
        }

        await this.fetchCarts()

        // 204 No Content - no body to parse
        if (response.status === 204) {
          return null
        }

        return await response.json()
      } catch (err) {
        if (err.name === 'AbortError') {
          throw err
        }
        this._showError(err)
      } finally {
        this.loading = false
      }
    },

    async deleteAllCart() {
      try {
        this.loading = true
        this.error = null
        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/carts/all`,
          {
            method: 'DELETE',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        if (!response.ok) {
          const result = await response.json().catch(() => ({}))
          throw new Error(result.message || 'Something went wrong!')
        }

        await this.fetchCarts()

        // 204 No Content - no body to parse
        if (response.status === 204) {
          return null
        }

        return await response.json()
      } catch (err) {
        if (err.name === 'AbortError') {
          throw err
        }
        this._showError(err)
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    totalItems: (state) =>
      state.carts.reduce((total, item) => total + (Number(item.quantity) || 0), 0)
  }
})
