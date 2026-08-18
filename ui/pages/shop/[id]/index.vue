<template>
  <main>
    <div class="container">
      <div class="flex flex-wrap mt-4">
        <div class="w-full">
          <!-- breadcrumb -->
          <nav aria-label="breadcrumb">
            <ol class="flex flex-wrap">
              <li class="inline-block text-blue-600 mr-2">
                <a href="/shop">
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
              <li class="inline-block text-blue-600 mr-2">
                <a :href="`/shop?category=${encodeURIComponent(featuredProduct.category_slug)}`">
                  {{ featuredProduct.category_name }}
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
                {{ featuredProduct.product_name }}
              </li>
            </ol>
          </nav>
        </div>
      </div>
    </div>

    <section class="my-10">
      <div class="container">
        <div class="flex flex-wrap">
          <!-- Loading Skeleton -->
          <div v-if="loading" class="w-full">
            <div class="flex flex-wrap">
              <div class="lg:w-1/2 pr-4 pl-4">
                <div class="skeleton skeleton-image" style="height: 400px" />
              </div>
              <div class="lg:w-1/2 pr-4 pl-4">
                <div class="flex flex-col gap-4">
                  <div class="skeleton skeleton-text w-32" />
                  <div class="skeleton skeleton-title w-full" />
                  <div class="skeleton skeleton-title w-3/4" />
                  <div class="skeleton skeleton-price" style="width: 150px" />
                  <div class="skeleton skeleton-text w-full" />
                  <div class="skeleton skeleton-text w-full" />
                  <div class="skeleton skeleton-button" style="width: 200px" />
                </div>
              </div>
            </div>
          </div>
          <!-- Product Detail -->
          <div v-else class="w-full">
            <div class="flex flex-wrap">
              <ProductSingleGallery :product="featuredProduct" />
              <div class="lg:w-1/2 pr-4 pl-4">
                <div class="lg:pl-10 mt-6 md:mt-0">
                  <div class="flex flex-col gap-4">
                    <!-- content -->
                    <a href="#!" class="block text-blue-600">{{ featuredProduct.category_name }}</a>
                    <!-- heading -->
                    <div class="flex flex-col">
                      <h1>{{ featuredProduct.product_name }}</h1>
                      <div class="flex flex-col gap-4">
                        <div class="text-md">
                          <span class="text-gray-900 font-semibold">
                            Rp.
                            {{
                              formatPrice(
                                selectedChild
                                  ? selectedChild.sale_price
                                  : featuredProduct.sale_price
                              )
                            }}
                          </span>
                          <span
                            v-if="
                              selectedChild
                                ? selectedChild.sale_price !== selectedChild.regular_price
                                : featuredProduct.sale_price !== featuredProduct.regular_price
                            "
                            class="line-through text-gray-500 ml-2"
                          >
                            Rp.
                            {{
                              formatPrice(
                                selectedChild
                                  ? selectedChild.regular_price
                                  : featuredProduct.regular_price
                              )
                            }}
                          </span>
                        </div>
                      </div>
                    </div>
                    <!-- hr -->
                    <div class="flex flex-col gap-6">
                      <hr />
                      <div>
                        <button
                          type="button"
                          :class="[
                            'btn inline-flex items-center gap-x-2 border-gray-300 border mr-2',
                            !selectedChild
                              ? 'bg-blue-600 text-white hover:text-white'
                              : 'bg-white text-gray-800'
                          ]"
                          @click="handleVariantSelect(null)"
                        >
                          {{ featuredProduct.weight }}{{ featuredProduct.unit }}
                        </button>
                        <button
                          v-for="childs in featuredProduct.childs"
                          :key="childs.id"
                          type="button"
                          :class="[
                            'btn inline-flex items-center gap-x-2 border-gray-300 border mr-2',
                            selectedChild?.id === childs.id
                              ? 'bg-blue-600 text-white hover:text-white'
                              : 'bg-white text-gray-800'
                          ]"
                          @click="handleVariantSelect(childs)"
                        >
                          {{ childs.weight }}{{ featuredProduct.unit }}
                        </button>
                      </div>
                      <div class="flex flex-wrap justify-start gap-2 items-center">
                        <div class="lg:w-1/3 md:w-2/5 w-full grid">
                          <button
                            type="button"
                            class="btn gap-x-1 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300 justify-center"
                            @click="
                              addToCart(
                                selectedChild === null ? featuredProduct.id : selectedChild.id
                              )
                            "
                          >
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
                              class="icon icon-tabler icons-tabler-outline icon-tabler-shopping-bag"
                            >
                              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                              <path
                                d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304z"
                              />
                              <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
                            </svg>
                            Add to cart
                          </button>
                        </div>
                        <div class="md:w-1/3 w-full" />
                      </div>
                      <!-- hr -->
                      <hr />
                    </div>
                    <div>
                      <!-- table -->
                      <table class="text-left w-full">
                        <tbody>
                          <tr>
                            <td class="px-6 py-3">Availability:</td>
                            <td class="px-6 py-3">
                              {{ selectedChild ? selectedChild.stock : featuredProduct.stock }}
                            </td>
                          </tr>
                          <tr>
                            <td class="px-6 py-3">Type:</td>
                            <td class="px-6 py-3">{{ featuredProduct.category_name }}</td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    <section class="mb-10">
      <div class="container">
        <div class="flex flex-wrap">
          <div class="w-full">
            <ul id="myTab" class="nav nav-line-bottom border-b border-gray-300 pl-0" role="tablist">
              <li class="nav-item" role="presentation">
                <button
                  id="product-tab"
                  class="inline-block py-3 font-semibold px-4 no-underline nav-link active"
                  data-bs-toggle="tab"
                  data-bs-target="#product-tab-pane"
                  type="button"
                  role="tab"
                  aria-controls="product-tab-pane"
                  aria-selected="true"
                >
                  Product Details
                </button>
              </li>
            </ul>
            <div id="myTabContent" class="tab-content">
              <div
                id="product-tab-pane"
                class="tab-pane active opacity-100 block"
                role="tabpanel"
                aria-labelledby="product-tab"
                tabindex="0"
              >
                <!-- Loading Skeleton for Description & Variant -->
                <div v-if="loading" class="my-8 flex flex-col gap-5">
                  <div class="flex flex-col gap-1">
                    <div class="skeleton skeleton-title w-64" />
                    <div class="skeleton skeleton-text w-full" />
                    <div class="skeleton skeleton-text w-full" />
                    <div class="skeleton skeleton-text w-3/4" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <div class="skeleton skeleton-title w-48" />
                    <div class="skeleton skeleton-text w-32" />
                  </div>
                  <div class="flex flex-col gap-1">
                    <div class="skeleton skeleton-title w-48" />
                    <div class="skeleton skeleton-text w-full" />
                    <div class="skeleton skeleton-text w-3/4" />
                  </div>
                </div>
                <!-- Actual Content -->
                <div v-else class="my-8 flex flex-col gap-5">
                  <div class="flex flex-col gap-1">
                    <h3 class="text-lg font-semibold">Nutrient Value & Benefits</h3>
                    <p>{{ featuredProduct.description }}</p>
                  </div>
                  <div class="flex flex-col gap-1">
                    <h4 class="text-md font-semibold">Variant</h4>
                    <p>
                      {{ featuredProduct.childs === null ? 1 : featuredProduct.childs?.length + 1 }}
                      variant
                    </p>
                  </div>
                  <div class="flex flex-col gap-1">
                    <h4 class="text-md font-semibold">Disclaimer</h4>
                    <p>
                      Image shown is a representation and may slightly vary from the actual product.
                      Every effort is made to maintain accuracy of all information displayed.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useProductStore } from '~/stores/product'
import { useCartStore } from '~/stores/cart'
import ProductSingleGallery from '~/components/home/ProductSingleGallery.vue'

const featuredProduct = ref({})
const productStore = useProductStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const productId = route.params.id
const loading = ref(true)
const { showError } = useErrorModal()

const selectedChild = ref(null)

const handleVariantSelect = (childs) => {
  if (childs === null) {
    selectedChild.value = null
    return
  }
  selectedChild.value = childs
}

onMounted(async () => {
  try {
    await productStore.fetchProductDetailHome(productId)
    featuredProduct.value = productStore.product
    selectedChild.value = null
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  productStore.cancelRequests()
  cartStore.cancelRequests()
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const addToCart = async (productId) => {
  if (!authStore.token) {
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
</script>
