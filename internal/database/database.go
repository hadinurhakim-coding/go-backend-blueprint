package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Open membuka koneksi ke PostgreSQL menggunakan DSN (Data Source Name).
// Contoh DSN: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	return db, nil
}

// Ping memeriksa bahwa koneksi ke database masih hidup. Dipanggil setelah Open untuk memastikan DB siap.
func Ping(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}

// Close menutup koneksi database. Biasanya dipanggil dengan defer setelah Open di main.
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
