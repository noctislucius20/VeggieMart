<template>
  <main>
    <div class="mt-4">
      <div class="container">
        <div class="flex flex-wrap">
          <div class="w-full">
            <!-- breadcrumb -->
            <nav aria-label="breadcrumb">
              <ol class="flex flex-wrap">
                <li class="inline-block text-blue-600 mr-2">
                  <a href="/">
                    Home
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="icon icon-tabler icon-tabler-chevron-right inline-block"
                      width="14"
                      height="14"
                      viewBox="0 0 24 24"
                      stroke-width="2"
                      stroke="currentColor"
                      fill="none"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path d="M9 6l6 6l-6 6" />
                    </svg>
                  </a>
                </li>
                <li class="inline-block text-blue-600 mr-2">
                  <a href="shop">
                    Shop

                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="icon icon-tabler icon-tabler-chevron-right inline-block"
                      width="14"
                      height="14"
                      viewBox="0 0 24 24"
                      stroke-width="2"
                      stroke="currentColor"
                      fill="none"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path d="M9 6l6 6l-6 6" />
                    </svg>
                  </a>
                </li>
                <li class="inline-block text-gray-500 active" aria-current="page">
                  {{ selectedCategory?.name || 'All Products' }}
                </li>
              </ol>
            </nav>
          </div>
        </div>
      </div>
    </div>
    <div class="my-10">
      <div class="container">
        <div class="flex lg:gap-8">
          <aside class="lg:w-1/4 mb-6 md:">
            <div
              id="offcanvasCategory"
              class="offcanvas offcanvas-left offcanvas-collapse"
              tabindex="-1"
              aria-labelledby="offcanvasCategoryLabel"
            >
              <div class="lg:invisible lg:hidden flex items-center p-4 justify-between">
                <h5 id="offcanvasCategoryLabel" class="offcanvas-title">Filter</h5>
                <button
                  type="button"
                  class="btn-close"
                  data-bs-dismiss="offcanvas"
                  aria-label="Close"
                  @click="closeOffcanvas"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon icon-tabler icon-tabler-x text-gray-700"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    stroke-width="2"
                    stroke="currentColor"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M18 6l-12 12" />
                    <path d="M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="offcanvas-body flex flex-col gap-8">
                <div class="flex flex-col gap-3">
                  <h5>Categories</h5>
                  <ul class="flex flex-wrap nav-category">
                    <CategoryMenuItem
                      v-for="(category, index) in categories"
                      :key="index"
                      :category="category"
                      :is-active="isCategoryActive(category)"
                      :active-slug="categoryFilter"
                      @select="handleCategoryFilter"
                    />
                  </ul>
                </div>
                <div class="flex flex-col gap-3">
                  <!-- price -->
                  <h5>Price</h5>
                  <div class="flex flex-col gap-3">
                    <!-- range -->
                    <div ref="priceRangeRef" />
                    <div class="flex flex-row gap-2 items-center">
                      <span class="text-gray-800">Price:</span>
                      <span class="text-xs">{{ priceRangeValue }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </aside>
          <section class="lg:w-3/4 w-full">
            <!-- card -->
            <div
              v-if="categoryFilter && selectedCategory"
              class="relative flex flex-col min-w-0 rounded-lg wrap-break-word bg-gray-100 p-8 mb-6"
            >
              <!-- card body -->
              <div class="flex-auto">
                <h1 class="text-xl">{{ selectedCategory?.name || 'All Products' }}</h1>
              </div>
            </div>
            <!-- list icon -->
            <div
              v-if="paginateProds.total_count > 0"
              class="flex flex-col md:flex-row justify-between lg:items-center mb-6 gap-3"
            >
              <div>
                <p>
                  <span class="text-gray-900">{{ paginateProds.total_count }}</span>
                  Products found
                </p>
              </div>
              <div class="flex flex-col md:flex-row justify-between md:items-center gap-3">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-1">
                    <!-- view mode icons -->
                    <button
                      type="button"
                      :class="[
                        'inline-flex items-center justify-center p-2 rounded-lg border transition-colors',
                        viewMode === 'grid'
                          ? 'bg-blue-600 text-white border-blue-600'
                          : 'bg-white text-gray-800 border-gray-300 hover:bg-gray-100'
                      ]"
                      title="Grid View"
                      @click="viewMode = 'grid'"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="icon icon-tabler icon-tabler-layout-grid"
                        width="18"
                        height="18"
                        viewBox="0 0 24 24"
                        stroke-width="2"
                        stroke="currentColor"
                        fill="none"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      >
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path
                          d="M4 4m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z"
                        />
                        <path
                          d="M14 4m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z"
                        />
                        <path
                          d="M4 14m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z"
                        />
                        <path
                          d="M14 14m0 1a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1z"
                        />
                      </svg>
                    </button>
                    <button
                      type="button"
                      :class="[
                        'inline-flex items-center justify-center p-2 rounded-lg border transition-colors',
                        viewMode === 'list'
                          ? 'bg-blue-600 text-white border-blue-600'
                          : 'bg-white text-gray-800 border-gray-300 hover:bg-gray-100'
                      ]"
                      title="List View"
                      @click="viewMode = 'list'"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="icon icon-tabler icon-tabler-list"
                        width="18"
                        height="18"
                        viewBox="0 0 24 24"
                        stroke-width="2"
                        stroke="currentColor"
                        fill="none"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      >
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path d="M9 6l11 0" />
                        <path d="M9 12l11 0" />
                        <path d="M9 18l11 0" />
                        <path d="M5 6l0 .01" />
                        <path d="M5 12l0 .01" />
                        <path d="M5 18l0 .01" />
                      </svg>
                    </button>
                  </div>
                  <div class="ml-3 lg:hidden">
                    <a
                      class="btn inline-flex items-center gap-x-2 bg-white text-gray-800 border-gray-300 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
                      data-bs-toggle="offcanvas"
                      href="#offcanvasCategory"
                      role="button"
                      aria-controls="offcanvasCategory"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="icon icon-tabler icon-tabler-filter inline-block"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        stroke-width="2"
                        stroke="currentColor"
                        fill="none"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      >
                        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                        <path
                          d="M4 4h16v2.172a2 2 0 0 1 -.586 1.414l-4.414 4.414v7l-6 2v-8.5l-4.48 -4.928a2 2 0 0 1 -.52 -1.345v-2.227z"
                        />
                      </svg>
                      Filters
                    </a>
                  </div>
                </div>

                <div class="flex gap-3">
                  <div class="grow">
                    <!-- select option -->
                    <select
                      v-model.number="limit"
                      class="text-base py-2 block w-full border-gray-300 rounded-lg focus:border-blue-600 focus:ring-blue-600 disabled:opacity-50 disabled:pointer-events-none"
                      @change="handleLimitChange"
                    >
                      <option :value="10">Show: 10</option>
                      <option :value="20">20</option>
                      <option :value="30">30</option>
                      <option :value="50">50</option>
                    </select>
                  </div>
                  <div>
                    <!-- select option -->
                    <select
                      v-model="sort"
                      class="text-base py-2 block w-full border-gray-300 rounded-lg focus:border-blue-600 focus:ring-blue-600 disabled:opacity-50 disabled:pointer-events-none"
                      @change="handleSortChange"
                    >
                      <option value="">Sort by: Featured</option>
                      <option value="price_asc">Price: Low to High</option>
                      <option value="price_desc">Price: High to Low</option>
                      <option value="newest">Release Date</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
            <!-- Skeleton Loading -->
            <div v-if="loading" class="grid lg:grid-cols-4 md:grid-cols-3 gap-4">
              <div v-for="i in 8" :key="i" class="border border-gray-300 bg-white rounded-lg p-4">
                <div class="skeleton skeleton-image w-full aspect-square rounded mb-4" />
                <div class="skeleton skeleton-text h-3 w-1/2 mb-2" />
                <div class="skeleton skeleton-title h-5 w-3/4 mb-3" />
                <div class="flex justify-between items-center">
                  <div class="skeleton skeleton-text h-5 w-1/3" />
                  <div class="skeleton skeleton-button h-8 w-16 rounded" />
                </div>
              </div>
            </div>

            <!-- Empty State -->
            <div
              v-else-if="featuredProducts.length === 0"
              class="flex flex-col items-center gap-3 text-center py-16"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="icon icon-tabler icon-tabler-package text-gray-400"
                width="48"
                height="48"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M12 3l8 4.5v9L12 21l-8-4.5v-9L12 3z" />
                <path d="M12 12l8 -4.5" />
                <path d="M12 12v9" />
                <path d="M12 12L4 7.5" />
              </svg>
              <p class="text-gray-500 text-lg">Tidak ada produk ditemukan</p>
              <p class="text-gray-400 text-sm">
                Coba pilih kategori lain atau perbarui filter pencarian
              </p>
            </div>

            <!-- Grid View -->
            <div v-else-if="viewMode === 'grid'" class="grid lg:grid-cols-4 md:grid-cols-3 gap-4">
              <div
                v-for="product in featuredProducts"
                :key="product.id"
                class="relative rounded-lg wrap-break-word border bg-white border-gray-300 card-product"
              >
                <div class="flex-auto p-4">
                  <div class="text-center relative flex justify-center">
                    <div class="absolute top-0 left-0">
                      <span
                        class="inline-block p-1 text-center font-semibold text-sm align-baseline leading-none rounded bg-red-600 text-white"
                        >Sale</span
                      >
                    </div>
                    <NuxtLink :to="`/shop/${product.id}`">
                      <img
                        :src="product.product_image"
                        alt="Grocery Ecommerce Template"
                        class="w-full h-auto"
                      />
                    </NuxtLink>
                  </div>
                  <div class="flex flex-col gap-3">
                    <a href="#!" class="text-decoration-none text-gray-500">
                      <small>{{ product.category_name }}</small></a
                    >
                    <div class="flex flex-col gap-2">
                      <h3 class="text-base truncate">
                        <NuxtLink :to="`/shop/${product.id}`">
                          {{ product.product_name }}
                        </NuxtLink>
                      </h3>
                    </div>
                    <div class="flex justify-between items-center">
                      <div>
                        <span class="text-gray-900 font-semibold"
                          >Rp. {{ formatPrice(product.sale_price) }}</span
                        >
                        <span
                          v-if="product.sale_price != product.regular_price"
                          class="line-through text-gray-500"
                          >Rp. {{ formatPrice(product.regular_price) }}</span
                        >
                      </div>
                      <div>
                        <button
                          type="button"
                          class="btn inline-flex items-center gap-x-1 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300 btn-sm"
                          @click="addToCart(product.id)"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="icon icon-tabler icon-tabler-plus"
                            width="14"
                            height="14"
                            viewBox="0 0 24 24"
                            stroke-width="3"
                            stroke="currentColor"
                            fill="none"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                          >
                            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                            <path d="M12 5l0 14" />
                            <path d="M5 12l14 0" />
                          </svg>
                          <span>Add</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- List View -->
            <div v-else class="flex flex-col gap-4">
              <div
                v-for="product in featuredProducts"
                :key="product.id"
                class="relative rounded-lg wrap-break-word border bg-white border-gray-300 card-product"
              >
                <div class="flex flex-col md:flex-row">
                  <div class="md:w-1/4 relative flex justify-center p-4">
                    <div class="absolute top-0 left-0">
                      <span
                        class="inline-block p-1 text-center font-semibold text-sm align-baseline leading-none rounded bg-red-600 text-white"
                        >Sale</span
                      >
                    </div>
                    <NuxtLink :to="`/shop/${product.id}`">
                      <img
                        :src="product.product_image"
                        alt="Grocery Ecommerce Template"
                        class="w-full h-auto"
                      />
                    </NuxtLink>
                  </div>
                  <div class="md:w-3/4 flex-auto p-4">
                    <div class="flex flex-col gap-3">
                      <a href="#!" class="text-decoration-none text-gray-500">
                        <small>{{ product.category_name }}</small></a
                      >
                      <div class="flex flex-col gap-2">
                        <h3 class="text-lg">
                          <NuxtLink :to="`/shop/${product.id}`">
                            {{ product.product_name }}
                          </NuxtLink>
                        </h3>
                      </div>
                      <div class="flex justify-between items-center">
                        <div>
                          <span class="text-gray-900 font-semibold"
                            >Rp. {{ formatPrice(product.sale_price) }}</span
                          >
                          <span
                            v-if="product.sale_price != product.regular_price"
                            class="line-through text-gray-500"
                            >Rp. {{ formatPrice(product.regular_price) }}</span
                          >
                        </div>
                        <div>
                          <button
                            type="button"
                            class="btn inline-flex items-center gap-x-1 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300 btn-sm"
                            @click="addToCart(product.id)"
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              class="icon icon-tabler icon-tabler-plus"
                              width="14"
                              height="14"
                              viewBox="0 0 24 24"
                              stroke-width="3"
                              stroke="currentColor"
                              fill="none"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                            >
                              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                              <path d="M12 5l0 14" />
                              <path d="M5 12l14 0" />
                            </svg>
                            <span>Add</span>
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!loading && featuredProducts.length > 0" class="flex my-8">
              <nav class="flex items-center gap-x-1">
                <button
                  type="button"
                  :disabled="paginateProds.page === 1"
                  class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
                  @click="handlePageChange(paginateProds.page - 1)"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon icon-tabler icon-tabler-chevron-left"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    stroke-width="2"
                    stroke="currentColor"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M15 6l-6 6l6 6" />
                  </svg>
                </button>
                <div class="flex items-center gap-x-1">
                  <button
                    v-for="page in paginateProds.total_page"
                    :key="page"
                    :class="[
                      'leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border',
                      page === paginateProds.page
                        ? 'text-white border bg-blue-600 border-blue-600 hover:bg-blue-600 focus:outline-none focus:bg-blue-600'
                        : 'bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300'
                    ]"
                    @click="handlePageChange(page)"
                  >
                    {{ page }}
                  </button>
                </div>
                <button
                  type="button"
                  :disabled="paginateProds.page === paginateProds.total_page"
                  class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
                  @click="handlePageChange(paginateProds.page + 1)"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon icon-tabler icon-tabler-chevron-right"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    stroke-width="2"
                    stroke="currentColor"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M9 6l6 6l-6 6" />
                  </svg>
                </button>
              </nav>
            </div>
          </section>
        </div>
      </div>
    </div>
    <LoginModal v-if="showLoginModal" @close="showLoginModal = false" />
  </main>
  <!-- <QuickView v-if="show" @close="showLoginModal = false" /> -->
