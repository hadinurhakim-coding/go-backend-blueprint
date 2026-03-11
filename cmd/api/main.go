package main

import (
	"log"
	"net/http"

	"go-backend-blueprint/internal/config"
	"go-backend-blueprint/internal/database"
	"go-backend-blueprint/internal/handler"
	"go-backend-blueprint/internal/migrate"
	"go-backend-blueprint/internal/store"
)

func main() {
	cfg := config.FromEnv()

	var itemStore store.ItemStore
	if cfg.DBDSN != "" {
		db, err := database.Open(cfg.DBDSN)
		if err != nil {
			log.Fatalf("database open: %v", err)
		}
		defer database.Close(db)
		if err := database.Ping(db); err != nil {
			log.Fatalf("database ping: %v", err)
		}
		if err := migrate.Run(db, "migrations"); err != nil {
			log.Fatalf("migrate: %v", err)
		}
		log.Println("database: connected, migrations ok")
		itemStore = store.NewPostgresStore(db)
	} else {
		itemStore = store.NewMemoryStore()
	}

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
