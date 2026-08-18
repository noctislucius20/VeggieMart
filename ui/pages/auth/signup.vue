<template>
  <section class="my-10">
    <div class="container">
      <div class="flex flex-wrap justify-center items-center gap-8 lg:gap-16">
        <div class="w-full md:w-1/3 xl:w-1/3 lg:order-1 order-2">
          <!-- img -->
          <img src="~/assets/images/svg-graphics/signup-g.svg" alt="" class="max-w-full h-auto" />
        </div>

        <div class="w-full md:w-1/2 lg:mx-1/6 xl:w-1/3 lg:order-2 order-1 flex flex-col gap-6">
          <div class="flex flex-col gap-1">
            <h1 class="text-xl">Get Start Shopping</h1>
            <p>Welcome to DailyMart! Enter your email to get started.</p>
          </div>
          <form class="needs-validation" novalidate @submit.prevent="handleSubmit">
            <div
              v-if="showSuccessMessage"
              class="mb-4 p-4 text-sm text-blue-800 rounded-lg bg-blue-50"
              role="alert"
            >
              <span class="font-medium">Berhasil!</span> Akun anda telah berhasil dibuat. Silahkan
              login.
            </div>

            <!-- Alert Error -->
            <div
              v-if="authStore.error"
              class="mb-4 p-4 text-sm text-red-800 rounded-lg bg-red-50"
              role="alert"
            >
              <span class="font-medium">Gagal!</span> {{ authStore.error }}
            </div>
            <div class="flex flex-wrap gap-5">
              <div class="flex flex-col gap-3">
                <div class="flex flex-row gap-3">
                  <div class="w-1/2">
                    <!-- input -->
                    <label for="formSignupfname" class="invisible hidden">Name</label>
                    <input
                      id="formSignupfname"
                      v-model="form.name"
                      :disabled="authStore.loading"
                      type="text"
                      class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                      :class="{ 'border-red-500': v$.name.$error && v$.name.$dirty }"
                      placeholder="Name"
                      required
                      @blur="v$.name.$touch()"
                    />
                    <div v-if="v$.name.$error && v$.name.$dirty" class="text-red-600 text-sm mt-1">
                      <template v-if="v$.name.required.$invalid">Nama harus diisi</template>
                      <template v-else-if="v$.name.minLength.$invalid"
                        >Nama minimal 3 karakter</template
                      >
                      <template v-else-if="v$.name.maxLength.$invalid"
                        >Nama maksimal 255 karakter</template
                      >
                    </div>
                  </div>
                  <div class="w-1/2">
                    <!-- input -->
                    <label for="formSignupEmail" class="invisible hidden">Email address</label>
                    <input
                      id="formSignupEmail"
                      v-model="form.email"
                      :disabled="authStore.loading"
                      type="email"
                      class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                      :class="{ 'border-red-500': v$.email.$error && v$.email.$dirty }"
                      placeholder="Email"
                      required
                      @blur="v$.email.$touch()"
                    />
                    <div
                      v-if="v$.email.$error && v$.email.$dirty"
                      class="text-red-600 text-sm mt-1"
                    >
                      <template v-if="v$.email.required.$invalid">Email harus diisi</template>
                      <template v-else-if="v$.email.email.$invalid"
                        >Format email tidak valid</template
                      >
                      <template v-else-if="v$.email.maxLength.$invalid"
                        >Email maksimal 255 karakter</template
                      >
                    </div>
                  </div>
                </div>
                <div class="w-full">
                  <div class="password-field relative">
                    <label for="formSignupPassword" class="invisible hidden">Password</label>
                    <div class="password-field relative">
                      <input
                        id="formSignupPassword"
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
                <div class="w-full">
                  <div class="password-field relative">
                    <label for="formSignupPasswordConfirm" class="invisible hidden"
                      >Password Confirmation</label
                    >
                    <div class="password-field relative">
                      <input
                        id="formSignupPasswordConfirm"
                        v-model="form.password_confirmation"
                        :disabled="authStore.loading"
                        :type="showPasswordConfirm ? 'text' : 'password'"
                        class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base fakePassword"
                        :class="{
                          'border-red-500':
                            v$.password_confirmation.$error && v$.password_confirmation.$dirty
                        }"
                        placeholder="*****"
                        required
                        @blur="v$.password_confirmation.$touch()"
                      />
                      <span @click="togglePasswordConfirm">
                        <i
                          :class="showPasswordConfirm ? 'ti ti-eye' : 'ti ti-eye-off'"
                          class="passwordToggler"
                        />
                      </span>
                      <div
                        v-if="v$.password_confirmation.$error && v$.password_confirmation.$dirty"
                        class="text-red-600 text-sm mt-1"
                      >
                        <template v-if="v$.password_confirmation.required.$invalid"
                          >Konfirmasi password harus diisi</template
                        >
                        <template v-else-if="v$.password_confirmation.minLength.$invalid"
                          >Konfirmasi password minimal 8 karakter</template
                        >
                        <template v-else-if="v$.password_confirmation.maxLength.$invalid"
                          >Konfirmasi password maksimal 255 karakter</template
                        >
                        <template v-else-if="v$.password_confirmation.sameAs.$invalid"
                          >Konfirmasi password harus sama dengan password</template
                        >
                      </div>
                    </div>
                  </div>
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
                  <span v-else>Register</span>
                </button>
              </div>

              <!-- text -->
            </div>
            <div class="mt-2">
              <p>
                <small>
                  By continuing, you agree to our
                  <a href="#!">Terms of Service</a>
                  &
                  <a href="#!" class="text-blue-600">Privacy Policy</a>
                </small>
              </p>
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
import { required, email, minLength, maxLength, sameAs } from '@vuelidate/validators'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const showPassword = ref(false)
const showPasswordConfirm = ref(false)
const showSuccessMessage = ref(false)
const { showError } = useErrorModal()

const form = reactive({
  name: '',
  email: '',
  password: '',
  password_confirmation: ''
})

const rules = computed(() => ({
  name: {
    required,
    minLength: minLength(3),
    maxLength: maxLength(255)
  },
  email: {
    required,
    email,
    maxLength: maxLength(255)
  },
  password: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(255)
  },
  password_confirmation: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(255),
    sameAs: sameAs(computed(() => form.password))
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
    const userData = {
      name: form.name,
      email: form.email,
      password: form.password,
      password_confirmation: form.password_confirmation
    }
    await authStore.signup(userData)
    showSuccessMessage.value = true

    form.name = ''
    form.email = ''
    form.password = ''
    form.password_confirmation = ''
    v$.value.$reset()

    setTimeout(() => {
      router.push('/auth/signin')
    }, 2000)
  } catch (error) {
    showError(error)
  }
}

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const togglePasswordConfirm = () => {
  showPasswordConfirm.value = !showPasswordConfirm.value
}
</script>
