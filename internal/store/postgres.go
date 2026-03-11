package store

import (
	"database/sql"

	"go-backend-blueprint/internal/entity"
)

// Pastikan PostgresStore mengimplementasi ItemStore.
var _ ItemStore = (*PostgresStore)(nil)

// PostgresStore menyimpan Item ke tabel items di PostgreSQL.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore membuat store yang memakai database db.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) List() ([]*entity.Item, error) {
	rows, err := s.db.Query("SELECT id, name, created_at FROM items ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.Item
	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []*entity.Item{}
	}
	return items, nil
}

func (s *PostgresStore) GetByID(id string) (*entity.Item, error) {
	var item entity.Item
	err := s.db.QueryRow("SELECT id, name, created_at FROM items WHERE id = $1", id).
		Scan(&item.ID, &item.Name, &item.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *PostgresStore) Create(name string) (*entity.Item, error) {
	var item entity.Item
	err := s.db.QueryRow(
		"INSERT INTO items (name, created_at) VALUES ($1, now()) RETURNING id, name, created_at",
		name,
	).Scan(&item.ID, &item.Name, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *PostgresStore) Update(id string, name string) (*entity.Item, error) {
	var item entity.Item
	err := s.db.QueryRow(
		"UPDATE items SET name = $1 WHERE id = $2 RETURNING id, name, created_at",
		name, id,
	).Scan(&item.ID, &item.Name, &item.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *PostgresStore) Delete(id string) error {
	res, err := s.db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
