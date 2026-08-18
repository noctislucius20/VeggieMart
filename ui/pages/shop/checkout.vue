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
                  <a href="../shop">
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
                <li class="inline-block text-gray-500 active" aria-current="page">Shop Checkout</li>
              </ol>
            </nav>
          </div>
        </div>
      </div>
    </div>
    <div class="my-10">
      <div class="container">
        <div class="flex flex-wrap">
          <div class="w-full mb-6">
            <!-- card -->

            <h1 class="text-xl">Checkout</h1>
          </div>
        </div>

        <div class="flex flex-wrap lg:flex-nowrap gap-10">
          <div class="lg:w-3/5 w-full">
            <!-- accordion -->
            <div id="accordionFlushExample" class="accordion accordion-flush">
              <!-- accordion item -->
              <div class="border-b border-gray-300 py-4">
                <div class="flex justify-between items-center">
                  <button
                    class="flex flex-row gap-2 items-center text-gray-900 text-md font-bold"
                    :aria-expanded="activeSection === 'address'"
                    @click="toggleSection('address')"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="icon icon-tabler icon-tabler-map-pin inline-block text-gray-500"
                      width="20"
                      height="20"
                      viewBox="0 0 24 24"
                      stroke-width="1.5"
                      stroke="currentColor"
                      fill="none"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path d="M9 11a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
                      <path
                        d="M17.657 16.657l-4.243 4.243a2 2 0 0 1 -2.827 0l-4.244 -4.243a8 8 0 1 1 11.314 0z"
                      />
                    </svg>
                    Delivery address
                  </button>
                </div>
                <div v-show="activeSection === 'address'" class="my-6">
                  <!-- Loading Skeleton -->
                  <div v-if="loading">
                    <div class="md:flex gap-6">
                      <div class="lg:w-1/2 w-full">
                        <div class="card card-body flex-col gap-4">
                          <div class="skeleton skeleton-text w-24 h-5" />
                          <div class="skeleton skeleton-text w-full" />
                          <div class="skeleton skeleton-text w-3/4" />
                        </div>
                      </div>
                    </div>
                    <div class="mt-5 flex justify-end gap-2">
                      <div class="skeleton skeleton-button" style="width: 100px; height: 44px" />
                      <div class="skeleton skeleton-button" style="width: 100px; height: 44px" />
                    </div>
                  </div>
                  <!-- Content -->
                  <div v-else class="md:flex gap-6">
                    <div class="lg:w-1/2 w-full">
                      <div class="card card-body flex-col gap-4">
                        <div class="relative flex items-center gap-2">
                          <input
                            id="homeRadio"
                            v-model="selectedAddress"
                            class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded-full focus:ring-blue-600 focus:outline-none focus:ring-2"
                            type="radio"
                            value="home"
                          />
                          <label class="text-gray-800 inline-block" for="homeRadio">Home</label>
                        </div>
                        <address class="not-italic">
                          {{ formData.address }}
                          <br />
                          <abbr title="Phone">{{ formData.phone }}</abbr>
                        </address>
                      </div>
                    </div>
                  </div>
                  <div v-if="!loading" class="mt-5 flex justify-end gap-2">
                    <button
                      :disabled="isProcessing"
                      class="btn inline-flex items-center gap-x-2 bg-white text-gray-800 border-gray-300 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
                      @click="prevSection"
                    >
                      Prev
                    </button>
                    <button
                      :disabled="isProcessing"
                      class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                      @click="nextSection"
                    >
                      Next
                    </button>
                  </div>
                </div>
              </div>

              <!-- Delivery Instructions Section -->
              <div class="border-b border-gray-300 py-4">
                <button
                  class="flex flex-row gap-2 items-center text-gray-900 text-md font-bold"
                  :aria-expanded="activeSection === 'instructions'"
                  @click="toggleSection('instructions')"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon icon-tabler icon-tabler-shopping-bag inline-block text-gray-500"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path
                      d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304z"
                    />
                    <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
                  </svg>
                  Delivery Method
                </button>
                <div v-show="activeSection === 'instructions'" class="my-6">
                  <!-- Loading Skeleton -->
                  <div v-if="loading">
                    <div class="mb-4">
                      <div class="skeleton skeleton-text w-32 mb-3" style="height: 18px" />
                      <div class="skeleton skeleton-text w-full" style="height: 44px" />
                    </div>
                    <div class="mb-4">
                      <div class="skeleton skeleton-text w-full" style="height: 80px" />
                    </div>
                    <div class="mt-5 flex justify-end gap-2">
                      <div class="skeleton skeleton-button" style="width: 100px; height: 44px" />
                      <div class="skeleton skeleton-button" style="width: 100px; height: 44px" />
                    </div>
                  </div>
                  <!-- Content -->
                  <div v-else class="mb-4">
                    <div class="flex gap-4">
                      <div class="flex items-center">
                        <input
                          id="delivery"
                          v-model="form.shipping_type"
                          type="radio"
                          value="DELIVERY"
                          class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded-full focus:ring-blue-600 focus:outline-none focus:ring-2"
                        />
                        <label for="delivery" class="ml-2 text-gray-800">Delivery</label>
                      </div>
                      <div class="flex items-center">
                        <input
                          id="pickup"
                          v-model="form.shipping_type"
                          type="radio"
                          value="PICKUP"
                          class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded-full focus:ring-blue-600 focus:outline-none focus:ring-2"
                        />
                        <label for="pickup" class="ml-2 text-gray-800">Pickup</label>
                      </div>
                    </div>
                    <div
                      v-if="v$.shipping_type.$error && v$.shipping_type.$dirty"
                      class="text-red-600 text-sm mt-1"
                    >
                      <template v-if="v$.shipping_type.required.$invalid"
                        >Metode pengiriman harus dipilih</template
                      >
                    </div>
                  </div>

                  <textarea
                    v-if="!loading"
                    v-model="form.remarks"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    rows="3"
                    :placeholder="
                      form.shipping_type === 'PICKUP'
                        ? 'Tulis catatan untuk pickup'
                        : 'Tulis instruksi pengiriman'
                    "
                  />
                  <p v-if="!loading" class="block mt-1 small text-gray-500">
                    {{
                      form.shipping_type === 'PICKUP'
                        ? 'Tambahkan catatan untuk pengambilan pesanan'
                        : 'Tambahkan instruksi untuk pengiriman pesanan Anda'
                    }}
                  </p>
                  <div v-if="!loading" class="mt-5 flex justify-end gap-2">
                    <button
                      :disabled="isProcessing"
                      class="btn inline-flex items-center gap-x-2 bg-white text-gray-800 border-gray-300 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
                      @click="prevSection"
                    >
                      Prev
                    </button>
                    <button
                      :disabled="isProcessing"
                      class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                      @click="nextSection"
                    >
                      Next
                    </button>
                  </div>
                </div>
              </div>

              <!-- Payment Method Section -->
              <div class="py-4">
                <button
                  class="flex flex-row gap-2 items-center text-gray-900 text-md font-bold"
                  :aria-expanded="activeSection === 'payment'"
                  @click="toggleSection('payment')"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="icon icon-tabler icon-tabler-credit-card inline-block text-gray-500"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path
                      d="M3 5m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v8a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"
                    />
                    <path d="M3 10l18 0" />
                    <path d="M7 15l.01 0" />
                    <path d="M11 15l2 0" />
                  </svg>
                  Payment Method
                </button>
                <div v-show="activeSection === 'payment'" class="mt-6 flex flex-col gap-4">
                  <!-- Loading Skeleton -->
                  <div v-if="loading">
                    <div class="card cursor-pointer mb-3">
                      <div class="flex-col flex p-6 gap-4">
                        <div class="flex gap-3">
                          <div class="skeleton w-4 h-4 rounded" />
                          <div class="flex flex-col gap-1">
                            <div class="skeleton skeleton-text w-32" style="height: 18px" />
                            <div class="skeleton skeleton-text w-48" />
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="card cursor-pointer mb-3">
                      <div class="flex-col flex p-6 gap-4">
                        <div class="flex gap-3">
                          <div class="skeleton w-4 h-4 rounded" />
                          <div class="flex flex-col gap-1">
                            <div class="skeleton skeleton-text w-32" style="height: 18px" />
                            <div class="skeleton skeleton-text w-48" />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <!-- Content -->
                  <template v-else>
                    <label class="card cursor-pointer">
                      <div class="flex-col flex p-6 gap-4">
                        <div class="flex gap-3">
                          <div class="relative flex mt-2">
                            <input
                              id="cashonDelivery"
                              v-model="form.payment_type"
                              class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded-full focus:ring-blue-600 focus:outline-none focus:ring-2"
                              type="radio"
                              value="COD"
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <h5 class="text-base">Cash on Delivery</h5>
                            <p class="text-sm">Pay with cash when your order is delivered.</p>
                          </div>
                        </div>
                      </div>
                    </label>
                    <label class="card cursor-pointer">
                      <div class="flex-col flex p-6 gap-4">
                        <div class="flex gap-3">
                          <div class="relative flex mt-2">
                            <input
                              id="payoneer"
                              v-model="form.payment_type"
                              class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded-full focus:ring-blue-600 focus:outline-none focus:ring-2"
                              type="radio"
                              value="TRANSFER"
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <h5 class="text-base">Pay with Midtrans</h5>
                            <p class="text-sm">
                              You will be redirected to Midtrans website to complete your purchase
                              securely.
                            </p>
                          </div>
                        </div>
                      </div>
                    </label>
                    <div
                      v-if="v$.payment_type.$error && v$.payment_type.$dirty"
                      class="text-red-600 text-sm"
                    >
                      <template v-if="v$.payment_type.required.$invalid"
                        >Metode pembayaran harus dipilih</template
                      >
                    </div>
                  </template>
                  <div v-if="!loading" class="flex justify-end">
                    <button
                      :disabled="isProcessing"
                      class="btn inline-flex items-center gap-x-2 bg-white text-gray-800 border-gray-300 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
                      @click="prevSection"
                    >
                      Prev
                    </button>
                    <button
                      :disabled="isProcessing"
                      class="ml-3 btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                      @click="placeOrder"
                    >
                      Place Order
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="w-full md:w-full lg:w-1/3 lg:ml-14">
            <div>
              <!-- Loading Skeleton for Cart Summary -->
              <div v-if="loading" class="card shadow-sm">
                <h5
                  class="px-6 py-4 border-b border-gray-300 skeleton skeleton-text w-40"
                  style="height: 24px"
                />
                <ul class="flex flex-col">
                  <li
                    v-for="i in 3"
                    :key="'skeleton-cart-' + i"
                    class="py-3 px-6 border-b border-gray-300"
                  >
                    <div class="flex flex-wrap items-center gap-2">
                      <div class="skeleton w-10 h-10 rounded" />
                      <div class="flex-1">
                        <div class="skeleton skeleton-text w-24 mb-2" />
                        <div class="skeleton skeleton-text w-16" />
                      </div>
                      <div class="skeleton skeleton-text w-8" />
                      <div class="skeleton skeleton-text w-20" />
                    </div>
                  </li>
                  <li class="py-3 px-6">
                    <div class="flex items-center justify-between mb-2">
                      <div class="skeleton skeleton-text w-32" />
                      <div class="skeleton skeleton-text w-24" />
                    </div>
                    <div class="flex items-center justify-between mb-2">
                      <div class="skeleton skeleton-text w-32" />
                      <div class="skeleton skeleton-text w-24" />
                    </div>
                    <div
                      class="flex items-center justify-between border-t border-gray-300 pt-2 mt-2"
                    >
                      <div class="skeleton skeleton-text w-32" style="height: 18px" />
                      <div class="skeleton skeleton-text w-24" style="height: 18px" />
                    </div>
                  </li>
                </ul>
              </div>
              <!-- Cart Summary Content -->
              <div v-else class="card shadow-sm">
                <h5 class="px-6 py-4 border-b border-gray-300">Detail Pesanan</h5>
                <ul class="flex flex-col">
                  <!-- list group item -->
                  <li
                    v-for="item in cartItems"
                    :key="item.id"
                    class="py-3 px-6 border-b border-gray-300"
                  >
                    <div class="flex flex-wrap items-center">
                      <div class="w-1/5 md:w-1/5">
                        <img :src="item.product_image" alt="Ecommerce" class="w-10" />
                      </div>
                      <div class="w-2/5 md:w-2/5 flex flex-col flex-wrap gap-1">
                        <h6>{{ item.product_name }}</h6>
                        <span class="text-gray-700 text-sm"
                          >{{ item.weight }} / {{ item.unit }}</span
                        >
                      </div>
                      <div class="w-1/5 md:w-1/5 text-center text-gray-700">
                        <span>{{ item.quantity }}</span>
                      </div>
                      <div class="w-1/5 text-center md:w-1/5">
                        <span class="font-bold text-gray-800"
                          >Rp. {{ formatPrice(item.sale_price) }}</span
                        >
                      </div>
                    </div>
                  </li>
                  <!-- list group item -->
                  <li class="py-3 px-6">
                    <div class="flex items-center justify-between text-gray-800 mb-1">
                      <div>Subtotal Pesanan</div>
                      <div>Rp. {{ formatPrice(subtotal) }}</div>
                    </div>
                    <div class="flex items-center justify-between text-gray-800">
                      <div>Biaya Pengiriman</div>
                      <div>
                        Rp.
                        {{
                          form.shipping_type.toLowerCase() == 'delivery'
                            ? formatPrice(5000)
                            : formatPrice(0)
                        }}
                      </div>
                    </div>
                    <div
                      class="flex items-center justify-between font-bold text-gray-800 border-t border-gray-300 pt-2 mt-2"
                    >
                      <div>Total Pembayaran</div>
                      <div>Rp. {{ formatPrice(total) }}</div>
                    </div>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- Processing Modal Overlay -->
    <div
      v-if="isProcessing"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    >
      <div class="bg-white rounded-lg p-8 shadow-xl flex flex-col items-center gap-4 mx-4">
        <div
          class="animate-spin rounded-full h-12 w-12 border-4 border-blue-600 border-t-transparent"
        />
        <p class="text-lg font-semibold text-gray-800">Memproses Pesanan</p>
        <p class="text-sm text-gray-500">Mohon tunggu sebentar...</p>
      </div>
    </div>
  </main>
