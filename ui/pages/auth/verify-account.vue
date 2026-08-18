<template>
  <section class="my-10">
    <div class="container">
      <div class="flex flex-wrap justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
          <!-- Alert Error -->
          <div
            v-if="authStore.error"
            class="mb-4 p-4 text-sm text-red-800 rounded-lg bg-red-50 relative"
            role="alert"
          >
            <div class="flex items-center pr-8">
              <span class="font-medium mr-2">Gagal!</span> {{ authStore.error }}
            </div>
          </div>

          <!-- Loading State -->
          <div v-if="authStore.loading" class="text-center">
            <p>Memverifikasi akun anda...</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
<script setup>
import { onBeforeUnmount } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useRouter, useRoute } from 'vue-router'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { showError } = useErrorModal()

const token = route.query.token

if (!token) {
  router.push('/auth/signin')
}

onMounted(async () => {
  try {
    await authStore.verifyAccount(token)

    if (authStore.user?.role === 'Super Admin') {
      router.push('/dashboard')
    } else {
      router.push('/')
    }
  } catch (error) {
    showError(error)
  }
})

onBeforeUnmount(() => {
  authStore.cancelRequests()
})
</script>
