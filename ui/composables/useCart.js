import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'

export const useCart = () => {
  const cartItems = ref([])
  const authStore = useAuthStore()

  // Load cart dari localStorage saat composable diinisialisasi
  const initCart = () => {
    if (process.client) {
      const savedCart = localStorage.getItem('cart')
      if (savedCart) {
        cartItems.value = JSON.parse(savedCart)
      }
    }
  }

  // Simpan cart ke localStorage
  const saveCart = () => {
    if (process.client) {
      localStorage.setItem('cart', JSON.stringify(cartItems.value))
    }
  }

  // Fungsi untuk menambah item ke cart
  const addToCart = async (product, quantity = 1) => {
    if (!authStore.isAuthenticated) {
      if (process.client) {
        // Simpan item yang ingin ditambahkan ke localStorage
        localStorage.setItem('pendingCartItem', JSON.stringify({ product, quantity }))
      }

      return {
        success: false,
        message: 'Silakan login terlebih dahulu',
        showLoginModal: true
      }
    }

    // Cek apakah produk sudah ada di cart
    const existingItem = cartItems.value.find(
      (item) => item.product.id === product.id && item.variant?.id === product.selectedVariant?.id
    )

    if (existingItem) {
      existingItem.quantity += quantity
    } else {
      cartItems.value.push({
        product,
        quantity,
        variant: product.selectedVariant || null
      })
    }

    // Simpan ke localStorage
    saveCart()

    return {
      success: true,
      message: 'Produk berhasil ditambahkan ke keranjang'
    }
  }

  // Fungsi untuk mengambil pending cart item setelah login
  const processPendingCartItem = () => {
    if (process.client) {
      const pendingItem = localStorage.getItem('pendingCartItem')
      if (pendingItem) {
        const { product, quantity } = JSON.parse(pendingItem)
        localStorage.removeItem('pendingCartItem')
        return addToCart(product, quantity)
      }
    }
    return null
  }

  // Fungsi untuk menghapus item dari cart
  const removeFromCart = (productId, variantId = null) => {
    cartItems.value = cartItems.value.filter(
      (item) => !(item.product.id === productId && item.variant?.id === variantId)
    )
    saveCart()
  }

  // Fungsi untuk mengupdate quantity
  const updateQuantity = (productId, variantId = null, quantity) => {
    const item = cartItems.value.find(
      (item) => item.product.id === productId && item.variant?.id === variantId
    )
    if (item) {
      item.quantity = quantity
      saveCart()
    }
  }

  // Inisialisasi cart hanya di client-side
  if (process.client) {
    initCart()
  }

  return {
    cartItems,
    addToCart,
    removeFromCart,
    updateQuantity,
    processPendingCartItem
  }
}
