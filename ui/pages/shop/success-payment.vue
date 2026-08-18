<template>
  <section>
    <div class="container">
      <div class="flex flex-wrap min-h-screen justify-center items-center">
        <div class="lg:mx-1/6 lg:w-4/5 flex flex-col gap-10">
          <!-- <div class="flex flex-col justify-center lg:justify-start w-full">
            <a href="/" class="flex justify-center lg:justify-start"
              ><img src="~/assets/images/logo/dailymart-logo.svg" alt=""
            /></a>
          </div> -->
          <div class="flex-col lg:flex-row flex justify-center items-center gap-8">
            <!-- content -->
            <div class="lg:w-1/2 w-full">
              <div
                class="flex flex-col gap-8 justify-center lg:justify-start text-center lg:text-left"
              >
                <div class="flex flex-col gap-3">
                  <img
                    src="~/assets/images/svg-graphics/200-success.svg"
                    alt=""
                    class="max-w-full h-auto"
                  />
                  <h1 class="leading-tight text-center">Payment Success</h1>
                </div>
                <div class="text-right">
                  <!-- btn -->
                  <a
                    href="/"
                    class="btn inline-flex items-center gap-x-2 bg-blue-600 text-white border-blue-600 disabled:opacity-50 disabled:pointer-events-none hover:text-white hover:bg-blue-700 hover:border-blue-700 active:bg-blue-700 active:border-blue-700 focus:outline-none focus:ring-4 focus:ring-blue-300"
                  >
                    Back to home
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
                      class="icon icon-tabler icons-tabler-outline icon-tabler-arrow-narrow-right"
                    >
                      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                      <path d="M5 12l14 0" />
                      <path d="M15 16l4 -4" />
                      <path d="M15 8l4 4" />
                    </svg>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  setTimeout(() => {
    router.push('/')
  }, 5000)
})
</script>

<style></style>
