<template>
  <div>
    <section class="mt-8">
      <div class="container">
        <Swiper
          :modules="[SwiperAutoplay, SwiperEffectFade, SwiperPagination, SwiperNavigation]"
          :effect="'fade'"
          :speed="400"
          :space-between="100"
          :pagination="{
            clickable: true,
            el: '.swiper-pagination'
          }"
          :autoplay="{
            delay: 3000,
            disableOnInteraction: false
          }"
          :breakpoints="{
            480: { slidesPerView: 1 },
            768: { slidesPerView: 1 },
            1024: { slidesPerView: 1 }
          }"
          class="swiper-container"
        >
          <!-- Slide 1 -->
          <SwiperSlide>
            <div class="relative rounded-lg overflow-hidden" style="min-height: 400px">
              <img
                src="~/assets/images/slider/slide-1.jpg"
                alt="Slider 1"
                class="absolute inset-0 w-full h-full object-cover"
              />
              <div class="relative lg:py-32 p-12 lg:pl-12 xl:w-2/5 md:w-3/5">
                <span
                  class="inline-block p-2 text-sm align-baseline leading-none rounded-lg bg-yellow-500 text-gray-900 font-semibold"
                >
                  Opening Sale Discount 50%
                </span>
                <div class="my-7 flex flex-col gap-2">
                  <h1 class="text-gray-900 text-xl lg:text-5xl font-bold leading-tight">
                    SuperMarket For Daily Grocery
                  </h1>
                  <p class="text-md font-light">
                    Introduced a new model for online grocery shopping and convenient home delivery.
                  </p>
                </div>
                <NuxtLink
                  to="/shop"
                  class="btn inline-flex items-center gap-x-2 bg-gray-800 text-white border-gray-800 hover:bg-gray-900 hover:border-gray-900 focus:ring-4 focus:ring-gray-300"
                >
                  Shop Now
                  <Icon name="tabler:arrow-right" class="inline-block" size="14" />
                </NuxtLink>
              </div>
            </div>
          </SwiperSlide>

          <!-- Slide 2 -->
          <SwiperSlide>
            <div class="relative rounded-lg overflow-hidden" style="min-height: 400px">
              <img
                src="/images/slider/slider-2.jpg"
                alt="Slider 2"
                class="absolute inset-0 w-full h-full object-cover"
              />
              <div class="relative lg:py-32 lg:pl-12 lg:pr-6 px-12 py-12 xl:w-2/5 md:w-3/5">
                <span
                  class="inline-block p-2 text-sm align-baseline leading-none rounded-lg bg-yellow-500 text-gray-900 font-semibold"
                >
                  Free Shipping - orders over $100
                </span>
                <div class="my-7 flex flex-col gap-2">
                  <h2 class="text-gray-900 text-xl lg:text-5xl font-bold leading-tight">
                    Free Shipping on
                    <br />
                    orders over
                    <span class="text-blue-600">$100</span>
                  </h2>
                  <p class="text-md font-light">
                    Free Shipping to First-Time Customers Only, After promotions and discounts are
                    applied.
                  </p>
                </div>
                <NuxtLink
                  to="/shop"
                  class="btn inline-flex items-center gap-x-2 bg-gray-800 text-white border-gray-800 hover:bg-gray-900 hover:border-gray-900 focus:ring-4 focus:ring-gray-300"
                >
                  Shop Now
                  <Icon name="tabler:arrow-right" class="inline-block" size="14" />
                </NuxtLink>
              </div>
            </div>
          </SwiperSlide>

          <!-- Pagination -->
          <div class="swiper-pagination bottom-14!" />
        </Swiper>
      </div>
    </section>
    <!-- Categories Section -->
    <section class="mt-8">
      <div class="container">
        <div class="flex flex-wrap">
          <div class="w-full">
            <h2 class="text-lg absolute z-10">Featured Categories</h2>
          </div>
        </div>

        <Swiper
          :modules="[SwiperAutoplay, SwiperNavigation]"
          :speed="400"
          :space-between="20"
          :navigation="true"
          :autoplay="{
            delay: 3000,
            disableOnInteraction: false
          }"
          :breakpoints="{
            480: { slidesPerView: 2 },
            768: { slidesPerView: 3 },
            1024: { slidesPerView: 6 }
          }"
          class="py-12"
        >
          <!-- Loading Skeleton -->
          <div v-if="loading" class="grid gap-4 grid-cols-1 md:grid-cols-2 lg:gap-4 xl:grid-cols-5">
            <div v-for="i in 5" :key="'skeleton-cat-' + i">
              <div class="relative rounded-lg wrap-break-word border bg-white border-gray-300">
                <div class="py-8 text-center">
                  <div class="skeleton skeleton-image w-28 h-28 mx-auto mb-3 rounded-full" />
                  <div class="skeleton skeleton-text mx-auto w-20" />
                </div>
              </div>
            </div>
          </div>
          <!-- Categories -->
          <SwiperSlide v-for="category in categories" :key="category.id">
            <NuxtLink :to="`/shop?category=${encodeURIComponent(category.slug)}`">
              <div
                class="relative rounded-lg wrap-break-word border bg-white border-gray-300 transition duration-75 hover:transition hover:duration-500 ease-in-out hover:border-blue-600 hover:shadow-md"
              >
                <div class="py-8 text-center">
                  <img :src="category.icon" :alt="category.name" class="w-28 h-28 mb-3 mx-auto" />
                  <div class="text-base">{{ category.name }}</div>
                </div>
              </div>
            </NuxtLink>
          </SwiperSlide>
        </Swiper>
      </div>
    </section>

    <!-- Featured Products -->
    <section class="lg:my-14 my-8">
      <div class="container">
        <div class="flex flex-wrap">
          <div class="w-full mb-6">
            <h2 class="text-lg">New Products</h2>
          </div>
        </div>

        <!-- Loading Skeleton -->
        <div v-if="loading" class="grid gap-4 grid-cols-1 md:grid-cols-2 lg:gap-4 xl:grid-cols-5">
          <div
            v-for="i in 5"
            :key="'skeleton-prod-' + i"
            class="relative rounded-lg wrap-break-word border bg-white border-gray-300"
          >
            <div class="flex-auto p-4">
              <div class="text-center relative flex justify-center">
                <div class="skeleton skeleton-image mb-3" />
              </div>
              <div class="flex flex-col gap-3">
                <div class="skeleton skeleton-text-sm w-16 mx-auto" />
                <div class="skeleton skeleton-title w-full" />
                <div class="flex justify-between items-center">
                  <div class="skeleton skeleton-price" />
                  <div class="skeleton skeleton-button" />
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- Products -->
        <div v-else class="grid gap-4 grid-cols-1 md:grid-cols-2 lg:gap-4 xl:grid-cols-5">
          <ProductCard v-for="product in featuredProducts" :key="product.id" :product="product" />
        </div>
      </div>
    </section>

    <!-- Why Choose Us -->
    <section class="py-10 bg-gray-50">
      <div class="container">
        <h2 class="text-2xl font-bold text-center mb-8">Mengapa Memilih Kami</h2>
        <div class="grid md:grid-cols-3 gap-6">
          <div class="text-center">
            <div
              class="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <Icon name="tabler:truck-delivery" class="text-blue-600" size="32" />
            </div>
            <h3 class="font-semibold mb-2">Pengiriman Cepat</h3>
            <p class="text-gray-600">Pengiriman di hari yang sama untuk area tertentu</p>
          </div>
          <div class="text-center">
            <div
              class="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <Icon name="tabler:leaf" class="text-blue-600" size="32" />
            </div>
            <h3 class="font-semibold mb-2">100% Segar</h3>
            <p class="text-gray-600">Langsung dari petani terpercaya</p>
          </div>
          <div class="text-center">
            <div
              class="bg-blue-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <Icon name="tabler:wallet" class="text-blue-600" size="32" />
            </div>
            <h3 class="font-semibold mb-2">Pembayaran Aman</h3>
            <p class="text-gray-600">Berbagai metode pembayaran tersedia</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ProductCard from '~/components/home/ProductCard.vue'
import { Swiper, SwiperSlide } from 'swiper/vue'
import { Autoplay, EffectFade, Pagination, Navigation } from 'swiper/modules'
import 'swiper/css'
import 'swiper/css/effect-fade'
import 'swiper/css/pagination'
import 'swiper/css/navigation'
import { useCategoryStore } from '~/stores/category'
import { useProductStore } from '~/stores/product'
const categories = ref([])
const categoryStore = useCategoryStore()
const featuredProducts = ref([])
const productStore = useProductStore()
const loading = ref(true)
const { showError } = useErrorModal()

// Definisikan modul Swiper yang akan digunakan
const SwiperAutoplay = Autoplay
const SwiperEffectFade = EffectFade
const SwiperPagination = Pagination
const SwiperNavigation = Navigation

onMounted(async () => {
  try {
    await Promise.all([categoryStore.fetchCategoriesHome(), productStore.fetchProductsHome()])
    categories.value = categoryStore.categories
    featuredProducts.value = productStore.products
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  categoryStore.cancelRequests()
  productStore.cancelRequests()
})
</script>

<style scoped>
.swiper-container {
  width: 100%;
  height: 100%;
}
</style>
