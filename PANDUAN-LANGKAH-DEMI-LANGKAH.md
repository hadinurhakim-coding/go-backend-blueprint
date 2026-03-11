# Panduan Belajar: Go Backend Blueprint — Langkah demi Langkah

Panduan ini untuk Anda yang masih sangat awam. Setiap langkah kecil akan dijelaskan: **apa yang kita lakukan**, **mengapa**, dan **bagaimana cara mengecek bahwa langkah berhasil**.

---

## Peta Besar (Apa Saja yang Akan Kita Lakukan)

| # | Langkah | Singkatnya | Status |
|---|--------|------------|--------|
| 1 | **Inisialisasi modul Go** | Memberi nama resmi ke proyek dan mengaktifkan manajemen dependensi. | ✅ |
| 2 | **Program Go paling sederhana** | Satu file yang hanya menampilkan teks di terminal. | ✅ |
| 3 | **Menjalankan program** | Memastikan Go terpasang dan kode bisa dijalankan. | ✅ |
| 4 | **Server HTTP minimal** | Satu endpoint yang menjawab "Halo" ketika diakses lewat browser. | ✅ |
| 5 | **Struktur folder** | Mengatur kode ke folder `cmd`, `internal`, dll. | ✅ |
| 6 | **Apa selanjutnya** | REST API, database, Redis, Docker, dll. (dilakukan nanti). | 📌 |
| 7 | **Endpoint REST pertama** | GET /health mengembalikan JSON dan status code 200. | ✅ |
| 8 | **CRUD satu entitas (Item)** | GET/POST /items, GET/PUT/DELETE /items/:id dengan in-memory store. | ✅ |
| 9 | **Persiapan koneksi database** | Config dari env (PORT, DB_DSN), package database (Open, Ping, Close). | ✅ |

Kita akan mengerjakan **satu langkah dalam satu waktu**. Jangan terburu-buru; pastikan Anda paham sebelum lanjut.

---

## Langkah 1: Inisialisasi Modul Go

### Apa yang kita lakukan?
Kita menjalankan perintah **`go mod init go-backend-blueprint`** di dalam folder proyek.

### Apa itu "modul"?
- Di Go, **modul** = satu proyek dengan nama dan versi dependensi (library) yang dipakai.
- File **`go.mod`** yang akan terbentuk berisi:
  - Nama modul (misalnya `go-backend-blueprint`).
  - Versi Go yang dipakai.
  - Daftar library dari luar (nanti akan bertambah saat kita pakai library).

### Mengapa ini langkah pertama?
Tanpa `go.mod`, Go tidak menganggap folder ini sebagai "proyek" yang bisa di-build atau dijalankan dengan rapi. Semua proyek Go modern dimulai dengan `go mod init`.

### Apa yang akan terbentuk?
Satu file baru: **`go.mod`** di akar proyek.

### ✅ Langkah 1 sudah dilakukan

Perintah yang dijalankan:
```bash
go mod init go-backend-blueprint
```

File **`go.mod`** yang terbentuk isinya kira-kira:
```text
module go-backend-blueprint

go 1.26.0
```

- **Baris 1:** `module go-backend-blueprint` = nama resmi proyek ini. Import path paket di dalam proyek akan diawali dengan nama ini (misalnya `go-backend-blueprint/cmd/api`).
- **Baris 3:** `go 1.26.0` = versi bahasa Go yang dipakai di proyek ini (diselaraskan dengan versi Go di komputer Anda).

**Cara mengecek:** Pastikan di folder proyek ada file `go.mod`. Buka dengan editor dan bandingkan dengan isi di atas.

---

## Langkah 2: Program Go Paling Sederhana

### Apa yang kita lakukan?
Kita membuat satu file bernama **`main.go`** di akar proyek. Isinya hanya: menjalankan fungsi `main` dan menampilkan teks **"Halo"** ke terminal.

### Mengapa pakai nama `main.go`?
- Di Go, program yang bisa dijalankan (bukan sekadar library) **harus** punya fungsi bernama `main` di dalam **package main**.
- Nama file bisa apa saja (boleh `main.go`, `server.go`, dll.), yang penting isinya ada `package main` dan `func main()`.

### Apa arti setiap baris?

