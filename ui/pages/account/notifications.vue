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
            <a class="nav-link gap-2" href="/account">
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
            <a class="nav-link gap-2 active" href="/account/notifications">
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
                :disabled="notificationStore.loading"
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
                  <a class="nav-link gap-2" href="./">
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
                  <a class="nav-link gap-2 active" href="./notifications">
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
                  <a class="nav-link gap-2" href="./setting">
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
              <h2 class="mb-6 text-lg">Notifications</h2>

              <!-- Loading Skeleton -->
              <div v-if="loading">
                <div
                  v-for="i in 5"
                  :key="'skeleton-notif-' + i"
                  class="px-5 py-4 border-b border-gray-300"
                >
                  <div class="flex">
                    <div class="skeleton w-10 h-10 rounded-full" />
                    <div class="ms-4 flex-1">
                      <div class="skeleton skeleton-text w-48 mb-2" />
                      <div class="skeleton skeleton-text w-full mb-2" />
                      <div class="skeleton skeleton-text w-32" />
                    </div>
                  </div>
                </div>
              </div>
              <!-- Notifications List -->
              <div v-else-if="notificationStore?.loading" class="text-center py-4">Loading...</div>
              <div
                v-else-if="notificationStore?.notifications?.length === 0"
                class="text-center py-4"
              >
                Tidak ada notifikasi
              </div>
              <ul v-else class="notification-list">
                <li
                  v-for="notification in notificationList"
                  :key="notification.id"
                  class="px-5 py-4 hover:bg-gray-100 border-b border-gray-300 cursor-pointer"
                  :class="{ 'bg-blue-50': !notification.is_read, 'bg-white': notification.is_read }"
                  @click="handleNotificationClick(notification)"
                >
                  <div class="flex">
                    <div class="">
                      <div
                        class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-600"
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
                          class="icon icon-tabler icons-tabler-outline icon-tabler-bell"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path
                            d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3l-1 4h-10l-1 -4v-3a7 7 0 0 1 4 -6z"
                          />
                        </svg>
                      </div>
                    </div>

                    <div class="ms-4">
                      <div class="mb-1 text-left">
                        <p class="text-gray-900 font-semibold">{{ notification.subject }}</p>
                        {{ notification.message }}
                      </div>
                      <div class="flex items-center">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="12"
                          height="12"
                          fill="currentColor"
                          class="bi bi-clock text-gray-500"
                          viewBox="0 0 16 16"
                        >
                          <path
                            d="M8 3.5a.5.5 0 0 0-1 0V9a.5.5 0 0 0 .252.434l3.5 2a.5.5 0 0 0 .496-.868L8 8.71V3.5z"
                          />
                          <path
                            d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16zm7-8A7 7 0 1 1 1 8a7 7 0 0 1 14 0z"
                          />
                        </svg>
                        <small class="ms-2">{{ formattedTime(notification.sent_at) }}</small>
                      </div>
                    </div>
                  </div>
                </li>
              </ul>

              <div
                v-if="paginateNotifications?.total_count > 0"
                class="border-t border-gray-300 flex flex-col md:flex-row justify-between items-center px-6 py-6 gap-3"
              >
                <span
                  >Showing {{ paginateNotifications.page }} to
                  {{ paginateNotifications.total_pages }} of
                  {{ paginateNotifications.total_count }} entries</span
                >
                <nav class="flex items-center gap-x-1">
                  <button
                    :disabled="paginateNotifications.page === 1"
                    type="button"
                    class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
                    @click="handleFilterChange(paginateNotifications.page - 1)"
                  >
                    Previous
                  </button>
                  <div class="flex items-center gap-x-1">
                    <button
                      v-for="page in paginateNotifications.total_pages"
                      :key="page"
                      :disabled="notificationStore.loading"
                      type="button"
                      :class="[
                        'leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border',
                        page === paginateNotifications.page
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
                    :disabled="paginateNotifications.page === paginateNotifications.total_page"
                    type="button"
                    class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
                    @click="handleFilterChange(paginateNotifications.page + 1)"
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
import { useNotificationStore } from '~/stores/notification'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

definePageMeta({
  middleware: ['auth']
})
const router = useRouter()
const authStore = useAuthStore()
const orderStore = useOrderStore()
const notificationStore = useNotificationStore()
const notificationList = ref([])
const paginateNotifications = ref({})
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

const syncNotificationsFromStore = () => {
  notificationList.value = notificationStore.notifications
  paginateNotifications.value = notificationStore.pagination
}

const formattedTime = (sentAt) => {
  if (!sentAt) return ''
  return dayjs(sentAt).fromNow()
}

const handleFilterChange = async (page) => {
  try {
    await notificationStore.fetchPushNotifications({ page: page, limit: 5 })
    syncNotificationsFromStore()
  } catch (error) {
    showError(error)
  }
}

const handleNotificationClick = async (notification) => {
  try {
    if (!notification.is_read) {
      await notificationStore.markAsRead(notification.id)
      syncNotificationsFromStore()
    }
  } catch (error) {
    showError(error)
  }
}

onMounted(async () => {
  try {
    let mounted = true
    await notificationStore.fetchPushNotifications({ page: 1, limit: 5 })
    if (mounted) syncNotificationsFromStore()
    return () => {
      mounted = false
    }
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  notificationStore.cancelRequests()
})
</script>
