<template>
  <div class="card h-full card-lg">
    <div class="card-body p-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <p class="text-sm text-gray-500 font-medium mb-1">Earnings</p>
          <h3 class="text-2xl font-bold text-gray-800">Rp {{ formatPrice(earnings) }}</h3>
          <p class="text-xs text-gray-400 mt-1">{{ periodLabel }}</p>
        </div>
        <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
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
            class="text-green-600"
          >
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path
              d="M16.7 8a3 3 0 0 0 -2.7 -2h-4a3 3 0 0 0 0 6h4a3 3 0 0 1 0 6h-4a3 3 0 0 1 -2.7 -2"
            />
            <path d="M12 3v3m0 12v3" />
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
  earnings: {
    type: Number,
    default: 0
  },
  percentageChange: {
    type: Number,
    default: null
  },
  periodLabel: {
    type: String,
    default: 'This Month'
  }
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>
