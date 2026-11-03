# Implementation Plan (PLAN)
## Aplikasi Photobooth Sekolah (Mastergin: Go + Gin + GORM + SQLite)

Dokumen ini berisi tahapan eksekusi dan roadmap implementasi menggunakan pondasi starterkit **Mastergin**.

---

## 1. Arsitektur & Struktur Direktori

```text
photobooth/
├── main.go                        # Entry point server Gin
├── core/                          # Framework layer (Config, DB GORM+SQLite, Router, Template Engine)
├── middleware/                    # Middleware (RequireAuth JWT cookie)
├── models/                        # Core models (User)
├── modules/                       # Domain Feature Modules
│   ├── auth/                      # Built-in Auth (Login, Logout, SeedAdmin)
│   ├── kiosk/                     # Photobooth Touchscreen Frontend & API Capture
│   │   ├── handler.go             # Layar Kiosk, Upload Render, Offline QR Code API
│   │   └── routes.go              # /kiosk, /api/kiosk/finish, /api/qr
│   ├── frames/                    # Frame & Template Management
│   │   ├── model.go               # Frame GORM Model
│   │   ├── handler.go             # Admin Frame CRUD & API List Active
│   │   ├── routes.go              # /admin/frames, /api/frames
│   │   └── service.go             # Default frames seeder (PNG Generator)
│   ├── sessions/                  # Photo Sessions, Gallery & Mobile Download
│   │   ├── model.go               # Session & SessionPhoto GORM Models
│   │   ├── handler.go             # Session Gallery & Public Mobile Download (/d/:code)
│   │   ├── routes.go              # /admin/sessions, /admin/sessions/export-zip, /d/:code
│   │   └── service.go             # Windows Print spooler & Image storage
│   ├── vouchers/                  # Ticket & Voucher Management (Fundraiser / Bazaar)
│   │   ├── model.go               # Voucher GORM Model
│   │   ├── handler.go             # Admin Vouchers, Batch Generator, Toggle Mode
│   │   ├── routes.go              # /admin/vouchers, /api/vouchers/validate
│   │   └── service.go             # Code generator & validation logic
│   └── settings/                  # Kiosk & Printer Settings
│       ├── model.go               # Setting GORM Model
│       ├── handler.go             # Admin Settings Form & API Config
│       ├── routes.go              # /admin/settings, /api/settings
│       └── service.go             # Local IP auto-detection & default seeder
├── template/meridian/             # Meridian UI Theme
│   ├── layouts/master.html        # Master admin layout (sidebar navigation)
│   └── pages/                     # Templates (frames, sessions, settings_app, vouchers, kiosk, download)
├── storage/                       # Local disk storage
│   ├── photos/raw/                # Foto mentah per shot
│   ├── photos/final/              # Foto komposit selesai render
│   └── frames/                    # File PNG template frame
├── run_kiosk.bat                  # One-click Kiosk launcher (Chrome/Edge Fullscreen)
├── .env                           # Konfigurasi aplikasi (PORT=8090, DB=app.db)
├── go.mod / go.sum
└── README.md                      # Dokumentasi & panduan penggunaan
```

---

## 2. Status Tahapan Eksekusi (Milestones)

### Phase 1: Adaptasi Fondasi Mastergin
- [x] Salin starterkit `mastergin` ke workspace `photobooth`.
- [x] Perbarui PRD dan PLAN di folder `context/`.
- [x] Konfigurasi database SQLite pure-go (WAL mode, foreign keys).
- [x] Buat direktori `storage/photos/raw`, `storage/photos/final`, `storage/frames`.

### Phase 2: Modul Frames & Settings (Admin Management)
- [x] Buat model `Frame` dan `Setting` dengan auto-migration GORM.
- [x] Tambahkan 3 frame bawaan: Classic 2-Strip (Funcam style), Korean 4-Cut Grid (4R), dan Polaroid Minimalist.
- [x] Buat halaman admin Meridian `frames.html` dengan **Preset Koordinat Otomatis**.
- [x] Buat halaman admin Meridian `settings_app.html` (nama acara, countdown, auto-print, copies, IP lokal).

### Phase 3: Modul Kiosk UI (Clean Minimalism Touchscreen)
- [x] Desain antarmuka standalone `/kiosk` fullscreen, modern clean minimalism.
- [x] Integrasi WebRTC Camera Feed dengan preview mirror selfie & **Camera Switcher**.
- [x] Alur interaktif lengkap:
  1. *Welcome screen* ("Sentuh untuk Memulai").
  2. *Pilihan Frame* (thumbnail grid interaktif dengan visual checklist).
  3. *Countdown timer visual & audio* (beep tone via Web Audio API + flash putih layar saat capture).
  4. *Multi-shot capture* (mengambil 1-4 foto berurutan dengan thumbnail progress bar).
  5. *Pilihan filter estetik* (Normal, B&W Classic, Warm Vintage, Cool Pastel).
  6. *Stiker Digital & Emojis* (drag-and-drop langsung di kanvas sentuh).
  7. *Client-side canvas compositing* presisi tinggi + watermark sekolah otomatis.

### Phase 4: Penyimpanan Sesi & QR Code Softfile Download (100% Offline)
- [x] API endpoint `/api/kiosk/finish` untuk menyimpan komposit JPEG dan foto mentah.
- [x] Deteksi otomatis IP lokal server (contoh: `192.168.1.4:8090`).
- [x] Generator QR Code murni Go (`/api/qr?text=...`) — **100% offline tanpa perlu koneksi internet**.
- [x] Laman unduh mobile `/d/:code` responsif dengan **Tab Foto Cetak** dan **Tab Boomerang Loop**.