</template>

<script setup>
import CategoryMenuItem from '~/components/home/CategoryMenuItem.vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { watchDebounced } from '@vueuse/core'
import noUiSlider from 'nouislider'
import wNumb from 'wnumb'
import { useCategoryStore } from '~/stores/category'
import { useProductStore } from '~/stores/product'
import { useAuthStore } from '~/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import LoginModal from '~/components/modals/LoginModal.vue'
import { useCartStore } from '~/stores/cart'

const categories = ref([])
const categoryStore = useCategoryStore()
const featuredProducts = ref([])
const productStore = useProductStore()
const authStore = useAuthStore()
const cartStore = useCartStore()
const showLoginModal = ref(false)
const loading = ref(true)
const { showError } = useErrorModal()

const priceRangeRef = ref(null)
const priceRangeValue = ref('')
const route = useRoute()
const paginateProds = ref({})
const limit = ref(10)
const sort = ref('')
const viewMode = ref('grid')
const router = useRouter()
const selectedCategory = ref(null)
const searchFilter = ref(String(route.query.search ?? ''))
const priceFilter = ref('')
const isReady = ref(false)
let priceSlider = null

const categoryFilter = computed(() => selectedCategory.value?.slug || '')

const findCategoryBySlug = (categoryList, slug) => {
  for (const category of categoryList) {
    if (category.slug === slug) {
      return category
    }

    const childCategory = findCategoryBySlug(category.childs || [], slug)
    if (childCategory) {
      return childCategory
    }
  }

  return null
}

