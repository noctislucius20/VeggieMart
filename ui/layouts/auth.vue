<script setup lang="ts">
import Footer from '~/components/common/Footer.vue'

const route = useRoute()
const isSignup = computed(() => route.path === '/auth/signin')
</script>

<template>
  <main>
    <div class="border-b shadow-sm">
      <nav class="relative md:flex flex-wrap w-full items-center py-3 text-center">
        <div class="container sm:flex-center md:justify-between">
          <div class="lg:w-1/5 md:w-1/4 sm:w-1/2">
            <a class="inline-block mr-4 text-lg whitespace-no-wrap" href="/">
              <img
                src="~/assets/images/logo/dailymart-logo.svg"
                alt=""
                class="inline-block align-text-top"
              />
            </a>
          </div>

          <span v-if="!isSignup" class="inline-block">
            Already have an account?
            <a href="signin" class="text-blue-600">Sign in</a>
          </span>
        </div>
      </nav>
    </div>

    <slot />
    <Footer />
  </main>
</template>
