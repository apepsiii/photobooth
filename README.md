# 📸 GoPhotobooth — Modern Offline Photobooth & Kiosk System

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
  <img src="https://img.shields.io/badge/Gin_Web_Framework-v1.9+-008ECF?style=for-the-badge&logo=gin&logoColor=white" alt="Gin" />
  <img src="https://img.shields.io/badge/SQLite-Pure--Go_(Zero_CGO)-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite Pure Go" />
  <img src="https://img.shields.io/badge/Design-Neobrutalism_&_Clean_Kiosk-FFE600?style=for-the-badge&logo=html5&logoColor=black" alt="Design" />
  <img src="https://img.shields.io/badge/Network-100%25_Offline_Local_LAN-28A745?style=for-the-badge&logo=wifi&logoColor=white" alt="Offline Ready" />
</p>

<p align="center">
  <strong>Aplikasi Photobooth interaktif, self-service, dan ultra-cepat untuk sekolah, bazaar, wisuda, wedding, dan event komunitas.</strong><br>
  <em>100% Mandiri tanpa ketergantungan internet luar — siap cetak seketika & distribusi softfile via QR code lokal!</em>
</p>

---

## ⚡ Sekilas Fitur Unggulan

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   🎪 Kiosk UI   │ ───►  │  🎨 Frame Studio│ ───►  │  📱 Softfile QR │
│ Touchscreen &   │       │ Drag & Drop     │       │ 100% Offline    │
│ Live Countdown  │       │ On-Screen Mapper│       │ Boomerang + High│
└─────────────────┘       └─────────────────┘       └─────────────────┘
         │                         │                         │
         ▼                         ▼                         ▼
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│ 🎟️ Tiket & PIN  │       │ 🖨️ Auto-Print   │       │ 🌟 Neobrutalism │
│ Kupon Bazaar    │       │ Windows Spooler │       │ Landing Page    │
│ & Akses Khusus  │       │ & One-Click ZIP │       │ & Live Showcase │
└─────────────────┘       └─────────────────┘       └─────────────────┘
```

- 📺 **Live Projector Slideshow Wall (`/live-wall`):**
  - **Proyektor / TV Antrean Venue:** Layar khusus full-screen resolusi tinggi untuk ditayangkan ke proyektor aula, videotron LED, atau monitor antrean.
  - **Live Auto-Sync & Polling:** Otomatis memperbarui foto secara real-time setiap kali ada pengunjung yang menyelesaikan sesi di kiosk.
  - **Celebration Fanfare Overlay:** Begitu foto baru terambil, proyektor otomatis memunculkan pop-up animasi *"FOTO BARU BARUSAN DIAMBIL! (ANTREAN #...)"* lengkap dengan efek audio chime perayaan.
  - **Dual Mode (Spotlight vs Mosaic Grid):** Mode sorotan utama dengan efek Ken Burns halus atau galeri dinding multi-foto (*shortcut: tekan G*).
  - **Pojok Scan QR Cepat:** Memudahkan tamu di aula memindai QR langsung dari layar panggung untuk mencari dan mengunduh foto mereka.
- 🎨 **Creative Neobrutalism Landing Page (`/`):** Tampilan retro-pop modern dengan marquee running ticker, bilah pencarian instan kode softfile (`/d/PB-XXXXX`), kartu fitur pop-out, dan galeri pameran hasil jepretan live.
- 📸 **Fullscreen Kiosk Layar Sentuh (`/kiosk`):**
  - **💳 Mode Self-Service Vending QRIS:** Kiosk mandiri tanpa operator; pengunjung memilih varian (Paket Single / Double Bestie), scan QRIS di layar via GoPay/OVO/BCA/Dana, dan sesi otomatis terbuka begitu terkonfirmasi.
  - **✋ Hands-Free Optical Gesture Shutter:** Pengunjung berpose dari jarak 1-2 meter dan cukup mengangkat tangan ✋ ke target HUD di sudut kamera untuk memicu hitungan mundur secara otomatis (atau via tombol remote Bluetooth/Spasi).
  - **🗣️ Audio Voice Assistant & Pemandu Pose Bahasa Indonesia:** Suara pemandu interaktif yang memandu pengunjung dari hitungan mundur ("3... 2... 1... Senyum!"), rekomendasi gaya pose bergantian di sela-sela jepretan, hingga instruksi pengambilan foto cetak dan scan QR.
  - **Live Webcam Switcher:** Deteksi multi-kamera (Webcam laptop, USB Cam, Cam Link DSLR) + saklar *Mirroring*.
  - **Audio & Visual Beep Countdown:** Hitungan mundur 3-2-1 disertai efek suara dan *flash screen* putih real-time.
  - **Freeze-Frame Preview per Shot:** Begitu foto terambil, hasil langsung tampil membeku selama 1.2 detik sebelum lanjut ke pose berikutnya!
  - **Single-Shot Retake:** Kurang puas dengan pose ke-2? Cukup klik tombol `↺ Retake Shot 2` tanpa harus mengulang dari awal.
  - **✏️ Freehand Doodle & Tanda Tangan Digital:** Goreskan kuas neon atau tanda tangan langsung di atas canvas foto dengan 8 warna cerah, 3 ukuran kuas, penghapus pintar, dan fitur Undo.
  - **10 Filter Estetik & K-Beauty Glow:** K-Glow peach smoothing, Korean Mono, Golden Hour, 90s Film Grain, Cyber Y2K, Noir Dramatic, Nordic Matte, dan Cinematic Cine.
  - **Universal Inactivity Watchdog:** Reset otomatis dengan countdown dialog jika booth ditinggalkan pengunjung, serta proteksi anti-tamper.
- 📊 **Laporan Finansial & Analitik Booth (`/admin/analytics`):**
  - **Financial Summary Cards:** Omzet kotor gabungan, pendapatan QRIS mandiri, pendapatan tiket/voucher tunai, serta estimasi margin laba bersih setelah HPP kertas foto.
  - **⏱️ Grafik Jam Ramai (Hourly Rush Hours):** Distribusi kepadatan sesi dari pukul 07:00 hingga 22:00 dalam bentuk bar chart interaktif untuk evaluasi alur antrean pengunjung.
  - **🏆 Leaderboard Template Frame Terlaris:** Peringkat frame favorit yang paling sering dipilih pengunjung.
  - **🎨 Leaderboard Filter Terpopuler:** Statistik pemakaian filter estetik (K-Glow, Monochrome, 90s Vintage, dll.).
  - **💳 Log Transaksi QRIS Mandiri:** Riwayat pesanan vending dengan ID transaksi, nominal, varian paket, status pembayaran, dan stempel waktu.
- 📱 **Remote Operator Mobile Web Controller (`/operator`):**
  - **Live Booth Status & Screen Monitor:** Memantau status kiosk secara nirkabel (Online/Offline, layar aktif saat ini, dan latency ping).
  - **📸 Remote Shutter Trigger:** Picu timer hitungan mundur kamera dari HP operator saat mengarahkan rombongan atau foto bersama.
  - **🔄 Emergency Force Reset:** Kembalikan kiosk yang macet atau ditinggalkan pengunjung ke layar awal Standby dengan satu sentuhan.
  - **🎟️ Pembuat Kupon Tiket Cepat:** Hasilkan kode tiket voucher baru dalam 1 tap langsung di meja resepsionis/kasir untuk dibagikan ke pengunjung.
  - **🖨️ Cetak Ulang Instan (Reprint):** Menampilkan thumbnail foto komposit sesi terakhir dan tombol spool cetak ulang tanpa perlu buka laptop.
  - **📜 Pemantau Stok Kertas & Ribbon:** Gauge visual kapasitas kertas tersisa dengan indikator peringatan warna kuning/merah saat kertas hampir habis, serta tombol reset roll baru.
- 📐 **Visual Frame Studio & Interactive Mapper (`/admin/frames`):**
  - **Direct On-Screen Drag & Resize:** Geser dan ubah ukuran slot foto langsung dengan mouse/touchpad tanpa perlu menebak koordinat pixel!
  - **8-Handle Bounding Box:** Kontrol sudut dan sisi box dengan presisi tinggi (60 FPS smooth rendering).
  - **Draw Marquee Tool (`✏️ Tarik & Gambar Slot`):** Buat slot baru cukup dengan klik-dan-tarik kotak di atas template frame.
  - **Action Toolbar & Keyboard Shortcuts:** Duplikat slot (`Ctrl+D`), hapus slot (`Delete`), geser presisi (`Tombol Panah`), serta Center X / Center Y otomatis.
- 📱 **Softfile Download 100% Offline Lokal (`/d/:code`):**
  - QR Code dibuat langsung secara internal oleh mesin Go (tanpa API pihak ketiga).
  - Pengunjung cukup scan QR dengan HP di jaringan WiFi yang sama untuk mengunduh:
    1. **Foto Cetak Komposit Final** (resolusi tinggi siap cetak/posting).
    2. **Animasi Boomerang GIF** (gerak bolak-balik dari shot foto untuk Instagram Stories / TikTok).
    3. **File Foto Mentah (*Raw Shots*)** per pose.
- 🎟️ **Sistem Tiket & Voucher Bazaar (`/admin/vouchers`):**
  - Generator kode kupon batch untuk tiket berbayar atau kuota foto di bazaar/event sekolah.
  - Mendukung batas pemakaian (*usage limit*) dan masa berlaku.
- 🖨️ **Integrasi Mesin Cetak Fisik:**
  - Terkoneksi ke Windows Print Spooler (Canon SELPHY, Epson L-Series, DNP RX1HS, HiTi, dsb.).
  - Mendukung *Auto-Print* otomatis saat sesi usai, serta tombol *Re-Print* darurat di dashboard operator.
- 📦 **Arsip 1-Click ZIP Export:**
  - Panitia/Dokumentasi dapat mengunduh seluruh komposit foto hari tersebut dalam satu file ZIP utuh langsung dari dashboard.

---

## 🛠️ Arsitektur & Teknologi

| Komponen | Spesifikasi | Keterangan |
|---|---|---|
| **Backend Engine** | **Go (Golang) 1.22+** | Kompilasi biner mandiri, performa tinggi, hemat memori (< 35 MB RAM). |
| **HTTP Framework** | **Gin Gonic v1.9** | Router berkecepatan tinggi dengan REST API modular. |
| **Database** | **Pure-Go SQLite** (`github.com/glebarez/sqlite`) | **Zero CGO / Tanpa GCC compiler** di Windows, single file `.db`, ACID compliant. |
| **ORM** | **GORM v1.25** | Migrasi skema otomatis (*AutoMigrate*). |
| **Templating** | **Meridian Multi-Template Engine** | Template Go terisolasi per halaman, mendukung inheritance `{{block "content"}}`. |
| **Desain UI** | **Neobrutalism + Clean Touch Kiosk** | High-contrast borders, bold typography, dan layout responsif sentuh. |
| **Kamera & Grafis** | **WebRTC + HTML5 Canvas + Pure-Go Imaging** | Pemrosesan komposit foto dan rendering QR Code sepenuhnya di server/klien lokal. |

---

## 🚀 Panduan Memulai Cepat

### 1. Kloning Repositori
```bash
git clone https://github.com/username/photobooth.git
cd photobooth
```

### 2. Salin Konfigurasi Environment
```bash
copy .env.example .env
```

Isi default `.env`:
```ini
APP_NAME="Photobooth Kiosk"
APP_PORT=8090
APP_ENV=local
GIN_MODE=debug
DB_PATH=photobooth.db
JWT_SECRET=super-secret-photobooth-key-12345
LOCAL_IP=192.168.1.100
```
> 💡 *Ubah `LOCAL_IP` sesuai alamat IP LAN / WiFi laptop Anda agar QR Code softfile langsung dapat dibuka dari smartphone pengunjung.*

### 3. Jalankan Aplikasi

#### Pilihan A: Mode Siap Pakai (Windows One-Click)
Klik dua kali file **`run_kiosk.bat`**. Skrip ini otomatis:
1. Menjalankan server Go di background.
2. Membuka Google Chrome / Edge dalam **Fullscreen Kiosk Mode** yang ramah layar sentuh.

#### Pilihan B: Jalankan via Terminal
```powershell
# Menggunakan biner terkompilasi:
.\photobooth.exe

