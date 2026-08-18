<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="store.isVisible" class="modal-overlay" @click.self="store.closeError()">
        <div class="modal-content">
          <div class="flex items-start gap-4">
            <div class="flex-shrink-0">
              <div class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center">
                <Icon name="tabler:alert-circle" class="text-red-600" size="28" />
              </div>
            </div>
            <div class="flex-1">
              <h3 class="font-bold text-gray-800 mb-1">Terjadi Kesalahan</h3>
              <p class="text-gray-600 text-sm">{{ store.message }}</p>
            </div>
            <button class="btn-close" @click="store.closeError()">
              <Icon name="tabler:x" class="text-gray-700" size="20" />
            </button>
          </div>
          <!-- Progress Bar -->
          <div class="mt-4 h-1.5 bg-gray-200 rounded-full overflow-hidden">
            <div
              class="h-full bg-red-500 transition-all duration-100 ease-linear rounded-full"
              :style="{ width: store.progress + '%' }"
            />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { useErrorModalStore } from '~/stores/errorModal'

const store = useErrorModalStore()
</script>

<style scoped>
/* .modal-container {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  pointer-events: none;
}


.modal-content {
  max-width: 420px;
  width: 90vw;
  background: white;
  border-radius: 12px;
  padding: 20px 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
  border: 1px solid #fecaca;
  pointer-events: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-10px);
} */

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: flex-start;
  z-index: 9999;
}

.modal-content {
  max-width: 420px;
  width: 90%;
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
