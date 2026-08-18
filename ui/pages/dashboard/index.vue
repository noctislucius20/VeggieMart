<template>
  <section class="container">
    <!-- row -->
    <div class="grid grid-cols-12 mb-8">
      <div class="col-span-12">
        <!-- card -->
        <div
          class="card bg-gray-100 border-0 rounded-2xl"
          style="
            background-image: url(/images/slider/slider-image-1.jpg);
            background-repeat: no-repeat;
            background-size: cover;
            background-position: right;
          "
        >
          <div class="card-body lg:p-12">
            <h1>Welcome back! DailyMart</h1>
            <p class="mb-6">DailyMart is simple & clean design for developer and designer.</p>
            <NuxtLink
              to="/dashboard/products/create"
              class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
            >
              Create Product
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- summary cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <EarningsCard
        :earnings="summaryData.earnings"
        :percentage-change="summaryData.earningsChange"
        period-label="This Month"
      />
      <OrdersSummaryCard
        :total-orders="summaryData.totalOrders"
        :new-orders="summaryData.newOrders"
        :percentage-change="summaryData.ordersChange"
      />
      <CustomersSummaryCard
        :total-customers="summaryData.totalCustomers"
        :new-customers="summaryData.newCustomers"
        :percentage-change="summaryData.customersChange"
      />
    </div>

    <!-- recent orders table -->
    <div class="grid grid-cols-1 mb-8">
      <RecentOrdersTable :orders="recentOrders" :loading="loading" />
    </div>
  </section>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useOrderStore } from '~/stores/orders'
import { useCustomerStore } from '~/stores/customer'
import EarningsCard from '~/components/dashboard/EarningsCard.vue'
import OrdersSummaryCard from '~/components/dashboard/OrdersSummaryCard.vue'
import CustomersSummaryCard from '~/components/dashboard/CustomersSummaryCard.vue'
import RecentOrdersTable from '~/components/dashboard/RecentOrdersTable.vue'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const orderStore = useOrderStore()
const customerStore = useCustomerStore()

const loading = ref(true)
const recentOrders = ref([])

const { showError } = useErrorModal()

const summaryData = reactive({
  earnings: 0,
  earningsChange: null,
  totalOrders: 0,
  newOrders: null,
  ordersChange: null,
  totalCustomers: 0,
  newCustomers: null,
  customersChange: null
})

const calculateSummary = () => {
  const orders = orderStore.orders || []
  const customers = customerStore.customers || []

  // Total orders
  summaryData.totalOrders = orderStore.pagination?.total_count || orders.length

  // Earnings: sum total_amount dari orders with status DONE or CONFIRMED
  const totalEarnings = orders
    .filter((o) => o.status === 'DONE' || o.status === 'CONFIRMED')
    .reduce((sum, o) => sum + Number(o.total_amount || 0), 0)
  summaryData.earnings = totalEarnings

  // New orders (PENDING in last 7 days)
  const oneWeekAgo = new Date()
  oneWeekAgo.setDate(oneWeekAgo.getDate() - 7)
  const newOrdersCount = orders.filter((o) => {
    const orderDate = new Date(o.order_datetime || o.created_at)
    return orderDate >= oneWeekAgo
  }).length
  summaryData.newOrders = newOrdersCount

  // Total customers
  summaryData.totalCustomers = customerStore.pagination?.total_count || customers.length

  // New customers in the last 7 days
  const newCustomersCount = customers.filter((c) => {
    const createdDate = new Date(c.created_at || c.registered_at)
    return createdDate >= oneWeekAgo
  }).length
  summaryData.newCustomers = newCustomersCount

  // Recent orders (latest 5)
  recentOrders.value = [...orders].slice(0, 5)
}

onMounted(async () => {
  try {
    await Promise.all([
      orderStore.fetchOrdersAdmin({ limit: 100 }),
      customerStore.fetchCustomers({ limit: 100 })
    ])
    calculateSummary()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  orderStore.cancelRequests()
  customerStore.cancelRequests()
})
</script>