# Atau jalankan langsung dengan Go:
go run .
```

Aplikasi siap diakses di port **8090**:
- 🌐 **Landing Page & Softfile Lookup:** [http://localhost:8090/](http://localhost:8090/)
- 📸 **Kiosk Layar Sentuh Siswa:** [http://localhost:8090/kiosk](http://localhost:8090/kiosk)
- 📐 **Visual Frame Studio:** [http://localhost:8090/admin/frames](http://localhost:8090/admin/frames)
- 📊 **Panel Sesi & Operator:** [http://localhost:8090/admin/sessions](http://localhost:8090/admin/sessions)
- 🎟️ **Manajemen Voucher:** [http://localhost:8090/admin/vouchers](http://localhost:8090/admin/vouchers)

---

## 🔐 Kredensial Default

| Akses | Alamat URL | Akun / Sandi |
|---|---|---|
| **Super Admin / Operator** | `/login` | `admin@meridian.com` / `admin123` |
| **PIN Cepat Kiosk** | Pojok kanan atas `/kiosk` | `1234` |
| **Portal Softfile Pengunjung** | `/d/:code` *(contoh: `/d/PB-SAMPLE`)* | Publik (via QR Code / Input Kode) |

---

## 🎮 Pintasan Keyboard (Visual Frame Studio)

Saat berada di studio mapping frame ([`/admin/frames`](http://localhost:8090/admin/frames)):

| Tombol Keyboard | Aksi |
|---|---|
| <kbd>Delete</kbd> / <kbd>Backspace</kbd> | Hapus slot foto yang sedang dipilih |
| <kbd>Ctrl</kbd> + <kbd>D</kbd> | Duplikasi slot foto aktif dengan dimensi sama |
| <kbd>▲</kbd> <kbd>▼</kbd> <kbd>◄</kbd> <kbd>►</kbd> | Geser posisi slot per 1 piksel (*nudge*) |
| <kbd>Shift</kbd> + Tombol Panah | Geser posisi slot per 10 piksel |
| <kbd>Drag Mouse</kbd> | Pindahkan posisi box slot |
| <kbd>Drag Handle Sudut</kbd> | Perbesar / perkecil ukuran slot foto secara proporsional |

---

## 📂 Struktur Direktori Proyek

```
photobooth/
├── .env.example               # Template environment configuration
├── main.go                    # Entry point: DB, Seed Admin, Routes, & HTTP Server
├── run_kiosk.bat              # One-click Windows Kiosk Launcher
├── core/                      # Fondasi arsitektur sistem
│   ├── config.go              # Env loader
│   ├── database.go            # GORM + Pure-Go SQLite singleton
│   ├── router.go              # Gin engine & static files
│   └── template.go            # Custom isolated multi-template renderer
├── middleware/                # JWT Cookie Authentication
├── modules/                   # Modul bisnis independen
│   ├── auth/                  # Otentikasi & Landing page Neobrutalism
│   ├── frames/                # Visual Studio, CRUD Frame & Slot Coordinate Mapping
│   ├── kiosk/                 # Antarmuka Kiosk, WebRTC, Filters & Stickers
│   ├── sessions/              # Log sesi foto, komposit, softfile download, & ZIP export
│   ├── settings/              # Konfigurasi nama event, countdown, & driver printer
│   └── vouchers/              # Tiket masuk, kupon bazaar, & validasi kode
├── storage/                   # Penyimpanan file lokal
│   ├── frames/                # Template PNG frame overlay
│   └── photos/                # Arsip raw shots & final composite photos
└── template/meridian/         # Aset tampilan & tema Meridian
    ├── assets/                # CSS, JS, Fonts, dan Icons
    ├── layouts/master.html    # Master dashboard layout
    └── pages/                 # Halaman HTML (landing, kiosk, frames, download, dsb.)
