<template>
  <header>
    <div class="border-b">
      <div class="pt-5">
        <div class="container">
          <div class="flex flex-wrap w-full items-center justify-between">
            <!-- Logo -->
            <div class="sm:w-full lg:w-1/6 md:w-1/5">
              <NuxtLink to="/" class="navbar-brand">
                <img
                  src="/images/logo/dailymart-logo.svg"
                  alt="TailwindCSS eCommerce HTML Template"
                />
              </NuxtLink>
            </div>

            <!-- Search Bar -->
            <div class="lg:w-2/5 hidden lg:block">
              <form @submit.prevent="handleSearch">
                <div class="relative">
                  <label for="searchProducts" class="invisible hidden">Search</label>
                  <input
                    id="searchProducts"
                    v-model="searchQuery"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    type="search"
                    placeholder="Search for products"
                  />
                  <button class="absolute right-0 top-0 p-3" type="submit">
                    <Icon name="tabler:search" size="16" />
                  </button>
                </div>
              </form>
            </div>

            <!-- User Actions -->
            <div class="lg:w-1/5 text-end md:w-1/2 w-3/5">
              <div class="flex items-center justify-end">
                <!-- User Menu -->
                <div v-if="!authStore.isAuthenticated">
                  <button class="text-gray-600" @click="showLoginModal = true">
                    <Icon name="tabler:user" size="22" />
                  </button>
                </div>

                <!-- Cart -->
                <div>
                  <button
                    type="button"
                    class="text-gray-600 h-12 w-12 rounded-full hover:bg-gray-100 hover:text-blue-600 relative"
                    @click="showCart = true"
                  >
                    <Icon name="tabler:shopping-bag" size="24" />
                    <span
                      v-if="cartCount > 0"
                      class="absolute top-0 -mt-1 left-full rounded-full h-5 w-5 -ml-3 bg-blue-600 text-white text-center font-semibold text-sm"
                    >
                      {{ cartCount }}
                      <span class="invisible">unread messages</span>
                    </span>
                  </button>
                </div>

                <!-- Notifications -->
                <div v-if="authStore.isAuthenticated">
                  <ul class="list-unstyled flex items-center gap-2">
                    <NotificationDropdown />
                    <ProfileDropdown />
                  </ul>
                </div>

                <!-- Mobile Menu -->
                <div class="lg:hidden leading-none">
                  <button
                    type="button"
                    class="p-2 text-gray-600 hover:text-gray-800 rounded-lg hover:bg-gray-100"
                    aria-label="Toggle navigation"
                    @click="showMobileMenu = true"
                  >
                    <Icon name="tabler:menu-2" class="text-gray-800" size="24" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav
        class="navbar relative navbar-expand-lg lg:flex lg:flex-wrap items-center content-between text-black navbar-default"
        aria-label="Main navigation"
      >
        <div class="container max-w-7xl mx-auto xl:px-4 lg:px-0">
          <div class="hidden lg:flex lg:items-center">
            <!-- Category Dropdown -->
            <div class="dropdown">
              <button
                class="mr-4 btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                @click="toggleCategoryMenu"
              >
                <Icon name="tabler:layout-grid" size="16" />
                All Category
              </button>
              <ul v-if="showCategoryMenu" class="dropdown-menu">
                <li v-for="category in categories" :key="category.slug">
                  <NuxtLink
                    :to="`/shop?category=${encodeURIComponent(category.slug)}`"
                    class="dropdown-item"
                    @click="showCategoryMenu = false"
                  >
                    {{ category.name }}
                  </NuxtLink>
                </li>
              </ul>
            </div>

            <!-- Main Menu -->
            <ul class="navbar-nav lg:flex gap-3 lg:items-center">
              <li class="nav-item w-full lg:w-auto">
                <NuxtLink to="/" class="nav-link">Home</NuxtLink>
              </li>
              <li class="nav-item w-full lg:w-auto">
                <NuxtLink to="/shop" class="nav-link">Shop</NuxtLink>
              </li>
              <li class="nav-item w-full lg:w-auto">
                <NuxtLink v-if="authStore.isAuthenticated" to="/account" class="nav-link"
                  >Account</NuxtLink
                >
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </div>
  </header>

  <!-- Login Modal -->
  <LoginModal v-if="showLoginModal" @close="showLoginModal = false" />

  <!-- Cart Sidebar -->
  <CartSidebar v-if="showCart" @close="showCart = false" />

  <!-- Mobile Menu -->
  <TheMobileMenu :show="showMobileMenu" :categories="categories" @close="showMobileMenu = false" />
</template>

<script setup>
import CartSidebar from '../modals/CartSidebar.vue'
import LoginModal from '../modals/LoginModal.vue'
import NotificationDropdown from '../auth/NotificationDropdown.vue'
import ProfileDropdown from '../auth/ProfileDropdown.vue'
import TheMobileMenu from './TheMobileMenu.vue'
import { useAuthStore } from '~/stores/auth'
import { useCategoryStore } from '~/stores/category'
import { useCartStore } from '~/stores/cart'
import { useRouter } from 'vue-router'

const cartStore = useCartStore()
const authStore = useAuthStore()

const searchQuery = ref('')
const showLoginModal = ref(false)
const showCart = ref(false)
const showMobileMenu = ref(false)
const showCategoryMenu = ref(false)
const categories = ref([])
const categoryStore = useCategoryStore()
const { showError } = useErrorModal()

const syncCart = async () => {
  if (!authStore.isAuthenticated) {
    cartStore.clearCarts()
    return
  }

  try {
    await cartStore.fetchCarts()
  } catch (error) {
    showError(error)
  }
}

onMounted(async () => {
  const response = await categoryStore.fetchCategoriesHome()
  categories.value = response.data
  await syncCart()
})

watch(
  () => authStore.isAuthenticated,
  () => {
    syncCart()
  }
)

watch(
  () => showCart.value,
  async (isOpen) => {
    if (isOpen) {
      await syncCart()
    }
  }
)

const router = useRouter()

const cartCount = computed(() => {
  return cartStore.totalItems
})

const handleSearch = () => {
  if (searchQuery.value) {
    router.push({
      path: '/shop',
      query: { search: searchQuery.value }
    })
  }
}

const toggleCategoryMenu = () => {
  showCategoryMenu.value = !showCategoryMenu.value
}

onBeforeUnmount(() => {
  categoryStore.cancelRequests()
  authStore.cancelRequests()
  cartStore.cancelRequests()
})
</script>
