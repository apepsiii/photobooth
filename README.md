# Aplikasi Photobooth Sekolah (Inspirasi: Funcam.id)

Aplikasi photobooth modern, sangat ringan (*ultra-lightweight*), berdesain *clean minimalism*, dan mudah dioperasikan secara mandiri (*self-service kiosk*) untuk lingkungan sekolah (Pentas Seni, Bazar, Wisuda/Kelulusan, Class Meeting, dan event OSIS).

Dibangun menggunakan **Go (Gin Framework)** + **GORM** + **Pure-Go SQLite** (tanpa CGO/GCC di Windows) + **Meridian Admin Dashboard**, dengan antarmuka kiosk fullscreen ramah layar sentuh, stiker interaktif, dan sistem QR Code softfile offline lokal.

---

## Fitur Utama & Pengalaman Pengguna

### 1. Layar Kiosk Fullscreen (`/kiosk` atau `/`):
- **Modern Clean Minimalism:** Tampilan gelap (*dark aesthetic*), tipografi bersih, tombol sentuh besar, bebas elemen rumit.
- **WebRTC Camera & Switcher:**
  - Mendukung webcam laptop maupun USB HD Webcam / DSLR via Cam Link.
  - Tombol **"Ganti Kamera"** otomatis muncul jika sistem mendeteksi lebih dari satu perangkat kamera.
  - Tampilan cermin (*mirroring*) natural untuk kenyamanan selfie.
- **Live Countdown Visual & Sound Effect:** Hitungan mundur 3-2-1 dinamis disertai efek audio beep dan *flash screen* putih saat pemotretan.
- **3 Template Frame Default:**
  - *Classic 2-Strip (Funcam Style)* — 3 foto vertikal (2x6 inci).
  - *Korean 4-Cut Grid (4R)* — 4 foto dalam grid 2x2 ala photobooth Korea (*Life Four Cuts*).
  - *Polaroid Minimalist (Wide)* — 1 foto square besar dengan border tanda tangan.
- **Filter Warna Estetik:** Normal (Original), B&W Classic (Monochrome), Warm Vintage, dan Cool Pastel.
- **Stiker Digital & Emojis Interaktif (Funcam Style):**
  - Siswa dapat menempelkan berbagai stiker (🎓 Topi Wisuda, 👑 Mahkota, 🎀 Pita, 💖 Love, ⭐ Bintang, 👓 Kacamata, 🐱 Kucing, dsb.).
  - Mendukung **drag & drop touch** langsung di atas kanvas foto untuk mengatur posisi stiker.
- **Watermark & Branding Sekolah Otomatis:**
  - Menempelkan identitas sekolah dan nama acara secara rapi di bagian bawah foto secara otomatis.
- **Auto-Reset Kiosk:** Timer 60 detik di layar akhir untuk mereset booth bagi antrean berikutnya.

### 2. Distribusi Softfile QR Code 100% Offline Lokal (`/d/:code`):
- Otomatis mendeteksi IP lokal server (contoh: `http://192.168.1.4:8090/d/PB-XXXXX`).
- QR code digenerate secara native oleh Go (`/api/qr`), **100% bekerja tanpa koneksi internet luar** (cukup sambungkan laptop dan HP ke WiFi/hotspot yang sama).
- **Laman Unduh Mobile:**
  - **Tab 1 (Foto Cetak):** Unduh foto komposit final beresolusi penuh.
  - **Tab 2 (Boomerang Loop):** Animasi bergerak bolak-balik dari shot foto yang siap di-upload ke Instagram Stories / TikTok!
  - **Foto Individual:** Unduh tiap-tiap foto mentah (*raw shot*).

### 3. Layanan Cetak (Print Service):**
- Terhubung langsung dengan spooler cetak Windows OS untuk mengirim berkas ke printer foto (Canon SELPHY, Epson L-series, DNP, dsb.).
- Opsi *Auto-print* (otomatis cetak tanpa konfirmasi) atau *Manual Confirmation*.
- Fitur *Re-print* di panel admin jika kertas printer macet.

### 4. Panel Admin & Operator ([`/admin/sessions`](http://localhost:8090/admin/sessions)):**
- **Galeri Sesi & Counter:** Pantau seluruh riwayat sesi foto hari ini, counter cetak, dan counter unduh.
- **1-Click ZIP Export (`📦 Unduh ZIP Dokumentasi`):** Download seluruh foto komposit hari itu dalam satu file ZIP untuk arsip panitia OSIS.
- **Manajemen Frame ([`/admin/frames`](http://localhost:8090/admin/frames)):** Upload template frame PNG baru dan atur koordinat slot foto.
- **Pengaturan & Printer ([`/admin/settings`](http://localhost:8090/admin/settings)):** Atur nama sekolah, nama acara, durasi countdown, PIN kiosk, dan printer target.

---

## Panduan Menjalankan

### Cara 1: Menggunakan Launcher Satu-Klik (Rekomendasi di Hari H)
Cukup klik dua kali file **`run_kiosk.bat`** di Windows. Skrip ini akan:
1. Memulai server backend Go di port 8090.
2. Membuka Google Chrome / Microsoft Edge dalam **Fullscreen Kiosk Mode** yang terkunci.

### Cara 2: Menjalankan via Terminal
```bash
# Jalankan biner langsung:
.\photobooth.exe

# Atau jalankan via Go:
go run .
```

Buka di browser:
- **Layar Kiosk Siswa:** [http://localhost:8090/kiosk](http://localhost:8090/kiosk) (atau [http://localhost:8090/](http://localhost:8090/))
- **Panel Admin:** [http://localhost:8090/admin/sessions](http://localhost:8090/admin/sessions)

---

## Kredensial Default

| Akses | Kredensial |
|---|---|
| **Login Admin/Operator** | `admin@meridian.com` / `admin123` |
| **PIN Operator Kiosk** | `1234` *(klik tombol gerigi di kanan atas layar kiosk)* |
| **Port Server** | `8090` *(dapat diubah via `.env`)* |
