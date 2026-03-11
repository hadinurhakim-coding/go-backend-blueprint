package main

import (
	// log dipakai untuk menulis informasi startup dan error fatal ke terminal.
	"log"
	// net/http adalah server HTTP bawaan Go (tanpa framework).
	"net/http"

	// config membaca konfigurasi dari environment variable (PORT, DB_DSN).
	"go-backend-blueprint/internal/config"
	// database berisi helper membuka koneksi PostgreSQL (Open/Ping/Close).
	"go-backend-blueprint/internal/database"
	// handler berisi fungsi/struct yang menangani request HTTP (endpoint).
	"go-backend-blueprint/internal/handler"
	// migrate menjalankan file SQL di folder migrations agar tabel siap dipakai.
	"go-backend-blueprint/internal/migrate"
	// store adalah layer akses data. Di sini kita bisa pilih penyimpanan in-memory atau PostgreSQL.
	"go-backend-blueprint/internal/store"
)

func main() {
	// Ambil konfigurasi dari env. Tujuannya: perilaku aplikasi bisa diubah tanpa mengubah kode.
	cfg := config.FromEnv()

	// itemStore adalah “backend penyimpanan data” untuk CRUD Item.
	// Kita pakai interface agar handler tidak perlu tahu detail: memory atau database.
	var itemStore store.ItemStore
	if cfg.DBDSN != "" {
		// Jika DB_DSN ada, artinya kita ingin pakai PostgreSQL (data tersimpan permanen).
		db, err := database.Open(cfg.DBDSN)
		if err != nil {
			// Fatal = stop aplikasi. Alasannya: tanpa DB yang valid, mode database tidak bisa jalan benar.
			log.Fatalf("database open: %v", err)
		}
		// Pastikan koneksi database ditutup saat aplikasi berhenti.
		defer database.Close(db)
		if err := database.Ping(db); err != nil {
			// Ping memastikan DB benar-benar bisa diakses (host/port/password/dbname benar).
			log.Fatalf("database ping: %v", err)
		}
		if err := migrate.Run(db, "migrations"); err != nil {
			// Migrasi membuat/menyesuaikan tabel (mis. items) sebelum aplikasi mulai melayani request.
			log.Fatalf("migrate: %v", err)
		}
		log.Println("database: connected, migrations ok")
		// Gunakan PostgresStore: implementasi ItemStore yang menjalankan query SQL ke tabel items.
		itemStore = store.NewPostgresStore(db)
	} else {
		// Jika DB_DSN kosong, pakai MemoryStore untuk mode belajar/dev cepat (data hilang saat server stop).
		itemStore = store.NewMemoryStore()
	}

	// ItemsHandler butuh ItemStore untuk melakukan operasi CRUD.
	itemsHandler := &handler.ItemsHandler{Store: itemStore}

	// Endpoint /health: dipakai untuk health check (monitoring / load balancer / Kubernetes).
	http.HandleFunc("/health", handler.Health)
	// CRUD Item:
	// - GET/POST /items
	// - GET/PUT/DELETE /items/:id
	http.HandleFunc("/items", itemsHandler.HandleItems)
	http.HandleFunc("/items/", itemsHandler.HandleItemByID)
	// Endpoint /: endpoint paling sederhana untuk memastikan server merespons.
	http.HandleFunc("/", handler.Halo)

	// Tentukan alamat listen dari PORT (default 8080). Format net/http: ":8080".
	addr := ":" + cfg.Port
	log.Printf("server listening on %s", addr)
	// ListenAndServe akan “blocking” (program berhenti di sini) sampai server dimatikan (Ctrl+C) atau error.
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("listen: %v", err) 
	}
}
