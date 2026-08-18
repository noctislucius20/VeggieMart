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
                  <hr class="pt-3 mt-5" />
                </li>
                <!-- nav item -->
                <li class="nav-item">
                  <button type="button" class="nav-link gap-2" @click="handleLogout">
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
              <!-- card -->
              <div class="card h-full card-lg">
                <!-- Loading Skeleton -->
                <div v-if="loading" class="card-body p-6">
                  <div class="flex flex-col md:flex-row gap-4 justify-between">
                    <div class="flex items-center">
                      <div class="skeleton skeleton-text w-40" style="height: 24px" />
                      <div class="skeleton skeleton-button" style="width: 100px; height: 28px" />
                    </div>
                    <div class="skeleton skeleton-button" style="width: 150px; height: 40px" />
                  </div>
                  <div class="mt-8">
                    <div class="grid grid-cols-12">
                      <div class="lg:col-span-4 md:col-span-4 col-span-12">
                        <div class="skeleton skeleton-text w-32 mb-4" style="height: 20px" />
                        <div class="skeleton skeleton-text w-full mb-2" />
                        <div class="skeleton skeleton-text w-3/4 mb-2" />
                        <div class="skeleton skeleton-text w-1/2" />
                      </div>
                      <div class="lg:col-span-4 md:col-span-4 col-span-12">
                        <div class="skeleton skeleton-text w-32 mb-4" style="height: 20px" />
                        <div class="skeleton skeleton-text w-full mb-2" />
                        <div class="skeleton skeleton-text w-3/4" />
                      </div>
                      <div class="lg:col-span-4 md:col-span-4 col-span-12">
                        <div class="skeleton skeleton-text w-32 mb-4" style="height: 20px" />
                        <div class="skeleton skeleton-text w-full mb-2" />
                        <div class="skeleton skeleton-text w-3/4 mb-2" />
                        <div class="skeleton skeleton-text w-1/2" />
                      </div>
                    </div>
                  </div>
                </div>
                <!-- Order Content -->
                <div v-else class="card-body p-6">
                  <div class="flex flex-col md:flex-row gap-4 justify-between">
                    <div class="flex items-center">
                      <h3 class="mb-0 mr-3">Order ID: #{{ orderData.order_code }}</h3>
                      <span
                        class="inline-block p-1 text-sm align-baseline leading-none rounded border font-semibold"
                        :class="getOrderStatusClass(orderData.status)"
                      >
                        {{ orderData.status }}
                      </span>
                    </div>
                    <div class="mt-4 md:mt-0">
                      <a
                        href="/account"
                        class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                      >
                        Back to orders
                      </a>
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
                              {{ orderData.order_code }}
                            </span>
                            <br />
                            Order Date:
                            <span class="text-gray-800">
                              {{ orderData.order_datetime }}
                            </span>
                            <br />
                            Order Total:
                            <span class="text-gray-800"
                              >Rp. {{ formatPrice(orderData.total_amount) }}</span
                            >
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <!-- Table Skeleton -->
                <div v-if="loading" class="card-body">
                  <div class="grid grid-cols-1">
                    <div class="relative overflow-x-auto">
                      <table class="text-left w-full whitespace-nowrap">
                        <thead class="bg-gray-200 text-gray-700">
                          <tr>
                            <th scope="col" class="px-6 py-3">
                              <div class="skeleton skeleton-text w-20" style="height: 16px" />
                            </th>
                            <th scope="col" class="px-6 py-3">
                              <div class="skeleton skeleton-text w-16" style="height: 16px" />
                            </th>
                            <th scope="col" class="px-6 py-3">
                              <div class="skeleton skeleton-text w-20" style="height: 16px" />
                            </th>
                            <th scope="col" class="px-6 py-3">
                              <div class="skeleton skeleton-text w-16" style="height: 16px" />
                            </th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="i in 3"
                            :key="'skeleton-row-' + i"
                            class="border-b border-gray-300"
                          >
                            <td class="py-3 px-6 text-left">
                              <div class="flex items-center gap-2">
                                <div class="skeleton w-12 h-12 rounded" />
                                <div class="skeleton skeleton-text w-32" />
                              </div>
                            </td>
                            <td class="py-3 px-6 text-left">
                              <div class="skeleton skeleton-text w-24" />
                            </td>
                            <td class="py-3 px-6 text-left">
                              <div class="skeleton skeleton-text w-16" />
                            </td>
                            <td class="py-3 px-6 text-left">
                              <div class="skeleton skeleton-text w-24" />
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
                <!-- Order Table -->
                <div v-else class="grid grid-cols-1 px-5">
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
                            Sub Total :
                          </td>
                          <td class="font-medium text-gray-800 border-b border-gray-300">
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
                            Shipping Cost
                          </td>
                          <td class="font-medium text-gray-800 border-b border-gray-300">
                            Rp. {{ formatPrice(orderData.shipping_fee) }}
                          </td>
                        </tr>
                        <tr class="">
                          <td class="py-3 px-6 text-left" />
                          <td class="py-3 px-6 text-left" />
                          <td
                            colspan="1"
                            class="font-semibold text-gray-800 py-3 px-6 text-left border-b border-gray-300"
                          >
                            Grand Total
                          </td>
                          <td class="font-semibold text-gray-800 py-3 border-b border-gray-300">
                            Rp. {{ formatPrice(orderData.shipping_fee + calculateSubtotal) }}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Payment & Notes Skeleton -->
                <div v-if="loading" class="card-body">
                  <div class="grid md:grid-cols-2 gap-4">
                    <div>
                      <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-text w-40" />
                    </div>
                    <div>
                      <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-text w-full" />
                      <div class="skeleton skeleton-text w-3/4" />
                    </div>
                  </div>
                  <div class="grid md:grid-cols-2 gap-4 mt-6">
                    <div>
                      <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-button" style="width: 120px; height: 32px" />
                    </div>
                  </div>
                </div>
                <!-- Payment & Notes Content -->
                <div v-else class="card-body">
                  <div class="grid md:grid-cols-2 gap-4">
                    <div class="">
                      <h6 class="mb-2">Payment Info</h6>
                      <span>{{
                        paymentData
                          ? paymentData?.payment_method?.toLowerCase() === 'transfer'
                            ? 'Transfer'
                            : 'Cash or Duel'
                          : ''
                      }}</span>
                    </div>
                    <div class="">
                      <h6 class="mb-3">Delivery Notes</h6>
                      <p class="text-gray-700">{{ orderData.remarks }}</p>
                    </div>
                  </div>
                  <div class="grid md:grid-cols-2 gap-4">
                    <div class="">
                      <h6 class="mb-2">Payment Status</h6>
                      <span
                        class="inline-block p-1 text-sm align-baseline leading-none rounded border font-semibold"
                        :class="getPaymentStatusClass(paymentData?.payment_status)"
                      >
                        {{ getPaymentStatusLabel(paymentData?.payment_status) }}
                      </span>
                      <div
                        v-if="paymentData?.payment_status?.toLowerCase() === 'pending'"
                        class="mt-3"
                      >
                        <button
                          class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                          @click="handlePayNow"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            class="icon icon-tabler icons-tabler-outline icon-tabler-external-link"
                          >
                            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                            <path
                              d="M12 6h-6a2 2 0 0 0 -2 2v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-6"
                            />
                            <path d="M11 13l9 -9" />
                            <path d="M15 4h5v5" />
                          </svg>
                          Bayar Sekarang
                        </button>
                      </div>
                    </div>
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
import { useAuthStore } from '~/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useOrderStore } from '~/stores/orders'

