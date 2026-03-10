package store

import (
	"errors"

	"go-backend-blueprint/internal/entity"
)

// ErrNotFound dipakai ketika resource tidak ditemukan (GetByID nil, Update/Delete tidak ada).
var ErrNotFound = errors.New("not found")

// ItemStore adalah interface untuk penyimpanan Item. Dengan interface ini
// kita bisa pakai in-memory sekarang dan ganti ke database nanti tanpa mengubah handler.
type ItemStore interface {
	List() ([]*entity.Item, error)
	GetByID(id string) (*entity.Item, error)
	Create(name string) (*entity.Item, error)
	Update(id string, name string) (*entity.Item, error)
	Delete(id string) error
}