const syncSelectedCategoryFromRoute = () => {
  const categorySlug = String(route.query.category ?? '')
  selectedCategory.value = categorySlug ? findCategoryBySlug(categories.value, categorySlug) : null
}

const syncProductsFromStore = () => {
  featuredProducts.value = productStore.products
  paginateProds.value = productStore.pagination
}

const buildProductParams = (page = 1) => ({
  search: searchFilter.value,
  category: selectedCategory.value?.id || '',
  page,
  limit: limit.value,
  orderBy: sort.value,
  price: priceFilter.value
})

const fetchProducts = async ({ page = 1 } = {}) => {
  loading.value = true

  try {
    await productStore.fetchProductsShop(buildProductParams(page))
    syncProductsFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    await categoryStore.fetchCategoriesHome()
    categories.value = categoryStore.categories
  } catch (error) {
    showError(error)
  }
}

watch(
  () => route.query.search,
  (newSearch) => {
    searchFilter.value = String(newSearch ?? '')
  }
)

watch(
  () => route.query.category,
  async () => {
    if (!isReady.value) {
      return
    }

    syncSelectedCategoryFromRoute()
    await fetchProducts({ page: 1 })
  }
)

watchDebounced(
  searchFilter,
  async () => {
    if (!isReady.value) {
      return
    }

    await fetchProducts({ page: 1 })
  },
  { debounce: 500, maxWait: 1000 }
)