| Baris | Kode | Arti singkat |
|-------|------|--------------|
| 1 | `package main` | Kode ini bagian dari paket bernama `main`. Hanya paket `main` yang boleh punya fungsi `main` dan bisa dijadikan program yang dijalankan. |
| 2 | *(kosong)* | Pemisah antara deklarasi paket dan import. |
| 3 | `import "fmt"` | Kita akan memakai paket **fmt** (format I/O dari standard library Go). Paket ini berisi fungsi seperti `Println` untuk mencetak teks. |
| 4 | *(kosong)* | Pemisah. |
| 5 | `func main() {` | Ini **titik masuk** program. Saat kita jalankan program, Go akan memanggil fungsi `main()` sekali. |
| 6 | `fmt.Println("Halo")` | Cetak teks **Halo** ke layar (terminal), lalu pindah baris. `Println` = "print line". |
| 7 | `}` | Akhir dari fungsi `main`. |

### ✅ Langkah 2 sudah dilakukan

File **`main.go`** yang dibuat isinya:

```go
package main

import "fmt"

func main() {
	fmt.Println("Halo")
}
```

**Cara mengecek:** Jalankan di folder proyek:
```bash
go run main.go
```
Jika berhasil, di terminal akan muncul satu baris: **Halo**. Artinya program Go kita sudah jalan. *(File `main.go` di akar sudah dipindahkan ke `cmd/api/main.go` di Langkah 5.)*

---

## Langkah 3: Menjalankan Program

### Apa yang kita lakukan?
Kita memastikan program Go bisa dijalankan dengan perintah **`go run`** dan output yang diharapkan muncul.

### Apa itu `go run`?
- **`go run main.go`** (atau **`go run ./cmd/api`**) = kompilasi kode Go lalu langsung jalankan, tanpa menyimpan file executable.
- Cocok untuk development. Untuk produksi nanti biasanya pakai **`go build`** untuk menghasilkan satu file binary.

### ✅ Langkah 3 sudah dilakukan

Kita sudah menjalankan program dan memastikan output **Halo** muncul. Setelah struktur pindah ke `cmd/api`, cara menjalankan menjadi:
```bash
go run ./cmd/api
```
Program akan jalan (dan untuk server HTTP akan mendengarkan di port 8080). Menghentikan: **Ctrl+C** di terminal.

---

## Langkah 4: Server HTTP Minimal

### Apa yang kita lakukan?
Program tidak lagi sekadar mencetak "Halo" ke terminal, tetapi **menjadi server** yang:
1. **Mendengarkan** di port **8080** (komputer siap menerima koneksi HTTP).
2. Ketika ada yang membuka **http://localhost:8080/** di browser (atau mengirim request GET ke `/`), server **membalas** dengan teks **Halo**.

