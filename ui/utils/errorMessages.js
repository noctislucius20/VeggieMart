// utils/errorMessages.js
// Mapping error code → UI Message (dari error_mapping.md)

// Mapping error code → UI Message
const ERROR_MESSAGES = {
  // 400 Bad Request
  ID_REQUIRED: 'ID wajib diisi.',
  SLUG_REQUIRED: 'Slug wajib diisi.',
  PRODUCT_ID_REQUIRED: 'ID produk wajib diisi.',
  IDEMPOTENCY_KEY_REQUIRED: 'Terjadi kesalahan. Silakan coba lagi.',
  ORDER_CODE_REQUIRED: 'Kode pesanan wajib diisi.',
  LAT_OR_LNG_REQUIRED: 'Lokasi wajib diisi.',

  // 401 Unauthorized
  SESSION_EXPIRED: 'Sesi Anda telah berakhir. Silakan masuk kembali.',
  TOKEN_INVALID: 'Sesi tidak valid. Silakan masuk kembali.',
  TOKEN_EXPIRED: 'Sesi Anda telah kedaluwarsa. Silakan masuk kembali.',
  LOGIN_INVALID: 'Email atau kata sandi salah.',

  // 403 Forbidden
  ACCESS_FORBIDDEN: 'Anda tidak memiliki akses ke halaman ini.',
  GATEWAY_REQUIRED: 'Akses ditolak.',
  GATEWAY_SECRET_INVALID: 'Akses ditolak.',
  SERVICE_NOT_ALLOWED: 'Akses ditolak.',
  SERVICE_SECRET_INVALID: 'Akses ditolak.',

  // 404 Not Found
  DATA_NOT_FOUND: 'Data tidak ditemukan.',

  // 409 Conflict
  EMAIL_ALREADY_EXISTS: 'Email sudah terdaftar.',
  EMAIL_NOT_VERIFIED: 'Email belum diverifikasi. Silakan cek email Anda.',
  DATA_ALREADY_EXISTS: 'Data sudah ada.',
  DATA_STILL_IN_USED: 'Data masih digunakan dan tidak dapat dihapus.',
  STOCK_UNAVAILABLE: 'Stok produk tidak mencukupi.',
  INVALID_STATUS_TRANSITION: 'Status tidak dapat diubah.',
  REQUEST_PROCESSING: 'Permintaan sedang diproses. Silakan tunggu.',

  // 422 Unprocessable Entity
  RELATION_DATA_NOT_FOUND: 'Data terkait tidak ditemukan.',
  ID_INVALID: 'ID tidak valid.',
  PRODUCT_ID_INVALID: 'ID produk tidak valid.',
  QUANTITY_INVALID: 'Jumlah produk tidak valid.',
  PRICE_RANGE_INVALID: 'Rentang harga tidak valid.',
  LAT_OR_LNG_INVALID: 'Lokasi tidak valid.',
  DISTANCE_TOO_FAR: 'Jarak pengiriman terlalu jauh.',
  INVALID_PAYMENT_METHOD: 'Metode pembayaran tidak valid.',
  INVALID_NOTIFICATION_METHOD: 'Metode notifikasi tidak valid.',

  // 429 Too Many Requests
  TOO_MANY_REQUESTS: 'Terlalu banyak permintaan. Silakan coba lagi nanti.',

  // 500 Internal Server Error
  INTERNAL_SERVER_ERROR: 'Terjadi kesalahan pada server. Silakan coba lagi.',

  // 503 Service Unavailable
  SERVICE_UNAVAILABLE: 'Layanan sedang tidak tersedia. Silakan coba lagi nanti.',

  // 504 Gateway Timeout
  TIMEOUT_LIMIT_EXCEEDED: 'Waktu permintaan habis. Silakan coba lagi.'
}

// Fallback message berdasarkan HTTP status code
const HTTP_STATUS_MESSAGES = {
  400: 'Permintaan tidak valid. Silakan periksa kembali input Anda.',
  401: 'Sesi tidak valid. Silakan masuk kembali.',
  403: 'Anda tidak memiliki akses ke halaman ini.',
  404: 'Data tidak ditemukan.',
  409: 'Terjadi konflik data. Silakan coba lagi.',
  422: 'Data yang Anda masukkan tidak valid.',
  429: 'Terlalu banyak permintaan. Silakan coba lagi nanti.',
  500: 'Terjadi kesalahan pada server. Silakan coba lagi.',
  502: 'Terjadi kesalahan pada server. Silakan coba lagi.',
  503: 'Layanan sedang tidak tersedia. Silakan coba lagi nanti.',
  504: 'Waktu permintaan habis. Silakan coba lagi.'
}

// Literal error messages yang langsung ditampilkan dari handler
const LITERAL_MESSAGES = {
  'New Password and Confirm Password do not match': 'Kata sandi baru dan konfirmasi tidak cocok.',
  'Password and Confirm Password do not match': 'Kata sandi dan konfirmasi tidak cocok.',
  'file size exceeds maximum limit of 1 MB': 'Ukuran file melebihi batas maksimal 1 MB.',
  'only JPG and PNG files are allowed': 'Hanya file JPG dan PNG yang diperbolehkan.',
  'invalid target url': 'Terjadi kesalahan pada server. Silakan coba lagi.',
  'failed to connect to backend WebSocket': 'Koneksi gagal. Silakan coba lagi.',
  'failed to read request body': 'Terjadi kesalahan pada server. Silakan coba lagi.',
  'Failed to create request': 'Terjadi kesalahan pada server. Silakan coba lagi.',
  'Failed to forward request': 'Terjadi kesalahan pada server. Silakan coba lagi.',
  'Failed to read response body': 'Terjadi kesalahan pada server. Silakan coba lagi.'
}

export const getErrorMessage = (message) => {
  if (!message) return 'Terjadi kesalahan. Silakan coba lagi.'

  // 1. Cek literal message terlebih dahulu
  if (LITERAL_MESSAGES[message]) {
    return LITERAL_MESSAGES[message]
  }

  // 2. Cek format "HTTP_CODE: ERROR_CODE" (contoh: "401: TOKEN_INVALID")
  const match = message.match(/^(\d{3}):\s*(.+)$/)
  if (match) {
    const httpCode = match[1]
    const errorCode = match[2]

    // Cari di mapping error code
    if (ERROR_MESSAGES[errorCode]) {
      return ERROR_MESSAGES[errorCode]
    }

    // Fallback ke HTTP status message
    if (HTTP_STATUS_MESSAGES[httpCode]) {
      return HTTP_STATUS_MESSAGES[httpCode]
    }
  }

  // 3. Fallback: jika message sudah berupa teks biasa, tampilkan apa adanya
  return message
}
