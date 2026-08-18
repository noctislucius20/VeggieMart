import { defineStore } from 'pinia'
import { useErrorModal } from '~/composables/useErrorModal'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    loading: false,
    error: null,
    user: null,
    abortControllers: []
  }),

  actions: {
    isSuperAdmin() {
      return this.user?.role === 'Super Admin'
    },

    cancelRequests() {
      this.abortControllers.forEach((controller) => controller.abort())
      this.abortControllers = []
    },

    _createSignal() {
      const controller = new AbortController()
      this.abortControllers.push(controller)
      // Cleanup: remove from array when request completes
      const signal = controller.signal
      signal.addEventListener('abort', () => {
        const index = this.abortControllers.indexOf(controller)
        if (index !== -1) {
          this.abortControllers.splice(index, 1)
        }
      })
      return signal
    },

    logout() {
      this.cancelRequests()
      this.user = null
      this.token = null
      const tokenCookie = useCookie('token')
      tokenCookie.value = null
      const userCookie = useCookie('user')
      userCookie.value = null
      // localStorage.removeItem('token')
    },

    checkAuth() {
      const tokenCookie = useCookie('token')
      const userCookie = useCookie('user')
      if (tokenCookie.value && userCookie.value) {
        this.user = userCookie.value
        this.token = tokenCookie.value
      }
    },

    async signin(email, password) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/signin`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password }),
            signal: this._createSignal()
          }
        )

        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.message || 'Something went wrong!')
        }

        this.user = {
          id: data.data.id,
          name: data.data.name,
          email: data.data.email,
          phone: data.data.phone,
          lat: data.data.lat,
          lng: data.data.lng,
          photo: data.data.photo,
          role: data.data.role
        }
        this.token = data.data.access_token

        const tokenCookie = useCookie('token')
        tokenCookie.value = this.token

        const userCookie = useCookie('user')
        userCookie.value = JSON.stringify(this.user)
        // localStorage.setItem('token', this.token)
        // localStorage.setItem('user', JSON.stringify(this.user))

        return data
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    },

    async signup(userData) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/signup`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData),
            signal: this._createSignal()
          }
        )

        const data = await response.json()

        if (!response.ok) {
          throw new Error(data.message || 'Terjadi kesalahan saat registrasi')
        }

        return data
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    },

    async forgotPassword(email) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/forgot-password`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email }),
            signal: this._createSignal()
          }
        )

        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.message || 'Something went wrong!')
        }

        return data
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    },

    async verifyAccount(token) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/activate-account?token=${token}`,
          {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json'
            },
            signal: this._createSignal()
          }
        )

        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.message || 'Something went wrong!')
        }

        return data
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updatePasswordNoAuth(password_new, password_confirmation, token) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/reset-password?token=${token}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              password_new,
              password_confirmation
            }),
            signal: this._createSignal()
          }
        )

        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.message || 'Something went wrong!')
        }

        return data
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    },

    async getProfile() {
      try {
        this.loading = true

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/profile`,
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

        this.user = result.data

        return result
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateProfile(userData) {
      try {
        this.loading = true
        this.error = null

        const response = await fetch(
          `${useRuntimeConfig().public.apiGatewayBaseUrl}/users/profile`,
          {
            method: 'PUT',
            headers: {
              Authorization: `Bearer ${useCookie('token').value}`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData),
            signal: this._createSignal()
          }
        )
        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.message || 'Something went wrong!')
        }

        this.user = {
          id: data.data.id,
          name: data.data.name,
          email: data.data.email,
          phone: data.data.phone,
          photo: data.data.photo,
          role: data.data.role,
          lat: data.data.lat,
          lng: data.data.lng
        }
        this.token = data.data.access_token

        const tokenCookie = useCookie('token')
        tokenCookie.value = this.token

        const userCookie = useCookie('user')
        userCookie.value = JSON.stringify(this.user)

        return data
      } catch (error) {
        if (error.name === 'AbortError') {
          throw error
        }
        this.error = error.message
        const { showError } = useErrorModal()
        showError(error.message)
        throw error
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    isAuthenticated: (state) => !!state.token,
    getUser: (state) => state.user,
    userRole: (state) => state.user?.role
  }
})
