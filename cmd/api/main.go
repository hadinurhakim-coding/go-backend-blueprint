package main

import (
	"log"
	"net/http"

	"go-backend-blueprint/internal/config"
	"go-backend-blueprint/internal/database"
	"go-backend-blueprint/internal/handler"
	"go-backend-blueprint/internal/store"
)

func main() {
	cfg := config.FromEnv()

	// Opsional: sambungkan ke database jika DB_DSN diset. Koneksi siap untuk dipakai nanti (misalnya oleh store).
	if cfg.DBDSN != "" {
		db, err := database.Open(cfg.DBDSN)
		if err != nil {
			log.Fatalf("database open: %v", err)
		}
		defer database.Close(db)
		if err := database.Ping(db); err != nil {
			log.Fatalf("database ping: %v", err)
		}
		log.Println("database: connected")
		_ = db // nantinya db bisa diinjeksi ke store (misalnya PostgresStore)
	}

	// Store in-memory untuk CRUD Item. Nantinya bisa diganti dengan store yang pakai database.
	itemStore := store.NewMemoryStore()
	itemsHandler := &handler.ItemsHandler{Store: itemStore}

	// Endpoint /health: untuk health check
	http.HandleFunc("/health", handler.Health)
	// CRUD Item: GET/POST /items, GET/PUT/DELETE /items/:id
	http.HandleFunc("/items", itemsHandler.HandleItems)
	http.HandleFunc("/items/", itemsHandler.HandleItemByID)
	// Endpoint /: menyapa dengan teks "Halo"
	http.HandleFunc("/", handler.Halo)

	addr := ":" + cfg.Port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