definePageMeta({
  middleware: ['auth']
})

const config = useRuntimeConfig()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const orderStore = useOrderStore()
const paymentStore = usePaymentStore()
const orderData = ref({})
const paymentData = ref(null)
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

const getOrderStatusClass = (productStatus) => {
  const statusClass = {
    PENDING: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    CONFIRMED: 'bg-blue-100 text-blue-800 border-blue-200',
    PROCESS: 'bg-orange-100 text-orange-800 border-orange-200',
    SENDING: 'bg-purple-100 text-purple-800 border-purple-200',
    DONE: 'bg-blue-100 text-blue-800 border-blue-200',
    CANCEL: 'bg-red-100 text-red-800 border-red-200'
  }

  return (
    statusClass[String(productStatus).toUpperCase()] || 'bg-gray-100 text-gray-800 border-gray-200'
  )
}

const getPaymentStatusClass = (status) => {
  const statusClass = {
    pending: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    success: 'bg-blue-100 text-blue-800 border-blue-200',
    failed: 'bg-red-100 text-red-800 border-red-200'
  }

  return statusClass[String(status).toLowerCase()] || 'bg-gray-100 text-gray-800 border-gray-200'
}

const getPaymentStatusLabel = (status) => {
  const labels = {
    pending: 'MENUNGGU PEMBAYARAN',
    success: 'LUNAS',
    failed: 'PEMBAYARAN GAGAL'
  }

  return labels[String(status).toLowerCase()] || status || ''
}

onMounted(async () => {
  try {
    await Promise.all([
      orderStore.getDetailOrders(route.params.id, false),
      paymentStore.getDetailPaymentByOrderId(route.params.id, false)
    ])
    orderData.value = orderStore.order
    paymentData.value = paymentStore.payment
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

const handlePayNow = () => {
  if (paymentData.value?.payment_token) {
    window.snap.pay(paymentData.value.payment_token, {
      onSuccess: function () {
        window.location.reload()
      },
      onPending: function () {
        // alert('Menunggu Pembayaran')
      },
      onError: function () {
        alert('Pembayaran gagal')
      },
      onClose: function () {
        alert('Anda menutup popup tanpa menyelesaikan pembayaran')
      }
    })
  }
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}

useScript({
  src: 'https://app.sandbox.midtrans.com/snap/snap.js',
  attrs: {
    'data-client-key': config.public.midtransClientKey
  },
  crossorigin: false
})
</script>
