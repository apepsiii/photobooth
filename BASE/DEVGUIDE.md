Berikut adalah draf lengkap untuk file **`devguide.md`**. File ini dirancang khusus sebagai "buku manual" atau *system prompt* bagi AI. Bahasa dan instruksinya dibuat sangat terstruktur, definitif, dan menggunakan pendekatan *step-by-step* agar AI tidak melakukan halusinasi kode atau merusak arsitektur dasar.

Silakan salin teks di bawah ini dan simpan sebagai `devguide.md` di dalam folder `BASE/`.

---

```markdown
# DEVGUIDE.md - AI Developer Guide & Implementation Manual
## Project: Go-Gin-SQLite-GORM Boilerplate Framework (Base Project)

Dokumen ini adalah panduan resmi bagi AI Agent (Co-Developer) dalam membangun, mengonfigurasi, dan mengembangkan fitur pada boilerplate ini. Semua instruksi di bawah ini bersifat **MUTLAK** dan harus diikuti setiap kali melakukan *code generation* atau modifikasi proyek.

---

## ATURAN DASAR UNTUK AI (AI CORE RULES)
1. **Modularity is King:** Jangan pernah mencampur logika bisnis ke dalam folder `core/` atau file utama `main.go`. Semua fitur baru wajib dibuat di dalam folder `modules/`.
2. **Do Not Break the Base:** Core framework (Database initialization, Router setup, Auth Middleware) yang sudah berjalan tidak boleh diubah ketika sedang mengerjakan PRD baru, kecuali diminta secara eksplisit oleh pengguna.
3. **Idiomatic Go:** Gunakan standar penulisan Go yang baik (PascalCase untuk public struct/func, camelCase untuk private). Kembalikan `error` sebagai *return value* alih-alih melakukan `panic()`.
4. **Context Awareness:** Selalu baca `ABOUT.md` di folder `BASE/` untuk memahami framework, lalu baca file di folder `CONTEXT/` untuk memahami logika bisnis spesifik proyek yang sedang dikerjakan.

---

## FASE 1: SETUP BOILERPLATE AWAL (CORE FRAMEWORK)
Jika pengguna meminta untuk "Build Boilerplate Awal", lakukan langkah-langkah setup berikut secara berurutan:

### 1. Inisialisasi Proyek
* Jalankan `go mod init <module_name>`
* Install dependencies wajib: 
  * `github.com/gin-gonic/gin`
  * `gorm.io/gorm`
  * `gorm.io/driver/sqlite`
  * `github.com/joho/godotenv`
  * `github.com/golang-jwt/jwt/v5` (atau session manager sejenis)

### 2. Struktur Direktori Standar
Buat struktur direktori persis seperti berikut:
```text
.
├── main.go
├── core/
│   ├── config.go      # Load .env
│   ├── database.go    # Setup GORM & SQLite
│   ├── router.go      # Setup Gin & Global Middleware
│   └── logger.go      # Setup Logging
├── middleware/
│   └── auth.go        # JWT / Session checker
├── modules/           # (Kosong pada tahap awal)
├── models/
│   └── user.go        # Model dasar untuk Autentikasi
└── template/
    └── meridian/      # Master template HTML

```

---

## FASE 2: IMPLEMENTASI FITUR INTI SECARA DETAIL

### A. Setup Database (GORM + SQLite)

* **File:** `core/database.go`
* **Instruksi AI:**
1. Buat fungsi `InitDB()` yang me-return `*gorm.DB`.
2. Gunakan `gorm.Open(sqlite.Open("app.db"), &gorm.Config{})`.
3. Aktifkan koneksi PRAGMA untuk performa SQLite: `db.Exec("PRAGMA foreign_keys = ON")` dan `db.Exec("PRAGMA journal_mode = WAL")`.
4. Jalankan `db.AutoMigrate(&models.User{})` di dalam fungsi inisialisasi.



### B. Autentikasi & Keamanan (Auth Setup)

* **Model:** `models/user.go` harus memiliki field standar: `ID`, `Username`, `PasswordHash`, `Role`, `CreatedAt`, `UpdatedAt`.
* **Flow Autentikasi (Server-Side HTML):**
1. Buat fungsi Login yang memverifikasi password (gunakan `bcrypt`).
2. Setelah sukses, buat JWT token atau Session dan simpan ke dalam HTTP-Only Cookie (karena ini aplikasi berbasis *template HTML*, cookie lebih aman dan otomatis terkirim).


* **Middleware:** `middleware/auth.go`
1. Buat fungsi `RequireAuth()` yang membaca cookie `auth_token`.
2. Jika token tidak ada atau tidak valid, redirect user ke `/login` (`c.Redirect(http.StatusFound, "/login")`).
3. Jika valid, simpan user ID ke *context* (`c.Set("userID", id)`).



### C. Integrasi Template Engine (Meridian UI)

* **File:** `core/router.go`
* **Instruksi AI:**
1. Inisialisasi router: `r := gin.Default()`.
2. Load semua template: `r.LoadHTMLGlob("template/meridian//*.html")`.
3. Load static assets: `r.Static("/assets", "./template/meridian/assets")`.
4. Atur grouping route. Pisahkan antara route publik (`/login`, `/register`) dan route privat yang diapit oleh `middleware.RequireAuth()` (seperti `/dashboard`).



---

## FASE 3: WORKFLOW PENGEMBANGAN PROYEK BARU (MODULES)

Ketika pengguna memberikan `CONTEXT/PRD.md` baru (misalnya fitur "Manajemen Produk"), AI wajib membuat fitur tersebut dengan pendekatan *Domain-Driven Design (DDD)* sederhana di dalam folder `modules/`.

### Struktur per Modul (Contoh: `modules/product/`)

Setiap modul baru **wajib** memiliki 3 layer:

1. **`model.go`**: Definisi `struct` GORM untuk tabel database.
2. **`service.go`**: Logika bisnis (CRUD, kalkulasi). Service harus menerima dependensi `*gorm.DB`.
3. **`handler.go`**: Gin HTTP handler. Bertugas menerima request, memanggil service, dan me-render HTML atau mengembalikan JSON.
4. **`routes.go`**: Fungsi untuk meregistrasi rute modul ke router utama. Contoh: `func RegisterProductRoutes(rg *gin.RouterGroup, db *gorm.DB)`.

### Langkah Eksekusi AI saat membuat Modul Baru:

1. **Analisis PRD:** Baca requirement dari PRD. Tentukan field apa saja yang dibutuhkan di database.
2. **Buat Model:** Tulis `model.go` dengan tag GORM yang tepat.
3. **Update AutoMigrate:** Tambahkan model baru tersebut ke dalam `AutoMigrate()` di `core/database.go` secara hati-hati.
4. **Buat Service & Handler:** Tulis logika CRUD. Pastikan jika me-render halaman, gunakan nama template yang sesuai (misal: `c.HTML(200, "product_list.html", gin.H{"data": products})`).
5. **Registrasi Route:** Panggil fungsi `Register[NamaModul]Routes` di dalam file `main.go` atau `core/router.go`.

```

```