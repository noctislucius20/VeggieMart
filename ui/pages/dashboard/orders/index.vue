<template>
  <section class="container">
    <!-- row -->
    <div class="grid grid-cols-1 mb-8">
      <div class="md:flex justify-between items-center">
        <div>
          <h2 class="text-xl">Orders</h2>
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

              <li class="inline-block text-gray-500 active" aria-current="page">Orders</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
      </div>
    </div>
    <div class="grid grid-cols-1">
      <!-- card -->
      <div class="card h-full card-lg">
        <div class="px-6 py-6">
          <div class="grid grid-cols-12 justify-between gap-2">
            <div class="lg:col-span-3 md:col-span-6 col-span-12">
              <!-- form -->
              <form class="flex" role="search">
                <input
                  v-model="searchFilter"
                  :disabled="useOrder.loading"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  type="search"
                  placeholder="Search"
                  aria-label="Search"
                />
              </form>
            </div>
            <!-- select option -->
            <div class="md:col-start-11 md:col-end-13 md:col-span-4 col-span-12">
              <select
                v-model="statusFilter"
                :disabled="useOrder.loading"
                class="text-base py-2 block w-full border-gray-300 rounded-lg focus:border-blue-600 focus:ring-blue-600 disabled:opacity-50 disabled:pointer-events-none"
                @change="handleFilterChange()"
              >
                <option value="" disabled>Pilih Status</option>
                <option value="PENDING">Pending</option>
                <option value="CONFIRMED">Confirmed</option>
                <option value="SENDING">Sending</option>
                <option value="DONE">Done</option>
                <option value="CANCEL">Cancel</option>
              </select>
            </div>
          </div>
        </div>
        <!-- card body -->
        <div class="card-body p-0">
          <!-- table -->
          <div class="relative overflow-x-auto">
            <table class="text-left w-full whitespace-nowrap table-with-checkbox table-hover">
              <thead class="bg-gray-200 text-gray-700">
                <tr class="border-transparent border-b-0!">
                  <th scope="col" class="px-6 py-3">Image</th>
                  <th scope="col" class="px-6 py-3">Order Name</th>
                  <th scope="col" class="px-6 py-3">Customer</th>
                  <th scope="col" class="px-6 py-3">Date & TIme</th>
                  <th scope="col" class="px-6 py-3">Status</th>
                  <th scope="col" class="px-6 py-3">Amount</th>
                  <th scope="col" class="px-6 py-3" />
                </tr>
              </thead>
              <tbody>
                <tr v-if="loading">
                  <td colspan="8" class="py-4">
                    <div
                      v-for="i in 5"
                      :key="'skeleton-dash-' + i"
                      class="flex items-center gap-4 py-3 border-b border-gray-100 px-6"
                    >
                      <div class="skeleton w-12 h-12 rounded" />
                      <div class="flex-1">
                        <div class="skeleton skeleton-text w-48 mb-2" />
                        <div class="skeleton skeleton-text w-32" />
                      </div>
                      <div class="skeleton skeleton-button" style="width: 24px; height: 24px" />
                    </div>
                  </td>
                </tr>
                <tr v-else-if="useOrder?.orders?.length === 0">
                  <td colspan="8" class="text-center py-4">No data available</td>
                </tr>
                <tr v-for="item in orderDatas" v-else :key="item.id">
                  <td class="py-3 px-6 text-left">
                    <a :href="`/dashboard/orders/detail/${item.id}`">
                      <img :src="item.product_image" alt="" class="h-10 w-10" />
                    </a>
                  </td>
                  <td class="py-3 px-6 text-left">
                    <a :href="`/dashboard/orders/detail/${item.id}`" class="text-inherit">{{
                      item.order_code
                    }}</a>
                  </td>
                  <td class="py-3 px-6 text-left">{{ item.customer_name }}</td>

                  <td class="py-3 px-6 text-left">{{ item.order_datetime }}</td>

                  <td class="py-3 px-6 text-left">
                    <span
                      class="inline-block p-1 text-sm align-baseline leading-none rounded font-semibold"
                      :class="getOrderStatusClass(item.status)"
                    >
                      {{ item.status }}
                    </span>
                  </td>
                  <td class="py-3 px-6 text-left">Rp {{ formatPrice(item.total_amount) }}</td>

                  <td class="py-3 px-6 text-left">
                    <div class="dropdown">
                      <a
                        href="#"
                        class="text-inherit"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="20"
                          height="20"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="icon icon-tabler icons-tabler-outline icon-tabler-dots-vertical"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path d="M12 12m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                          <path d="M12 19m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                          <path d="M12 5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                        </svg>
                      </a>
                      <ul class="dropdown-menu">
                        <li>
                          <a :href="`/dashboard/orders/detail/${item.id}`" class="dropdown-item">
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
                              class="icon icon-tabler icons-tabler-outline icon-tabler-edit"
                            >
                              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                              <path
                                d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"
                              />
                              <path
                                d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"
                              />
                              <path d="M16 5l3 3" />
                            </svg>
                            Edit
                          </a>
                        </li>
                        <li>
                          <a class="dropdown-item" href="#">
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
                              class="icon icon-tabler icons-tabler-outline icon-tabler-trash"
                            >
                              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                              <path d="M4 7l16 0" />
                              <path d="M10 11l0 6" />
                              <path d="M14 11l0 6" />
                              <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" />
                              <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" />
                            </svg>
                            Delete
                          </a>
                        </li>
                      </ul>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div
          v-if="useOrder.pagination?.total_count > 0"
          class="border-t border-gray-300 flex flex-col md:flex-row justify-between items-center px-6 py-6 gap-3"
        >
          <span
            >Showing {{ paginateOrders.page }} to {{ paginateOrders.total_page }} of
            {{ paginateOrders.total_count }} entries</span
          >
          <nav class="flex items-center gap-x-1">
            <button
              :disabled="paginateOrders.page === 1"
              type="button"
              class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
              @click="handleFilterChange(paginateOrders.page - 1)"
            >
              Previous
            </button>
            <div class="flex items-center gap-x-1">
              <button
                v-for="page in paginateOrders.total_page"
                :key="page"
                :disabled="useOrder.loading"
                type="button"
                :class="[
                  'leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border',
                  page === paginateOrders.page
                    ? 'text-white border bg-blue-600 border-blue-600 hover:bg-blue-600 focus:outline-none focus:bg-blue-600'
                    : 'bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300'
                ]"
                aria-current="page"
                @click="handleFilterChange(page)"
              >
                {{ page }}
              </button>
            </div>
            <button
              :disabled="paginateOrders.page === paginateOrders.total_page"
              type="button"
              class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
              @click="handleFilterChange(paginateOrders.page + 1)"
            >
              Next
            </button>
          </nav>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { useOrderStore } from '~/stores/orders'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { watchDebounced } from '@vueuse/core'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const useOrder = useOrderStore()
