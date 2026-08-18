<template>
  <div class="card h-full card-lg">
    <div class="px-6 py-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <h5 class="text-lg font-semibold text-gray-800">Recent Orders</h5>
        <a
          href="/dashboard/orders"
          class="text-sm text-blue-600 hover:text-blue-700 hover:underline"
        >
          View All
        </a>
      </div>
    </div>
    <div class="card-body p-0">
      <div class="relative overflow-x-auto">
        <table class="text-left w-full whitespace-nowrap table-hover">
          <thead class="bg-gray-200 text-gray-700">
            <tr class="border-transparent border-b-0!">
              <th scope="col" class="px-6 py-3">Order Number</th>
              <th scope="col" class="px-6 py-3">Order Date</th>
              <th scope="col" class="px-6 py-3">Total Price</th>
              <th scope="col" class="px-6 py-3">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="5" class="py-4">
                <div
                  v-for="i in 5"
                  :key="'skeleton-' + i"
                  class="flex items-center gap-4 py-3 border-b border-gray-100"
                >
                  <div class="flex-1">
                    <div class="skeleton skeleton-text w-48 mb-2" />
                    <div class="skeleton skeleton-text w-32" />
                  </div>
                </div>
              </td>
            </tr>
            <tr v-else-if="orders.length === 0">
              <td colspan="5" class="text-center py-8 text-gray-500">No recent orders</td>
            </tr>
            <tr
              v-for="item in orders"
              v-else
              :key="item.id"
              class="border-b border-gray-100 hover:bg-gray-50"
            >
              <td class="py-3 px-6 text-left">
                <a
                  :href="`/dashboard/orders/detail/${item.id}`"
                  class="text-blue-600 hover:text-blue-700"
                >
                  {{ item.order_code || '#' + item.id }}
                </a>
              </td>
              <td class="py-3 px-6 text-left text-gray-600">{{ item.order_datetime }}</td>
              <td class="py-3 px-6 text-left text-gray-800 font-medium">
                Rp {{ formatPrice(item.total_amount) }}
              </td>
              <td class="py-3 px-6 text-left">
                <span
                  class="inline-block p-1 text-sm align-baseline leading-none rounded font-semibold"
                  :class="getOrderStatusClass(item.status)"
                >
                  {{ item.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  orders: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

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

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>
