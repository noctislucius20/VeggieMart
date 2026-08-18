<template>
  <div class="container">
    <div class="grid grid-cols-1 mb-8">
      <!-- page header -->
      <div class="md:flex justify-between items-center">
        <div>
          <h2 class="text-xl">Update Product</h2>
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
                <a href="/dashboard/products">
                  Products
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
              <li class="inline-block text-gray-500 active" aria-current="page">Update Product</li>
            </ol>
          </nav>
        </div>
        <!-- button -->
        <div>
          <a
            href="/dashboard/products"
            class="btn inline-flex items-center gap-x-2 bg-gray-300 text-black border-gray-300 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-100"
            >Back to Product</a
          >
        </div>
      </div>
    </div>
    <!-- row -->
    <div class="grid grid-cols-1 mb-5">
      <!-- Loading Skeleton -->
      <div v-if="loading">
        <div class="card card-lg border-0 row-span-4 col-span-2">
          <div class="card-body flex flex-col gap-8 p-7">
            <div class="skeleton skeleton-title w-48 mb-4" style="height: 24px" />
            <div class="grid grid-cols-12 gap-6">
              <div class="lg:col-span-6 col-span-12">
                <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                <div class="skeleton skeleton-text w-full" style="height: 44px" />
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                <div class="skeleton skeleton-text w-full" style="height: 44px" />
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                <div class="skeleton skeleton-text w-full" style="height: 44px" />
              </div>
              <div class="lg:col-span-12 col-span-12">
                <div class="skeleton skeleton-text w-40 mb-2" style="height: 18px" />
                <div class="skeleton skeleton-text w-full" style="height: 100px" />
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div class="skeleton skeleton-text w-32 mb-2" style="height: 18px" />
                <div class="skeleton skeleton-text w-full" style="height: 44px" />
              </div>
              <div class="col-span-12 mt-3">
                <div class="skeleton skeleton-button" style="width: 150px; height: 44px" />
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- Form Content -->
      <div v-else class="card card-lg border-0 row-span-4 col-span-2">
        <div class="card-body flex flex-col gap-8 p-7">
          <div class="flex flex-col gap-4">
            <h3 class="mb-0 text-md">Product Information</h3>
            <form
              class="grid grid-cols-12 gap-6 needs-validation"
              novalidate
              @submit.prevent="handleSubmit"
            >
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <label
                    for="creatCustomerName"
                    class="inline-block text-gray-800 font-medium mb-2"
                  >
                    Product Name
                    <span class="text-red-600">*</span>
                  </label>
                  <input
                    id="creatCustomerName"
                    v-model="form.name"
                    :disabled="useProduct.loading"
                    type="text"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.name.$error && v$.name.$dirty }"
                    placeholder="Product Name"
                    required
                    @blur="v$.name.$touch()"
                  />
                  <div v-if="v$.name.$error && v$.name.$dirty" class="text-red-600 text-sm mt-1">
                    <template v-if="v$.name.required.$invalid">Nama produk harus diisi</template>
                    <template v-else-if="v$.name.minLength.$invalid"
                      >Nama produk minimal 3 karakter</template
                    >
                    <template v-else-if="v$.name.maxLength.$invalid"
                      >Nama produk maksimal 100 karakter</template
                    >
                  </div>
                </div>
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <label class="inline-block text-gray-800 font-medium mb-2" for="selectOne"
                    >Product Category
                    <span class="text-red-600">*</span>
                  </label>
                  <select
                    v-model="form.category_slug"
                    :disabled="useProduct.loading"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{
                      'border-red-500': v$.category_slug.$error && v$.category_slug.$dirty
                    }"
                    aria-label="Default select example"
                    required
                    @blur="v$.category_slug.$touch()"
                  >
                    <option value="" disabled>Open this select category</option>
                    <option
                      v-for="category in categoryDatas"
                      :key="category.slug"
                      :value="category.slug"
                    >
                      {{ category.name }}
                    </option>
                  </select>
                  <div
                    v-if="v$.category_slug.$error && v$.category_slug.$dirty"
                    class="text-red-600 text-sm mt-1"
                  >
                    <template v-if="v$.category_slug.required.$invalid"
                      >Kategori harus dipilih</template
                    >
                    <template v-else-if="v$.category_slug.maxLength.$invalid"
                      >Kategori maksimal 120 karakter</template
                    >
                  </div>
                </div>
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <label
                    for="creatCustomerName"
                    class="inline-block text-gray-800 font-medium mb-2"
                  >
                    Unit
                    <span class="text-red-600">*</span>
                  </label>
                  <select
                    v-model="form.unit"
                    :disabled="useProduct.loading"
                    required
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.unit.$error && v$.unit.$dirty }"
                    aria-label="Default select example"
                    @blur="v$.unit.$touch()"
                  >
                    <option value="" disabled>Open this select unit</option>
                    <option value="gr">Gram</option>
                    <option value="kg">Kilogram</option>
                  </select>
                  <div v-if="v$.unit.$error && v$.unit.$dirty" class="text-red-600 text-sm mt-1">
                    <template v-if="v$.unit.required.$invalid">Unit harus dipilih</template>
                    <template v-else-if="v$.unit.maxLength.$invalid"
                      >Unit maksimal 120 karakter</template
                    >
                  </div>
                </div>
              </div>
              <div class="lg:col-span-12 col-span-12">
                <div>
                  <label class="inline-block text-gray-800 font-medium mb-2"
                    >Product Description
                    <span class="text-red-600">*</span>
                  </label>
                  <textarea
                    id="product_description"
                    v-model="form.description"
                    :disabled="useProduct.loading"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{
                      'border-red-500': v$.description.$error && v$.description.$dirty
                    }"
                    name="product_description"
                    required
                    @blur="v$.description.$touch()"
                  />
                  <div
                    v-if="v$.description.$error && v$.description.$dirty"
                    class="text-red-600 text-sm mt-1"
                  >
                    <template v-if="v$.description.required.$invalid"
                      >Deskripsi produk harus diisi</template
                    >
                  </div>
                </div>
              </div>
              <div class="lg:col-span-6 col-span-12">
                <div>
                  <label
                    for="creatCustomerName"
                    class="inline-block text-gray-800 font-medium mb-2"
                  >
                    Product Status
                    <span class="text-red-600">*</span>
                  </label>
                  <select
                    v-model="form.status"
                    :disabled="useProduct.loading"
                    class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                    :class="{ 'border-red-500': v$.status.$error && v$.status.$dirty }"
                    aria-label="Default select example"
                    required
                    @blur="v$.status.$touch()"
                  >
                    <option value="" disabled>Open this select status</option>
                    <option value="DRAFT">Draft</option>
                    <option value="ACTIVE">Active</option>
                    <option value="INACTIVE">Inactive</option>
                  </select>
                  <div
                    v-if="v$.status.$error && v$.status.$dirty"
                    class="text-red-600 text-sm mt-1"
                  >
                    <template v-if="v$.status.required.$invalid">Status harus dipilih</template>
                  </div>
                </div>
              </div>

              <div
                v-if="v$.variant.$error && v$.variant.$dirty"
                class="col-span-12 text-red-600 text-sm"
              >
                <template v-if="v$.variant.required.$invalid || v$.variant.gte.$invalid"
                  >Minimal 1 variant harus ditambahkan</template
                >
              </div>

              <div class="col-span-12 mt-3">
                <div class="flex flex-col md:flex-row gap-2">
                  <button
                    :disabled="useProduct.loading"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                    type="submit"
                  >
                    <span v-if="useProduct.loading">Loading...</span>
                    <span v-else>Update Product</span>
                  </button>
                  <a
                    href="/dashboard/products"
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
    <!-- Loading Skeleton for Variant -->
    <div v-if="loading">
      <div class="card h-full card-lg">
        <div class="p-6">
          <div class="flex justify-between flex-row items-center mb-4">
            <div class="skeleton skeleton-title w-32" style="height: 24px" />
            <div class="skeleton skeleton-button" style="width: 120px; height: 40px" />
          </div>
        </div>
        <div class="card-body p-0">
          <div class="relative overflow-x-auto">
            <table class="text-left w-full whitespace-nowrap table-with-checkbox">
              <thead class="bg-gray-200 text-gray-700">
                <tr class="border-transparent border-b-0!">
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton w-4 h-4 rounded" />
                  </th>
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton skeleton-text w-16" style="height: 16px" />
                  </th>
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton skeleton-text w-20" style="height: 16px" />
                  </th>
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton skeleton-text w-28" style="height: 16px" />
                  </th>
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton skeleton-text w-24" style="height: 16px" />
                  </th>
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton skeleton-text w-24" style="height: 16px" />
                  </th>
                  <th scope="col" class="px-6 py-3">
                    <div class="skeleton skeleton-button" style="width: 24px; height: 24px" />
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="i in 3"
                  :key="'skeleton-variant-' + i"
                  class="border-transparent border-b-0!"
                >
                  <td class="py-3 px-6"><div class="skeleton w-4 h-4 rounded" /></td>
                  <td class="py-3 px-6"><div class="skeleton skeleton-text w-16" /></td>
                  <td class="py-3 px-6"><div class="skeleton skeleton-text w-20" /></td>
                  <td class="py-3 px-6"><div class="skeleton skeleton-text w-28" /></td>
                  <td class="py-3 px-6"><div class="skeleton skeleton-text w-24" /></td>
                  <td class="py-3 px-6"><div class="skeleton w-12 h-12 rounded" /></td>
                  <td class="py-3 px-6">
                    <div class="skeleton skeleton-button" style="width: 24px; height: 24px" />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
    <!-- Variant Content -->
    <div v-else class="card h-full card-lg">
      <div class="p-6">
        <div class="flex justify-between flex-row items-center">
          <div>
            <h3 class="text-md">Variant</h3>
          </div>
          <div>
            <button
              type="button"
              class="btn inline-flex items-center gap-x-2 bg-cyan-600 text-white border-cyan-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-cyan-700 hover:border-cyan-700 active:bg-cyan-700 active:border-cyan-700 focus:outline-none focus:ring-4 focus:ring-cyan-100"
              data-bs-toggle="modal"
              data-bs-target="#variantModal"
              :disabled="useProduct.loading"
            >
              Add Variant
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
                class="icon icon-tabler icons-tabler-outline icon-tabler-plus"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M12 5l0 14" />
                <path d="M5 12l14 0" />
              </svg>
            </button>
          </div>
        </div>
      </div>
      <div class="card-body p-0">
        <div class="relative overflow-x-auto">
          <table class="text-left w-full whitespace-nowrap table-with-checkbox">
            <thead class="bg-gray-200 text-gray-700">
              <tr class="border-transparent border-b-0!">
                <th scope="col" class="px-6 py-3">
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
                </th>
                <th scope="col" class="px-6 py-3">Stock</th>
                <th scope="col" class="px-6 py-3">Weight</th>
                <th scope="col" class="px-6 py-3">Reguler Price</th>
                <th scope="col" class="px-6 py-3">Sale Price</th>
                <th scope="col" class="px-6 py-3">Product Image</th>

                <th scope="col" class="px-6 py-3">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="variant in variants"
                :key="variant.id"
                class="border-transparent border-b-0!"
              >
                <td class="py-3 px-6 text-left">
                  <div class="flex items-center">
                    <input
                      :disabled="useProduct.loading"
                      class="w-4 h-4 text-blue-600 bg-white border-gray-300 rounded focus:ring-blue-600 focus:outline-none focus:ring-2"
                      type="checkbox"
                    />
                  </div>
                </td>
                <td class="py-3 px-6 text-left">{{ variant.stock }}</td>
                <td class="py-3 px-6 text-left">{{ variant.weight }}</td>
                <td class="py-3 px-6 text-left">Rp {{ variant.regular_price }}</td>
                <td class="py-3 px-6 text-left">Rp {{ variant.sale_price }}</td>
                <td class="py-3 px-6 text-left">
                  <img
                    :src="variant.product_image"
                    class="h-12 w-12 object-cover rounded"
                    alt="Variant Image"
                  />
                </td>
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
                        <a
                          class="dropdown-item"
                          href="#"
                          @click.prevent="deleteVariant(variant.id)"
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
    <div
      id="variantModal"
      class="modal fade"
      tabindex="-1"
      aria-labelledby="variantModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog modal-lg modal-dialog-centered">
        <div class="modal-content p-6 flex flex-col gap-6">
          <div class="flex flex-row items-center justify-between">
            <h5 id="addressLabel" class="modal-title">Add Variant</h5>
            <button
              :disabled="useProduct.loading"
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="icon icon-tabler icon-tabler-x text-gray-700"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                stroke-width="2"
                stroke="currentColor"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M18 6l-12 12" />
                <path d="M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="modal-body p-0">
            <form
              class="grid grid-cols-12 needs-validation gap-3"
              novalidate
              @submit.prevent="addVariant"
            >
              <div class="lg:col-span-6 col-span-12">
                <label for="customerEditAdd" class="inline-block text-gray-800 font-medium mb-2">
                  Stock <span class="text-red-600">*</span>
                </label>
                <input
                  id="customerEditAdd"
                  v-model="variantForm.stock"
                  :disabled="useProduct.loading"
                  type="number"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  :class="{ 'border-red-500': vv$.stock.$error && vv$.stock.$dirty }"
                  placeholder="10"
                  required
                  @blur="vv$.stock.$touch()"
                />
                <div v-if="vv$.stock.$error && vv$.stock.$dirty" class="text-red-600 text-sm mt-1">
                  <template v-if="vv$.stock.required.$invalid">Stock harus diisi</template>
                  <template v-else-if="vv$.stock.gte.$invalid">Stock minimal 0</template>
                </div>
              </div>

              <div class="lg:col-span-6 col-span-12">
                <label for="customerZip" class="inline-block text-gray-800 font-medium mb-2"
                  >Weight <span class="text-red-600">*</span></label
                >
                <input
                  id="customerZip"
                  v-model="variantForm.weight"
                  :disabled="useProduct.loading"
                  type="number"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  :class="{ 'border-red-500': vv$.weight.$error && vv$.weight.$dirty }"
                  placeholder="100"
                  required
                  @blur="vv$.weight.$touch()"
                />
                <div
                  v-if="vv$.weight.$error && vv$.weight.$dirty"
                  class="text-red-600 text-sm mt-1"
                >
                  <template v-if="vv$.weight.required.$invalid">Berat harus diisi</template>
                  <template v-else-if="vv$.weight.gte.$invalid">Berat minimal 0</template>
                </div>
              </div>

              <div class="lg:col-span-6 col-span-12">
                <label for="customerCity" class="inline-block text-gray-800 font-medium mb-2"
                  >Regular Price <span class="text-red-600">*</span></label
                >
                <input
                  id="customerCity"
                  v-model="variantForm.regular_price"
                  :disabled="useProduct.loading"
                  type="number"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  :class="{
                    'border-red-500': vv$.regular_price.$error && vv$.regular_price.$dirty
                  }"
                  placeholder="10000"
                  required
                  @blur="vv$.regular_price.$touch()"
                />
                <div
                  v-if="vv$.regular_price.$error && vv$.regular_price.$dirty"
                  class="text-red-600 text-sm mt-1"
                >
                  <template v-if="vv$.regular_price.required.$invalid"
                    >Harga reguler harus diisi</template
                  >
                  <template v-else-if="vv$.regular_price.gte.$invalid">Harga minimal 0</template>
                </div>
              </div>

              <div class="lg:col-span-6 col-span-12">
                <label for="customerState" class="inline-block text-gray-800 font-medium mb-2"
                  >Sale Price <span class="text-red-600">*</span></label
                >
                <input
                  id="customerState"
                  v-model="variantForm.sale_price"
                  :disabled="useProduct.loading"
                  type="number"
                  class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                  :class="{ 'border-red-500': vv$.sale_price.$error && vv$.sale_price.$dirty }"
                  placeholder="9000"
                  required
                  @blur="vv$.sale_price.$touch()"
                />
                <div
                  v-if="vv$.sale_price.$error && vv$.sale_price.$dirty"
                  class="text-red-600 text-sm mt-1"
                >
                  <template v-if="vv$.sale_price.required.$invalid"
                    >Harga jual harus diisi</template
                  >
                  <template v-else-if="vv$.sale_price.gte.$invalid">Harga jual minimal 0</template>
                </div>
              </div>

              <div class="lg:col-span-8 col-span-12">
                <label class="inline-block text-gray-800 font-medium mb-2">Image Product</label>
                <input
                  id="variantProductImage"
                  :disabled="useProduct.loading"
                  type="file"
                  class="sr-only"
                  accept="image/jpeg,image/png"
                  @change="handleFileUpload"
                />
                <label
                  for="variantProductImage"
                  class="flex min-h-[120px] cursor-pointer items-center justify-center border-2 border-dashed border-gray-400 bg-gray-50 px-4 py-8 text-center text-gray-600 transition hover:border-blue-600 hover:bg-blue-50"
                  :class="[
                    {
                      'border-red-500': vv$.product_image.$error && vv$.product_image.$dirty,
                      'pointer-events-none opacity-75': variantImageUploading
                    }
                  ]"
                  @dragover.prevent
                  @drop.prevent="handleFileDrop"
                >
                  <!-- Skeleton loading saat proses upload gambar berlangsung -->
                  <span
                    v-if="variantImageUploading"
                    class="flex w-full flex-col items-center gap-2"
                  >
                    <div
                      class="skeleton skeleton-image"
                      style="height: 64px; width: 64px; border-radius: 0.5rem"
                    />
                    <div class="skeleton skeleton-text w-48" style="margin-bottom: 0" />
                    <div class="skeleton skeleton-text-sm w-32" style="margin-bottom: 0" />
                  </span>
                  <span v-else-if="!variantForm.product_image">Drop file here to upload</span>
                  <img
                    v-else
                    :src="variantForm.product_image"
                    class="h-24 w-24 rounded object-cover"
                    alt="Variant image preview"
                  />
                </label>
                <div v-if="variantImageError" class="text-red-600 text-sm mt-1">
                  {{ variantImageError }}
                </div>
                <div
                  v-if="vv$.product_image.$error && vv$.product_image.$dirty"
                  class="text-red-600 text-sm mt-1"
                >
                  <template v-if="vv$.product_image.required.$invalid"
                    >Gambar produk harus diupload</template
                  >
                  <template v-else-if="vv$.product_image.url.$invalid"
                    >URL gambar tidak valid</template
                  >
                  <template v-else-if="vv$.product_image.maxLength.$invalid"
                    >URL gambar maksimal 255 karakter</template
                  >
                </div>
              </div>
            </form>
          </div>
          <div class="flex flex-row gap-3">
            <button
              type="button"
              :disabled="useProduct.loading"
              class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
              @click="addVariant"
            >
              Add
            </button>
            <button
              :disabled="useProduct.loading"
              type="button"
              class="btn inline-flex items-center gap-x-2 bg-gray-200 text-gray-800 border-gray-200 border disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useProductStore } from '~/stores/product'
