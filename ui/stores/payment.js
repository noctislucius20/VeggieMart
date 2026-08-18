import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const usePaymentStore = defineStore('payment', {
  state: () => ({
    token: null,
    loading: false,
    error: null,
    payment: null,
    midToken: null,
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

    async createPayment(payload) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(`${useRuntimeConfig().public.apiGatewayBaseUrl}/payments`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${useCookie('token').value}`
          },
          body: JSON.stringify(payload),
          signal: this._createSignal()
        })
        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.payment = result.data

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

    async getDetailPaymentByOrderId(orderId, isAdmin = false) {
      try {
        this.loading = true
        this.error = null

        let url = `${useRuntimeConfig().public.apiGatewayBaseUrl}/payments/order/${orderId}`
        if (isAdmin) {
          url = `${useRuntimeConfig().public.apiGatewayBaseUrl}/payments/order/${orderId}/admin`
        }

        const response = await fetch(url, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${useCookie('token').value}`
          },
          signal: this._createSignal()
        })
        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.payment = result.data

        return result
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this._showError(error)
      } finally {
        this.loading = false
      }
    }
  }
})
