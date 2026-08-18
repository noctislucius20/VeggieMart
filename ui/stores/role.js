import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useRoleStore = defineStore('roles', {
  state: () => ({
    roles: [],
    role: {},
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

    async fetchRoles() {
      try {
        this.loading = true

        const response = await fetch(`${useRuntimeConfig().public.apiGatewayBaseUrl}/users/roles`, {
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

        this.roles = result.data

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

    async fetchRoleID(roleID) {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/roles/${roleID}`,
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
          throw new Error(result.message || 'Something went wrong!')
        }

        this.role = result.data

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

    async createRole(name) {
      try {
        this.loading = true

        const response = await fetch(`${useRuntimeConfig().public.apiGatewayBaseUrl}/users/roles`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${useCookie('token').value}`
          },
          body: JSON.stringify({
            name: name
          }),
          signal: this._createSignal()
        })

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

    async updateRole(id, name) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/roles/${id}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${useCookie('token').value}`
            },
            body: JSON.stringify({
              name: name
            }),
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

    async deleteRole(id) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/roles/${id}`,
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
