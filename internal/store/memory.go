package store

import (
	"fmt"
	"sync"
	"sync/atomic"

	"go-backend-blueprint/internal/entity"
)

// Pastikan MemoryStore mengimplementasi ItemStore.
var _ ItemStore = (*MemoryStore)(nil)

// MemoryStore menyimpan Item di memori (map). Cocok untuk development dan testing.
// Nantinya bisa ditambah ItemStore lain (misalnya PostgresStore) yang implement interface yang sama.
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]*entity.Item
	next  atomic.Uint64
}

// NewMemoryStore membuat in-memory store baru.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]*entity.Item),
	}
}

func (s *MemoryStore) List() ([]*entity.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*entity.Item, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out, nil
}

func (s *MemoryStore) GetByID(id string) (*entity.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return nil, nil
	}
	return item, nil
}

func (s *MemoryStore) Create(name string) (*entity.Item, error) {
	id := fmt.Sprintf("%d", s.next.Add(1))
	item := &entity.Item{
		ID:        id,
		Name:      name,
		CreatedAt: entity.NowFunc(),
	}
	s.mu.Lock()
	s.items[id] = item
	s.mu.Unlock()
	return item, nil
}

func (s *MemoryStore) Update(id string, name string) (*entity.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return nil, nil
	}
	item.Name = name
	return item, nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
