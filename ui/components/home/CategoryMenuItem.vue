<template>
  <li class="border-b w-full">
    <a
      :class="{ active: isActive }"
      href="#"
      class="flex justify-between py-2 items-center text-gray-800 hover:text-blue-600"
      :aria-expanded="isOpen"
      :aria-controls="category.slug"
      @click.prevent="toggle"
    >
      {{ category.name }}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="icon icon-tabler icon-tabler-chevron-right inline-block transition-transform duration-200"
        :class="{ 'rotate-90': isOpen }"
        width="14"
        height="14"
        viewBox="0 0 24 24"
        stroke-width="2"
        stroke="currentColor"
        fill="none"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
        <path d="M9 6l6 6l-6 6" />
      </svg>
    </a>
    <div v-show="isOpen" :id="category.slug" class="text-gray-800">
      <ul class="flex flex-wrap flex-col ml-3">
        <li v-for="(item, index) in category.childs" :key="index">
          <a
            href="#!"
            class="inline-block py-2 no-underline hover:text-blue-600"
            :class="{ active: item.slug === activeChildSlug }"
            @click.prevent="$emit('select', item.slug)"
          >
            {{ item.name }}
          </a>
        </li>
      </ul>
    </div>
  </li>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  category: {
    type: Object,
    required: true
  },
  isActive: {
    type: Boolean,
    default: false
  },
  activeSlug: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['select'])

const isOpen = ref(false)

// Otomatis buka menu jika kategori (atau salah satu child) sedang aktif
watch(
  () => props.isActive,
  (active) => {
    if (active) {
      isOpen.value = true
    }
  },
  { immediate: true }
)

const toggle = () => {
  const wasOpen = isOpen.value
  isOpen.value = !wasOpen

  if (props.isActive) {
    // Kategori sedang dipilih → klik untuk menghapus filter (kembali ke All Products)
    emit('select', '')
  } else if (!wasOpen) {
    // Buka submenu → filter ke kategori ini
    emit('select', props.category.slug)
  }
}

const activeChildSlug = computed(() => {
  if (!props.activeSlug || !props.category.childs) {
    return null
  }

  const activeChild = props.category.childs.find((child) => child.slug === props.activeSlug)
  return activeChild?.slug
})
</script>
