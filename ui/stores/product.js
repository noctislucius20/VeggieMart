import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useProductStore = defineStore('products', {
  state: () => ({
    products: [],
    product: {},
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

    async fetchProductsHome() {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/home`,
          {
            headers: {
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.products = result.data

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

    async fetchProductsShop(params = {}) {
      try {
        this.loading = true

        const { search, page, limit, price, orderBy, category } = params

        const query = new URLSearchParams()
        if (search) query.append('search', search)
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)
        if (price) query.append('price', price)
        if (orderBy) query.append('order_by', orderBy)
        if (category) query.append('category', category)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/shop?${query.toString()}`,
          {
            headers: {
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.products = result.data
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

    async fetchProductDetailHome(productId) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/home/${productId}`,
          {
            headers: {
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.product = result.data

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

    async fetchProductsAdmin(params = {}) {
      try {
        this.loading = true

        const { search, page, limit, status, orderBy, category } = params

        const query = new URLSearchParams()
        if (search) query.append('search', search)
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)
        if (status) query.append('status', status)
        if (orderBy) query.append('order_by', orderBy)
        if (category) query.append('category', category)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products?${query.toString()}`,
          {
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${useCookie('token').value}`
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.products = result.data
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

    async createProductAdmin(formData) {
      try {
        this.loading = true

        const response = await fetch(`${useRuntimeConfig().public.apiGatewayBaseUrl}/products`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${useCookie('token').value}`
          },
          body: JSON.stringify(formData),
          signal: this._createSignal()
        })

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
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

    async fetchProductDetail(id) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/${id}`,
          {
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`
            },
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.product = result.data

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

    async updateProduct(id, data) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/${id}`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(data),
            signal: this._createSignal()
          }
        )

        const result = await response.json()

        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
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
    },

    async deleteProductAdmin(productId) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/${productId}`,
          {
            method: 'DELETE',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${useCookie('token').value}`
            },
            signal: this._createSignal()
          }
        )

        if (!response.ok) {
          const result = await response.json().catch(() => ({}))
          throw new Error(result.message || 'Terjadi kesalahan server')
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
