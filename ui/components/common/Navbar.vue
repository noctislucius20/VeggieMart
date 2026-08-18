<template>
  <nav
    class="navbar relative navbar-expand-lg lg:flex lg:flex-wrap items-center content-between text-black navbar-default"
    aria-label="Navigasi utama"
  >
    <div class="container max-w-7xl mx-auto w-full xl:px-4 lg:px-0">
      <!-- Desktop Navigation -->
      <div class="hidden lg:flex lg:items-center">
        <!-- Category Dropdown -->
        <div class="dropdown relative">
          <button
            class="mr-4 btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 hover:bg-blue-700 hover:border-blue-700"
            @click="toggleCategoryMenu"
          >
            <Icon name="tabler:layout-grid" size="16" />
            Semua Kategori
          </button>
          <ul
            v-if="showCategoryMenu"
            class="dropdown-menu absolute z-10 mt-2 w-48 bg-white rounded-lg shadow-lg py-2"
          >
            <li v-for="category in categories" :key="category.id">
              <NuxtLink
                :to="`/shop?category=${category.slug}`"
                class="block px-4 py-2 text-gray-800 hover:bg-blue-50"
              >
                {{ category.name }}
              </NuxtLink>
            </li>
          </ul>
        </div>

        <!-- Main Menu -->
        <ul class="navbar-nav lg:flex gap-3 lg:items-center">
          <li v-for="menu in mainMenus" :key="menu.path" class="nav-item">
            <NuxtLink :to="menu.path" class="nav-link px-4 py-2 text-gray-700 hover:text-blue-600">
              {{ menu.name }}
            </NuxtLink>
          </li>
        </ul>
      </div>

      <!-- Mobile Navigation -->
      <div
        id="navbar-default"
        class="lg:hidden offcanvas offcanvas-left"
        :class="{ show: isMobileMenuOpen }"
        tabindex="-1"
      >
        <div class="offcanvas-header pb-1">
          <NuxtLink to="/">
            <img src="/images/logo/dailymart-logo.svg" alt="Logo Sayur" />
          </NuxtLink>
          <button type="button" class="btn-close" aria-label="Tutup" @click="closeMobileMenu">
            <Icon name="tabler:x" class="text-gray-700" size="24" />
          </button>
        </div>

        <div class="offcanvas-body lg:flex lg:items-center">
          <!-- Mobile Search -->
          <div class="block lg:hidden mb-4">
            <form @submit.prevent="handleSearch">
              <div class="relative">
                <label for="searchMobile" class="invisible hidden">Cari</label>
                <input
                  id="searchMobile"
                  v-model="searchQuery"
                  type="search"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2"
                  placeholder="Cari produk"
                />
                <button type="submit" class="absolute right-0 top-0 p-3">
                  <Icon name="tabler:search" size="16" />
                </button>
              </div>
            </form>
          </div>

          <!-- Mobile Category Button -->
          <div class="block lg:hidden mb-4">
            <button
              class="btn w-full inline-flex items-center gap-x-2 bg-blue-600 text-white"
              @click="toggleMobileCategoryMenu"
            >
              <Icon name="tabler:layout-grid" size="16" />
              Semua Kategori
            </button>
            <div v-if="showMobileCategoryMenu" class="mt-2">
              <ul class="list-none">
                <li v-for="category in categories" :key="category.id">
                  <NuxtLink
                    :to="`/shop?category=${category.slug}`"
                    class="block py-2 px-3 text-gray-800 hover:bg-gray-100"
                  >
                    {{ category.name }}
                  </NuxtLink>
                </li>
              </ul>
            </div>
          </div>

          <!-- Mobile Menu -->
          <ul class="navbar-nav">
            <li v-for="menu in mainMenus" :key="menu.path" class="nav-item w-full">
              <NuxtLink
                :to="menu.path"
                class="nav-link block py-2 text-gray-800 hover:text-blue-600"
                @click="closeMobileMenu"
              >
                {{ menu.name }}
              </NuxtLink>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
const searchQuery = ref('')
const showCategoryMenu = ref(false)
const showMobileCategoryMenu = ref(false)
const isMobileMenuOpen = ref(false)

const categories = ref([
  { id: 1, name: 'Sayuran & Telur', slug: 'sayuran-telur' },
  { id: 2, name: 'Cemilan', slug: 'cemilan' },
  { id: 3, name: 'Buah-buahan', slug: 'buah-buahan' },
  { id: 4, name: 'Minuman', slug: 'minuman' },
  { id: 5, name: 'Makanan Instan', slug: 'makanan-instan' },
  { id: 6, name: 'Roti & Kue', slug: 'roti-kue' },
  { id: 7, name: 'Daging & Ikan', slug: 'daging-ikan' }
])

const mainMenus = ref([
  { name: 'Beranda', path: '/' },
  { name: 'Belanja', path: '/shop' },
  { name: 'Akun Saya', path: '/account' }
])

const toggleCategoryMenu = () => {
  showCategoryMenu.value = !showCategoryMenu.value
}

const toggleMobileCategoryMenu = () => {
  showMobileCategoryMenu.value = !showMobileCategoryMenu.value
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
  showMobileCategoryMenu.value = false
}

const handleSearch = () => {
  // Implementasi logika pencarian
  console.log('Mencari:', searchQuery.value)
}
</script>
