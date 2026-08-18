// composables/useErrorModal.js
import { computed } from 'vue'
import { useErrorModalStore } from '~/stores/errorModal'

export const useErrorModal = () => {
  const store = useErrorModalStore()

  return {
    isVisible: computed(() => store.isVisible),
    message: computed(() => store.message),
    progress: computed(() => store.progress),
    showError: (message) => store.showError(message),
    closeError: () => store.closeError()
  }
}
