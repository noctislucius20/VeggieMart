# 🥬 VeggieMart - Aplikasi E-commerce Sayuran

Aplikasi e-commerce sayuran berbasis **microservices architecture** dengan **frontend Nuxt 3** dan **backend Go** yang terdiri dari 5 layanan terpisah, API Gateway, serta didukung oleh message broker, cache, dan search engine.

## 🏗️ Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend (Nuxt 3)                   │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   API Gateway (Go/Echo)                    │
│            JWT Auth • Rate Limit • CORS • Logging          │
└────────────────────────────────────────────────────────────┘
        │           │           │           │           │
        ▼           ▼           ▼           ▼           ▼
   ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
   │  User   │ │ Product │ │  Order  │ │ Payment │ │Notific. │
   │ Service │ │ Service │ │ Service │ │ Service │ │ Service │
   └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘
        │           │           │           │           │
        ▼           ▼           ▼           ▼           ▼
   ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
   │Postgres │ │Postgres │ │Postgres │ │Postgres │ │Postgres │
   └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘

   Infrastruktur:  Redis • RabbitMQ • Elasticsearch
```

## ✨ Fitur Utama

- 🔐 Autentikasi & otorisasi pengguna (JWT)
- 👤 Manajemen profil pengguna
- 🥦 Katalog produk & kategori
- 🛒 Keranjang belanja
- 💳 Checkout & pembayaran
- 📦 Riwayat pesanan
- 🔔 Notifikasi real-time (WebSocket)
- 📍 Geolokasi untuk pengiriman
- 🛡️ API Gateway dengan rate limiting & CORS

## 🛠️ Teknologi

### Frontend

| Teknologi    | Keterangan           |
| ------------ | -------------------- |
| Nuxt 3       | Vue.js framework     |
| Vue 3        | JavaScript framework |
| Pinia        | State management     |
| Tailwind CSS | CSS framework        |
| Bootstrap 5  | UI framework         |
| Swiper       | Carousel/slider      |
| Quill        | Rich text editor     |
| Vuelidate    | Form validation      |
| Day.js       | Date utilities       |

### Backend

| Teknologi     | Keterangan             |
| ------------- | ---------------------- |
| Go            | Bahasa pemrograman     |
| Echo          | Web framework          |
| GORM          | ORM                    |
| PostgreSQL    | Database               |
| Redis         | Cache & rate limiting  |
| RabbitMQ      | Message broker         |
| Elasticsearch | Search engine          |
| WebSocket     | Real-time notification |

## 📁 Struktur Project

```
VeggieMart/
├── api/                              # Backend services
│   ├── docker-compose.yml            # Docker Compose untuk semua services
│   ├── api-gateway/                  # API Gateway
│   ├── user-service/                 # User service
│   ├── product-service/              # Product service
│   ├── order-service/                # Order service
│   ├── payment-service/              # Payment service
│   └── notification-service/         # Notification service
├── ui/                               # Frontend Nuxt 3
├── k6-performance/                   # Performance testing scripts
└── Jenkinsfile                       # CI/CD pipeline
```

## 🚀 Menjalankan dengan Docker Compose

### Prasyarat

- [Docker](https://www.docker.com/) & Docker Compose
- [Node.js](https://nodejs.org/) (v18+)
- [Bun](https://bun.sh/) (opsional, untuk frontend)

### 1. Jalankan Backend (Semua Services)

```bash
cd api

# Salin file environment
cp .env.example .env

# Sesuaikan konfigurasi di .env
# DB_USER, DB_PASS, RBMQ_USER, RBMQ_PASS, dll.

# Jalankan semua services
docker-compose up --build
```

### 2. Jalankan Frontend

```bash
cd ui

# Install dependencies
bun install
# atau
npm install

# Salin file environment
cp .env.example .env

# Jalankan development server
bun dev
# atau
npm run dev
```

## 🐳 Layanan

| Layanan              | Deskripsi             |
| -------------------- | --------------------- |
| Frontend (Nuxt)      | Aplikasi web          |
| API Gateway          | Entry point API       |
| User Service         | Manajemen user & auth |
| Product Service      | Katalog produk        |
| Order Service        | Transaksi pesanan     |
| Payment Service      | Pembayaran            |
| Notification Service | Notifikasi            |
| PostgreSQL           | Database per service  |
| Redis                | Cache & rate limit    |
| RabbitMQ             | Message broker        |
| Elasticsearch        | Search engine         |

## 🧪 Performance Testing

Proyek ini menggunakan [k6](https://k6.io/) untuk performance testing:

```bash
cd k6-performance/veggiemart-app-k6
k6 run user-service/script.js
k6 run product-test/script.js
k6 run order-service/script.js
```

## 🔄 CI/CD

Proyek ini menggunakan **Jenkins** untuk CI/CD pipeline. Lihat file `Jenkinsfile` untuk detail konfigurasi.

## 📝 Lisensi

[MIT License](LICENSE)
