<template>
  <li
    :class="[
      'px-5 py-4 hover:bg-gray-100 border-b border-gray-300',
      { 'bg-gray-100 active': notification.is_read === false }
    ]"
  >
    <a href="javascript:void(0)" class="text-gray-500" @click.prevent="handleClick">
      <div class="flex">
        <div class="">
          <div
            class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-600"
          >
            <Icon name="tabler:bell" size="20" />
          </div>
        </div>

        <div class="ms-4">
          <div class="mb-1 text-left">
            <p class="text-gray-900 font-semibold">{{ notification.subject }}</p>
            {{ truncatedMessage }}
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
              <path d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16zm7-8A7 7 0 1 1 1 8a7 7 0 0 1 14 0z" />
            </svg>
            <small class="ms-2">{{ formattedTime }}</small>
          </div>
        </div>
      </div>
    </a>
  </li>
</template>

<script setup>
import { computed } from 'vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { useNotificationStore } from '~/stores/notification'
import { useRouter } from 'vue-router'

dayjs.extend(relativeTime)

const props = defineProps({
  notification: {
    type: Object,
    default: () => ({
      id: '',
      subject: '',
      message: '',
      sent_at: '',
      is_read: false
    })
  }
})

const router = useRouter()
const notificationStore = useNotificationStore()
const authStore = useAuthStore()

const truncatedMessage = computed(() => {
  const msg = props.notification.message || ''
  if (msg.length <= 60) return msg
  return msg.substring(0, 60) + '...'
})

const formattedTime = computed(() => {
  if (!props.notification.sent_at) return ''
  return dayjs(props.notification.sent_at).fromNow()
})

const emit = defineEmits(['close'])

const handleClick = async () => {
  if (props.notification.notification_type === 'orders') {
    if (authStore.userRole.toLowerCase() === 'super admin') {
      router.push(`/dashboard/orders/detail/${props.notification.notification_type_id}`)
    } else {
      router.push(`/account/orders/detail/${props.notification.notification_type_id}`)
    }
  }

  emit('close')

  if (!props.notification.is_read) {
    await notificationStore.markAsRead(props.notification.id)
  }
}
</script>
