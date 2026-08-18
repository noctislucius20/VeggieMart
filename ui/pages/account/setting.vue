<template>
  <main>
    <!-- Offcanvas for mobile -->
    <div
      id="offcanvasAccount"
      class="offcanvas offcanvas-start"
      tabindex="-1"
      aria-labelledby="offcanvasAccountLabel"
    >
      <div class="offcanvas-header border-bottom border-gray-300">
        <div class="flex items-center gap-2 w-full">
          <h5 id="offcanvasAccountLabel" class="offcanvas-title mb-0">Menu</h5>
          <button
            type="button"
            class="btn btn-sm ms-auto p-0 border-0 shadow-none text-gray-600 hover:text-gray-900"
            data-bs-dismiss="offcanvas"
            aria-label="Close"
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
              class="icon icon-tabler icons-tabler-outline icon-tabler-x"
            >
              <path stroke="none" d="M0 0h24v24H0z" fill="none" />
              <path d="M18 6l-12 12" />
              <path d="M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      <div class="offcanvas-body p-0">
        <ul class="nav flex-col nav-pills nav-pills-dark p-4">
          <!-- nav item -->
          <li class="nav-item">
            <a class="nav-link gap-2" href="/account">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="icon icon-tabler icons-tabler-outline icon-tabler-shopping-bag"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path
                  d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304z"
                />
                <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
              </svg>
              Your Orders
            </a>
          </li>

          <!-- nav item -->
          <li class="nav-item">
            <a class="nav-link gap-2" href="/account/notifications">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="icon icon-tabler icons-tabler-outline icon-tabler-bell"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3l-1 4h-10l-1 -4v-3a7 7 0 0 1 4 -6z" />
              </svg>
              Notifications
            </a>
          </li>

          <!-- nav item -->
          <li class="nav-item">
            <a class="nav-link gap-2 active" href="/account/setting">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="icon icon-tabler icons-tabler-outline icon-tabler-settings"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path
                  d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"
                />
                <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
              </svg>
              Settings
            </a>
          </li>
          <!-- nav item -->

          <li class="nav-item">
            <hr class="pt-3 mt-5" />
          </li>
          <!-- nav item -->
          <li class="nav-item">
            <button
              :disabled="authStore.loading"
              type="button"
              class="nav-link gap-2"
              @click="handleLogout"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="icon icon-tabler icons-tabler-outline icon-tabler-logout"
              >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path
                  d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2"
                />
                <path d="M9 12h12l-3 -3" />
                <path d="M18 15l3 -3" />
              </svg>
              Log out
            </button>
          </li>
        </ul>
      </div>
    </div>
    <!-- section -->
    <section>
      <div class="container">
        <div class="grid grid-cols-12">
          <div class="col-span-12">
            <div class="flex justify-between items-center md:hidden py-4">
              <!-- heading -->
              <h3 class="text-md mb-0">Account Setting</h3>
              <!-- button -->
              <button
                :disabled="authStore.loading"
                class="btn inline-flex items-center gap-x-2 bg-transparent text-gray-600 border-gray-300 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300 md:hidden btn-sm ms-3"
                type="button"
                data-bs-toggle="offcanvas"
                data-bs-target="#offcanvasAccount"
                aria-controls="offcanvasAccount"
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
                  class="icon icon-tabler icons-tabler-outline icon-tabler-indent-increase"
                >
                  <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                  <path d="M20 6l-11 0" />
                  <path d="M20 12l-7 0" />
                  <path d="M20 18l-11 0" />
                  <path d="M4 8l4 4l-4 4" />
                </svg>
              </button>
            </div>
          </div>
          <div class="lg:col-span-2 md:col-span-4 col-span-12 border-r hidden md:block">
            <div class="pt-10 md:pr-10">
              <!-- nav -->
              <ul class="nav flex-col nav-pills nav-pills-dark">
                <!-- nav item -->
                <li class="nav-item">
                  <a class="nav-link gap-2" href="./">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      class="icon icon-tabler icons-tabler-outline icon-tabler-shopping-bag"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M6.331 8h11.339a2 2 0 0 1 1.977 2.304l-1.255 8.152a3 3 0 0 1 -2.966 2.544h-6.852a3 3 0 0 1 -2.965 -2.544l-1.255 -8.152a2 2 0 0 1 1.977 -2.304z"
                      />
                      <path d="M9 11v-5a3 3 0 0 1 6 0v5" />
                    </svg>
                    Your Orders
                  </a>
                </li>

                <!-- nav item -->
                <li class="nav-item">
                  <a class="nav-link gap-2" href="./notifications">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      class="icon icon-tabler icons-tabler-outline icon-tabler-bell"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3l-1 4h-10l-1 -4v-3a7 7 0 0 1 4 -6z"
                      />
                    </svg>
                    Notifications
                  </a>
                </li>

                <!-- nav item -->
                <li class="nav-item">
                  <a class="nav-link gap-2 active" href="./setting">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      class="icon icon-tabler icons-tabler-outline icon-tabler-settings"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"
                      />
                      <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
                    </svg>
                    Settings
                  </a>
                </li>
                <!-- nav item -->

                <li class="nav-item">
                  <hr class="pt-3 mt-5" />
                </li>
                <!-- nav item -->
                <li class="nav-item">
                  <button
                    :disabled="authStore.loading"
                    type="button"
                    class="nav-link gap-2"
                    @click="handleLogout"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      class="icon icon-tabler icons-tabler-outline icon-tabler-logout"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path
                        d="M14 8v-2a2 2 0 0 0 -2 -2h-7a2 2 0 0 0 -2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2 -2v-2"
                      />
                      <path d="M9 12h12l-3 -3" />
                      <path d="M18 15l3 -3" />
                    </svg>
                    Log out
                  </button>
                </li>
              </ul>
            </div>
          </div>
          <div class="lg:col-span-10 md:col-span-8 col-span-12">
            <div class="py-6 md:p-6 lg:p-10">
              <div class="mb-6">
                <!-- heading -->
                <h2 class="mb-0 text-lg">Account Setting</h2>
              </div>

              <!-- Loading Skeleton -->
              <div v-if="loading">
                <div class="grid grid-cols-12">
                  <div class="col-span-5">
                    <div class="skeleton skeleton-title w-48 mb-4" style="height: 24px" />
                    <div class="mb-3">
                      <div class="skeleton skeleton-text w-20 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-text w-full" style="height: 44px" />
                    </div>
                    <div class="mb-3">
                      <div class="skeleton skeleton-text w-20 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-text w-full" style="height: 44px" />
                    </div>
                    <div class="mb-5">
                      <div class="skeleton skeleton-text w-20 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-text w-full" style="height: 44px" />
                    </div>
                    <div class="mb-5">
                      <div class="skeleton skeleton-text w-20 mb-2" style="height: 18px" />
                      <div class="skeleton skeleton-text w-full" style="height: 100px" />
                    </div>
                    <div class="mb-12">
                      <div class="skeleton skeleton-text w-20 mb-2" style="height: 18px" />
                      <div class="flex items-center gap-2">
                        <div class="skeleton rounded-lg" style="width: 150px; height: 150px" />
                        <div class="skeleton skeleton-button" style="width: 120px; height: 44px" />
                      </div>
                    </div>
                    <div class="mb-3">
                      <div class="skeleton skeleton-button" style="width: 150px; height: 44px" />
                    </div>
                  </div>
                </div>
              </div>
              <!-- Form Content -->
              <div v-else>
                <!-- heading -->
                <h5 class="mb-4">Account details</h5>
                <div class="grid grid-cols-12">
                  <div class="sm:col-span-12 xl:col-span-5">
                    <!-- form -->
                    <form enctype="multipart/form-data" @submit.prevent="handleSubmit">
                      <!-- input -->
                      <div class="mb-3">
                        <label class="inline-block text-gray-800 font-medium mb-2">Name</label>
                        <input
                          v-model="formData.name"
                          :disabled="authStore.loading"
                          type="text"
                          class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                          :class="{ 'border-red-500': v$.name.$error && v$.name.$dirty }"
                          placeholder="Name"
                          @blur="v$.name.$touch()"
                        />
                        <div
                          v-if="v$.name.$error && v$.name.$dirty"
                          class="text-red-600 text-sm mt-1"
                        >
                          <template v-if="v$.name.required.$invalid">Nama harus diisi</template>
                          <template v-else-if="v$.name.minLength.$invalid"
                            >Nama minimal 3 karakter</template
                          >
                          <template v-else-if="v$.name.maxLength.$invalid"
                            >Nama maksimal 255 karakter</template
                          >
                        </div>
                      </div>
                      <!-- input -->
                      <div class="mb-3">
                        <label class="inline-block text-gray-800 font-medium mb-2">Email</label>
                        <input
                          v-model="formData.email"
                          :disabled="authStore.loading"
                          type="email"
                          class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                          :class="{ 'border-red-500': v$.email.$error && v$.email.$dirty }"
                          placeholder="Email"
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
                      <!-- input -->
                      <div class="mb-5">
                        <label class="inline-block text-gray-800 font-medium mb-2">Phone</label>
                        <input
                          v-model="formData.phone"
                          :disabled="authStore.loading"
                          type="text"
                          class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                          :class="{ 'border-red-500': v$.phone.$error && v$.phone.$dirty }"
                          placeholder="Phone number"
                          @blur="v$.phone.$touch()"
                        />
                        <div
                          v-if="v$.phone.$error && v$.phone.$dirty"
                          class="text-red-600 text-sm mt-1"
                        >
                          <template v-if="v$.phone.maxLength.$invalid"
                            >Nomor telepon maksimal 17 karakter</template
                          >
                          <template v-if="v$.phone.numeric.$invalid"
                            >Nomor telepon harus berupa angka</template
                          >
                        </div>
                      </div>
                      <div class="mb-5">
                        <label class="inline-block text-gray-800 font-medium mb-2">Address</label>
                        <textarea
                          v-model="formData.address"
                          :disabled="authStore.loading"
                          class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                          placeholder="Address"
                        />
                      </div>
                      <div class="mb-12">
                        <label class="inline-block text-gray-800 font-medium mb-2">Photo</label>
                        <div
                          v-if="v$.photo.$error && v$.photo.$dirty"
                          class="text-red-600 text-sm mb-2"
                        >
                          <template v-if="v$.photo.maxLength.$invalid"
                            >Nama file foto maksimal 255 karakter</template
                          >
                        </div>
                        <div v-if="photoError" class="text-red-600 text-sm mb-2">
                          {{ photoError }}
                        </div>
                        <div class="flex items-center gap-4">
                          <div>
                            <img
                              v-if="imageUrl"
                              style="width: 150px; height: 150px; object-fit: cover"
                              class="image rounded-lg"
                              :src="imageUrl"
                              alt="Image"
                            />
                            <img
                              v-else
                              style="width: 150px; height: 150px; object-fit: cover"
                              class="image rounded-lg"
                              src="../../assets/images/avatar/avatar.jpg"
                            />
                          </div>
                          <label
                            :aria-disabled="authStore.loading"
                            :class="[
                              'file-upload btn inline-flex items-center gap-x-2 bg-gray-200 text-gray-800 border-gray-200 border hover:text-white hover:bg-gray-700 hover:border-gray-700 active:bg-gray-700 active:border-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-300 cursor-pointer p-3 rounded-lg',
                              authStore.loading ? 'opacity-50 pointer-events-none' : ''
                            ]"
                          >
                            <input
                              type="file"
                              class="file-input hidden"
                              accept="image/jpeg,image/png"
                              :disabled="authStore.loading"
                              @change="handleFileUpload"
                            />
                            Upload Photo
                          </label>
                        </div>
                      </div>
                      <!-- button -->
                      <div class="mb-3">
                        <button
                          :disabled="authStore.loading"
                          type="submit"
                          class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100"
                        >
                          <span v-if="authStore.loading">Loading...</span>
                          <span v-else>Save</span>
                        </button>
                      </div>
                    </form>
                  </div>
                </div>
              </div>
              <!-- <hr class="my-10">
              <div class="lg:pr-14">
                <h5 class="mb-4">Password</h5>
                <form class="grid grid-cols-1 lg:grid-cols-2 gap-x-4">
                  <div class="mb-3">
                    <label class="inline-block text-gray-800 font-medium mb-2">New Password</label>
                    <input :disabled="authStore.loading" type="password"
                      class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                      placeholder="**********">
                  </div>
                  <div class="mb-3">
                    <label class="inline-block text-gray-800 font-medium mb-2">Current
                      Password</label>
                    <input :disabled="authStore.loading" type="password"
                      class="border border-gray-300 text-gray-900 rounded-lg focus:shadow-[0_0_0_.25rem_rgba(37,99,235,.25)] focus:ring-blue-600 focus:ring-0 focus:border-blue-600 block p-2 px-3 disabled:opacity-50 disabled:pointer-events-none w-full text-base"
                      placeholder="**********">
                  </div>
                  <div class="">
                    <p class="mb-4">
                      Can't remember your current password?
                      <a href="#" class="text-blue-600">Reset your password.</a>
                    </p>
                    <a href="#"
                      class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-100">
                      Save Password
                    </a>
                  </div>
                </form>
              </div>
              <hr class="my-10">
              <div>
                <h5 class="mb-4">Delete Account</h5>
                <p class="mb-2">Would you like to delete your account?</p>
                <p class="mb-5">This account contain 12 orders, Deleting your account will remove all
                  the order details associated with it.</p>
                <a href="#"
                  class="btn inline-flex items-center gap-x-2 bg-transparent text-red-600 border-red-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-red-600 hover:border-red-600 active:bg-red-600 active:border-red-600 focus:outline-none focus:ring-4 focus:ring-red-100">
                  I want to delete my account
                </a>
              </div> -->
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useCustomerStore } from '~/stores/customer'
import { useRouter } from 'vue-router'
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength, maxLength } from '@vuelidate/validators'
import { numeric } from '~/utils/validators'

