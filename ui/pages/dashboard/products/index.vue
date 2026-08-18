<template>
  <div class="container">
    <div class="grid grid-cols-1 mb-8">
      <!-- page header -->
      <div class="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
        <div>
          <h2 class="text-xl">Products</h2>
          <!-- breacrumb -->
          <nav aria-label="breadcrumb">
            <ol class="flex flex-wrap">
              <li class="inline-block text-blue-600">
                <a href="#!">
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

              <li class="inline-block text-gray-500 active" aria-current="page">Products</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
        <div>
          <a
            href="/dashboard/products/create"
            class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
          >
            Add New Product
          </a>
        </div>
      </div>
    </div>
    <!-- row -->
    <div class="grid grid-cols-1">
      <!-- card -->
      <div class="card h-full card-lg">
        <div class="px-6 py-6">
          <div class="grid grid-cols-12 justify-between">
            <!-- form -->
            <div class="lg:col-span-3 md:col-span-6 col-span-12 mb-2 lg:mb-0">
              <form class="flex" role="search">
                <input
                  v-model="searchFilter"
                  :disabled="useProduct.loading"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  type="search"
                  placeholder="Search Products"
                  aria-label="Search"
                />
              </form>
            </div>
            <!-- select option -->
            <div class="md:col-start-11 md:col-end-13 md:col-span-4 col-span-12">
              <select
                v-model="status"
                :disabled="useProduct.loading"
                class="text-base py-2 block w-full border-gray-300 rounded-lg focus:border-blue-600 focus:ring-blue-600 disabled:opacity-50 disabled:pointer-events-none"
                @change="handleFilterChange"
              >
                <option value="" disabled>Pilih Status</option>
                <option value="ACTIVE">Active</option>
                <option value="INACTIVE">Inactive</option>
                <option value="DRAFT">Draft</option>
              </select>
            </div>
          </div>
        </div>
        <!-- card body -->
        <div class="card-body p-0">
          <!-- table -->
          <div class="relative overflow-x-auto">
            <table class="text-left w-full whitespace-nowrap table-with-checkbox table-hover">
              <thead class="bg-gray-200 text-gray-700">
                <tr class="border-transparent border-b-0!">
                  <!-- <th scope="col" class="px-6 py-3">
                    <div class="flex items-center">
                      <input
                        id="checkAll"
                        :disabled="useProduct.loading"
                        class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded focus:ring-blue-600 focus:outline-none focus:ring-2"
                        type="checkbox"
                        value=""
                      />
                      <label class="text-gray-800 ms-3" for="checkAll" />
                    </div>
                  </th> -->
                  <th scope="col" class="px-6 py-3">Image</th>
                  <th scope="col" class="px-6 py-3">Product Name</th>
                  <th scope="col" class="px-6 py-3">Category</th>
                  <th scope="col" class="px-6 py-3">Status</th>
                  <th scope="col" class="px-6 py-3">Price</th>
                  <th scope="col" class="px-6 py-3">Create at</th>
                  <th scope="col" class="px-6 py-3" />
                </tr>
              </thead>
              <tbody class="divide-y">
                <!-- Loading Skeleton -->
                <tr v-if="loading">
                  <td colspan="8">
                    <div class="py-4">
                      <div
                        v-for="i in 5"
                        :key="'skeleton-dash-' + i"
                        class="flex items-center gap-4 py-3 border-b border-gray-300 px-6"
                      >
                        <!-- <div class="skeleton w-4 h-4 rounded" /> -->
                        <div class="skeleton w-12 h-12 rounded" />
                        <div class="flex-1">
                          <div class="skeleton skeleton-text w-48 mb-2" />
                          <div class="skeleton skeleton-text w-32" />
                        </div>
                        <div class="skeleton skeleton-button" style="width: 24px; height: 24px" />
                      </div>
                    </div>
                  </td>
                </tr>
                <tr v-else-if="useProduct?.loading">
                  <td colspan="8" class="text-center py-4">Loading...</td>
                </tr>
                <tr v-else-if="useProduct?.products?.length === 0">
                  <td colspan="8" class="text-center py-4">No data available</td>
                </tr>
                <tr
                  v-for="item in productDatas"
                  v-else
                  :key="item.id"
                  class="border-transparent border-b-0!"
                >
                  <!-- <td class="py-3 px-6 text-center">
                    <div class="flex items-center">
                      <input
                        id="productOne"
                        :disabled="useProduct.loading"
                        class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded focus:ring-blue-600 focus:outline-none focus:ring-2"
                        type="checkbox"
                        value=""
                      />
                      <label class="text-gray-800 ms-3" for="productOne" />
                    </div>
                  </td> -->
                  <td class="py-3 px-6 text-left">
                    <NuxtLink :to="`/dashboard/products/edit/${item.id}`">
                      <img :src="item.product_image" alt="" class="h-12 w-12" />
                    </NuxtLink>
                  </td>
                  <td class="py-3 px-6 text-left">
                    <NuxtLink :to="`/dashboard/products/edit/${item.id}`" class="text-inherit">
                      {{ item.product_name }}
                    </NuxtLink>
                  </td>
                  <td class="py-3 px-6 text-left">{{ item.category_name }}</td>

                  <td class="py-3 px-6 text-left">
                    <span
                      class="inline-block p-1 text-sm align-baseline leading-none rounded border font-semibold"
                      :class="getProductStatusClass(item.product_status)"
                    >
                      {{ item.product_status }}
                    </span>
                  </td>
                  <td class="py-3 px-6 text-left">Rp. {{ item.sale_price }}</td>
                  <td class="py-3 px-6 text-left">{{ formatDate(item.created_at) }}</td>
                  <td class="py-3 px-6 text-left">
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
                          <NuxtLink
                            :to="`/dashboard/products/edit/${item.id}`"
                            class="dropdown-item"
                          >
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
                          <a class="dropdown-item" href="#" @click="handleDelete(item.id)">
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
        <div
          v-if="useProduct.pagination?.total_count > 0"
          class="border-t border-gray-300 flex flex-col md:flex-row justify-between items-center px-6 py-6 gap-3"
        >
          <span
            >Showing {{ paginateProds.page }} to {{ paginateProds.total_page }} of
            {{ paginateProds.total_count }} entries</span
          >
          <nav class="flex items-center gap-x-1">
            <button
              :disabled="paginateProds.page === 1"
              type="button"
              class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
              @click="handlePageChange(paginateProds.current_page - 1)"
            >
              Previous
            </button>
            <div class="flex items-center gap-x-1">
              <button
                v-for="page in paginateProds.total_page"
                :key="page"
                :disabled="useProduct.loading"
                type="button"
                :class="[
                  'leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border',
                  page === paginateProds.page
                    ? 'text-white border bg-blue-600 border-blue-600 hover:bg-blue-600 focus:outline-none focus:bg-blue-600'
                    : 'bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300'
                ]"
                aria-current="page"
                @click="handlePageChange(page)"
              >
                {{ page }}
              </button>
            </div>
            <button
              :disabled="paginateProds.page === paginateProds.total_page"
              type="button"
              class="leading-none min-h-[36px] min-w-[36px] py-2 px-2.5 inline-flex justify-center items-center gap-x-1.5 rounded-md border bg-white border-gray-300 text-gray-800 hover:bg-gray-300 focus:outline-none focus:bg-gray-300 disabled:opacity-50 disabled:pointer-events-none"
              @click="handlePageChange(paginateProds.page + 1)"
            >
              Next
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { useProductStore } from '~/stores/product'
import { useErrorModal } from '~/composables/useErrorModal'
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { watchDebounced } from '@vueuse/core'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})
const { formatDate } = useDateTime()

