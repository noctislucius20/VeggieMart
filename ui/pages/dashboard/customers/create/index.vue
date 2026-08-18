<template>
  <div class="container">
    <div class="grid grid-cols-1 mb-8">
      <!-- page header -->
      <div class="md:flex justify-between items-center">
        <div>
          <h2 class="text-xl">Create Customer</h2>
          <!-- breacrumb -->
          <nav aria-label="breadcrumb">
            <ol class="flex flex-wrap">
              <li class="inline-block text-blue-600">
                <a href="/dashboard">
                  Dashboard
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="icon icon-tabler icons-tabler-outline icon-tabler-slash inline-block mx-2"
                  >
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <path d="M17 5l-10 14" />
                  </svg>
                </a>
              </li>

              <li class="inline-block text-gray-500 active" aria-current="page">Create Customer</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
      </div>
    </div>
    <!-- row -->
    <div class="grid grid-cols-1">
      <div class="card card-lg border-0">
        <div class="card-body flex flex-col gap-8 p-7">
          <div class="flex flex-col md:flex-row items-center mb-4 file-input-wrapper gap-2">
            <div>
              <!-- Skeleton loading saat proses upload gambar berlangsung -->
              <div
                v-if="photoUploading"
                class="skeleton skeleton-image"
                style="height: 64px; width: 64px; border-radius: 0.5rem"
              />
              <img
                v-else-if="imageUrl"
                class="image h-16 w-16 rounded-lg"
                :src="imageUrl"
                alt="Image"
              />
              <!-- Kotak kosong placeholder saat foto belum diupload -->
              <div
                v-else
                class="h-16 w-16 rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 flex items-center justify-center"
              >
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
                  class="text-gray-400"
                >
                  <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                  <path d="M15 8h.01" />
                  <path
                    d="M3 6a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v12a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3v-12z"
                  />
                  <path d="M3 16l5 -5c.928 -.893 2.072 -.893 3 0l5 5" />
                  <path d="M14 14l1 -1c.928 -.893 2.072 -.893 3 0l3 3" />
                </svg>
              </div>
            </div>

            <div
              class="file-upload btn inline-flex items-center gap-x-2 bg-gray-200 text-gray-800 border-gray-200 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300 md:ml-4"
            >
              <input
                :disabled="customerStore.loading"
                type="file"
                class="file-input opacity-0"
                accept="image/jpeg,image/png"
                @change="handleFileUpload"
              />
              Upload Photo
            </div>

            <span class="ms-2">JPG, PNG. 1MB Max.</span>
          </div>
          <div v-if="photoError" class="text-red-600 text-sm -mt-4">{{ photoError }}</div>
          <div v-if="v$.photo.$error && v$.photo.$dirty" class="text-red-600 text-sm -mt-4">
            <template v-if="v$.photo.required.$invalid">Foto harus diupload</template>
            <template v-else-if="v$.photo.url.$invalid">URL foto tidak valid</template>
            <template v-else-if="v$.photo.maxLength.$invalid"
              >URL foto maksimal 255 karakter</template
            >
          </div>
          <div class="flex flex-col gap-4">
            <h3 class="mb-0 text-md">Customer Information</h3>
            <form
              class="grid grid-cols-12 gap-6 needs-validation"
              novalidate
              @submit.prevent="handlerSubmit"
            >
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <!-- input -->
                  <label
                    for="creatCustomerName"
                    class="inline-block text-gray-800 font-medium mb-2"
                  >
                    Name
                    <span class="text-red-600">*</span>
                  </label>
                  <input
                    id="creatCustomerName"
                    v-model="form.name"
                    :disabled="customerStore.loading"
                    type="text"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.name.$error && v$.name.$dirty }"
                    placeholder="Customer Name"
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
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <!-- input -->
                  <label
                    for="creatCustomerEmail"
                    class="inline-block text-gray-800 font-medium mb-2"
                  >
                    Email
                    <span class="text-red-600">*</span>
                  </label>
                  <input
                    id="creatCustomerEmail"
                    v-model="form.email"
                    :disabled="customerStore.loading"
                    type="email"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.email.$error && v$.email.$dirty }"
                    placeholder="Email Address"
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
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <!-- input -->
                  <label
                    for="creatCustomerPhone"
                    class="inline-block text-gray-800 font-medium mb-2"
                    >Phone</label
                  >
                  <input
                    id="creatCustomerPhone"
                    v-model="form.phone"
                    :disabled="customerStore.loading"
                    type="text"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.phone.$error && v$.phone.$dirty }"
                    placeholder="Number"
                    @blur="v$.phone.$touch()"
                  />
                  <div v-if="v$.phone.$error && v$.phone.$dirty" class="text-red-600 text-sm mt-1">
                    <template v-if="v$.phone.maxLength.$invalid"
                      >Nomor telepon maksimal 17 karakter</template
                    >
                    <template v-if="v$.phone.numeric.$invalid"
                      >Nomor telepon harus berupa angka</template
                    >
                  </div>
                </div>
              </div>

              <div class="lg:col-span-6 col-span-12">
                <label class="inline-block text-gray-800 font-medium mb-2" for="address"
                  >Address</label
                >
                <input
                  id="address"
                  v-model="form.address"
                  :disabled="customerStore.loading"
                  type="text"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base flatpickr"
                />
              </div>

              <div class="lg:col-span-6 col-span-12">
                <div class="password-field relative">
                  <label class="inline-block text-gray-800 font-medium mb-2" for="password"
                    >Password <span class="text-red-600">*</span></label
                  >
                  <div class="password-field relative">
                    <input
                      id="password"
                      v-model="form.password"
                      :disabled="customerStore.loading"
                      :type="showPassword ? 'text' : 'password'"
                      class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base flatpickr"
                      :class="{ 'border-red-500': v$.password.$error && v$.password.$dirty }"
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
                      <template v-if="v$.password.required.$invalid">Password harus diisi</template>
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

              <div class="lg:col-span-6 col-span-12">
                <div class="password-field relative">
                  <label
                    class="inline-block text-gray-800 font-medium mb-2"
                    for="password_confirmation"
                    >Password Confirmation <span class="text-red-600">*</span></label
                  >
                  <div class="password-field relative">
                    <input
                      id="password_confirmation"
                      v-model="form.password_confirmation"
                      :disabled="customerStore.loading"
                      :type="showPasswordConfirm ? 'text' : 'password'"
                      class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base flatpickr"
                      :class="{
                        'border-red-500':
                          v$.password_confirmation.$error && v$.password_confirmation.$dirty
                      }"
                      @blur="v$.password_confirmation.$touch()"
                    />
                    <span @click="togglePasswordConfirmation">
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

              <div class="col-span-12 mt-3">
                <div class="flex flex-col md:flex-row gap-2">
                  <button
                    :disabled="customerStore.loading"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                    type="submit"
                  >
                    <span v-if="customerStore.loading">Loading...</span>
                    <span v-else>Create New Customer</span>
                  </button>
                  <a
                    href="/dashboard/customers"
                    class="btn inline-flex items-center gap-x-2 bg-gray-200 text-gray-800 border-gray-200 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
                  >
                    Cancel
                  </a>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useCustomerStore } from '~/stores/customer'
