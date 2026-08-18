<template>
  <div class="relative rounded-lg wrap-break-word border bg-white border-gray-300 card-product">
    <div class="flex-auto p-4">
      <div class="text-center relative flex justify-center">
        <div v-if="product.reguler_price != product.sale_price" class="absolute top-0 left-0">
          <span
            class="inline-block p-1 text-center font-semibold text-sm align-baseline leading-none rounded bg-red-600 text-white"
            >Sale</span
          >
        </div>
        <NuxtLink :to="`/shop/${product.id}`">
          <img :src="product.product_image" :alt="product.product_name" class="w-full h-auto" />
        </NuxtLink>

        <!-- <div class="absolute w-full bottom-[15%] opacity-0 invisible card-product-action">
          <button type="button"
            class="h-[34px] w-[34px] leading-[34px] bg-white shadow inline-flex items-center justify-center rounded-lg hover:bg-blue-600 hover:text-white"
            data-bs-toggle="tooltip" data-bs-target="#quickViewModal" data-bs-html="true" title="Quick View">
          
            <span data-bs-toggle="modal" @click="showQickViewModal = true">
              <svg xmlns="http://www.w3.org/2000/svg"
                class="icon icon-tabler icon-tabler-eye" width="24" height="24"
                viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"
                fill="none" stroke-linecap="round" stroke-linejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
                <path
                  d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6" />
              </svg>
            </span>
          </button>
        </div> -->
      </div>
      <div class="flex flex-col gap-3">
        <NuxtLink :to="`/shop/${product.id}`" class="text-decoration-none text-gray-500"
          ><small>{{ product.category_name }}</small></NuxtLink
        >
        <div class="flex flex-col gap-2">
          <h3 class="text-base truncate">
            <NuxtLink :to="`/shop/${product.id}`">{{ product.product_name }}</NuxtLink>
          </h3>
        </div>
        <div class="flex justify-between items-center">
          <div>
            <span class="text-gray-900 font-semibold"
              >Rp. {{ formatPrice(product.sale_price) }}</span
            >
            <span
              v-if="product.regular_price != product.sale_price"
              class="line-through text-gray-500"
              >Rp. {{ formatPrice(product.regular_price) }}</span
            >
          </div>
          <div>
            <button
              type="button"
              class="btn inline-flex items-center gap-x-1 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300 btn-sm"
              @click="addToCart"
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
  <QuickView v-if="showQickViewModal" @close="showQickViewModal = false" />
</template>

<script setup>
import QuickView from '../modals/QuickView.vue'
import { useCartStore } from '~/stores/cart'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

const cartStore = useCartStore()
const authStore = useAuthStore()
const router = useRouter()

const props = defineProps({
  product: {
    type: Object,
    default: () => ({
      id: '',
      name: '',
      price: 0,
      image: '',
      category: '',
      unit: '',
      discount: 0
    })
  }
})

// const isWishlist = ref(false)
const showQickViewModal = ref(false)
// const emit = defineEmits(['add-to-cart'])

const { showError } = useErrorModal()

// Menghitung harga setelah diskon
// const calculateDiscountedPrice = computed(() => {
//   if (!props.product.discount) return props.product.price

//   const discount = props.product.price * (props.product.discount / 100)
//   return props.product.price - discount
// })

// Format harga ke format Rupiah
const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

// Toggle wishlist
// const toggleWishlist = () => {
//   isWishlist.value = !isWishlist.value
// }

// const showQuickView = () => {
//   showQuickViewModal.value = true
// }

// Tambah ke keranjang
const addToCart = async () => {
  if (!authStore.token) {
    // Jika belum login, arahkan ke halaman login
    alert('Silakan login terlebih dahulu')
    router.push('/auth/signin')
    return
  }
  try {
    await cartStore.addToCart(props.product.id, 1)
    alert('Produk berhasil ditambahkan ke keranjang')
  } catch (error) {
    showError(error)
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
