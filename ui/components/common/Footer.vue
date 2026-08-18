<template>
  <footer class="bg-gray-200 py-8">
    <div class="container">
      <!-- Categories and Links Section -->
      <div class="flex flex-wrap md:gap-4 lg:gap-0 py-4 mb-6">
        <!-- Categories Column -->
        <div class="w-full md:w-full lg:w-1/3 flex flex-col gap-4 mb-6">
          <h6>Categories</h6>
          <div class="flex flex-wrap">
            <div class="w-1/2">
              <ul class="flex flex-col gap-2">
                <li v-for="(category, index) in categoriesLeft" :key="index">
                  <NuxtLink :to="'/shop'" class="inline-block hover:text-blue-600">
                    {{ category }}
                  </NuxtLink>
                </li>
              </ul>
            </div>
            <div class="w-1/2">
              <ul class="flex flex-col gap-2">
                <li v-for="(category, index) in categoriesRight" :key="index">
                  <NuxtLink :to="'/shop'" class="inline-block hover:text-blue-600">
                    {{ category }}
                  </NuxtLink>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Other Links Columns -->
        <div class="w-full md:w-full lg:w-2/3">
          <div class="flex flex-wrap">
            <div
              v-for="(section, index) in footerSections"
              :key="index"
              class="w-1/2 sm:w-1/2 md:w-1/4 flex flex-col gap-4 mb-6"
            >
              <h6>{{ section.title }}</h6>
              <ul class="flex flex-col gap-2">
                <li v-for="(link, linkIndex) in section.links" :key="linkIndex">
                  <NuxtLink :to="link.url || '#'" class="inline-block hover:text-blue-600">
                    {{ link.text }}
                  </NuxtLink>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Payment Partners Section -->
      <div class="border-t py-4 border-gray-300">
        <div class="gap-y-4 flex flex-wrap items-center justify-center lg:justify-start">
          <div class="lg:w-2/5 lg:text-left text-center">
            <div class="flex md:flex-row flex-col gap-3 md:gap-6 items-center">
              <div class="text-gray-900">Payment Partners</div>
              <ul class="flex items-center flex-row gap-4">
                <li v-for="(payment, index) in paymentMethods" :key="index">
                  <img :src="`/images/payment/${payment}.svg`" :alt="payment" class="h-6" />
                </li>
              </ul>
            </div>
          </div>
          <div class="lg:w-3/5 flex justify-end">
            <div class="flex flex-col md:flex-row items-center gap-3 md:gap-6">
              <div class="text-gray-900">Get deliveries with DailyMart</div>
              <ul class="flex flex-row gap-2">
                <li>
                  <img
                    src="~/assets/images/appbutton/appstore-btn.svg"
                    alt="App Store"
                    class="h-10"
                  />
                </li>
                <li>
                  <img
                    src="~/assets/images/appbutton/googleplay-btn.svg"
                    alt="Google Play"
                    class="h-10"
                  />
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Copyright Section -->
      <div class="border-t py-4 border-gray-300">
        <div class="flex flex-col md:flex-row items-center gap-3">
          <div class="md:w-1/2">
            <span class="text-sm text-gray-500">
              © {{ new Date().getFullYear() }} DailyMart TailwindCSS eCommerce HTML Template.
              Powered by
              <a href="https://codescandy.com/" target="_blank" class="text-blue-600">Codescandy</a
              >.
            </span>
          </div>
          <div class="md:w-1/2 flex md:justify-end items-center">
            <div class="flex flex-row gap-5 items-center">
              <div class="text-gray-500">Follow us on</div>
              <ul class="flex items-center justify-end text-sm gap-1">
                <li v-for="(social, index) in socialMedia" :key="index">
                  <a
                    :href="social.url"
                    class="inline-flex justify-center items-center align-middle text-center select-none border font-normal whitespace-no-wrap rounded leading-normal no-underline h-8 w-8 border-gray-300 hover:border-blue-600 hover:text-blue-600 transition ease-in-out"
                  >
                    <Icon :name="social.icon" size="16" />
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
const categoriesLeft = [
  'Vegetables & Fruits',
  'Breakfast & instant food',
  'Bakery & Biscuits',
  'Atta, rice & dal',
  'Sauces & spreads',
  'Organic & gourmet',
  'Baby care',
  'Cleaning essentials',
  'Personal care'
]

const categoriesRight = [
  'Dairy, bread & eggs',
  'Cold drinks & juices',
  'Tea, coffee & drinks',
  'Masala, oil & more',
  'Chicken, meat & fish',
  'Paan corner',
  'Pharma & wellness',
  'Home & office',
  'Pet care'
]

const footerSections = [
  {
    title: 'Get to know us',
    links: [
      // { text: 'Company', url: '/about' },
      // { text: 'About', url: '/about' },
      // { text: 'Blog', url: '/blog' },
      // { text: 'Help Center', url: '/help' },
      // { text: 'Our Value', url: '/values' }
      { text: 'Company', url: '#' },
      { text: 'About', url: '#' },
      { text: 'Blog', url: '#' },
      { text: 'Help Center', url: '#' },
      { text: 'Our Value', url: '#' }
    ]
  },
  {
    title: 'For Consumers',
    links: [
      // { text: 'Payments', url: '/payments' },
      // { text: 'Shipping', url: '/shipping' },
      // { text: 'Product Returns', url: '/returns' },
      // { text: 'FAQ', url: '/faq' },
      // { text: 'Shop Checkout', url: '/checkout' }
      { text: 'Payments', url: '#' },
      { text: 'Shipping', url: '#' },
      { text: 'Product Returns', url: '#' },
      { text: 'FAQ', url: '#' },
      { text: 'Shop Checkout', url: '#' }
    ]
  },
  {
    title: 'Become a Shopper',
    links: [
      // { text: 'Shopper Opportunities', url: '/shopper' },
      // { text: 'Become a Shopper', url: '/shopper/register' },
      // { text: 'Earnings', url: '/shopper/earnings' },
      // { text: 'Ideas & Guides', url: '/guides' },
      // { text: 'New Retailers', url: '/retailers' }
      { text: 'Shopper Opportunities', url: '#' },
      { text: 'Become a Shopper', url: '#' },
      { text: 'Earnings', url: '#' },
      { text: 'Ideas & Guides', url: '#' },
      { text: 'New Retailers', url: '#' }
    ]
  },
  {
    title: 'Dailymart programs',
    links: [
      // { text: 'Dailymart programs', url: '/programs' },
      // { text: 'Gift Cards', url: '/gift-cards' },
      // { text: 'Promos & Coupons', url: '/coupons' },
      // { text: 'Dailymart Ads', url: '/ads' },
      // { text: 'Careers', url: '/careers' }
      { text: 'Dailymart programs', url: '#' },
      { text: 'Gift Cards', url: '#' },
      { text: 'Promos & Coupons', url: '#' },
      { text: 'Dailymart Ads', url: '#' },
      { text: 'Careers', url: '#' }
    ]
  }
]

const paymentMethods = ['amazonpay', 'american-express', 'mastercard', 'paypal', 'visa']

const socialMedia = [
  { icon: 'tabler:brand-facebook', url: '#' },
  { icon: 'tabler:brand-x', url: '#' },
  { icon: 'tabler:brand-instagram', url: '#' }
]
</script>
