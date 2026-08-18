<template>
  <div class="cart-sidebar fixed right-0 top-0 h-full bg-white shadow-lg w-96 z-50">
    <div class="h-full flex flex-col">
      <div class="border-b p-4 flex justify-between items-center">
        <h5>Keranjang Belanja</h5>
        <button type="button" class="text-gray-700" @click="$emit('close')">
          <Icon name="tabler:x" size="24" />
        </button>
      </div>

      <div class="grow h-full overflow-y-auto flex flex-col">
        <!-- Empty Cart State -->
        <div
          v-if="cartItems.length === 0"
          class="flex flex-col items-center justify-center flex-1 p-4"
        >
          <div class="text-center">
            <Icon name="tabler:shopping-cart-off" size="64" class="text-gray-300 mx-auto mb-4" />
            <h5 class="text-gray-700 font-semibold mb-2">Keranjang Belanja Kosong</h5>
            <p class="text-gray-500 text-sm mb-6">
              Belum ada produk di keranjang Anda. Mulai belanja sekarang!
            </p>
            <NuxtLink
              to="/shop"
              class="inline-block px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition text-sm font-medium"
              @click="$emit('close')"
            >
              Mulai Belanja
            </NuxtLink>
          </div>
        </div>

        <!-- Cart Items -->
        <div v-else class="p-4">
          <div class="text-right mb-3">
            <NuxtLink to="/shop/cart" class="text-blue-600">Lihat Semua</NuxtLink>
          </div>

          <ul class="list-none">
            <li v-for="item in cartItems" :key="item.id" class="py-3 border-t">
              <div class="flex items-center">
                <div class="w-1/2 md:w-1/2 lg:w-3/5">
                  <div class="flex">
                    <img
                      :src="item.product_image"
                      :alt="item.name"
                      class="w-16 h-16 object-cover"
                    />
                    <div class="ml-3">
                      <NuxtLink :to="`/shop/${item.id}`" class="text-inherit">
                        <h6>{{ item.product_name }}</h6>
                      </NuxtLink>
                      <span
                        ><small class="text-gray-500"
                          >{{ item.weight }} {{ item.unit }}</small
                        ></span
                      >
                      <div class="mt-2 small leading-none">
                        <button
                          class="text-blue-600 flex items-center"
                          @click="removeItem(item.id)"
                        >
                          <Icon name="tabler:trash" size="14" class="mr-1" />
                          <span class="text-gray-500 text-sm">Hapus</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="w-1/3 md:w-1/4 lg:w-1/5">
                  <div class="input-spinner rounded-lg flex justify-between items-center">
                    <button
                      class="w-8 py-1 border-r cursor-pointer border-gray-300"
                      @click="handleQuantityChange(item, 'decrement')"
                    >
                      -
                    </button>
                    <input
                      v-model="item.quantity"
                      type="number"
                      class="w-9 px-2 text-center h-7 border-0 bg-transparent"
                      min="1"
                      max="10"
                      @change="() => handleQuantityChange(item, '')"
                    />
                    <button
                      class="w-8 py-1 border-l cursor-pointer border-gray-300"
                      @click="handleQuantityChange(item, 'increment')"
                    >
                      +
                    </button>
                  </div>
                </div>

                <div class="w-1/5 text-center">
                  <span class="font-bold text-gray-800">Rp {{ formatPrice(item.sale_price) }}</span>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>

      <div v-if="cartItems.length > 0" class="border-t p-4">
        <div class="flex justify-between mb-3">
          <span>Subtotal:</span>
          <span class="font-bold">Rp {{ formatPrice(subtotal) }}</span>
        </div>
        <NuxtLink
          :to="`/shop/checkout`"
          class="btn w-full bg-blue-600 text-white hover:bg-blue-700 focus:ring-4 focus:ring-blue-300"
          @click="$emit('close')"
        >
          Checkout
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useCartStore } from '~/stores/cart'
import { onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'

defineEmits(['close'])

const cartStore = useCartStore()
const authStore = useAuthStore()
const { showError } = useErrorModal()
const cartItems = computed(() => cartStore.carts.slice(0, 5))

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    cartStore.clearCarts()
    return
  }

  await cartStore.fetchCarts()
})

const subtotal = computed(() => {
  return cartItems.value.reduce(
    (total, item) => total + Number(item.sale_price) * Number(item.quantity),
    0
  )
})

const handleQuantityChange = async (item, type) => {
  try {
    let quantity = item.quantity
    if (type === 'increment') {
      quantity += 1
    } else if (type === 'decrement') {
      quantity -= 1
    }

    if (quantity <= 0) {
      const isConfirmed = window.confirm('Apakah Anda yakin ingin menghapus produk dari keranjang?')

      if (isConfirmed) {
        await cartStore.deleteCart(item.id)
      }

      return
    }

    await cartStore.addToCart(parseInt(item.id), parseInt(quantity))
  } catch (error) {
    showError(error)
  }
}

const removeItem = async (itemId) => {
  if (!authStore.token) {
    // Jika belum login, arahkan ke halaman login
    alert('Silakan login terlebih dahulu')
    return
  }
  try {
    const isConfirmed = window.confirm('Apakah Anda yakin ingin menghapus produk dari keranjang?')

    if (isConfirmed) {
      await cartStore.deleteCart(itemId)
    }
  } catch (error) {
    showError(error)
  }
}

onBeforeUnmount(() => {
  cartStore.cancelRequests()
  authStore.cancelRequests()
})
</script>
