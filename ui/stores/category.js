import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useCategoryStore = defineStore('categories', {
  state: () => ({
    categories: [],
    category: {},
    loading: false,
    error: null,
    imageUrl: null,
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

    async fetchCategoriesHome() {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories/home`,
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

        this.categories = result.data

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

    async fetchCategoriesShop() {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories/shop`,
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

        this.categories = result.data

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

    async fetchCategoriesAdmin(params = {}) {
      try {
        this.loading = true

        const { search, page, limit, orderBy, status } = params

        const query = new URLSearchParams()
        if (search) query.append('search', search)
        if (page) query.append('page', page)
        if (limit) query.append('limit', limit)
        if (orderBy) query.append('orderBy', orderBy)
        if (status) query.append('status', status)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories?${query.toString()}`,
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
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.categories = result.data
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

    async fetchCategoryByIDAdmin(id) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories/${id}`,
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
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

        this.category = result.data

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

    async uploadImage(file) {
      try {
        this.loading = true
        this.error = null

        if (!file) {
          throw new Error('File tidak ditemukan')
        }

        const formData = new FormData()
        formData.append('image', file)

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/image-upload`,
          {
            method: 'POST',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`
            },
            body: formData,
            signal: this._createSignal()
          }
        )

        const result = await response.json()
        if (!response.ok) {
          throw new Error(result.message || 'Terjadi kesalahan server')
        }
        this.imageUrl = result.data.imageUrl

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

    async createCategory(formData) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories`,
          {
            method: 'POST',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData),
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

    async editCategory(formData, id) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories/${id}`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData),
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

    async deleteCategory(id) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/products/categories/${id}`,
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
          throw new Error(result.message || 'Terjadi kesalahan server')
        }

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
  }
})
