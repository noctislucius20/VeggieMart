<template>
  <div class="card h-full card-lg">
    <div class="card-body p-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <p class="text-sm text-gray-500 font-medium mb-1">Total Customers</p>
          <h3 class="text-2xl font-bold text-gray-800">{{ totalCustomers }}</h3>
          <p v-if="newCustomers !== null" class="text-xs text-gray-400 mt-1">
            {{ newCustomers }} new in this week
          </p>
        </div>
        <div class="w-12 h-12 rounded-full bg-yellow-100 flex items-center justify-center">
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
            class="text-yellow-600"
          >
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M9 7m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0" />
            <path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" />
            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            <path d="M21 21v-2a4 4 0 0 0 -3 -3.85" />
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
  totalCustomers: {
    type: Number,
    default: 0
  },
  newCustomers: {
    type: Number,
    default: null
  },
  percentageChange: {
    type: Number,
    default: null
  }
})
</script>
