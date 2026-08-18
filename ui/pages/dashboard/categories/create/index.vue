<template>
  <div class="container">
    <div class="grid grid-cols-1 mb-8">
      <!-- page header -->
      <div class="md:flex justify-between items-center">
        <div>
          <h2 class="text-xl">Create Category</h2>
          <!-- breacrumb -->
          <nav aria-label="breadcrumb">
            <ol class="flex flex-wrap">
              <li class="inline-block text-blue-600">
                <a href="/dashboard/categories">
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
              <li class="inline-block text-blue-600">
                <a href="../dashboard/categories.html">
                  Categories
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

              <li class="inline-block text-gray-500 active" aria-current="page">Create Category</li>
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
          <div class="flex flex-col md:flex-row items-center file-input-wrapper gap-2">
            <div>
              <!-- Skeleton loading saat proses upload gambar berlangsung -->
              <div
                v-if="iconUploading"
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
                id="icon"
                :disabled="useCategory.loading"
                type="file"
                class="file-input opacity-0"
                accept="image/jpeg,image/png"
                @change="handleFileUpload"
              />
              Upload Photo
            </div>

            <span class="ms-2">JPG, PNG. 1MB Max.</span>
          </div>
          <div v-if="iconError" class="text-red-600 text-sm -mt-4">{{ iconError }}</div>
          <div v-if="v$.icon.$error && v$.icon.$dirty" class="text-red-600 text-sm -mt-4">
            <template v-if="v$.icon.required.$invalid">Icon harus diupload</template>
            <template v-else-if="v$.icon.url.$invalid">URL icon tidak valid</template>
            <template v-else-if="v$.icon.maxLength.$invalid"
              >URL icon maksimal 255 karakter</template
            >
          </div>
          <div class="flex flex-col gap-4">
            <h3 class="mb-0 text-md">Category Information</h3>
            <form
              class="grid grid-cols-12 gap-6 needs-validation"
              novalidate
              @submit.prevent="handleSubmit"
            >
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <!-- input -->
                  <label
                    for="creatCustomerName"
                    class="inline-block text-gray-800 font-medium mb-2"
                  >
                    Category Name
                    <span class="text-red-600">*</span>
                  </label>
                  <input
                    id="creatCustomerName"
                    v-model="form.name"
                    :disabled="useCategory.loading"
                    type="text"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.name.$error && v$.name.$dirty }"
                    placeholder="Category Name"
                    required
                    @blur="v$.name.$touch()"
                  />
                  <div v-if="v$.name.$error && v$.name.$dirty" class="text-red-600 text-sm mt-1">
                    <template v-if="v$.name.required.$invalid">Nama kategori harus diisi</template>
                    <template v-else-if="v$.name.minLength.$invalid"
                      >Nama kategori minimal 3 karakter</template
                    >
                    <template v-else-if="v$.name.maxLength.$invalid"
                      >Nama kategori maksimal 100 karakter</template
                    >
                  </div>
                </div>
              </div>
              <div class="lg:col-span-12 col-span-12">
                <div>
                  <!-- input -->
                  <label
                    for="creatCustomerPhone"
                    class="inline-block text-gray-800 font-medium mb-2"
                    >Description</label
                  >
                  <textarea
                    id=""
                    v-model="form.description"
                    :disabled="useCategory.loading"
                    name="description"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  />
                </div>
              </div>

              <div class="mb-3 col-lg-12">
                <label class="form-label" for="creatCustomerDate">Status</label>
                <div class="flex space-x-6">
                  <!-- Active Radio Button -->
                  <label class="inline-flex items-center">
                    <input
                      v-model="form.status"
                      :disabled="useCategory.loading"
                      type="radio"
                      class="form-radio focus:outline-none text-blue-600"
                      name="status"
                      :value="1"
                    />
                    <span class="ml-2">Active</span>
                  </label>

                  <!-- Disabled Radio Button -->
                  <label class="inline-flex items-center">
                    <input
                      v-model="form.status"
                      :disabled="useCategory.loading"
                      type="radio"
                      class="form-radio focus:outline-none text-blue-600"
                      name="status"
                      :value="0"
                    />
                    <span class="ml-2">Disabled</span>
                  </label>
                </div>
                <div v-if="v$.status.$error && v$.status.$dirty" class="text-red-600 text-sm mt-1">
                  <template v-if="v$.status.required.$invalid">Status harus dipilih</template>
                </div>
              </div>

              <div class="col-span-12 mt-3">
                <div class="flex flex-col md:flex-row gap-2">
                  <button
                    :disabled="useCategory.loading"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                    type="submit"
                  >
                    <span v-if="useCategory.loading">Loading...</span>
                    <span v-else>Create New Category</span>
                  </button>
                  <a
                    href="/dashboard/categories"
                    class="btn inline-flex items-center gap-x-2 bg-gray-200 text-gray-800 border-gray-200 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
                    type="submit"
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
import { useCategoryStore } from '~/stores/category'
import { ref, reactive, computed, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { required, minLength, maxLength, url } from '@vuelidate/validators'
import { oneof } from '~/utils/validators'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const router = useRouter()
const imageUrl = ref('')
const iconError = ref('')
const iconUploading = ref(false)
const useCategory = useCategoryStore()
const { showError } = useErrorModal()

const form = reactive({
  name: '',
  description: '',
  status: 1,
  icon: '',
  parent_id: ''
})

const rules = computed(() => ({
  name: {
    required,
    minLength: minLength(3),
    maxLength: maxLength(100)
  },
  icon: {
    required,
    url,
    maxLength: maxLength(255)
  },
  description: {},
  status: {
    required,
    oneof: oneof([1, 0])
  },
  parent_id: {}
}))

const v$ = useVuelidate(rules, form)

const handleFileUpload = async (event) => {
  try {
    const selectedFile = event.target.files[0]
    iconError.value = ''
    if (selectedFile) {
      // Validasi ukuran file (max 1MB)
      if (selectedFile.size > 1024 * 1024) {
        iconError.value = 'Ukuran file maksimal 1MB'
        return
      }

      // Validasi tipe file
      if (!['image/jpeg', 'image/png'].includes(selectedFile.type)) {
        iconError.value = 'File harus berupa JPG atau PNG'
        return
      }

      // Upload file
      iconUploading.value = true
      const result = await useCategory.uploadImage(selectedFile)

      // Update preview dan simpan URL
      imageUrl.value = result.data.image_url || result.data.imageUrl
      form.icon = imageUrl.value
    }
  } catch (error) {
    showError(error)
  } finally {
    iconUploading.value = false
  }
}

onBeforeUnmount(() => {
  useCategory.cancelRequests()
})

const handleSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    const category = {
      name: form.name,
      description: form.description,
      status: form.status === 0 ? 'Unpublished' : 'Published',
      icon: form.icon,
      parent_id: null
    }

    await useCategory.createCategory(category)

    router.push('/dashboard/categories')
  } catch (error) {
    showError(error)
  }
}
</script>

<style></style>
