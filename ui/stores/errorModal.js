import { defineStore } from 'pinia'
import { getErrorMessage } from '~/utils/errorMessages'

export const useErrorModalStore = defineStore('errorModal', {
  state: () => ({
    isVisible: false,
    message: '',
    progress: 100,
    timer: null,
    interval: null
  }),

  actions: {
    showError(rawMessage) {
      // Terjemahkan pesan error ke UI Message
      this.message = getErrorMessage(rawMessage)
      this.isVisible = true
      this.progress = 100

      // Bersihkan timer lama jika ada
      this.clearTimers()

      const DURATION = 3000 // 3 detik

      // Progress bar: 100% → 0% dalam 3 detik
      const step = 100 / (DURATION / 30) // update setiap 30ms
      this.interval = setInterval(() => {
        this.progress = Math.max(0, this.progress - step)
      }, 30)

      // Auto-close setelah 3 detik
      this.timer = setTimeout(() => {
        this.closeError()
      }, DURATION)
    },

    closeError() {
      this.isVisible = false
      this.message = ''
      this.progress = 100
      this.clearTimers()
    },

    clearTimers() {
      if (this.interval) {
        clearInterval(this.interval)
        this.interval = null
      }
      if (this.timer) {
        clearTimeout(this.timer)
        this.timer = null
      }
    }
  }
})