definePageMeta({
  middleware: ['auth']
})

const router = useRouter()
const authStore = useAuthStore()
const customerStore = useCustomerStore()

const latitude = ref(0)
const longitude = ref(0)
const imageUrl = ref('')
const loading = ref(true)
const photoError = ref('')

const { showError } = useErrorModal()

const formData = reactive({
  name: '',
  email: '',
  phone: '',
  address: '',
  lat: null,
  lng: null,
  photo: ''
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
  phone: {
    maxLength: maxLength(17),
    numeric
  },
  photo: {
    maxLength: maxLength(255)
  }
}))

const v$ = useVuelidate(rules, formData)

const handleFileUpload = async (event) => {
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
    const result = await customerStore.uploadImage(selectedFile)

    // Update preview dan simpan URL
    imageUrl.value = result.data.image_url
    formData.photo = result.data.image_url
  }
}

const handleLogout = () => {
  try {
    authStore.logout()

    useAuthStore().$reset()
    useCartStore().$reset()
    useCategoryStore().$reset()
    useCustomerStore().$reset()
    useNotificationStore().$reset()
    useOrderStore().$reset()
    usePaymentStore().$reset()
    useProductStore().$reset()
    useRoleStore().$reset()

    router.push('/auth/signin')
  } catch (error) {
    showError(error)
  }
}

