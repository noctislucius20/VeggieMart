<template>
  <section class="my-10">
    <div class="container">
      <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
        <div class="w-full md:w-1/3 xl:w-1/3 lg:order-1 order-2">
          <!-- img -->
          <img src="~/assets/images/svg-graphics/signin-g.svg" alt="" class="max-w-full h-auto" />
        </div>
        <div class="w-full md:w-1/2 lg:mx-1/6 xl:w-1/4 lg:order-2 order-1 flex flex-col gap-6">
          <div class="flex flex-col gap-1">
            <h1 class="text-xl">Sign in to DailyMart</h1>
            <p>Welcome back to DailyMart! Enter your email to get started.</p>
          </div>
          <form class="needs-validation" novalidate @submit.prevent="handleSubmit">
            <div class="flex flex-col gap-5">
              <div class="flex flex-col gap-3">
                <div class="w-full">
                  <!-- input -->
                  <label for="formSigninEmail" class="invisible hidden">Email address</label>
                  <input
                    id="formSigninEmail"
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
                    <template v-else-if="v$.email.email.$invalid"
                      >Format email tidak valid</template
                    >
                    <template v-else-if="v$.email.maxLength.$invalid"
                      >Email maksimal 255 karakter</template
                    >
                  </div>
                </div>
                <div class="w-full">
                  <!-- input -->
                  <div class="password-field relative">
                    <label for="formSigninPassword" class="invisible hidden">Password</label>
                    <div class="password-field relative">
                      <input
                        id="formSigninPassword"
                        v-model="form.password"
                        :disabled="authStore.loading"
                        :type="showPassword ? 'text' : 'password'"
                        class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base fakePassword"
                        :class="{ 'border-red-500': v$.password.$error && v$.password.$dirty }"
                        placeholder="*****"
                        required
                        @blur="v$.password.$touch()"
                      />
                      <span @click="togglePassword">
                        <i
                          :class="showPassword ? 'ti ti-eye' : 'ti ti-eye-off'"
                          class="passwordToggler"
                        />
                      </span>
                      <div
                        v-if="v$.password.$error && v$.password.$dirty"
                        class="text-red-600 text-sm mt-1"
                      >
                        <template v-if="v$.password.required.$invalid"
                          >Password harus diisi</template
                        >
                        <template v-else-if="v$.password.minLength.$invalid"
                          >Password minimal 8 karakter</template
                        >
                        <template v-else-if="v$.password.maxLength.$invalid"
                          >Password maksimal 255 karakter</template
                        >
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="flex flex-col gap-4">
                <div class="flex justify-between w-full">
                  <div>
                    Forgot password?
                    <a href="/auth/forgot-password" class="text-blue-600">Reset It</a>
                  </div>
                </div>
                <!-- btn -->
                <div class="w-full grid">
                  <button
                    :disabled="authStore.loading"
                    type="submit"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                  >
                    <span v-if="authStore.loading">Loading...</span>
                    <span v-else>Sign In</span>
                  </button>
                </div>
              </div>
            </div>
            <!-- link -->
            <div class="mt-2">
              Don&rsquo;t have an account?
              <a href="signup" class="text-blue-600">Sign Up</a>
            </div>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive, computed, onBeforeUnmount } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength, maxLength } from '@vuelidate/validators'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const { showError } = useErrorModal()

const showPassword = ref(false)

const form = reactive({
  email: '',
  password: ''
})

const rules = computed(() => ({
  email: {
    required,
    email,
    maxLength: maxLength(255)
  },
  password: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(255)
  }
}))

const v$ = useVuelidate(rules, form)

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

onBeforeUnmount(() => {
  authStore.cancelRequests()
})

const handleSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    await authStore.signin(form.email, form.password)
    // Redirect ke halaman dashboard setelah login berhasil

    if (authStore.user.role === 'Super Admin') {
      router.push('/dashboard')
    } else {
      router.push('/')
    }
  } catch (error) {
    // Error sudah ditangani di store
    showError(error)
  }
}
</script>