const initPriceRange = () => {
  if (!priceRangeRef.value || priceSlider) {
    return
  }

  priceSlider = noUiSlider.create(priceRangeRef.value, {
    start: [0, 100000],
    connect: true,
    range: {
      min: 0,
      max: 2000000
    },
    format: wNumb({
      decimals: 0,
      thousand: '.',
      prefix: 'Rp. '
    })
  })

  priceSlider.on('update', (values) => {
    priceRangeValue.value = values.join(' - ')
  })

  priceSlider.on('change', async (values) => {
    const min = parseInt(values[0].replace('Rp. ', '').replace(/\./g, ''))
    const max = parseInt(values[1].replace('Rp. ', '').replace(/\./g, ''))

    priceFilter.value = `${min} - ${max}`
    await fetchProducts({ page: 1 })
  })
}

onMounted(async () => {
  try {
    initPriceRange()
    await fetchCategories()
    syncSelectedCategoryFromRoute()
    await fetchProducts()
    isReady.value = true
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  priceSlider?.destroy()
  priceSlider = null
  categoryStore.cancelRequests()
  productStore.cancelRequests()
  authStore.cancelRequests()
  cartStore.cancelRequests()
})

const handlePageChange = async (page) => {
  if (!page || page === paginateProds.value?.page) {
    return
  }

  await fetchProducts({ page })
}

const handleLimitChange = async () => {
  await fetchProducts({ page: 1 })
}

const handleSortChange = async () => {
  await fetchProducts({ page: paginateProds.value?.page || 1 })
}

const handleCategoryFilter = async (categorySlug) => {
  selectedCategory.value = findCategoryBySlug(categories.value, categorySlug)

  const query = { ...route.query }
  if (selectedCategory.value?.slug) {
    query.category = selectedCategory.value.slug
  } else {
    delete query.category
  }

  closeOffcanvas()
  await router.push({ path: '/shop', query })
}

const isCategoryActive = (category) => {
  return (
    categoryFilter.value === category.slug ||
    (category.childs || []).some((child) => child.slug === categoryFilter.value)
  )
}

const addToCart = async (productId) => {
  if (!authStore.token) {
    // Jika belum login, arahkan ke halaman login
    alert('Silakan login terlebih dahulu')
    router.push('/auth/signin')
    return
  }
  try {
    const existing = cartStore.carts.find((c) => Number(c.id) === Number(productId))
    if (existing) {
      await cartStore.addToCart(productId, existing.quantity + 1)
    } else {
      await cartStore.addToCart(productId, 1)
    }
    alert('Produk berhasil ditambahkan ke keranjang')
  } catch (error) {
    showError(error)
  }
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const closeOffcanvas = () => {
  const { $bootstrap } = useNuxtApp()
  const offcanvasEl = document.getElementById('offcanvasCategory')

  if (offcanvasEl) {
    const offcanvas = $bootstrap.Offcanvas.getOrCreateInstance(offcanvasEl)
    offcanvas.hide()
  }
}
</script>

<style scoped>
.card-product:hover .card-product-action {
  opacity: 1 !important;
  visibility: visible !important;
  transition: all 0.3s ease;
}

.card-product .card-product-action {
  transition: all 0.3s ease;
}
</style>
