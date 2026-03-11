# Penjelasan Detail: PostgresStore + Migrasi (Database PostgreSQL Lokal)

Dokumen ini menjelaskan **secara teknis tapi sederhana** apa itu PostgresStore, apa itu migrasi, dan bagaimana kita memakai PostgreSQL di komputer sendiri (lokal).

---

## 1. Konteks: Di Mana Posisi Kita Sekarang?

Saat ini aplikasi punya:

- **ItemStore** = “tempat menyimpan data Item” (interface).
- **MemoryStore** = implementasi yang menyimpan di **RAM** (memori). Data hanya ada selama program jalan; begitu server dimatikan, data hilang.

Kita akan menambah:

- **PostgresStore** = implementasi yang menyimpan di **PostgreSQL** (database di disk). Data tetap ada meskipun server dimatikan.
- **Migrasi** = skrip SQL yang dipakai untuk “membuat atau mengubah struktur tabel” di database. Dijalankan sekali (atau per versi) agar tabel siap dipakai.

---

## 2. Apa Itu PostgreSQL (Postgres)?

- **PostgreSQL** = salah satu **database relasional** (data disimpan dalam bentuk **tabel**: baris dan kolom).
- **Lokal** = PostgreSQL di-install dan jalan di **komputer Anda**, bukan di cloud. Kita sambung ke sana dengan connection string (alamat, user, password, nama database).

Analogi singkat:

- **MemoryStore** = catatan di whiteboard: cepat, tapi begitu listrik mati (program berhenti), tulisan hilang.
- **PostgreSQL** = catatan di buku (file di disk): tetap ada setelah program berhenti, dan bisa dipakai lagi saat program jalan lagi.

---

## 3. Konsep Penting: Tabel, Kolom, Baris

Di database kita punya **tabel**. Satu tabel punya:

- **Nama tabel** — misalnya `items`.
- **Kolom (field)** — misalnya `id`, `name`, `created_at`. Setiap kolom punya **tipe data** (teks, angka, tanggal, dll.).
- **Baris (row)** — satu baris = satu “Item”. Setiap baris punya nilai untuk tiap kolom.

Contoh tabel `items`:

| id | name   | created_at          |
|----|--------|---------------------|
| 1  | Buku   | 2025-03-10 10:00:00 |
| 2  | Pensil | 2025-03-10 10:01:00 |

- **id** = pengenal unik (biasanya primary key).
- **name** = nama item.
- **created_at** = waktu pembuatan.

Agar tabel ini ada di database, kita harus “membuat”nya. Caranya: jalankan **perintah SQL** (CREATE TABLE ...). Kumpulan perintah SQL itulah yang nanti kita simpan dalam **file migrasi**.

---

## 4. Apa Itu Migrasi (Migration)?

- **Migrasi** = satu atau beberapa **file** berisi perintah SQL yang:
  - **Membuat** tabel baru (CREATE TABLE), atau
  - **Mengubah** tabel yang sudah ada (ALTER TABLE, ADD COLUMN, dll.).

Kenapa pakai file terpisah (bukan hardcode di kode Go)?

- Riwayat perubahan skema tercatat (siapa tambah kolom kapan).
- Bisa dijalankan ulang di lingkungan lain (dev, staging, production) dengan urutan yang sama.
- Rollback atau penambahan migrasi baru bisa dikelola rapi.

Urutan umum:

1. **Migrasi pertama** — misalnya: “Buat tabel `items` dengan kolom id, name, created_at.”
2. Nanti kalau butuh kolom baru (misalnya `updated_at`), kita buat **file migrasi kedua**: “Tambahkan kolom `updated_at` ke tabel `items`.”

Saat aplikasi (atau script) jalan, kita **jalankan migrasi yang belum pernah dijalankan** (biasanya dicatat di tabel khusus, misalnya `schema_migrations`). Jadi yang dijalankan hanya yang baru.

Untuk awal, kita bisa sederhana: **satu file migrasi** yang isinya cuma **CREATE TABLE items (...)**. Tidak wajib pakai library migrasi dulu; yang penting konsepnya: “ada file SQL yang membuat struktur tabel.”

---

## 5. Apa Itu PostgresStore?

- **PostgresStore** = struct di Go yang **mengimplementasi interface ItemStore** (sama seperti MemoryStore), tapi:
  - **List()** → jalankan query SQL `SELECT * FROM items`, baca hasil ke slice `[]*entity.Item`.
  - **GetByID(id)** → `SELECT * FROM items WHERE id = $1`, return satu item atau nil.
  - **Create(name)** → `INSERT INTO items (name, created_at) VALUES ($1, $2) RETURNING id, name, created_at`, lalu return item yang baru dibuat.
  - **Update(id, name)** → `UPDATE items SET name = $1 WHERE id = $2`, cek apakah ada baris ter-update; kalau tidak, return nil (not found).
  - **Delete(id)** → `DELETE FROM items WHERE id = $1`, cek jumlah baris terhapus; kalau 0, return ErrNotFound.

