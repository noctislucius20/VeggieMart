<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import Navbar from '~/components/admin/Navbar.vue'
import AdminSidebarNav from '~/components/admin/SidebarNav.vue'
import { useAuthStore } from '~/stores/auth'
import { useNotificationStore } from '~/stores/notification'
import { useWebSocket } from '~/composables/useWebSocket'

const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { connect, disconnect } = useWebSocket()

const sidebarOpen = ref(false)

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

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeSidebar = () => {
  sidebarOpen.value = false
}
</script>

<template>
  <div>
    <Navbar :sidebar-open="sidebarOpen" @toggle-sidebar="toggleSidebar" />
    <div class="main-wrapper">
      <!-- Overlay untuk mobile ketika sidebar terbuka -->
      <Transition name="fade">
        <div
          v-if="sidebarOpen"
          class="fixed inset-0 bg-black/50 z-40 xl:hidden"
          @click="closeSidebar"
        />
      </Transition>
      <AdminSidebarNav :open="sidebarOpen" @close="closeSidebar" />
      <main class="main-content-wrapper" @click="closeSidebar">
        <slot />
      </main>
    </div>
  </div>
</template>

<style></style>
