<template>
  <section class="container">
    <!-- row -->
    <div class="grid grid-cols-1 mb-8">
      <div class="md:flex justify-between items-center">
        <div>
          <h2 class="text-xl">Order Detail</h2>
          <!-- breacrumb -->
          <nav aria-label="breadcrumb">
            <ol class="flex flex-wrap">
              <li class="inline-block text-blue-600">
                <a href="/dashboard">
                  Dashboard
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="icon icon-tabler icons-tabler-outline icon-tabler-slash inline-block mx-2"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M17 5l-10 14" />
                  </svg>
                </a>
              </li>
              <li class="inline-block text-gray-500" aria-current="page">Order Detail</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
        <div class="mt-4 md:mt-0">
          <a
            href="/dashboard/orders"
            class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
          >
            Back to all orders
          </a>
        </div>
      </div>
    </div>
    <div class="grid grid-cols-1">
      <!-- Skeleton Loading -->
      <div v-if="loading" class="card h-full card-lg p-6">
        <div class="flex flex-col md:flex-row gap-4 justify-between mb-8">
          <div class="skeleton skeleton-button h-6 w-20 rounded" />
          <div class="skeleton skeleton-button h-10 w-40 rounded" />
        </div>
        <div class="grid grid-cols-12 gap-4 mb-8">
          <div v-for="i in 3" :key="i" class="lg:col-span-4 md:col-span-4 col-span-12">
            <div class="skeleton skeleton-title h-5 w-24 mb-3" />
            <div class="skeleton skeleton-text h-4 w-full mb-1" />
            <div class="skeleton skeleton-text h-4 w-3/4 mb-1" />
            <div class="skeleton skeleton-text h-4 w-1/2" />
          </div>
        </div>
        <div class="overflow-x-auto mb-6">
          <table class="text-left w-full whitespace-nowrap">
            <tbody>
              <tr>
                <td colspan="8">
                  <div class="py-4">
                    <div
                      v-for="i in 5"
                      :key="'skeleton-dash-' + i"
                      class="flex items-center gap-4 py-3 border-b border-gray-300"
                    >
                      <div class="skeleton w-12 h-12 rounded" />
                      <div class="flex-1">
                        <div class="skeleton skeleton-text w-48 mb-2" />
                        <div class="skeleton skeleton-text w-32" />
                      </div>
                    </div>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="grid md:grid-cols-2 gap-4 mb-4">
          <div>
            <div class="skeleton skeleton-title h-5 w-24 mb-3" />
            <div class="skeleton skeleton-text h-5 w-32" />
          </div>
          <div>
            <div class="skeleton skeleton-title h-5 w-24 mb-3" />
            <div class="skeleton skeleton-text h-5 w-48" />
          </div>
        </div>
        <div class="grid md:grid-cols-2 gap-4 mt-4">
          <div>
            <div class="skeleton skeleton-title h-5 w-16 mb-3" />
            <div class="skeleton skeleton-text h-20 w-full mb-3" />
            <div class="skeleton skeleton-button h-10 w-24 rounded" />
          </div>
        </div>
      </div>

      <!-- card -->
      <div v-if="!loading" class="card h-full card-lg">
        <div class="card-body p-6">
          <div class="flex flex-col md:flex-row gap-4 justify-between">
            <div class="flex items-center">
              <h3 class="mb-0 mr-3">Order ID: #{{ orderData?.order_code }}</h3>
              <span
                class="inline-block p-1 text-sm align-baseline leading-none rounded border font-semibold"
                :class="getOrderStatusClass(orderData?.status)"
              >
                {{ orderData?.status }}
              </span>
            </div>
            <!-- select option -->
            <div class="flex flex-col md:flex-row gap-3 md:items-center">
              <div class="">
                <select
                  v-model="form.status"
                  class="text-base py-2 block w-full border-gray-300 rounded-lg focus:border-blue-600 focus:ring-blue-600 disabled:opacity-50 disabled:pointer-events-none"
                  :class="{ 'border-red-500': v$.status.$error && v$.status.$dirty }"
                  @blur="v$.status.$touch()"
                >
                  <option value="" disabled>Pilih Status</option>
                  <option value="PENDING">Pending</option>
                  <option value="CONFIRMED">Confirmed</option>
                  <option value="PROCESSING">Processing</option>
                  <option value="SHIPPED">Shipped</option>
                  <option value="DELIVERED">Delivered</option>
                  <option value="CANCELLED">Cancelled</option>
                </select>
                <div v-if="v$.status.$error && v$.status.$dirty" class="text-red-600 text-sm mt-1">
                  <template v-if="v$.status.required.$invalid">Status harus dipilih</template>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-8">
            <div class="grid grid-cols-12">
              <!-- address -->
              <div class="lg:col-span-4 md:col-span-4 col-span-12">
                <div class="mb-6">
                  <h6>Customer Details</h6>
                  <p class="mb-2 leading-relaxed">
                    {{ orderData?.customer?.customer_name }}
                    <br />
                    {{ orderData?.customer?.customer_email }}
                    <br />
                    {{ orderData?.customer?.customer_phone }}
                  </p>
                  <!-- <a href="#" class="text-blue-600">View Profile</a> -->
                </div>
              </div>
              <!-- address -->
              <div class="lg:col-span-4 md:col-span-4 col-span-12">
                <div class="mb-6">
                  <h6>Shipping Address</h6>
                  <p class="mb-1 leading-relaxed">
                    {{ orderData?.customer?.customer_address }}
                    <br />
                    Contact No. {{ orderData?.customer?.customer_phone }}
                  </p>
                </div>
              </div>
              <!-- address -->
              <div class="lg:col-span-4 md:col-span-4 col-span-12">
                <div class="mb-6">
                  <h6>Order Details</h6>
                  <p class="mb-1 leading-relaxed">
                    Order ID:
                    <span class="text-gray-800">
                      {{ orderData?.order_code }}
                    </span>
                    <br />
                    Order Date:
                    <span class="text-gray-800">
                      {{ orderData?.order_datetime }}
                    </span>
                    <br />
                    Order Total:
                    <span class="text-gray-800"
                      >Rp. {{ formatPrice(orderData?.total_amount) }}</span
                    >
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 px-5">
          <div class="relative overflow-x-auto">
            <table class="text-left w-full whitespace-nowrap">
              <thead class="bg-gray-200 text-gray-700">
                <tr class="border-transparent border-b-0!">
                  <th scope="col" class="px-6 py-3">Products</th>
                  <th scope="col" class="px-6 py-3">Price</th>
                  <th scope="col" class="px-6 py-3">Quantity</th>
                  <th scope="col" class="px-6 py-3">Total</th>
                </tr>
              </thead>
              <!-- tbody -->
              <tbody>
                <tr
                  v-for="item in orderData?.order_items"
                  :key="item.id"
                  class="border-b border-gray-300"
                >
                  <td class="py-3 px-6 text-left">
                    <a href="#" class="text-inherit">
                      <div class="flex items-center gap-2 w-full">
                        <div class="flex-0">
                          <img :src="item.product_image" alt="" class="h-12 w-12" />
                        </div>
                        <div class="">
                          <h5 class="mb-0 text-base">{{ item.product_name }}</h5>
                        </div>
                      </div>
                    </a>
                  </td>
                  <td class="py-3 px-6 text-left">
                    <span class="">Rp. {{ formatPrice(item.product_price) }}</span>
                  </td>
                  <td class="py-3 px-6 text-left">{{ item.quantity }}</td>
                  <td class="py-3 px-6 text-left">
                    Rp. {{ formatPrice(item.product_price * item.quantity) }}
                  </td>
                </tr>
                <tr class="">
                  <td class="pb-0 py-3 px-6 text-left" />
                  <td class="pb-0 py-3 px-6 text-left" />
                  <td
                    colspan="1"
                    class="font-medium text-gray-800 py-3 px-6 text-left border-b border-gray-300"
                  >
                    <!-- text -->
                    Sub Total :
                  </td>
                  <td class="font-medium text-gray-800 border-b border-gray-300">
                    <!-- text -->
                    Rp. {{ formatPrice(calculateSubtotal) }}
                  </td>
                </tr>
                <tr class="">
                  <td class="pb-0 py-3 px-6 text-left" />
                  <td class="pb-0 py-3 px-6 text-left" />
                  <td
                    colspan="1"
                    class="font-medium text-gray-800 py-3 px-6 text-left border-b border-gray-300"
                  >
                    <!-- text -->
                    Shipping Cost
                  </td>
                  <td class="font-medium text-gray-800 border-b border-gray-300">
                    <!-- text -->
                    Rp. {{ formatPrice(orderData?.shipping_fee) }}
                  </td>
                </tr>

                <tr class="">
                  <td class="py-3 px-6 text-left" />
                  <td class="py-3 px-6 text-left" />
                  <td
                    colspan="1"
                    class="font-semibold text-gray-800 py-3 px-6 text-left border-b border-gray-300"
                  >
                    <!-- text -->
                    Grand Total
                  </td>
                  <td class="font-semibold text-gray-800 py-3 border-b border-gray-300">
                    <!-- text -->
                    Rp. {{ formatPrice(orderData?.shipping_fee + calculateSubtotal) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="card-body">
          <div class="grid md:grid-cols-2 gap-4">
            <div class="">
              <h6 class="mb-2">Payment Info</h6>
              <span>{{
                paymentData
                  ? paymentData?.payment_method?.toLowerCase() === 'midtrans'
                    ? 'Transfer'
                    : 'Cash or Duel'
                  : '-'
              }}</span>
            </div>
            <div class="">
              <h6 class="mb-3">Delivery Notes</h6>
              <p class="text-gray-700">{{ orderData ? orderData.remarks : '-' }}</p>
            </div>
          </div>
          <div class="grid md:grid-cols-2 gap-4">
            <div class="">
              <h6 class="mb-2">Payment Status</h6>
              <span
                class="inline-block text-sm align-baseline leading-none font-semibold"
                :class="getPaymentStatusClass(paymentData?.payment_status)"
              >
                {{ paymentData ? getPaymentStatusLabel(paymentData?.payment_status) : '-' }}
              </span>
            </div>
          </div>
          <div class="grid md:grid-cols-2 gap-4 mt-4">
            <div class="">
              <h5 class="mb-3">Notes</h5>
              <textarea
                v-model="form.remarks"
                class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base mb-3"
                rows="3"
                placeholder="Write note for order"
              />
              <button
                type="button"
                class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                @click="handleUpdateStatus(orderData.id)"
              >
                Submit
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useOrderStore } from '~/stores/orders'
import { usePaymentStore } from '~/stores/payment'
import { useVuelidate } from '@vuelidate/core'
import { required, maxLength } from '@vuelidate/validators'
import { oneof } from '~/utils/validators'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const router = useRouter()
const route = useRoute()
const orderStore = useOrderStore()
const paymentStore = usePaymentStore()
const orderData = ref(null)
const paymentData = ref(null)
const { showError } = useErrorModal()

const form = reactive({
  status: '',
  remarks: ''
})

const rules = computed(() => ({
  status: {
    required,
    maxLength: maxLength(20),
    oneof: oneof(['PENDING', 'CONFIRMED', 'PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED'])
  },
  remarks: {}
}))

const v$ = useVuelidate(rules, form)

const loading = ref(true)

const getPaymentStatusClass = (status) => {
  const statusClass = {
    pending: 'p-1 rounded border bg-yellow-100 text-yellow-800 border-yellow-200',
    success: 'p-1 rounded border bg-blue-100 text-blue-800 border-blue-200',
    failed: 'p-1 rounded border bg-red-100 text-red-800 border-red-200'
  }

  return statusClass[String(status).toLowerCase()]
}

const getPaymentStatusLabel = (status) => {
  const labels = {
    pending: 'MENUNGGU PEMBAYARAN',
    success: 'LUNAS',
    failed: 'PEMBAYARAN GAGAL'
  }

  return labels[String(status).toLowerCase()] || status || ''
}

const getOrderStatusClass = (productStatus) => {
  const statusClass = {
    PENDING: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    CONFIRMED: 'bg-blue-100 text-blue-800 border-blue-200',
    PROCESSING: 'bg-purple-100 text-purple-800 border-purple-200',
    SHIPPED: 'bg-cyan-100 text-cyan-800 border-cyan-200',
    DELIVERED: 'bg-green-100 text-green-800 border-green-200',
    CANCELLED: 'bg-red-100 text-red-800 border-red-200'
  }

  return (
    statusClass[String(productStatus).toUpperCase()] || 'bg-gray-100 text-gray-800 border-gray-200'
  )
}

onMounted(async () => {
  try {
    await orderStore.getDetailOrders(route.params.id, true)
    await paymentStore.getDetailPaymentByOrderId(route.params.id, true)

    orderData.value = orderStore.order
    paymentData.value = paymentStore.payment
    form.status = orderStore.order.status
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  orderStore.cancelRequests()
  paymentStore.cancelRequests()
})

const calculateSubtotal = computed(() => {
  if (!orderData.value?.order_items) return 0

  return orderData.value.order_items.reduce((total, item) => {
    return total + item.product_price * item.quantity
  }, 0)
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

const handleUpdateStatus = async (id) => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    await orderStore.updateStatusOrder(id, {
      status: form.status,
      remarks: form.remarks
    })
    router.push('/dashboard/orders')
  } catch (error) {
    showError(error)
  }
}
</script>

<style></style>
