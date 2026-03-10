## go-backend-blueprint

Blueprint proyek backend Golang yang **production-ready**, **scalable**, dan **tanpa framework besar** (mengandalkan standard library + utilitas kecil yang benar‑benar diperlukan). README ini dirancang untuk sekaligus:

- **Menunjukkan kemampuan Anda** sebagai Backend Engineer Go (sesuai poin tanggung jawab &amp; kualifikasi).
- **Menjadi template nyata** untuk membangun layanan backend baru dengan pola yang rapi dan idiomatik.

---

## Tujuan Proyek

- **Fondasi backend Golang yang scalable**: mudah dikembangkan menjadi microservices atau modul terpisah.
- **API berperforma tinggi**: dukungan RESTful dan gRPC.
- **Desain data yang kuat**: skema SQL/NoSQL yang rapi, query optimal, dan integritas data terjaga.
- **Kualitas kode terukur**: testing menyeluruh (unit, integration, table‑driven).
- **Siap produksi**: sudah dipikirkan untuk Docker, Kubernetes, dan integrasi CI/CD.

---

## Fitur Utama

- **Go murni tanpa framework besar**
  - Routing, middleware, dan struktur aplikasi dibangun di atas `net/http`, `context`, dan package standard lain.
  - Memudahkan Anda menunjukkan pemahaman mendalam terhadap core Go, bukan sekadar framework.

- **RESTful API**
  - Desain resource‑oriented (endpoint yang konsisten, penggunaan HTTP method &amp; status code yang tepat).
  - Validasi request, error handling terstruktur, dan response JSON yang seragam (envelope / metadata).

- **gRPC API**
  - Definisi service &amp; message menggunakan `.proto`.
  - Mendukung komunikasi antar‑layanan berperforma tinggi (service‑to‑service).

- **Database SQL &amp; NoSQL**
  - Contoh integrasi ke database relasional (misalnya PostgreSQL/MySQL) dengan:
    - Connection pooling.
    - Transaction (begin/commit/rollback).
    - Indexing dasar dan query yang dioptimasi.
  - Opsi NoSQL (misalnya MongoDB) bila dibutuhkan untuk use case tertentu.

- **Redis sebagai cache &amp; penyimpanan sementara**
  - Caching query berat dan hasil komputasi.
  - Menyimpan session/token/short‑lived data.

- **Object Storage (S3/GCS Compatible)**
  - Upload, download, dan mengelola file (misalnya gambar/dokumen) ke S3/GCS atau layanan kompatibel (MinIO, dsb).

- **Docker &amp; Kubernetes Ready**
  - `Dockerfile` untuk build image yang kecil dan efisien.
  - Manifest/deployment Kubernetes dasar (Deployment, Service, ConfigMap/Secret).
  - Dirancang agar mudah dihubungkan ke pipeline CI/CD.

- **Git Workflow Friendly**
  - Struktur repositori yang cocok untuk branching, PR, dan code review.
  - Memudahkan diskusi dan review standar idiomatic Go.

---

## Struktur Proyek (Contoh)

Struktur direktori dapat disesuaikan, namun pola umumnya:

```text
go-backend-blueprint/
  cmd/
    api/              # entrypoint REST/gRPC server
  internal/
    http/             # handler, middleware, routing
    grpc/             # gRPC server &amp; interceptors
    config/           # konfigurasi (file/env)
    domain/           # business logic &amp; entity/domain model
    repository/       # akses ke DB (SQL/NoSQL), transaction, query
    cache/            # integrasi Redis
    storage/          # integrasi S3/GCS/MinIO
    tests/            # test utilities, integration setup
  migrations/         # file migrasi database
  proto/              # definisi .proto untuk gRPC
  deployments/
    docker/           # Dockerfile &amp; docker-compose (opsional)
    k8s/              # manifest Kubernetes
  .github/workflows/  # pipeline CI/CD (opsional)
  go.mod
  README.md
```

Struktur ini membantu memisahkan **domain**, **transport (HTTP/gRPC)**, dan **infrastruktur (DB, cache, storage)** agar mudah di‑test dan di‑refactor.

---

## Menjalankan Proyek

### 1. Prasyarat

- Go (minimal 1.21 atau yang Anda gunakan).
- Docker (untuk menjalankan DB, Redis, dan service pendukung).
- Make (opsional, tapi memudahkan eksekusi perintah).

### 2. Menyiapkan Lingkungan Lokal

Contoh dengan `docker-compose` (bila disertakan dalam proyek):

```bash
docker-compose up -d
```

Layanan tipikal:

- Database SQL (PostgreSQL/MySQL).
- Redis.
- MinIO (bila ingin menguji S3 compatible storage).

### 3. Menjalankan Server

```bash
go run ./cmd/api
```

Server akan expose:

- **REST API** di port misalnya `:8080`.
- **gRPC** di port misalnya `:9090`.

---

## Desain RESTful API

- **Konvensi Endpoint**
  - `GET /v1/resources`
  - `GET /v1/resources/{id}`
  - `POST /v1/resources`
  - `PUT /v1/resources/{id}`
  - `DELETE /v1/resources/{id}`

- **Prinsip Utama**
  - Menggunakan HTTP method sesuai semantik (GET, POST, PUT, PATCH, DELETE).
  - Mengembalikan status code yang tepat (200, 201, 400, 401, 404, 409, 422, 500, dll).
  - Response JSON konsisten dengan field `data`, `error`, dan/atau `meta`.