const orderDatas = ref([])
const paginateOrders = ref({})
const searchFilter = ref('')
const statusFilter = ref('')
const { showError } = useErrorModal()

const loading = ref(true)

const getOrderStatusClass = (productStatus) => {
  const statusClass = {
    PENDING: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    CONFIRMED: 'bg-blue-100 text-blue-800 border-blue-200',
    SENDING: 'bg-purple-100 text-purple-800 border-purple-200',
    DONE: 'bg-blue-100 text-blue-800 border-blue-200',
    CANCEL: 'bg-red-100 text-red-800 border-red-200'
  }

  return (
    statusClass[String(productStatus).toUpperCase()] || 'bg-gray-100 text-gray-800 border-gray-200'
  )
}

onMounted(async () => {
  try {
    await useOrder.fetchOrdersAdmin()
    syncOrdersFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  useOrder.cancelRequests()
})

const syncOrdersFromStore = () => {
  orderDatas.value = useOrder.orders
  paginateOrders.value = useOrder.pagination
}

const handleFilterChange = async (page) => {
  loading.value = true

  try {
    await useOrder.fetchOrdersAdmin({
      status: statusFilter.value,
      page: page,
      search: searchFilter.value
    })
    syncOrdersFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

watchDebounced(
  searchFilter,
  () => {
    handleFilterChange()
  },
  { debounce: 500, maxWait: 1000 }
  // maxWait (opsional): Memaksa API dipanggil setiap 1 detik jika user mengetik tanpa henti
)

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>

<style></style>
