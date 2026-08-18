<template>
  <section class="my-10">
    <div class="container">
      <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
        <div class="w-full md:w-1/3 xl:w-1/3 lg:order-1 order-2">
          <!-- img -->
          <img src="~/assets/images/svg-graphics/fp-g.svg" alt="" class="max-w-full h-auto" />
        </div>

        <div class="w-full md:w-1/2 lg:mx-1/6 xl:w-1/3 lg:order-2 order-1 flex flex-col gap-6">
          <div class="flex flex-col gap-1">
            <h1 class="text-xl">Forgot your password?</h1>
            <p>
              Please enter the email address associated with your account and We will email you a
              link to reset your password.
            </p>
          </div>

          <form class="needs-validation" novalidate @submit.prevent="handleSubmit">
            <!-- Alert Error -->
            <div
              v-if="authStore.error"
              class="mb-4 p-4 text-sm text-red-800 rounded-lg bg-red-50 relative"
              role="alert"
            >
              <button
                :disabled="authStore.loading"
                class="absolute top-2 right-0 p-1.5 inline-flex items-center justify-center text-red-800 bg-red-50 hover:bg-red-200 rounded-full"
                @click="authStore.error = null"
              >
                <span class="sr-only">Tutup</span>
                <i class="ti ti-x text-lg" />
              </button>
              <div class="flex items-center pr-8">
                <span class="font-medium mr-2">Gagal!</span> {{ authStore.error }}
              </div>
            </div>
            <!-- Alert Success -->

            <div class="flex flex-wrap gap-4">
              <div class="w-full">
                <!-- input -->
                <label for="formForgotEmail" class="invisible hidden">Email address</label>
                <input
                  id="formForgotEmail"
                  v-model="form.email"
                  :disabled="authStore.loading"
                  type="email"
                  class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  :class="{ 'border-red-500': v$.email.$error && v$.email.$dirty }"
                  placeholder="Email"
                  required
                  @blur="v$.email.$touch()"
                />
                <div v-if="v$.email.$error && v$.email.$dirty" class="text-red-600 text-sm mt-1">
                  <template v-if="v$.email.required.$invalid">Email harus diisi</template>
                  <template v-else-if="v$.email.email.$invalid">Format email tidak valid</template>
                  <template v-else-if="v$.email.maxLength.$invalid"
                    >Email maksimal 255 karakter</template
                  >
                </div>
              </div>
              <div class="flex flex-wrap gap-2 flex-col w-full">
                <!-- btn -->
                <div class="w-full grid">
                  <button
                    :disabled="authStore.loading"
                    type="submit"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                  >
                    <span v-if="authStore.loading">Loading...</span>
                    <span v-else>Submit</span>
                  </button>
                </div>
                <div class="w-full grid">
                  <a
                    href="/auth/signin"
                    class="btn inline-flex items-center gap-x-2 bg-gray-200 text-gray-700 border-gray-200 disabled:opacity-50 disabled:pointer-events-none hover:text-gray-700 hover:bg-gray-300 hover:border-gray-300 active:bg-gray-300 active:border-gray-300 focus:outline-none focus:ring-4 focus:ring-gray-300"
                  >
                    Back
                  </a>
                </div>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { reactive, computed, onBeforeUnmount } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, email, maxLength } from '@vuelidate/validators'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const { showError } = useErrorModal()

const form = reactive({
  email: ''
})

const rules = computed(() => ({
  email: {
    required,
    email,
    maxLength: maxLength(255)
  }
}))

const v$ = useVuelidate(rules, form)

onBeforeUnmount(() => {
  authStore.cancelRequests()
})

const handleSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    await authStore.forgotPassword(form.email)

    router.push('/auth/signin')
  } catch (error) {
    showError(error)
  }
}
</script>