const useProduct = useProductStore()
const { showError } = useErrorModal()
const productDatas = ref([])
const paginateProds = ref({})
const searchFilter = ref('')
const status = ref('')
const loading = ref(true)

const getProductStatusClass = (productStatus) => {
  const statusClass = {
    ACTIVE: 'bg-blue-100 text-blue-800 border-blue-200',
    INACTIVE: 'bg-red-100 text-red-800 border-red-200',
    DRAFT: 'bg-yellow-100 text-yellow-800 border-yellow-200'
  }

  return (
    statusClass[String(productStatus).toUpperCase()] || 'bg-gray-100 text-gray-800 border-gray-200'
  )
}

const syncProductsFromStore = () => {
  productDatas.value = useProduct.products
  paginateProds.value = useProduct.pagination
}

const handleFilterChange = async () => {
  loading.value = true
  try {
    await useProduct.fetchProductsAdmin({
      search: searchFilter.value,
      status: status.value,
      page: 1
    })
    syncProductsFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

watchDebounced(
  searchFilter,
  () => {
    handleFilterChange()
  },
  { debounce: 500, maxWait: 1000 }
  // maxWait (opsional): Memaksa API dipanggil setiap 1 detik jika user mengetik tanpa henti
)

onMounted(async () => {
  try {
    await useProduct.fetchProductsAdmin()
    syncProductsFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  useProduct.cancelRequests()
})

const handlePageChange = async (page) => {
  loading.value = true

  try {
    await useProduct.fetchProductsAdmin({ page: page })
    syncProductsFromStore()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id) => {
  try {
    const isConfirmed = window.confirm('Apakah Anda yakin ingin menghapus kategori ini?')

    if (isConfirmed) {
      await useProduct.deleteProductAdmin(id)
      await useProduct.fetchProductsAdmin()
      syncProductsFromStore()
    }
  } catch (error) {
    showError(error)
  }
}
</script>
