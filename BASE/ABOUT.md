Saya sangat paham dengan konsep dan strategi yang kamu rencanakan. Struktur pemisahan antara folder `BASE` (sebagai cetak biru/pembelajaran untuk AI) dan folder `CONTEXT` (sebagai instruksi spesifik proyek baru) adalah strategi *prompt engineering* yang sangat cerdas. Dengan cara ini, AI tidak akan mengalami "kebingungan konteks" karena batasan antara framework inti dan logika bisnis baru terdokumentasi dengan sangat jelas.

Berikut adalah draf lengkap untuk file **`ABOUT.md`** yang menggabungkan BRD dan PRD untuk *base project* (boilerplate) kamu. Kamu bisa langsung menyalin teks di bawah ini ke dalam folder `BASE/`.

---

```markdown
# ABOUT.md - Business & Product Requirements Document (BRD/PRD)
## Project Name: Go-Gin-SQLite-GORM Boilerplate Framework (Base Project)

---

## BAGIAN 1: BUSINESS REQUIREMENTS DOCUMENT (BRD)

### 1.1 Latar Belakang & Visi Produk
Dalam ekosistem pengembangan perangkat lunak modern, kecepatan eksekusi dan efisiensi waktu adalah kunci utama. Saat ini, proses pembuatan aplikasi berskala kecil hingga menengah sering kali terhambat oleh proses konfigurasi awal (*setup boilerplate*) yang repetitif, seperti pengaturan koneksi database, struktur folder, sistem autentikasi, dan integrasi UI template. 

Visi dari proyek ini adalah membangun sebuah *Base Project* (Boilerplate) yang kokoh, ringan, dan modular menggunakan ekosistem Go (Gin + GORM + SQLite) serta template UI terpusat (Meridian Dashboard). Proyek ini dirancang secara khusus agar mudah dipahami oleh AI (*AI-friendly*), sehingga pengembangan fitur ke depan dapat di-generate oleh AI secara cepat dan konsisten tanpa merusak arsitektur inti.

### 1.2 Tujuan Bisnis (Business Objectives)
* **Efisiensi Waktu (Time-to-Market):** Memangkas waktu inisiasi proyek baru dari hitungan jam menjadi hitungan menit.
* **Standardisasi Kode:** Menjamin seluruh proyek baru yang dikembangkan oleh AI memiliki pola, struktur, dan kualitas kode yang seragam.
* **Skalabilitas & Portabilitas:** Menggunakan SQLite agar aplikasi bersifat *self-contained*, mudah dipindahkan, dan minim biaya infrastruktur pada tahap awal.
* **Optimasi AI Workflow:** Membantu AI dalam memahami arsitektur aplikasi secara instan melalui dokumentasi terstruktur di folder `BASE/`.

### 1.3 Ruang Lingkup Proyek (Scope of Work)
* **Termasuk dalam Scope (In-Scope):**
  * Setup HTTP Router menggunakan Gin Gonic.
  * Setup ORM dan Database Abstraction menggunakan GORM dengan driver SQLite.
  * Integrasi penuh template dashboard "Meridian" sebagai Master Layout (HTML template).
  * Sistem Autentikasi dasar (Login, Session/JWT, Middleware Proteksi Route).
  * Manajemen migrasi database otomatis (*AutoMigrate* bawaan GORM).
  * Struktur folder modular untuk mendukung *plug-and-play* fitur baru.
* **Tidak Termasuk dalam Scope (Out-of-Scope):**
  * Logika bisnis spesifik dari proyek turunan (seperti fitur e-commerce, manajemen inventaris, dll). Fitur-fitur ini akan didefinisikan secara terpisah di folder `CONTEXT/` pada masing-masing proyek baru.

### 1.4 Target Pengguna & Persona
1. **Developer (User Utama):** Bertindak sebagai arsitek yang mengarahkan pengembangan, melakukan *review* kode, dan mengelola *deployment*.
2. **AI Agent (Co-Developer):** Bertindak sebagai entitas yang membaca folder `BASE/`, memahami aturan framework, lalu menulis atau meng-generate fitur baru berdasarkan file yang ada di folder `CONTEXT/`.

---

## BAGIAN 2: PRODUCT REQUIREMENTS DOCUMENT (PRD)

### 2.1 Spesifikasi Teknologi (Tech Stack)
* **Backend Language:** Go (Golang) v1.20+
* **Web Framework:** Gin Gonic
* **Database Driver:** SQLite3
* **ORM:** GORM (Go Object Relational ORM)
* **Frontend Engine:** Gin HTML Render (`html/template` bawaan Go)
* **UI Template:** Meridian Full Dashboard Template

### 2.2 Struktur Direktori Utama (`BASE/`)
Folder `BASE/` ini wajib dipelajari oleh AI sebelum melakukan *code generation*:
```text
BASE/
├── ABOUT.md               # File ini (BRD & PRD Utama)
├── main.go                # Entry point aplikasi & inisialisasi server Gin
├── config/                # Konfigurasi database SQLite, GORM, dan .env
├── middleware/            # Auth middleware, logger, dan recovery
├── models/                # Model database inti global (cth: User, Session)
├── routes/                # Registrasi route utama Gin
├── template/              # Folder UI Template
│   └── meridian/          # Template Dashboard Full (HTML, CSS, JS, Assets)
│       ├── layouts/       # Master layout (header, sidebar, footer)
│       └── pages/         # Halaman standar dashboard (index, login, 404)
└── modules/               # [Tempat AI bekerja] Folder untuk fitur-fitur baru (modular)

