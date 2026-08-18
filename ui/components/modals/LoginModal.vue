<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content p-8 bg-white rounded-lg">
      <div class="flex justify-between items-center mb-4">
        <h3 class="font-bold text-gray-800">Sign Up</h3>
        <button
          :disabled="authStore.loading"
          type="button"
          class="btn-close"
          @click="$emit('close')"
        >
          <Icon name="tabler:x" class="text-gray-700" size="24" />
        </button>
      </div>

      <div class="modal-body">
        <form @submit.prevent="handleSubmit">
          <div
            v-if="authStore.error"
            class="mb-4 p-4 text-sm text-red-800 rounded-lg bg-red-50 relative"
            role="alert"
          >
            <button
              :disabled="authStore.loading"
              class="absolute -top-2 -right-2 p-1.5 inline-flex items-center justify-center text-red-800 bg-red-50 hover:bg-red-200 rounded-full"
              @click="authStore.error = null"
            >
              <span class="sr-only">Tutup</span>
              <i class="ti ti-x text-lg" />
            </button>
            <div class="flex items-center pr-8">
              <span class="font-medium mr-2">Gagal!</span> {{ authStore.error }}
            </div>
          </div>
          <div class="mb-3">
            <label for="fullName" class="mb-2 block text-gray-800">Nama</label>
            <input
              id="fullName"
              v-model="form.fullName"
              :disabled="authStore.loading"
              type="text"
              class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
              :class="{ 'border-red-500': v$.fullName.$error && v$.fullName.$dirty }"
              placeholder="Masukkan Nama Anda"
              required
              @blur="v$.fullName.$touch()"
            />
            <div v-if="v$.fullName.$error && v$.fullName.$dirty" class="text-red-600 text-sm mt-1">
              <template v-if="v$.fullName.required.$invalid">Nama harus diisi</template>
              <template v-else-if="v$.fullName.minLength.$invalid"
                >Nama minimal 3 karakter</template
              >
              <template v-else-if="v$.fullName.maxLength.$invalid"
                >Nama maksimal 255 karakter</template
              >
            </div>
          </div>

          <div class="mb-3">
            <label for="email" class="mb-2 block text-gray-800">Email</label>
            <input
              id="email"
              v-model="form.email"
              :disabled="authStore.loading"
              type="email"
              class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
              :class="{ 'border-red-500': v$.email.$error && v$.email.$dirty }"
              placeholder="Masukkan alamat email"
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

          <div class="mb-5">
            <label for="password" class="mb-2 block text-gray-800">Password</label>
            <input
              id="password"
              v-model="form.password"
              :disabled="authStore.loading"
              type="password"
              class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
              :class="{ 'border-red-500': v$.password.$error && v$.password.$dirty }"
              placeholder="Masukkan Password"
              required
              @blur="v$.password.$touch()"
            />
            <div v-if="v$.password.$error && v$.password.$dirty" class="text-red-600 text-sm mt-1">
              <template v-if="v$.password.required.$invalid">Password harus diisi</template>
              <template v-else-if="v$.password.minLength.$invalid"
                >Password minimal 8 karakter</template
              >
              <template v-else-if="v$.password.maxLength.$invalid"
                >Password maksimal 255 karakter</template
              >
            </div>
          </div>

          <div class="mb-5">
            <label for="passwordConfirmation" class="mb-2 block text-gray-800"
              >Password Confirmation</label
            >
            <input
              id="passwordConfirmation"
              v-model="form.password_confirmation"
              :disabled="authStore.loading"
              type="password"
              class="form-control border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
              :class="{
                'border-red-500': v$.password_confirmation.$error && v$.password_confirmation.$dirty
              }"
              placeholder="Masukkan Password Konfirmasi"
              required
              @blur="v$.password_confirmation.$touch()"
            />
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
            <span class="block mt-1 text-sm text-gray-500">
              Dengan mendaftar, Anda menyetujui
              <NuxtLink to="#" class="text-blue-600">Syarat dan Ketentuan</NuxtLink>
              &
              <NuxtLink to="#" class="text-blue-600">Kebijakan Privasi</NuxtLink>
            </span>
          </div>

          <button
            :disabled="authStore.loading"
            type="submit"
            class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300 justify-center"
          >
            <span v-if="authStore.loading">Loading...</span>
            <span v-else>Daftar</span>
          </button>
        </form>
      </div>

      <div class="modal-footer flex border-0 justify-center mt-3">
        Sudah punya akun?
        <NuxtLink to="/auth/signin" class="text-blue-600 ml-1">Masuk</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength, maxLength, sameAs } from '@vuelidate/validators'
import { useAuthStore } from '~/stores/auth'

const emit = defineEmits(['close'])

const authStore = useAuthStore()

const { showError } = useErrorModal()

const form = reactive({
  fullName: '',
  email: '',
  password: '',
  password_confirmation: ''
})

const rules = computed(() => ({
  fullName: {
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

const handleSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    const userData = {
      name: form.fullName,
      email: form.email,
      password: form.password,
      password_confirmation: form.password_confirmation
    }

    await authStore.signup(userData)

    form.fullName = ''
    form.email = ''
    form.password = ''
    form.password_confirmation = ''
    v$.value.$reset()

    emit('close')
    navigateTo('/auth/signin')
  } catch (error) {
    showError(error)
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  max-width: 500px;
  width: 90%;
}
</style>