</template>

<script setup>
import { useCartStore } from '~/stores/cart'
import { onMounted, onBeforeUnmount, ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useOrderStore } from '~/stores/orders'
import { usePaymentStore } from '~/stores/payment'
import { useVuelidate } from '@vuelidate/core'
import { required, maxLength } from '@vuelidate/validators'
import { oneof, dateFormat, timeFormat } from '~/utils/validators'

const config = useRuntimeConfig()
const cartStore = useCartStore()
const authStore = useAuthStore()
const router = useRouter()
const orderStore = useOrderStore()
const paymentStore = usePaymentStore()
const cartItems = ref([])
const loading = ref(true)
const isProcessing = ref(false)
const { showError } = useErrorModal()

const sections = ['address', 'instructions', 'payment']
const activeSection = ref('address')
const selectedAddress = ref('home')

const formData = ref({
  name: '',
  email: '',
  phone: '',
  address: '',
  lat: null,
  lng: null,
  photo: ''
})

const form = reactive({
  shipping_type: '',
  payment_type: '',
  remarks: '',
  order_date: '',
  order_time: ''
})

const rules = computed(() => ({
  shipping_type: {
    required,
    maxLength: maxLength(20),
    oneof: oneof(['PICKUP', 'DELIVERY'])
  },
  payment_type: {
    required,
    maxLength: maxLength(50),
    oneof: oneof(['COD', 'TRANSFER'])
  },
  remarks: {},
  order_date: {
    required,
    dateFormat
  },
  order_time: {
    required,
    timeFormat
  }
}))

