<template>
  <nav
    class="navbar-vertical-nav"
    :class="[
      open
        ? 'fixed inset-y-0 left-0 z-2000 translate-x-0'
        : 'fixed inset-y-0 left-0 -translate-x-full',
      'xl:block xl:translate-x-0 xl:relative',
      'transition-transform duration-300 ease-in-out'
    ]"
  >
    <div class="navbar-vertical">
      <div class="px-8 py-6">
        <NuxtLink to="/" class="navbar-brand">
          <img src="~/assets/images/logo/dailymart-logo.svg" alt="DailyMart Logo" />
        </NuxtLink>
      </div>
      <div class="navbar-vertical-content grow">
        <ul id="sideNavbar" class="navbar-nav flex-col">
          <li class="nav-item">
            <NuxtLink
              class="nav-link"
              :class="{ active: $route.path === '/dashboard' }"
              to="/dashboard"
              @click="handleNavClick"
            >
              <div class="flex items-center">
                <span class="nav-link-icon">
                  <IconHome />
                </span>
                <span class="nav-link-text">Dashboard</span>
              </div>
            </NuxtLink>
          </li>

          <li class="nav-item mt-6 mb-3">
            <span class="nav-label">Store Managements</span>
          </li>

          <li v-for="(item, index) in storeManagementItems" :key="index" class="nav-item">
            <NuxtLink
              class="nav-link"
              :class="{ active: $route.path === item.path }"
              :to="item.path"
              @click="handleNavClick"
            >
              <div class="flex items-center">
                <span class="nav-link-icon">
                  <component :is="item.icon" />
                </span>
                <span class="nav-link-text">{{ item.label }}</span>
              </div>
            </NuxtLink>
          </li>

          <li class="nav-item mt-6 mb-3">
            <span class="nav-label">Site Settings</span>
          </li>

          <li v-for="(item, index) in settingsItems" :key="index" class="nav-item">
            <NuxtLink
              class="nav-link"
              :class="{ active: $route.path === item.path }"
              :to="item.path"
              @click="handleNavClick"
            >
              <div class="flex items-center">
                <span class="nav-link-icon">
                  <component :is="item.icon" />
                </span>
                <span class="nav-link-text">{{ item.label }}</span>
              </div>
            </NuxtLink>
          </li>
          <li class="nav-item" style="cursor: pointer">
            <NuxtLink class="nav-link" @click="handleLogout">
              <div class="flex items-center">
                <span class="nav-link-icon">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="icon icon-tabler icons-tabler-outline icon-tabler-logout"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path
                      d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2"
                    />
                    <path d="M9 12h12l-3 -3" />
                    <path d="M18 15l3 -3" />
                  </svg>
                </span>
                <span class="nav-link-text">Log out</span>
              </div>
            </NuxtLink>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup>
import IconHome from './IconHome.vue'
import IconStore from './IconStore.vue'
import IconList from './IconList.vue'
import IconShoppingBag from './IconShoppingBag.vue'
import IconUsers from './IconUsers.vue'
import IconUserCog from './IconUserCog.vue'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

defineProps({
  open: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const authStore = useAuthStore()
const router = useRouter()
const { showError } = useErrorModal()

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

const handleNavClick = () => {
  emit('close')
}

const storeManagementItems = [
  {
    label: 'Products',
    path: '/dashboard/products',
    icon: IconStore
  },
  {
    label: 'Categories',
    path: '/dashboard/categories',
    icon: IconList
  },
  {
    label: 'Orders',
    path: '/dashboard/orders',
    icon: IconShoppingBag
  },
  {
    label: 'Customers',
    path: '/dashboard/customers',
    icon: IconUsers
  }
]

const settingsItems = [
  {
    label: 'Roles',
    path: '/dashboard/roles',
    icon: IconUserCog
  }
]
</script>

<style scoped></style>