Kita menyimpan **pointer ke database** (`*sql.DB`) di dalam PostgresStore. Semua operasi baca/tulis memakai `db.Query`, `db.QueryRow`, `db.Exec`, dll. dari package `database/sql`.

**Parameterized query** (`$1`, `$2`) = nilai dari variabel (id, name) kita kirim terpisah, bukan digabung ke string SQL. Ini mencegah **SQL injection** dan cara yang benar di Go.

---

## 6. Alur Kerja Secara Keseluruhan

```
1. PostgreSQL sudah jalan di PC Anda (service / proses).
2. Anda punya satu database (misalnya nama: blueprint).
3. Jalankan migrasi (sekali): file SQL membuat tabel items.
4. Aplikasi start → baca DB_DSN dari env → Open(dsn) → Ping().
5. Buat PostgresStore(db). Di main, pakai PostgresStore sebagai ItemStore (bukan MemoryStore).
6. Setiap request CRUD → handler memanggil store → PostgresStore menjalankan SQL ke tabel items.
7. Data tersimpan di disk (PostgreSQL). Restart aplikasi → data tetap ada.
```

---

## 7. Yang Perlu Disiapkan di Komputer Anda (PostgreSQL Lokal)

### 7.1 Install PostgreSQL

- **Windows:** download installer dari [postgresql.org](https://www.postgresql.org/download/windows/) atau pakai winget: `winget install PostgreSQL.PostgreSQL`.
- Saat install, Anda akan diminta set **password** untuk user **postgres**. Simpan password ini; dipakai di connection string.

### 7.2 Pastikan PostgreSQL Jalan

- Setelah install, biasanya PostgreSQL jalan sebagai **service**. Di Windows: Services → cari “PostgreSQL”, status harus Running.
- Atau dari command line (jika `psql` ada di PATH):  
  `psql -U postgres -c "SELECT 1"`  
  Kalau tidak error, koneksi ke server lokal berhasil.

### 7.3 Buat Database Khusus untuk Proyek Ini

Kita pakai satu database terpisah (bukan database default `postgres`) agar rapi:

- Buka **pgAdmin** (GUI yang ikut terinstall) atau pakai **psql**.
- Contoh dengan psql:
  - Login: `psql -U postgres`
  - Buat database: `CREATE DATABASE blueprint;`
  - Keluar: `\q`

Nanti **connection string (DSN)** kita arahkan ke database ini.

### 7.4 Connection String (DSN)

Format umum:

```
postgres://USER:PASSWORD@HOST:PORT/DATABASE?sslmode=disable
```

Untuk **lokal**:

- **USER** = `postgres` (atau user yang Anda buat).
- **PASSWORD** = password yang Anda set saat install.
- **HOST** = `localhost` (karena di PC sendiri).
- **PORT** = biasanya `5432` (default PostgreSQL).
- **DATABASE** = `blueprint` (nama database yang kita buat).
- **sslmode=disable** = untuk development lokal, SSL tidak dipakai.

Contoh (ganti `PASSWORD` dengan password Anda):

```
postgres://postgres:PASSWORD@localhost:5432/blueprint?sslmode=disable
```

Di aplikasi kita, nilai ini disimpan di **environment variable** `DB_DSN`. Jadi kita tidak menaruh password di kode, hanya di env (atau di file .env yang tidak di-commit).

---

## 8. Ringkasan Istilah

| Istilah        | Arti singkat |
|----------------|--------------|
| PostgreSQL     | Database relasional (data dalam tabel); bisa jalan di PC Anda (lokal). |
| Tabel          | Struktur data dengan kolom dan baris (misalnya tabel `items`). |
| Migrasi        | File SQL untuk membuat/mengubah struktur tabel; dijalankan berurutan. |
| DSN            | Connection string: user, password, host, port, nama database. |
| PostgresStore  | Implementasi ItemStore yang menyimpan data ke tabel di PostgreSQL. |
| DB_DSN         | Env berisi connection string; dibaca aplikasi untuk sambung ke PostgreSQL. |

---

## 9. Urutan Implementasi (Preview)

Kalau kita lanjut ke implementasi, urutan yang masuk akal:

1. **Folder dan file migrasi** — misalnya `migrations/001_create_items_table.sql` berisi `CREATE TABLE items (...);`.
2. **Cara menjalankan migrasi** — script kecil (Go atau SQL) yang jalan sekali saat development, atau dipanggil dari main saat startup (opsional).
3. **PostgresStore** — `internal/store/postgres.go`: struct yang pegang `*sql.DB`, implementasi List, GetByID, Create, Update, Delete dengan query SQL.
4. **Main** — jika `DB_DSN` ada: jalankan migrasi (jika belum), buat `PostgresStore(db)`, pakai sebagai `ItemStore`. Jika `DB_DSN` kosong: pakai `MemoryStore` seperti sekarang.

Dengan begitu Anda bisa mencoba PostgreSQL lokal dengan langkah yang jelas dan kode yang mengikuti pola yang sudah ada (ItemStore).

Jika Anda mau, langkah berikutnya bisa kita lakukan **implementasi nyata** (file migrasi + PostgresStore + perubahan di main) langkah demi langkah seperti sebelumnya.