const v$ = useVuelidate(rules, form)

const generateIdempotencyKey = () => {
  const raw = [
    authStore.user.id,
    JSON.stringify(cartItems.value.map((i) => `${i.id}:${i.quantity}`)),
    form.payment_type,
    form.shipping_type,
    form.remarks
  ].join('|')

  return btoa(raw).replace(/=+$/, '')
}

const toggleSection = (section) => {
  activeSection.value = section
}

const nextSection = () => {
  const currentIndex = sections.indexOf(activeSection.value)

  // Validate current section before moving forward
  if (activeSection.value === 'instructions') {
    v$.value.shipping_type.$touch()
    if (v$.value.shipping_type.$invalid) {
      return
    }
  }

  if (currentIndex < sections.length - 1) {
    activeSection.value = sections[currentIndex + 1]
  }
}

const prevSection = () => {
  const currentIndex = sections.indexOf(activeSection.value)
  if (currentIndex > 0) {
    activeSection.value = sections[currentIndex - 1]
  }
}

const placeOrder = async () => {
  if (isProcessing.value) return

  // Set auto-generated values
  const currentTime = new Date()
  form.order_date = currentTime.toISOString().split('T')[0]
  form.order_time = currentTime.toTimeString().split(' ')[0]

  // Validate all fields
  v$.value.$touch()
  if (v$.value.$invalid) {
    // Scroll to the first section with error
    if (v$.value.shipping_type.$invalid) {
      activeSection.value = 'instructions'
    } else if (v$.value.payment_type.$invalid) {
      activeSection.value = 'payment'
    }
    return
  }

  // Validate cart items
  if (!cartItems.value || cartItems.value.length === 0) {
    alert('Keranjang belanja kosong')
    return
  }

  isProcessing.value = true
  try {
    const idempotencyKey = generateIdempotencyKey()

    const orderData = {
      order_date: form.order_date,
      shipping_type: form.shipping_type,
      order_time: form.order_time,
      payment_type: form.payment_type,
      remarks: form.remarks,
      order_details: cartItems.value.map((item) => ({
        product_id: item.id,
        quantity: item.quantity
      }))
    }

    await orderStore.createOrders(orderData, formData.value.lat, formData.value.lng, idempotencyKey)
    await cartStore.deleteAllCart()

    const paymentData = {
      order_id: orderStore.order.order_id,
      remarks: form.remarks,
      payment_method: form.payment_type
    }

    await paymentStore.createPayment(paymentData)
    if (paymentStore.payment.payment_token === null) {
      router.push('/shop/success-payment')
      return
    }

    if (window.snap) {
      window.snap.pay(paymentStore.payment.payment_token, {
        onSuccess: function () {
          router.push('/shop/success-payment')
        },
        onPending: function () {
          alert('Menunggu Pembayaran')
        },
        onError: function () {
          alert('Pembayaran gagal')
        },
        onClose: function () {
          alert('Anda menutup popup tanpa menyelesaikan pembayaran')
          router.push('/account')
        }
      })
    }
  } catch (error) {
    alert(`${error}`)
  } finally {
    isProcessing.value = false
  }
}

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/auth/signin')
    return
  }

  try {
    await cartStore.fetchCarts()
    cartItems.value = cartStore.carts

    if (cartItems.value.length == 0) {
      router.push('/shop')
      return
    }

    await authStore.getProfile()
    if (authStore.user.address == '') {
      router.push('/account/setting')
      return
    }

    await getLocation()

    formData.value = {
      name: authStore.user.name || '',
      email: authStore.user.email || '',
      phone: authStore.user.phone || '',
      address: authStore.user.address || '',
      lat: authStore.user.lat || null,
      lng: authStore.user.lng || null,
      photo: authStore.user.photo || ''
    }
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  cartStore.cancelRequests()
  authStore.cancelRequests()
  orderStore.cancelRequests()
  paymentStore.cancelRequests()
})

const getCoordinates = () => {
  return new Promise((resolve, reject) => {
    navigator.geolocation.getCurrentPosition(resolve, reject)
  })
}

const getLocation = async () => {
  try {
    if (!(import.meta.client && 'geolocation' in navigator)) {
      throw new Error('Geolocation not supported or not in browser')
    }

    const position = await getCoordinates()

    authStore.user.lat = position.coords.latitude
    authStore.user.lng = position.coords.longitude
  } catch (error) {
    alert(`Error: ${error}`)
  }
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const subtotal = computed(() => {
  return cartItems.value.reduce((total, item) => total + item.sale_price * item.quantity, 0)
})

const shippingCost = computed(() => {
  return form.shipping_type.toLowerCase() == 'delivery' ? 5000 : 0
})

const total = computed(() => {
  return subtotal.value + shippingCost.value
})

useScript({
  src: 'https://app.sandbox.midtrans.com/snap/snap.js',
  attrs: {
    'data-client-key': config.public.midtransClientKey
  },
  crossorigin: false
})
</script>
