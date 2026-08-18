<template>
  <main>
    <!-- Offcanvas for mobile -->
    <div
      id="offcanvasAccount"
      class="offcanvas offcanvas-start"
      tabindex="-1"
      aria-labelledby="offcanvasAccountLabel"
    >
      <div class="offcanvas-header border-bottom border-gray-300">
        <div class="flex items-center gap-2 w-full">
          <h5 id="offcanvasAccountLabel" class="offcanvas-title mb-0">Menu</h5>
          <button
            type="button"
            class="btn btn-sm ms-auto p-0 border-0 shadow-none text-gray-600 hover:text-gray-900"
            data-bs-dismiss="offcanvas"
            aria-label="Close"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="icon icon-tabler icons-tabler-outline icon-tabler-x"
            >
              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
              <path d="M18 6l-12 12" />
              <path d="M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      <div class="offcanvas-body p-0">
        <ul class="nav flex-col nav-pills nav-pills-dark p-4">
          <!-- nav item -->
          <li class="nav-item">
            <a class="nav-link gap-2 active" href="/account">
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
              Your Orders
            </a>
          </li>

          <!-- nav item -->
          <li class="nav-item">
            <a class="nav-link gap-2" href="/account/notifications">
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
                class="icon icon-tabler icons-tabler-outline icon-tabler-bell"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3l-1 4h-10l-1 -4v-3a7 7 0 0 1 4 -6z" />
              </svg>
              Notifications
            </a>
          </li>

          <!-- nav item -->
          <li class="nav-item">
            <a class="nav-link gap-2" href="/account/setting">
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
                class="icon icon-tabler icons-tabler-outline icon-tabler-settings"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path
                  d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"
                />
                <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
              </svg>
              Settings
            </a>
          </li>
          <!-- nav item -->

          <li class="nav-item">
            <hr class="pt-3 mt-5" />
          </li>
          <!-- nav item -->
          <li class="nav-item">
            <button
              :disabled="authStore.loading"
              type="button"
              class="nav-link gap-2"
              @click="handleLogout"
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
                class="icon icon-tabler icons-tabler-outline icon-tabler-logout"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path
                  d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2"
                />
                <path d="M9 12h12l-3 -3" />
                <path d="M18 15l3 -3" />
              </svg>
              Log out
            </button>
          </li>
        </ul>
      </div>
    </div>
    <!-- section -->
    <section>
      <div class="container">
        <div class="grid grid-cols-12">
          <div class="col-span-12">
            <div class="flex justify-between items-center md:hidden py-4">
              <!-- heading -->
              <h3 class="text-md mb-0">Account Setting</h3>
              <!-- button -->
              <button
                :disabled="orderStore.loading"
                class="btn inline-flex items-center gap-x-2 bg-transparent text-gray-600 border-gray-300 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300 md:hidden btn-sm ms-3"
                type="button"
                data-bs-toggle="offcanvas"
                data-bs-target="#offcanvasAccount"
                aria-controls="offcanvasAccount"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="icon icon-tabler icons-tabler-outline icon-tabler-indent-increase"
                >
                  <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                  <path d="M20 6l-11 0" />
                  <path d="M20 12l-7 0" />
                  <path d="M20 18l-11 0" />
                  <path d="M4 8l4 4l-4 4" />
                </svg>
              </button>
            </div>
          </div>
          <div class="lg:col-span-2 md:col-span-4 col-span-12 border-r hidden md:block">
            <div class="pt-10 md:pr-10">
              <!-- nav -->
              <ul class="nav flex-col nav-pills nav-pills-dark">
                <!-- nav item -->
                <li class="nav-item">
                  <a class="nav-link active gap-2" href="/account">
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
                    Your Orders
                  </a>
                </li>
                <!-- nav item -->
                <li class="nav-item">
                  <a class="nav-link gap-2" href="/account/notifications">
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
                      class="icon icon-tabler icons-tabler-outline icon-tabler-bell"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3l-1 4h-10l-1 -4v-3a7 7 0 0 1 4 -6z"
                      />
                    </svg>
                    Notifications
                  </a>
                </li>
                <!-- nav item -->
                <li class="nav-item">
                  <a class="nav-link gap-2" href="/account/setting">
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
                      class="icon icon-tabler icons-tabler-outline icon-tabler-settings"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"
                      />
                      <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
                    </svg>
                    Settings
                  </a>
                </li>
                <!-- nav item -->
                <li class="nav-item">
                  <hr class="pt-3 mt-3" />
                </li>
                <!-- nav item -->
                <li class="nav-item">
                  <button
                    :disabled="orderStore.loading"
                    type="button"
                    class="nav-link gap-2"
                    @click="handleLogout"
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
                      class="icon icon-tabler icons-tabler-outline icon-tabler-logout"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2"
                      />
                      <path d="M9 12h12l-3 -3" />
                      <path d="M18 15l3 -3" />
                    </svg>
                    Log out
                  </button>
                </li>
              </ul>
            </div>
          </div>
          <div class="lg:col-span-10 md:col-span-8 col-span-12">
            <div class="py-6 md:p-6 lg:p-10">
              <!-- heading -->
              <h2 class="mb-6 text-lg">Your Orders</h2>

              <div class="block w-full overflow-auto scrolling-touch">
                <!-- Table -->
                <table
                  class="w-full max-w-full bg-transparent whitespace-nowrap table-with-checkbox"
                >
                  <!-- Table Head -->
                  <thead class="bg-gray-100">
                    <tr class="border-gray-300 border-b">
                      <th class="py-3 pl-6 text-left">&nbsp;</th>
                      <th class="py-3 pl-6 text-left">Product</th>
                      <th class="py-3 pl-6 text-left">Order</th>
                      <th class="py-3 pl-6 text-left">Date</th>
                      <th class="py-3 pl-6 text-left">Items</th>
                      <th class="py-3 pl-6 text-left">Status</th>
                      <th class="py-3 pl-6 text-left">Amount</th>

                      <th />
                    </tr>
                  </thead>
                  <tbody>
                    <!-- Loading Skeleton -->
                    <tr v-if="loading">
                      <td colspan="8">
                        <div class="py-4">
                          <div
                            v-for="i in 5"
                            :key="'skeleton-order-' + i"
                            class="flex items-center gap-4 py-3 border-b border-gray-300 px-6"
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
                    <!-- Table body -->
                    <tr v-else-if="orderStore?.orders?.length === 0">
                      <td colspan="8" class="text-center py-4">Tidak ada order</td>
                    </tr>
                    <tr
                      v-for="order in orderList"
                      v-else
                      :key="order.id"
                      class="border-gray-300 border-b"
                    >
                      <td class="align-middle px-6 py-3">
                        <a :href="`/account/orders/detail/${order.id}`" class="">
                          <img
                            :src="order.product_image"
                            :alt="order.order_code"
                            class="w-16 h-16 max-w-screen-lg"
                          />
                        </a>
                      </td>
                      <td class="align-middle px-6 py-3">
                        <a
                          :href="`/account/orders/detail/${order.id}`"
                          class="font-semibold text-inherit"
                        >
                          <h6 class="mb-0">{{ order.product_name }}</h6>
                        </a>
                        <span
                          ><small class="text-gray-500"
                            >{{ order.weight }} {{ order.unit }}</small
                          ></span
                        >
                      </td>
                      <td class="align-middle px-6 py-3">
                        <a :href="`/account/orders/detail/${order.id}`" class="text-inherit">{{
                          order.order_code
                        }}</a>
                      </td>
                      <td class="align-middle px-6 py-3">{{ order.order_datetime }}</td>
                      <td class="align-middle px-6 py-3">{{ order.quantity }}</td>
                      <td class="align-middle px-6 py-3">
                        <span :class="getStatusClass(order.status)">{{ order.status }}</span>
                      </td>
                      <td class="align-middle px-6 py-3">Rp {{ order.total_amount }}</td>
                      <td class="text-gray-500 align-middle px-6 py-3">
                        <a
                          :href="`/account/orders/detail/${order.id}`"
                          class="text-inherit"
                          data-bs-toggle="tooltip"
                          data-bs-placement="top"
                          data-bs-title="View"
                          ><i class="feather-icon icon-eye"
                        /></a>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div
                v-if="paginateOrders?.total_count > 0"
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
                      :disabled="orderStore.loading"
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
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'
import { onMounted, onBeforeUnmount } from 'vue'
import { useOrderStore } from '~/stores/orders'

