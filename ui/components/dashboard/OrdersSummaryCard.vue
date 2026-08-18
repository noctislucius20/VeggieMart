<template>
  <div class="card h-full card-lg">
    <div class="card-body p-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <p class="text-sm text-gray-500 font-medium mb-1">Total Orders</p>
          <h3 class="text-2xl font-bold text-gray-800">{{ totalOrders }}</h3>
          <p v-if="newOrders !== null" class="text-xs text-gray-400 mt-1">
            {{ newOrders }} new orders
          </p>
        </div>
        <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
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
            class="text-blue-600"
          >
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M6.5 7h11" />
            <path d="M9.5 11h6" />
            <path d="M9.5 15h6" />
            <path
              d="M5 3h14a1 1 0 0 1 1 1v16a1 1 0 0 1 -1 1h-14a1 1 0 0 1 -1 -1v-16a1 1 0 0 1 1 -1z"
            />
            <circle cx="8" cy="7" r="0.5" fill="currentColor" />
            <circle cx="8" cy="11" r="0.5" fill="currentColor" />
            <circle cx="8" cy="15" r="0.5" fill="currentColor" />
          </svg>
        </div>
      </div>
      <div v-if="percentageChange !== null" class="flex items-center gap-1">
        <span
          :class="[
            'inline-flex items-center text-xs font-medium px-2 py-0.5 rounded',
            percentageChange >= 0 ? 'text-green-700 bg-green-100' : 'text-red-700 bg-red-100'
          ]"
        >
          <svg
            v-if="percentageChange >= 0"
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="me-1"
          >
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M12 5l0 14" />
            <path d="M16 9l-4 -4" />
            <path d="M8 9l4 -4" />
          </svg>
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="me-1"
          >
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M12 5l0 14" />
            <path d="M16 15l-4 4" />
            <path d="M8 15l4 4" />
          </svg>
          {{ Math.abs(percentageChange) }}%
        </span>
        <span class="text-xs text-gray-400">vs last month</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  totalOrders: {
    type: Number,
    default: 0
  },
  newOrders: {
    type: Number,
    default: null
  },
  percentageChange: {
    type: Number,
    default: null
  }
})
</script>
