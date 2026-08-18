import { helpers } from '@vuelidate/validators'

/**
 * Custom validator: nilai harus salah satu dari daftar yang diizinkan
 * @param {string[]} allowedValues - Daftar nilai yang diizinkan
 */
export const oneof = (allowedValues) =>
  helpers.withMessage(`Nilai harus salah satu dari: ${allowedValues.join(', ')}`, (value) => {
    if (!value) return true
    return allowedValues.includes(value)
  })

/**
 * Custom validator: format tanggal YYYY-MM-DD
 */
export const dateFormat = helpers.withMessage('Format tanggal harus YYYY-MM-DD', (value) => {
  if (!value) return true
  return /^\d{4}-\d{2}-\d{2}$/.test(value)
})

/**
 * Custom validator: format waktu HH:MM:SS
 */
export const timeFormat = helpers.withMessage('Format waktu harus HH:MM:SS', (value) => {
  if (!value) return true
  return /^\d{2}:\d{2}:\d{2}$/.test(value)
})

/**
 * Custom validator: nomor telepon harus angka (omitempty)
 */
export const numeric = helpers.withMessage('Nomor telepon harus berupa angka', (value) => {
  if (!value || value === '') return true // omitempty
  return /^\d+$/.test(value)
})

/**
 * Custom validator: nilai harus lebih besar atau sama dengan minimum
 * @param {number} min - Nilai minimum yang diizinkan
 */
export const gte = (min) =>
  helpers.withMessage(`Nilai harus lebih besar atau sama dengan ${min}`, (value) => {
    if (value === null || value === undefined || value === '') return true // omitempty
    return Number(value) >= min
  })