definePageMeta({
  middleware: ['auth']
})
const router = useRouter()
const authStore = useAuthStore()
const orderStore = useOrderStore()
const orderList = ref([])
const paginateOrders = ref({})
const loading = ref(true)
const { showError } = useErrorModal()

const handleLogout = () => {
  try {
    authStore.logout()

    useAuthStore().$reset()
    useCartStore().$reset()
    useCategoryStore().$reset()
    useCustomerStore().$reset()
    useNotificationStore().$reset()
    useOrderStore().$reset()
    usePaymentStore().$reset()
    useProductStore().$reset()
    useRoleStore().$reset()

    router.push('/auth/signin')
  } catch (error) {
    showError(error)
  }
}

const getStatusClass = (status) => {
  const classes = {
    pending:
      'inline-block p-1 text-sm align-baseline leading-none rounded bg-gray-100 text-gray-800 border border-gray-200 font-semibold',
    confirmed:
      'inline-block p-1 text-sm align-baseline leading-none rounded bg-blue-100 text-blue-800 border border-blue-200 font-semibold',
    processing:
      'inline-block p-1 text-sm align-baseline leading-none rounded bg-yellow-100 text-yellow-800 border border-yellow-200 font-semibold',
    sending:
      'inline-block p-1 text-sm align-baseline leading-none rounded bg-purple-100 text-purple-800 border border-purple-200 font-semibold',
    done: 'inline-block p-1 text-sm align-baseline leading-none rounded bg-green-100 text-green-800 border border-green-200 font-semibold',
    canceled:
      'inline-block p-1 text-sm align-baseline leading-none rounded bg-red-100 text-red-800 border border-red-200 font-semibold'
  }
  return classes[status.toLowerCase()] || classes['processing']
}

const syncOrdersFromStore = () => {
  orderList.value = orderStore.orders
  paginateOrders.value = orderStore.pagination
}

const handleFilterChange = async (page) => {
  loading.value = true

  try {
    await orderStore.fetchOrders({ page: page })
    syncOrdersFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    await orderStore.fetchOrders()
    syncOrdersFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  orderStore.cancelRequests()
})
</script>