import { ref, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useCategoryStore } from '~/stores/category'
import { useRouter, useRoute } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { required, minLength, maxLength, url } from '@vuelidate/validators'
import { oneof, gte } from '~/utils/validators'

definePageMeta({
  middleware: ['admin', 'auth'],
  layout: 'dashboard'
})

const router = useRouter()
const route = useRoute()
const useProduct = useProductStore()
const useCategory = useCategoryStore()
const categoryDatas = ref([])
const productId = route.params.id
const loading = ref(true)
const variantImageError = ref('')
const variantImageUploading = ref(false)
const { showError } = useErrorModal()

const variants = reactive([])

const form = reactive({
  name: '',
  category_slug: '',
  unit: '',
  description: '',
  status: '',
  variant: 0
})

// Sinkronkan jumlah variant ke form agar Vuelidate tetap reaktif
watch(
  () => variants.length,
  (val) => {
    form.variant = val
  },
  { immediate: true }
)

const variantForm = reactive({
  stock: 0,
  weight: 0,
  regular_price: 0,
  sale_price: 0,
  product_image: ''
})

const productRules = computed(() => ({
  name: {
    required,
    minLength: minLength(3),
    maxLength: maxLength(100)
  },
  category_slug: {
    required,
    maxLength: maxLength(120)
  },
  unit: {
    required,
    maxLength: maxLength(120)
  },
  description: {
    required
  },
  status: {
    required,
    maxLength: maxLength(120),
    oneof: oneof(['DRAFT', 'ACTIVE', 'INACTIVE'])
  },
  variant: {
    required,
    gte: gte(1)
  }
}))