```

---

## 💡 Panduan Penggunaan di Hari Acara (Event Checklist)

1. **Setup Jaringan Lokal:** Nyalakan router WiFi portabel atau hotspot HP (tidak perlu ada kuota internet). Hubungkan laptop photobooth dan perangkat pengunjung ke WiFi tersebut.
2. **Cek Alamat IP:** Buka `cmd`, ketik `ipconfig`, catat alamat IPv4 (contoh `192.168.1.15`), lalu pastikan kolom `LOCAL_IP` di menu **Settings** sudah sesuai.
3. **Posisikan Kamera:** Pasang webcam / kamera setinggi mata, sesuaikan pencahayaan (ring light disarankan).
4. **Buka Kiosk Mode:** Klik dua kali `run_kiosk.bat` untuk mengunci layar ke antarmuka siswa.
5. **Siapkan Printer:** Pastikan kertas foto dan tinta/ribbon terpasang. Aktifkan opsi *Auto-Print* di menu Settings jika ingin cetak instan setiap kali foto selesai.

---

## 📜 Lisensi & Kontribusi
Dibuat dengan ❤️ untuk kemudahan dokumentasi sekolah, kampus, dan komunitas kreatif. Bebas dimodifikasi dan dikembangkan untuk keperluan acara Anda!
