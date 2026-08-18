<template>
  <div>
    <Header />
    <slot />
    <Footer />
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, watch } from 'vue'
import Footer from '~/components/common/Footer.vue'
import Header from '~/components/common/Header.vue'
import { useAuthStore } from '~/stores/auth'
import { useNotificationStore } from '~/stores/notification'
import { useWebSocket } from '~/composables/useWebSocket'

const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { connect, disconnect } = useWebSocket()

authStore.checkAuth()

const initWebSocket = () => {
  if (authStore.isAuthenticated && authStore.user && authStore.userRole) {
    connect(
      async (payload) => {
        const notification = JSON.parse(payload)

        notificationStore.addNotification(notification)
        await notificationStore.markAsSent(notification.id)
      },
      async () => {
        await notificationStore.fetchPushNotifications()
      }
    )
  } else {
    disconnect()
  }
}

onMounted(() => {
  initWebSocket()
})

watch(
  () => authStore.isAuthenticated,
  () => {
    initWebSocket()
  }
)

onUnmounted(() => {
  disconnect()
})
</script>
