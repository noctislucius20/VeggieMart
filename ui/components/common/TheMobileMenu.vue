<template>
  <div>
    <!-- Backdrop -->
    <div
      class="fixed inset-0 bg-black/50 z-40 transition-opacity duration-300"
      :class="show ? 'opacity-100' : 'opacity-0 pointer-events-none'"
      @click="$emit('close')"
    />

    <!-- Offcanvas Menu -->
    <div
      class="fixed top-0 left-0 h-full w-[85%] max-w-sm bg-white z-50 shadow-xl transition-transform duration-300 ease-in-out overflow-y-auto"
      :class="show ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200">
        <NuxtLink to="/" @click="$emit('close')">
          <img src="/images/logo/dailymart-logo.svg" alt="Logo Sayur" class="h-8" />
        </NuxtLink>
        <button
          type="button"
          class="p-2 text-gray-500 hover:text-gray-700 rounded-lg hover:bg-gray-100"
          aria-label="Tutup menu"
          @click="$emit('close')"
        >
          <Icon name="tabler:x" size="24" />
        </button>
      </div>

      <!-- Search -->
      <div class="p-4 border-b border-gray-200">
        <form @submit.prevent="handleSearch">
          <div class="relative">
            <label for="mobileSearch" class="invisible hidden">Cari produk</label>
            <input
              id="mobileSearch"
              v-model="searchQuery"
              type="search"
              class="w-full border border-gray-300 text-gray-900 rounded-lg focus:ring-blue-600 focus:border-blue-600 block p-2.5 text-sm"
              placeholder="Cari produk..."
            />
            <button type="submit" class="absolute right-0 top-0 p-3" aria-label="Cari">
              <Icon name="tabler:search" size="16" />
            </button>
          </div>
        </form>
      </div>

      <!-- Categories Dropdown -->
      <div class="px-4 pt-4 border-t border-gray-200">
        <button
          type="button"
          class="w-full flex items-center justify-between px-3 py-2.5 rounded-lg font-medium bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
          :aria-expanded="showCategories"
          @click="toggleCategories"
        >
          <span class="flex items-center gap-3">
            <Icon name="tabler:layout-grid" size="20" />
            Semua Kategori
          </span>
          <Icon
            name="tabler:chevron-down"
            size="20"
            class="transition-transform duration-300 ease-in-out"
            :class="showCategories ? 'rotate-180' : ''"
          />
        </button>
        <Transition name="categories">
          <div v-if="showCategories" class="categories-dropdown mt-2">
            <div class="categories-dropdown__inner">
              <ul class="categories-stagger space-y-1">
                <li
                  v-for="(category, index) in categories"
                  :key="category.slug"
                  :style="{ '--stagger-index': index }"
                >
                  <NuxtLink
                    :to="`/shop?category=${category.slug}`"
                    class="flex items-center gap-3 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
                    @click="$emit('close')"
                  >
                    <Icon name="tabler:point" size="16" class="text-blue-600" />
                    {{ category.name }}
                  </NuxtLink>
                </li>
              </ul>
            </div>
          </div>
        </Transition>
      </div>

      <!-- Navigation Menu -->
      <div class="p-4">
        <ul class="space-y-1">
          <li>
            <NuxtLink
              to="/"
              class="flex items-center gap-3 px-3 py-2.5 text-gray-800 hover:bg-gray-100 rounded-lg font-medium"
              @click="$emit('close')"
            >
              <Icon name="tabler:home" size="20" />
              Home
            </NuxtLink>
          </li>
          <li>
            <NuxtLink
              to="/shop"
              class="flex items-center gap-3 px-3 py-2.5 text-gray-800 hover:bg-gray-100 rounded-lg font-medium"
              @click="$emit('close')"
            >
              <Icon name="tabler:shopping-bag" size="20" />
              Shop
            </NuxtLink>
          </li>
          <li v-if="authStore.isAuthenticated">
            <NuxtLink
              to="/account"
              class="flex items-center gap-3 px-3 py-2.5 text-gray-800 hover:bg-gray-100 rounded-lg font-medium"
              @click="$emit('close')"
            >
              <Icon name="tabler:user" size="20" />
              Account
            </NuxtLink>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

defineProps({
  show: {
    type: Boolean,
    default: false
  },
  categories: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close'])

const authStore = useAuthStore()
const router = useRouter()
const searchQuery = ref('')
const showCategories = ref(false)

const toggleCategories = () => {
  showCategories.value = !showCategories.value
}

const handleSearch = () => {
  if (searchQuery.value) {
    router.push({
      path: '/shop',
      query: { search: searchQuery.value }
    })
    emit('close')
  }
}
</script>

<style scoped>
/* Smooth expand/collapse using grid-template-rows 0fr -> 1fr
   so the height animates naturally regardless of category count */
.categories-dropdown {
  display: grid;
  grid-template-rows: 1fr;
}

.categories-dropdown__inner {
  overflow: hidden;
  min-height: 0;
}

.categories-enter-active,
.categories-leave-active {
  transition:
    grid-template-rows 0.35s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.25s ease,
    transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  will-change: grid-template-rows, opacity, transform;
}

.categories-enter-from,
.categories-leave-to {
  grid-template-rows: 0fr;
  opacity: 0;
  transform: translateY(-10px);
}

.categories-enter-to,
.categories-leave-from {
  grid-template-rows: 1fr;
  opacity: 1;
  transform: translateY(0);
}

/* Cascade items in one-by-one for a polished, smooth feel */
.categories-enter-active .categories-stagger > li {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
  transition-delay: calc(var(--stagger-index) * 30ms);
}

.categories-enter-from .categories-stagger > li {
  opacity: 0;
  transform: translateY(6px);
}

.categories-leave-active .categories-stagger > li {
  transition: opacity 0.15s ease;
}
</style>
