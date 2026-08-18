<template>
  <div class="container">
    <div class="grid grid-cols-1 mb-8">
      <!-- page header -->
      <div class="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
        <div>
          <h2 class="text-xl">Roles</h2>
          <!-- breacrumb -->
          <nav aria-label="breadcrumb">
            <ol class="flex flex-wrap">
              <li class="inline-block text-blue-600">
                <a href="../../dashboard">
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

              <li class="inline-block text-gray-500 active" aria-current="page">Roles</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
        <div>
          <a
            href="/dashboard/roles/create"
            class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
          >
            Add New Role
          </a>
        </div>
      </div>
    </div>
    <!-- row -->
    <div class="grid grid-cols-1">
      <!-- card -->
      <div class="card h-full card-lg">
        <!-- card body -->
        <div class="card-body p-0">
          <!-- table -->
          <div class="relative overflow-x-auto">
            <table class="text-left w-full whitespace-nowrap table-with-checkbox table-hover">
              <thead class="bg-gray-200 text-gray-700">
                <tr class="border-transparent border-b-0!">
                  <th scope="col" class="px-6 py-3">Name</th>
                  <th scope="col" class="px-6 py-3" />
                </tr>
              </thead>
              <tbody class="divide-y">
                <tr v-if="loading">
                  <td colspan="8" class="py-4">
                    <div
                      v-for="i in 5"
                      :key="'skeleton-dash-' + i"
                      class="flex items-center gap-4 py-3 border-b border-gray-100 px-6"
                    >
                      <div class="flex-1">
                        <div class="skeleton skeleton-text w-48 mb-2" />
                        <div class="skeleton skeleton-text w-32" />
                      </div>
                      <div class="skeleton skeleton-button" style="width: 24px; height: 24px" />
                    </div>
                  </td>
                </tr>
                <tr v-else-if="useRole?.roles?.length === 0">
                  <td colspan="8" class="text-center py-4">No data available</td>
                </tr>
                <tr
                  v-for="role in rolesData"
                  v-else
                  :key="role.id"
                  class="border-transparent border-b-0!"
                >
                  <td class="py-3 px-6 text-left">
                    <NuxtLink :to="`/dashboard/roles/edit/${role.id}`" class="text-inherit">
                      {{ role.name }}
                    </NuxtLink>
                  </td>
                  <td class="py-3 px-6 text-right flex justify-end items-center gap-2">
                    <div class="dropdown">
                      <a
                        href="#"
                        class="text-inherit"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="20"
                          height="20"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="icon icon-tabler icons-tabler-outline icon-tabler-dots-vertical"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path d="M12 12m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                          <path d="M12 19m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                          <path d="M12 5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                        </svg>
                      </a>
                      <ul class="dropdown-menu">
                        <li>
                          <NuxtLink class="dropdown-item" :to="`/dashboard/roles/edit/${role.id}`">
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
                              class="icon icon-tabler icons-tabler-outline icon-tabler-edit"
                            >
                              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                              <path
                                d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"
                              />
                              <path
                                d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"
                              />
                              <path d="M16 5l3 3" />
                            </svg>
                            Edit
                          </NuxtLink>
                        </li>
                        <li>
                          <a class="dropdown-item" href="#" @click.prevent="handleDelete(role.id)">
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
                              class="icon icon-tabler icons-tabler-outline icon-tabler-trash"
                            >
                              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                              <path d="M4 7l16 0" />
                              <path d="M10 11l0 6" />
                              <path d="M14 11l0 6" />
                              <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" />
                              <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" />
                            </svg>
                            Delete
                          </a>
                        </li>
                      </ul>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRoleStore } from '~/stores/role'
import { useErrorModal } from '~/composables/useErrorModal'
import { ref, onMounted, onBeforeUnmount } from 'vue'

const useRole = useRoleStore()
const { showError } = useErrorModal()
const rolesData = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    await useRole.fetchRoles()
    rolesData.value = useRole.roles
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  useRole.cancelRequests()
})

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const handleDelete = async (id) => {
  try {
    const isConfirmed = window.confirm('Apakah Anda yakin ingin menghapus role ini?')

    if (isConfirmed) {
      await useRole.deleteRole(id)
      await useRole.fetchRoles()
      rolesData.value = useRole.roles
    }
  } catch (error) {
    showError(error)
  }
}
</script>

<style></style>