README ini dapat diperluas dengan dokumentasi OpenAPI/Swagger bila dibutuhkan.

---

## Desain gRPC Service

- Definisi service berada di direktori `proto/`.
- Menggunakan pola:
  - `service UserService { rpc GetUser (GetUserRequest) returns (GetUserResponse); }`
  - Mendukung unary RPC, dan dapat diperluas ke streaming bila dibutuhkan.
- Generator code gRPC dijalankan melalui `make proto` (atau perintah lain yang Anda definisikan).

---

## Database (SQL &amp; NoSQL)

- **SQL**
  - Menggunakan migrasi skema (di folder `migrations/`).
  - Mengoptimasi query melalui:
    - Penggunaan index dasar.
    - Penggunaan `JOIN` yang wajar dan jelas.
    - Transaction untuk operasi multi‑step (misalnya transfer saldo, update multi tabel).
  - Repository pattern agar mudah dipindah dari satu DB ke DB lain.

- **NoSQL (Opsional)**
  - Untuk data yang cocok disimpan sebagai dokumen atau key/value.
  - Dipisahkan pada package khusus agar tidak mencampur logika SQL &amp; NoSQL.

---

## Redis (Caching &amp; Penyimpanan Sementara)

- Menyimpan:
  - Session, token, OTP, dan data yang sifatnya sementara.
  - Cache hasil query/operasi berat agar response API lebih cepat.
- Mengatur TTL (time‑to‑live) dengan jelas sesuai kebutuhan business.
- Menggunakan kunci yang terstruktur, misalnya:
  - `user:profile:{id}`
  - `session:{token}`

---

## Object Storage (S3 / GCS / MinIO)

- Mendukung operasi:
  - Upload file.
  - Download file.
  - Menghapus file.
- Endpoint REST/gRPC dapat mengemas:
  - Presigned URL (jika diperlukan).
  - Metainformasi file (size, MIME type, dsb).

---

## Testing (Unit, Integration, Table‑Driven)

Proyek ini didesain untuk:

- **Unit Test**
  - Menggunakan `testing` package built‑in Go.
  - Fokus pada fungsi/komponen kecil (business logic, helper, dll).

- **Table‑Driven Test**
  - Menggunakan pola:
    - `tests := []struct{name string; input X; want Y}{ ... }`
    - `for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { ... }) }`
  - Memudahkan penambahan kasus baru tanpa mengulang banyak boilerplate.

- **Integration Test**
  - Menguji alur end‑to‑end (request ke API, DB, cache, dsb).
  - Dapat dijalankan dengan environment terpisah (menggunakan Docker).

Contoh perintah:

```bash
go test ./...
```

---

## Praktik Kode Idiomatic Go &amp; Code Review

Proyek ini mendorong:

- **Penulisan kode idiomatic**
  - Mengikuti guideline resmi Go (`Effective Go`, `Go Code Review Comments`).
  - Penamaan variabel, error handling, dan struktur package yang konsisten.

- **Code Review**
  - Setiap perubahan signifikan masuk lewat Pull Request.
  - Reviewer memeriksa:
    - Kesesuaian dengan standar kode.
    - Kualitas desain API &amp; domain.
    - Cakupan testing.

- **Linting &amp; Formatting**
  - Menggunakan `go fmt`, `go vet`, dan (opsional) linter tambahan seperti `golangci-lint`.

---

## Docker, Kubernetes &amp; CI/CD

- **Docker**
  - Image multi‑stage: build binary Go, lalu copy ke image runtime yang minimal.
  - Contoh perintah:

    ```bash
    docker build -t go-backend-blueprint .
    docker run --rm -p 8080:8080 go-backend-blueprint
    ```

- **Kubernetes**
  - Manifest contoh:
    - `Deployment` untuk aplikasi.
    - `Service` untuk expose di dalam cluster.
    - (Opsional) `Ingress` untuk expose ke publik.
    - `ConfigMap` &amp; `Secret` untuk konfigurasi.

- **CI/CD**
  - Pipeline minimal idealnya mencakup:
    - `go test ./...`
    - `go vet ./...`
    - Build Docker image &amp; push ke registry.
    - Deploy ke cluster (staging/production) via GitOps atau direct deploy.

---

## Kualifikasi yang Tercermin dalam Proyek Ini

README ini secara eksplisit mencerminkan poin kualifikasi Anda:

- **Pengalaman 2–3 tahun Backend Go**
  - Ditunjukkan melalui penggunaan idiomatik Go dan desain arsitektur yang terstruktur.
- **Penguasaan SQL (join, transaction, indexing)**
  - Didukung oleh layer repository dan migrasi database.
- **Penggunaan Redis**
  - Sebagai cache &amp; penyimpanan data sementara.
- **Penggunaan Object Storage (S3/GCS/MinIO)**
  - Untuk manajemen file yang scalable.
- **REST API yang kuat**
  - Struktur endpoint, error handling, dan konvensi yang jelas.
- **Git Workflow**
  - Struktur repo yang ramah branching, commit, merge, dan pull request.
- **Tanpa framework besar**
  - Ditekankan bahwa proyek ini dibangun di atas Go standard library, dengan utilitas tambahan seperlunya.

Dengan demikian, `go-backend-blueprint` dapat digunakan sebagai:

- **Blueprint teknis** ketika membangun layanan backend baru.
- **Portfolio profesional** untuk menunjukkan skill Anda sebagai Backend Engineer Golang.

