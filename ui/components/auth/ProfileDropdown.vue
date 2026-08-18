<template>
  <li ref="dropdownRef" class="dropdown">
    <button class="flex items-center" @click.prevent="toggleDropdown">
      <!-- Loading Skeleton -->
      <div v-if="loading" class="skeleton w-12 h-12 rounded-full" />
      <!-- Profile Image -->
      <img
        v-else
        :src="profile?.photo ? profile?.photo : `/images/avatar/avatar.jpg`"
        class="h-12 w-12 rounded-full"
      />
    </button>

    <div v-if="isOpen" class="dropdown-menu dropdown-menu-end !p-0 show" style="min-width: 220px">
      <div class="leading-snug px-5 py-4 border-b border-gray-300 text-left">
        <h5 class="mb-1 text-base">{{ profile?.name }}</h5>
        <small>{{ profile?.email }}</small>
      </div>

      <ul class="list-unstyled px-2 py-3">
        <li v-for="item in menuItems" :key="item.path">
          <NuxtLink :to="item.path" class="dropdown-item">
            {{ item.label }}
          </NuxtLink>
        </li>
      </ul>
      <div class="border-t px-5 py-3">
        <button @click="handleLogout">Log Out</button>
      </div>
    </div>
  </li>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useRouter } from 'vue-router'

const isOpen = ref(false)
const authStore = useAuthStore()
const dropdownRef = ref(null)
const router = useRouter()
const profile = ref(null)
const loading = ref(true)
const { showError } = useErrorModal()

const menuItems = computed(() => {
  const items = [
    { label: 'Home', path: '/' },
    { label: 'Settings', path: '/account/setting' }
  ]

  if (profile.value?.role && profile.value.role.toLowerCase().includes('super admin')) {
    items.unshift({ label: 'Dashboard', path: '/dashboard' })
  }

  return items
})

onMounted(() => {
  try {
    authStore.checkAuth()
    profile.value = authStore.user
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

const handleLogout = () => {
  try {
    authStore.logout()

    useAuthStore().$reset()
    useCartStore().$reset()
    useCategoryStore().$reset()
    useCustomerStore().$reset()
    useNotificationStore().$reset()
    useOrderStore().$reset()
    usePaymentStore().$reset()
    useProductStore().$reset()
    useRoleStore().$reset()

    router.push('/auth/signin')
  } catch (error) {
    showError(error)
  }
}

onClickOutside(dropdownRef, () => {
  isOpen.value = false
})
</script>
