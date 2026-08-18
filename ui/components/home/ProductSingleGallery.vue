<template>
  <div class="lg:w-1/2">
    <!-- Image slide -->
    <div class="product mb-6">
      <Swiper
        :modules="swiperModules"
        :navigation="true"
        :pagination="{ clickable: true }"
        :thumbs="{ swiper: thumbsSwiper }"
        class="product-images-slider"
        @swiper="setMainSwiper"
      >
        <SwiperSlide>
          <div
            class="zoom"
            :style="{ backgroundImage: `url(${defaultImage})` }"
            @mousemove="handleZoom"
            @mouseleave="handleMouseLeave"
          >
            <img
              :src="product?.image || defaultImage"
              :alt="product?.product_name || 'Product Image'"
              class="w-full"
            />
          </div>
        </SwiperSlide>
        <SwiperSlide v-for="(child, index) in product?.childs" :key="index">
          <div
            class="zoom"
            :style="{ backgroundImage: `url(${child.image || defaultImage})` }"
            @mousemove="handleZoom"
            @mouseleave="handleMouseLeave"
          >
            <img
              :src="child?.image || defaultImage"
              :alt="`${product?.product_name || 'Product'} - ${child.weight}${product?.unit}`"
              class="w-full"
            />
          </div>
        </SwiperSlide>
      </Swiper>
    </div>

    <!-- Product thumbnails -->
    <div class="product-tools">
      <Swiper
        :modules="swiperModules"
        :slides-per-view="4"
        :space-between="10"
        :free-mode="true"
        :watch-slides-progress="true"
        class="product-images-slider-thumbs"
        @swiper="setThumbsSwiper"
      >
        <SwiperSlide class="cursor-pointer">
          <img :src="defaultImage" :alt="product?.product_name || 'Product Image'" />
        </SwiperSlide>
        <SwiperSlide v-for="(child, index) in product?.childs" :key="index" class="cursor-pointer">
          <img
            :src="child?.image || defaultImage"
            :alt="`${product?.product_name || 'Product'} - ${child.weight}${product?.unit}`"
            class="w-full"
          />
        </SwiperSlide>
      </Swiper>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Swiper, SwiperSlide } from 'swiper/vue'
import { Navigation, Pagination, FreeMode, Thumbs } from 'swiper/modules'
import 'swiper/css'
import 'swiper/css/navigation'
import 'swiper/css/pagination'
import 'swiper/css/free-mode'
import 'swiper/css/thumbs'

const props = defineProps({
  product: {
    type: Object,
    required: false,
    default: () => ({
      product_name: '',
      child: [],
      unit: ''
    })
  }
})

const swiperModules = [Navigation, Pagination, FreeMode, Thumbs]
const thumbsSwiper = ref(null)
const mainSwiper = ref(null)

const setThumbsSwiper = (swiper) => {
  thumbsSwiper.value = swiper
}

const setMainSwiper = (swiper) => {
  mainSwiper.value = swiper
}

const handleZoom = (event) => {
  const zoomer = event.currentTarget
  const rect = zoomer.getBoundingClientRect()
  const offsetX = event.clientX - rect.left
  const offsetY = event.clientY - rect.top
  const x = (offsetX / rect.width) * 100
  const y = (offsetY / rect.height) * 100
  zoomer.style.backgroundPosition = `${x}% ${y}%`
  zoomer.style.backgroundSize = '200%'
}

const handleMouseLeave = (event) => {
  const zoomer = event.currentTarget
  zoomer.style.backgroundPosition = '50% 50%'
  zoomer.style.backgroundSize = 'cover'
}

const defaultImage = computed(() => {
  return props.product?.image || '/images/docs/placeholder-img.jpg'
})
</script>

<style scoped>
.zoom {
  position: relative;
  overflow: hidden;
  background-position: 50% 50%;
  background-size: cover;
  background-repeat: no-repeat;
  cursor: zoom-in;
  transition: background-size 0.1s ease;
}

.zoom:hover img {
  opacity: 0;
}

.zoom img {
  transition: opacity 0.3s ease;
  display: block;
  width: 100%;
  pointer-events: none;
}

:deep(.swiper-button-next),
:deep(.swiper-button-prev) {
  color: #16a34a;
  background-color: white;
  padding: 20px;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

:deep(.swiper-pagination-bullet-active) {
  background-color: #16a34a;
}

.product-images-slider-thumbs {
  padding: 10px 0;
}

.product-images-slider-thumbs :deep(.swiper-slide) {
  opacity: 0.5;
  transition: 0.3s;
}

.product-images-slider-thumbs :deep(.swiper-slide-thumb-active) {
  opacity: 1;
}
</style>
