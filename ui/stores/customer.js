import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useCustomerStore = defineStore('customers', {
  state: () => ({
    customers: [],
    loading: false,
    error: null,
    customer: {},
    abortControllers: [],
    pagination: {
      page: 1,
      total_count: 0,
      per_page: 10,
      total_pages: 0
    },
    imageUrl: null
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

    async fetchCustomers(params = {}) {
      try {
        this.loading = true
        const { search, page, limit, orderBy } = params

        const query = new URLSearchParams()
        if (search) query.append('search', search)
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)
        if (orderBy) query.append('orderBy', orderBy)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/customers?${query.toString()}`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${useCookie('token').value}`
            },
            signal: this._createSignal()
          }
        )
        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.customers = result.data
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

    async fetchCustomerByID(id) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/customers/${id}`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${useCookie('token').value}`
            },
            signal: this._createSignal()
          }
        )
        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Something went wrong!')
        }

        this.customer = result.data

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

    async uploadImage(file) {
      try {
        this.loading = true
        this.error = null

        if (!file) {
          throw new Error('File tidak ditemukan')
        }

        const formData = new FormData()
        formData.append('photo', file)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/profile/image-upload`,
          {
            method: 'POST',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`
            },
            body: formData,
            signal: this._createSignal()
          }
        )

        if (!response.ok) {
          const error = await response.json()
          throw new Error(error.message || 'Gagal mengupload gambar')
        }

        const result = await response.json()
        this.imageUrl = result.data.image_url

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

    async createCustomer(customerData) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/customers`,
          {
            method: 'POST',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(customerData),
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
      } finally {
        this.loading = false
      }
    },

    async updateCustomer(customerData, id) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/customers/${id}`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(customerData),
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
      } finally {
        this.loading = false
      }
    },

    async deleteCustomer(id) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/customers/${id}`,
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

        // 204 No Content - no body to parse
        if (response.status === 204) {
          return null
        }

        return await response.json()
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
