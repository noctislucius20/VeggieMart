<template>
  <nav class="navbar navbar-expand-lg navbar-glass">
    <div class="container-fluid">
      <div class="flex justify-between items-center w-full">
        <div class="flex items-center">
          <a
            class="text-inherit mr-3"
            :class="{ 'hidden xl:block': sidebarOpen, 'block xl:hidden': !sidebarOpen }"
            @click="toggleSidebar"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              fill="currentColor"
              class="bi bi-text-indent-right"
              viewBox="0 0 16 16"
            >
              <path
                d="M2 3.5a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm10.646 2.146a.5.5 0 0 1 .708.708L11.707 8l1.647 1.646a.5.5 0 0 1-.708.708l-2-2a.5.5 0 0 1 0-.708l2-2zM2 6.5a.5.5 0 0 1 .5-.5h6a.5.5 0 0 1 0 1h-6a.5.5 0 0 1-.5-.5zm0 3a.5.5 0 0 1 .5-.5h6a.5.5 0 0 1 0 1h-6a.5.5 0 0 1-.5-.5zm0 3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z"
              />
            </svg>
          </a>
        </div>
        <div>
          <ul class="list-unstyled flex items-center gap-2">
            <NotificationDropdown />
            <ProfileDropdown />
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import NotificationDropdown from '../auth/NotificationDropdown.vue'
import ProfileDropdown from '../auth/ProfileDropdown.vue'

defineProps({
  sidebarOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['toggle-sidebar'])

const toggleSidebar = () => {
  emit('toggle-sidebar')
}
</script>