### Konsep singkat
- **Port** = nomor yang mengidentifikasi "pintu" mana di komputer yang dipakai untuk layanan tertentu (8080 sering dipakai untuk development).
- **Endpoint /** = path utama. Request ke `http://localhost:8080/` akan memanggil handler yang kita daftarkan untuk `/`.
- **Handler** = fungsi yang menerima request dan menulis response (di sini: menulis "Halo" ke response).

### ✅ Langkah 4 sudah dilakukan

- **`internal/handler/handler.go`** berisi fungsi **`Halo`** yang menulis "Halo" ke `http.ResponseWriter`.
- **`cmd/api/main.go`** mendaftarkan `handler.Halo` untuk path `/` dan memanggil **`http.ListenAndServe(":8080", nil)`** agar server mendengarkan di port 8080.

**Cara mengecek:** Jalankan `go run ./cmd/api`, lalu buka browser ke **http://localhost:8080/** — halaman harus menampilkan **Halo**.

---

## Langkah 5: Struktur Folder (cmd/api + internal)

### Apa yang kita lakukan?
Kita mengatur ulang kode ke struktur yang rapi:
- **`cmd/api/`** = tempat **entrypoint** program (main package) untuk API server. Satu proyek bisa punya beberapa `cmd/...` (misalnya `cmd/api`, `cmd/worker`).
- **`internal/`** = kode yang **hanya dipakai di dalam proyek ini** (tidak diekspor ke proyek lain). Berisi handler, logic, akses DB, dll.

### Mengapa struktur ini?
- **cmd/** = konvensi di ekosistem Go: setiap "program yang bisa dijalankan" punya subfolder di `cmd/`.
- **internal/** = folder khusus di Go: paket di dalam `internal/` tidak bisa di-import oleh proyek lain. Jadi kode kita aman dari ketergantungan pihak luar.

### ✅ Langkah 5 sudah dilakukan

Struktur saat ini:

```text
go-backend-blueprint/
  cmd/
    api/
      main.go          ← entrypoint: daftar handler, ListenAndServe
  internal/
    handler/
      handler.go       ← fungsi Halo (handler HTTP)
  go.mod
  README.md
  PANDUAN-LANGKAH-DEMI-LANGKAH.md
```

- **`main.go`** di akar proyek sudah **dihapus**; titik masuk sekarang **`cmd/api/main.go`**.
- Handler dipindah ke **`internal/handler`** agar bisa dipakai ulang dan di-test terpisah.

**Cara mengecek:** Pastikan folder `cmd/api` dan `internal/handler` ada, lalu jalankan `go run ./cmd/api` — server harus jalan dan membalas "Halo" di http://localhost:8080/.

---

## Langkah 6: Apa Selanjutnya?

Langkah 1–5 sudah selesai. Proyek sekarang punya:
- Modul Go yang terdefinisi (`go.mod`),
- Server HTTP minimal yang merespons "Halo" di `/`,
- Struktur folder yang siap dikembangkan (`cmd/api`, `internal/handler`).

**Langkah berikutnya (nanti)** bisa mencakup:
- Menambah endpoint REST (GET/POST resource, status code, JSON).
- Menghubungkan database (SQL/NoSQL), Redis, object storage (S3/GCS).
- Docker & Kubernetes, serta pipeline CI/CD.

Setiap topik itu bisa kita lakukan lagi **tahap demi tahap** dengan penjelasan yang sama rincinya.

---

## Langkah 7: Endpoint REST Pertama (GET /health)

### Apa yang kita lakukan?
Kita menambah satu endpoint baru: **GET /health**. Ketika client memanggil `http://localhost:8080/health` dengan method **GET**, server mengembalikan:
- **Body:** JSON `{"status":"ok"}`
- **Header:** `Content-Type: application/json`
- **Status code:** 200 (OK)

Jika client memakai method selain GET (misalnya POST), server mengembalikan **405 Method Not Allowed**.

### Mengapa endpoint "health"?
- Di dunia nyata, **health check** dipakai oleh load balancer, Kubernetes, atau monitoring untuk memastikan aplikasi masih hidup.
- Endpoint ini sederhana (tidak baca database) dan mengembalikan JSON — pola yang sama dipakai untuk banyak endpoint REST lain.

### Apa yang berubah di kode?

**1. `internal/handler/handler.go`**
- Struct **`healthResponse`** dengan field `Status`. Tag **`json:"status"`** membuat field itu jadi key `"status"` di JSON.
- Fungsi **`Health(w, r)`**:
  - Cek `r.Method != http.MethodGet` → jika bukan GET, kirim 405 dan return.
  - Set header `Content-Type: application/json`.
  - Set status code 200 dengan `w.WriteHeader(http.StatusOK)`.
  - Encode struct ke JSON lalu tulis ke `w` dengan `json.NewEncoder(w).Encode(...)`.

**2. `cmd/api/main.go`**
- Daftarkan **`http.HandleFunc("/health", handler.Health)`** sebelum **`"/", handler.Halo`**, agar request ke `/health` tidak tertangkap oleh handler `/`.

### ✅ Langkah 7 sudah dilakukan

**Cara mengecek:**
1. Jalankan server: `go run ./cmd/api`
2. Di browser buka **http://localhost:8080/health** → harus tampil `{"status":"ok"}`.
3. Di browser buka **http://localhost:8080/** → tetap tampil **Halo**.

---

## Langkah 8: CRUD Satu Entitas (Item)

### Apa yang kita lakukan?
Kita menambah **resource REST** bernama **Item** dengan operasi CRUD lengkap:
- **GET /items** → daftar semua item (response JSON array).
- **POST /items** → buat item baru (body JSON `{"name":"..."}`), response 201 Created.
- **GET /items/:id** → ambil satu item; 404 jika tidak ada.
- **PUT /items/:id** → update item (body `{"name":"..."}`); 404 jika tidak ada.
- **DELETE /items/:id** → hapus item; 204 No Content; 404 jika tidak ada.

Data sementara disimpan **di memori** (in-memory store). Nantinya bisa diganti dengan store yang pakai database tanpa mengubah handler, berkat **interface** `ItemStore`.

### Konsep penting
- **Entity** = struktur data bisnis (di sini: `Item` dengan ID, Name, CreatedAt).
- **Store (interface)** = kontrak untuk “menyimpan/mengambil” data. Handler hanya bergantung pada interface, bukan implementasi. Jadi kita bisa ganti in-memory dengan database nanti.
- **In-memory store** = pakai `map` + mutex; ID dibuat otomatis (counter). Cocok untuk development dan testing.

### File yang ditambah
- **`internal/entity/item.go`** — struct `Item` dan `NowFunc` (untuk waktu, bisa di-mock saat testing).
- **`internal/store/store.go`** — interface `ItemStore` (List, GetByID, Create, Update, Delete) dan `ErrNotFound`.
- **`internal/store/memory.go`** — implementasi in-memory (map, mutex, ID increment).
- **`internal/handler/items.go`** — `ItemsHandler` dengan `HandleItems` (GET/POST /items) dan `HandleItemByID` (GET/PUT/DELETE /items/:id). Parsing ID dari path dengan `strings.TrimPrefix(r.URL.Path, "/items/")`.

### ✅ Langkah 8 sudah dilakukan

**Cara mengecek (tanpa DB):**
1. Jalankan: `go run ./cmd/api`
2. **GET** http://localhost:8080/items → `[]`
3. **POST** http://localhost:8080/items dengan body `{"name":"Buku"}` → 201, response berisi item dengan `id`, `name`, `created_at`
4. **GET** http://localhost:8080/items → array berisi satu item
5. **GET** http://localhost:8080/items/1 → satu item
6. **PUT** http://localhost:8080/items/1 dengan body `{"name":"Buku Baru"}` → 200, item ter-update
7. **DELETE** http://localhost:8080/items/1 → 204; **GET** /items/1 → 404

Untuk POST/PUT bisa pakai Postman, Insomnia, atau PowerShell:  
`Invoke-RestMethod -Uri http://localhost:8080/items -Method Post -ContentType "application/json" -Body '{"name":"Buku"}'`

---

## Langkah 9: Persiapan Koneksi Database

### Apa yang kita lakukan?
Kita **menyiapkan** koneksi database (belum dipakai untuk CRUD; data Item masih in-memory):
- **Config dari environment** — baca `PORT` (default 8080) dan `DB_DSN` (connection string PostgreSQL). Jika `DB_DSN` kosong, aplikasi tetap jalan tanpa DB.
- **Package `internal/database`** — fungsi `Open(dsn)`, `Ping(db)`, `Close(db)` menggunakan `database/sql` dan driver **PostgreSQL** (`github.com/lib/pq`). Di `main`, jika `DB_DSN` diset: buka koneksi, ping, dan `defer Close(db)` agar koneksi siap untuk langkah berikutnya (misalnya implementasi `PostgresStore`).

### Mengapa pakai interface ItemStore?
Agar nanti kita bisa buat **`internal/store/postgres.go`** yang mengimplementasi `ItemStore` dengan query ke database, lalu di `main` pilih: pakai `MemoryStore` atau `PostgresStore(db)`. Handler tidak berubah.

### File yang ditambah
- **`internal/config/config.go`** — struct `Config`, fungsi `FromEnv()` membaca `PORT` dan `DB_DSN`.
- **`internal/database/database.go`** — `Open(dsn)`, `Ping(db)`, `Close(db)`; import `_ "github.com/lib/pq"` untuk register driver.
- **`cmd/api/main.go`** — panggil `config.FromEnv()`, listen di `":" + cfg.Port`, dan jika `cfg.DBDSN != ""` maka `database.Open`, `Ping`, `defer Close(db)`.

### ✅ Langkah 9 sudah dilakukan

**Cara mengecek:**
- Tanpa DB: jalankan `go run ./cmd/api` seperti biasa. Server listen di port dari env `PORT` (default 8080). CRUD Item tetap pakai in-memory.
- Dengan DB: set env `DB_DSN` ke connection string PostgreSQL, misalnya  
  `postgres://user:password@localhost:5432/mydb?sslmode=disable`  
  Lalu jalankan lagi; di log harus muncul **"database: connected"**. Jika PostgreSQL belum jalan atau DSN salah, aplikasi akan exit dengan error (sesuai design: fail fast saat startup).

**Langkah berikutnya:** Lihat **POSTGRESSTORE-DAN-MIGRASI.md** untuk penjelasan detail dan **implementasi PostgresStore + migrasi** (sudah ada di proyek). Set env `DB_DSN` ke connection string PostgreSQL lalu jalankan `go run ./cmd/api` — data Item akan tersimpan di database. Setelah itu bisa lanjut: testing, middleware, Docker, dll.