const handleSubmit = async () => {
  v$.value.$touch()
  if (v$.value.$invalid) {
    return
  }

  try {
    const userData = {
      name: formData.name,
      email: formData.email,
      phone: formData.phone,
      address: formData.address,
      lat: latitude.value,
      lng: longitude.value,
      photo: imageUrl.value
    }

    await authStore.updateProfile(userData)

    // router.push('/account/setting')
  } catch (error) {
    showError(error)
  }
}

onMounted(async () => {
  try {
    await authStore.getProfile()

    if (import.meta.client && 'geolocation' in navigator) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          latitude.value = position.coords.latitude
          longitude.value = position.coords.longitude
        },
        (err) => {
          showError(err.message)
        }
      )
    } else {
      showError('Geolocation tidak didukung oleh browser')
    }

    if (authStore.user) {
      formData.name = authStore.user.name || ''
      formData.email = authStore.user.email || ''
      formData.phone = authStore.user.phone || ''
      formData.address = authStore.user.address || ''
      formData.lat = authStore.user.lat || 0
      formData.lng = authStore.user.lng || 0
      formData.photo = authStore.user.photo || ''

      imageUrl.value = formData.photo
    }
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  authStore.cancelRequests()
  customerStore.cancelRequests()
})
</script>

<style></style>