const variantRules = computed(() => ({
  stock: {
    required,
    gte: gte(0)
  },
  weight: {
    required,
    gte: gte(0)
  },
  regular_price: {
    required,
    gte: gte(0)
  },
  sale_price: {
    required,
    gte: gte(0)
  },
  product_image: {
    required,
    url,
    maxLength: maxLength(255)
  }
}))

const v$ = useVuelidate(productRules, form)
const vv$ = useVuelidate(variantRules, variantForm)

onMounted(async () => {
  try {
    // Fetch categories
    await useCategory.fetchCategoriesAdmin()
    categoryDatas.value = useCategory.categories

    // Fetch product detail
    await useProduct.fetchProductDetail(productId)

    // Populate form with product data
    form.name = useProduct.product.product_name
    form.category_slug =
      categoryDatas.value.find((cat) => cat.name === useProduct.product.category_name)?.slug || ''
    form.unit = useProduct.product.unit || ''
    form.description = useProduct.product.product_description
    form.status = useProduct.product.product_status || 'DRAFT'

    variants.push({
      id: useProduct.product.id,
      stock: useProduct.product.stock,
      weight: useProduct.product.weight,
      regular_price: useProduct.product.regular_price,
      sale_price: useProduct.product.sale_price,
      product_image: useProduct.product.product_image
    })

    if (useProduct.product.child.length > 0) {
      useProduct.product.child.forEach((child) => {
        variants.push({
          stock: child.stock,
          weight: child.weight,
          regular_price: child.regular_price,
          sale_price: child.sale_price,
          product_image: child.product_image,
          id: child.id
        })
      })
    }
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  useProduct.cancelRequests()
  useCategory.cancelRequests()
})

const addVariant = async () => {
  vv$.value.$touch()
  if (vv$.value.$invalid) {
    return
  }

  variants.push({
    stock: variantForm.stock,
    weight: variantForm.weight,
    regular_price: variantForm.regular_price,
    sale_price: variantForm.sale_price,
    product_image: variantForm.product_image,
    id: Date.now() // temporary id
  })

  // Reset form
  variantForm.stock = 0
  variantForm.weight = 0
  variantForm.regular_price = 0
  variantForm.sale_price = 0
  variantForm.product_image = ''
  vv$.value.$reset()
}

const deleteVariant = (id) => {
  const isConfirmed = window.confirm('Apakah Anda yakin ingin menghapus variant ini?')
  if (isConfirmed) {
    const index = variants.findIndex((variant) => variant.id === id)
    if (index !== -1) {
      variants.splice(index, 1)
    }
  }
}

const uploadVariantImage = async (file) => {
  try {
    if (!file) return

    variantImageError.value = ''

    if (file.size > 1024 * 1024) {
      variantImageError.value = 'Ukuran file maksimal 1MB'
      return
    }

    if (!['image/jpeg', 'image/png'].includes(file.type)) {
      variantImageError.value = 'File harus berupa JPG atau PNG'
      return
    }

    variantImageUploading.value = true
    const result = await useCategory.uploadImage(file)
    variantForm.product_image = result.data.image_url || result.data.imageUrl
  } catch (error) {
    showError(error)
  } finally {
    variantImageUploading.value = false
  }
}

const handleFileUpload = async (event) => {
  await uploadVariantImage(event.target.files[0])
  const input = document.getElementById('variantProductImage')
  if (input) input.value = ''
}

const handleFileDrop = async (event) => {
  await uploadVariantImage(event.dataTransfer.files[0])
  const input = document.getElementById('variantProductImage')
  if (input) input.value = ''
}

const handleSubmit = async () => {
  v$.value.variant.$touch()
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    const productData = {
      product_name: form.name,
      category_slug: form.category_slug,
      unit: form.unit,
      product_description: form.description,
      status: form.status,
      variant_detail: variants,
      variant: variants.length
    }

    await useProduct.updateProduct(productId, productData)

    router.push('/dashboard/products')
  } catch (error) {
    showError(error)
  }
}
</script>
