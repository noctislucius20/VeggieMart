<template>
  <div class="container">
    <div class="grid grid-cols-1 mb-8">
      <!-- page header -->
      <div class="md:flex justify-between items-center">
        <div>
          <h2 class="text-xl">Update Role</h2>
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
              <li class="inline-block text-blue-600">
                <a href="/dashboard/roles">
                  Roles
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
              <li class="inline-block text-gray-500 active" aria-current="page">Update Role</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
      </div>
    </div>
    <!-- row -->
    <div class="grid grid-cols-1">
      <!-- Skeleton Loading -->
      <div v-if="loading" class="card card-lg border-0">
        <div class="card-body flex flex-col gap-8 p-7">
          <div class="flex flex-col gap-4">
            <div class="skeleton skeleton-title h-5 w-16 mb-3" />
            <div class="grid grid-cols-12 gap-6">
              <div class="lg:col-span-6 col-span-12">
                <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                <div class="skeleton skeleton-text w-full" style="height: 44px" />
              </div>
              <div class="col-span-12 mt-3">
                <div class="flex flex-col md:flex-row gap-2">
                  <div class="skeleton skeleton-button h-10 w-28 rounded" />
                  <div class="skeleton skeleton-button h-10 w-24 rounded" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!loading" class="card card-lg border-0">
        <div class="card-body flex flex-col gap-8 p-7">
          <div class="flex flex-col gap-4">
            <h3 class="mb-0 text-md">Role</h3>
            <form
              class="grid grid-cols-12 gap-6 needs-validation"
              novalidate
              @submit.prevent="handleSubmit"
            >
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <!-- input -->
                  <label for="editRoleName" class="inline-block text-gray-800 font-medium mb-2">
                    Role Name
                    <span class="text-red-600">*</span>
                  </label>
                  <input
                    id="editRoleName"
                    v-model="form.name"
                    :disabled="useRole.loading"
                    type="text"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.name.$error && v$.name.$dirty }"
                    placeholder="Role Name"
                    required
                    @blur="v$.name.$touch()"
                  />
                  <div v-if="v$.name.$error && v$.name.$dirty" class="text-red-600 text-sm mt-1">
                    <template v-if="v$.name.required.$invalid">Nama role harus diisi</template>
                    <template v-else-if="v$.name.minLength.$invalid"
                      >Nama role minimal 3 karakter</template
                    >
                    <template v-else-if="v$.name.maxLength.$invalid"
                      >Nama role maksimal 50 karakter</template
                    >
                  </div>
                </div>
              </div>
              <div class="col-span-12 mt-3">
                <div class="flex flex-col md:flex-row gap-2">
                  <button
                    :disabled="useRole.loading"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                    type="submit"
                  >
                    <span v-if="useRole.loading">Loading...</span>
                    <span v-else>Update Role</span>
                  </button>
                  <a
                    href="/dashboard/roles"
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
import { useRoleStore } from '~/stores/role'
import { useErrorModal } from '~/composables/useErrorModal'
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { required, minLength, maxLength } from '@vuelidate/validators'

const useRole = useRoleStore()
const { showError } = useErrorModal()

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const route = useRoute()
const router = useRouter()
const roleId = route.params.id

const form = reactive({
  name: ''
})

const rules = computed(() => ({
  name: {
    required,
    minLength: minLength(3),
    maxLength: maxLength(50)
  }
}))

const v$ = useVuelidate(rules, form)

const loading = ref(true)

onMounted(async () => {
  try {
    await useRole.fetchRoleID(roleId)
    form.name = useRole.role.name
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  useRole.cancelRequests()
})

const handleSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    await useRole.updateRole(roleId, form.name)
    router.push('/dashboard/roles')
  } catch (error) {
    showError(error)
  }
}
</script>

<style></style>
