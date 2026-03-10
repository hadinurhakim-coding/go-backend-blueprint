package config

import "os"

// Config menyimpan konfigurasi aplikasi (dari environment variable).
// Nantinya bisa ditambah field lain (Redis URL, S3 bucket, dll.).
type Config struct {
	// Port untuk HTTP server (misalnya "8080").
	Port string
	// DBDSN adalah connection string database (PostgreSQL). Kosong = tidak connect ke DB.
	// Contoh: "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
	DBDSN string
}

// FromEnv membaca konfigurasi dari environment variable.
// PORT → Port (default "8080"), DB_DSN → DBDSN.
func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		Port:  port,
		DBDSN: os.Getenv("DB_DSN"),
	}
}
