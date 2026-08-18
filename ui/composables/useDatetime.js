export const useDateTime = () => {
  const formatDate = (date) => {
    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'long',
      year: 'numeric'
    }).format(new Date(date))
  }

  return {
    formatDate
  }
}

export const formatDatetime = () => {
  const formatDate = (date) => {
    // Add null check
    if (!date) return '-'

    const [datePart, timePart] = date.split(' ')
    const [year, month, day] = datePart.split('-').map(Number)
    const [hour, minute] = timePart.split(':').map(Number)

    const dateData = new Date(year, month - 1, day, hour, minute)

    const monthNames = [
      'Januari',
      'Februari',
      'Maret',
      'April',
      'Mei',
      'Juni',
      'Juli',
      'Agustus',
      'September',
      'Oktober',
      'November',
      'Desember'
    ]

    const formatted = `${String(day).padStart(2, '0')} ${monthNames[month - 1]} ${year} `

    return formatted
  }

  return {
    formatDate
  }
}