import { ref, reactive, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength, maxLength, sameAs } from '@vuelidate/validators'
import { numeric } from '~/utils/validators'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const customerStore = useCustomerStore()
const router = useRouter()
const showPassword = ref(false)
const showPasswordConfirm = ref(false)
const imageUrl = ref('')
const photoError = ref('')
const photoUploading = ref(false)
const { showError } = useErrorModal()

const form = reactive({
  name: '',
  email: '',
  phone: '',
  photo: '',
  password: '',
  address: '',
  role_id: 2,
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
  },
  phone: {
    maxLength: maxLength(17),
    numeric
  },
  address: {},
  photo: {
    maxLength: maxLength(255)
  }
}))

const v$ = useVuelidate(rules, form)

const handleFileUpload = async (event) => {
  try {
    const selectedFile = event.target.files[0]
    photoError.value = ''
    if (selectedFile) {
      // Validasi ukuran file (max 1MB)
      if (selectedFile.size > 1024 * 1024) {
        photoError.value = 'Ukuran file maksimal 1MB'
        return
      }

      // Validasi tipe file
      if (!['image/jpeg', 'image/png'].includes(selectedFile.type)) {
        photoError.value = 'File harus berupa JPG atau PNG'
        return
      }

      // Upload file
      photoUploading.value = true
      const result = await customerStore.uploadImage(selectedFile)

      // Update preview dan simpan URL
      imageUrl.value = result.data.image_url || result.data.imageUrl
      form.photo = imageUrl.value
    }
  } catch (error) {
    showError(error)
  } finally {
    photoUploading.value = false
  }
}

onBeforeUnmount(() => {
  customerStore.cancelRequests()
})

const handlerSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    const customerData = {
      name: form.name,
      email: form.email,
      phone: form.phone,
      photo: form.photo,
      password: form.password,
      address: form.address,
      password_confirmation: form.password_confirmation,
      role_id: 2
    }
    await customerStore.createCustomer(customerData)

    router.push('/dashboard/customers')
  } catch (error) {
    showError(error)
  }
}

const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const togglePasswordConfirmation = () => {
  showPasswordConfirm.value = !showPasswordConfirm.value
}
</script>

<style></style>