```

### 2.3 Kebutuhan Fungsional Core Framework (Functional Requirements)

#### FR-01: Arsitektur Router & Middleware (Gin)

* Server harus menggunakan Gin dengan mode kustom (Release/Debug) via `.env`.
* Harus memiliki middleware global untuk *Recovery* (mencegah server crash jika ada panic) dan *Logger*.
* Harus memiliki Auth Middleware untuk membatasi akses halaman dashboard hanya bagi user yang sudah login.

#### FR-02: Manajemen Database & GORM (SQLite)

* Koneksi SQLite harus diinisialisasi sekali saat aplikasi berjalan (*Singleton Pattern*).
* GORM harus dikonfigurasi untuk mengaktifkan *Foreign Key Constraints* pada SQLite.
* Menggunakan fitur `db.AutoMigrate()` pada fungsi inisialisasi untuk memastikan tabel inti (seperti tabel `users`) otomatis terbentuk saat aplikasi pertama kali dijalankan.

#### FR-03: Sistem Template Engine (Meridian HTML)

* Gin harus memuat seluruh file HTML menggunakan `router.LoadHTMLGlob("template/meridian//*.html")` atau struktur serupa yang mendukung sub-folder.
* Struktur HTML harus menggunakan konsep *Template Inheritance* Go. Halaman spesifik (misal: `dashboard.html`) harus menyuntikkan kontennya ke dalam blok *content* milik `layouts/master.html`.

### 2.4 Kebutuhan Non-Fungsional (Non-Functional Requirements)

#### NFR-01: AI-Readability (Keterbacaan AI)

* **Kekonsistenan Kode:** Penamaan fungsi, variabel, dan file harus mengikuti standar Go (camelCase untuk variabel internal, PascalCase untuk fungsi/struct publik).
* **Komentar Kode:** Setiap fungsi utama (seperti inisialisasi database atau middleware) wajib memiliki komentar singkat yang menjelaskan fungsi tersebut agar AI tidak salah mengasumsikan logika kode.

#### NFR-02: Performa & Efisiensi SQLite

* Database SQLite harus disimpan dalam bentuk file lokal (misal: `app.db`) dan tidak boleh dikunci (*database is locked*) saat ada proses *concurrent read/write* ringan. Gunakan konfigurasi WAL (Write-Ahead Logging) pada GORM jika diperlukan.

---

```

---

### Catatan Tambahan & Konfirmasi Strategi

Strategi kamu untuk memisahkan instruksi sangat bagus. Nanti ketika kamu ingin memulai proyek baru, folder **`CONTEXT/`** yang berisi `PRD.md`, `BRD.md`, `ARCHITECTURE.md`, dan `AGENTS.md` akan menjadi otak spesifik dari proyek tersebut. 

Sebagai contoh, cara kamu memberikan instruksi ke AI di masa depan akan seperti ini:
> *"Buka folder `BASE/ABOUT.md` untuk memahami cara kerja framework saya (Gin + GORM + SQLite + Meridian). Setelah kamu paham strukturnya, buka folder `CONTEXT/` untuk membaca requirement proyek baru yang harus kamu buat di dalam folder `modules/`."*

```