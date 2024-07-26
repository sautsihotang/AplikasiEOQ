# AplikasiEOQ

AplikasiEOQ adalah aplikasi berbasis Go yang dirancang untuk mengelola data terkait dengan Economic Order Quantity (EOQ) dan berbagai entitas terkait dalam sistem manajemen inventaris. Aplikasi ini menggunakan GORM sebagai ORM dan PostgreSQL sebagai database.

## Fitur

- **Manajemen Barang:** Menyimpan dan mengelola informasi barang.
- **Pemesanan:** Mengelola data pemesanan termasuk biaya terkait dan tanggal pemesanan.
- **Penyimpanan:** Menyimpan informasi mengenai biaya penyimpanan barang.
- **EOQ:** Menghitung dan menyimpan nilai EOQ untuk barang.
- **Pengguna:** Mengelola data pengguna sistem.
- **Supplier:** Menyimpan informasi mengenai supplier.
- **Penjualan:** Mencatat transaksi penjualan barang.

## Prerequisites

Sebelum menjalankan aplikasi, pastikan Anda memiliki:

- **Go (Golang)** versi 1.18 atau lebih baru.
- **PostgreSQL** sebagai database.
- **GORM** dan **Godotenv** untuk pengelolaan database dan variabel lingkungan.

## Instalasi

1. **Clone Repository:**

   ```bash
   git clone https://github.com/username/AplikasiEOQ.git
   cd AplikasiEOQ

2. **Instalasi Dependencies:**
    Instal GORM dan driver PostgreSQL menggunakan Go Modules:

    ```bash
    go mod tidy

3. **Konfigurasi Databas**
    Buat file .env di direktori root proyek Anda dengan konfigurasi berikut:
    
    ```bash
    DB_USER=your_username
    DB_PASSWORD=your_password
    DB_NAME=your_database
    DB_HOST=localhost
    DB_PORT=5432

4. **Menjalankan Aplikasi:**
    
    ```bash
    go run main.go