### Phase 5: Layanan Cetak (Print Service Integration)
- [x] Implementasi fungsi cetak Windows OS PowerShell spooler ke default atau target printer.
- [x] Mode konfigurasi cetak: Otomatis (*Auto Print*), Konfirmasi Manual (*Button*), atau Digital Only.
- [x] Tombol *Re-print* di galeri sesi admin jika kertas printer habis atau macet.

### Phase 6: Modul Tiket & Kupon Bazaar OSIS
- [x] Buat modul `vouchers` (model, generator batch, validator).
- [x] Halaman admin `/admin/vouchers` (buat kupon, cetak lembar tiket, toggle mode gratis/berbayar).
- [x] Modal input tiket dengan **Keypad Virtual Touchscreen** di layar Kiosk.

### Phase 7: Kiosk Launcher, Dokumentasi & Pengujian
- [x] Buat skrip launcher otomatis `run_kiosk.bat` (Google Chrome / Microsoft Edge `--kiosk`).
- [x] Tombol `📦 Unduh ZIP Dokumentasi` di galeri admin untuk arsip panitia OSIS.
- [x] Panduan instalasi dan penggunaan lengkap di `README.md`.
- [x] Pengujian menyeluruh endpoint & rendering: All HTTP 200 OK.

### Phase 8: Peningkatan Admin Gallery & Studio UX
- [x] Pencarian langsung (*Instant Live Search*) di Galeri Sesi (`sessions.html`).
- [x] Modal Lightbox Pratinjau Resolusi Penuh (*Full-res Image Lightbox*).
- [x] Manajemen penghapusan sesi foto beserta pembersihan file fisik di disk (`/admin/sessions/:id/delete`).
- [x] Penambahan 8 warna latar belakang bingkai (*Pastel & Dark color swatches*) dengan watermark otomatis menyesuaikan kontras.
- [x] Penambahan 6 filter estetik kamera (Normal, Korean Mono, Golden Hour, Tokyo Cool, Soft Glow, Vintage 90s).
- [x] Katalog stiker drag-and-drop 4 kategori (Sekolah, Wisuda, Lucu, Pesta, Bestie).

### Phase 9: Peningkatan Hardware Tooling, Perforated Tickets & Social Share
- [x] **Visual Slot Inspector** di Frame Management (`frames.html`) untuk pratinjau kanvas koordinat slot foto.
- [x] **Lembar Cetak Kupon A4 Perforasi** (`vouchers.html` via `@media print`) dengan garis gunting siap potong untuk panitia OSIS.
- [x] **Uji Coba Cetak Printer Langsung** (`settings_app.html` via `/admin/sessions/test-print`) untuk kalibrasi kertas & printer hardware.
- [x] **Web Share API & Salin Tautan** di laman unduh softfile mobile (`download.html`).

### Phase 10: Server-Side Boomerang GIF & Smart Mirror Camera
- [x] **Generator GIF Animasi Murni Go (`/d/:code/gif`)**: Downsampling, kuantisasi palet warna web, dan streaming animasi loop Boomerang bolak-balik tanpa pustaka eksternal CGO.
- [x] **Tombol Unduh File GIF**: Siswa dapat menyimpan berkas `.gif` asli langsung ke galeri smartphone untuk diunggah ke medsos.
- [x] **Toggle Cermin Kamera (Mirror/Un-mirror)**: Pengunjung dapat mengaktifkan atau menonaktifkan efek cermin horizontal secara instan agar teks selempang/baju tidak terbalik.

### Phase 11: Retake Single Shot & Export CSV Pelaporan
- [x] **Foto Ulang Shot Tertentu (*Retake Single Shot*)**: Pada Studio Editor Kiosk (`kiosk.html`), pengunjung dapat mengulang hanya satu shot pose tertentu tanpa mengulang seluruh sesi.
- [x] **Unduh Laporan CSV Sesi (`/admin/sessions/export-csv`)**: Panitia OSIS dapat mengunduh rekapitulasi data sesi foto berformat Excel/CSV lengkap dengan UTF-8 BOM untuk pelaporan keuangan bazaar dan pertanggungjawaban acara.
- [x] **Web Audio Resumption**: Penjaminan efek audio shutter dan beep countdown langsung aktif di semua browser saat sentuhan pertama di layar Kiosk.

### Phase 12: Custom Text Stamps, Live Admin Monitoring & Quick Ticket Copy
- [x] **Tab 4: ✍️ Teks (Custom Text Stamp)** pada Studio Editor Kiosk: Siswa dapat mengetik nama kelas, geng, atau ucapan wisuda dengan pilihan 5 palet warna dan menggesernya (*drag & drop*) di atas foto.
- [x] **Auto-Refresh Live Switch** di Galeri Admin (`sessions.html`): Switch live toggle auto-refresh 10 detik dengan persistensi `localStorage` agar operator dapat memantau sesi foto yang masuk secara real-time tanpa refresh manual.
- [x] **Salin Kode Tersedia Instan** di Manajemen Tiket (`vouchers.html`): Tombol 1-klik untuk menyalin seluruh kode tiket yang belum terpakai ke clipboard panitia pendaftaran.
