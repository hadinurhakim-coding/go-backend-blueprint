package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// schemaMigrationsTable dipakai untuk mencatat migrasi yang sudah dijalankan.
const schemaMigrationsTable = `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`

// Run menjalankan semua file .sql di dir yang belum tercatat di schema_migrations.
// File diurutkan berdasarkan nama (001_..., 002_..., dst). Version = bagian sebelum underscore pertama.
func Run(db *sql.DB, dir string) error {
	if _, err := db.Exec(schemaMigrationsTable); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		version := strings.TrimSuffix(name, ".sql")
		if idx := strings.Index(version, "_"); idx > 0 {
			version = version[:idx]
		}

		var applied string
		err := db.QueryRow("SELECT version FROM schema_migrations WHERE version = $1", version).Scan(&applied)
		if err == nil {
			continue // sudah dijalankan
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("check migration %s: %w", version, err)
		}

		path := filepath.Join(dir, name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		sqlText := strings.TrimSpace(string(sqlBytes))
		if sqlText == "" {
			continue
		}

		if _, err := db.Exec(sqlText); err != nil {
			return fmt.Errorf("run %s: %w", name, err)
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			return fmt.Errorf("record migration %s: %w", version, err)
		}
	}
	return nil
}
