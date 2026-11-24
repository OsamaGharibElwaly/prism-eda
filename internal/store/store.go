// Package store keeps analysed datasets in memory for the lifetime of the server.
package store

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"prismeda/internal/eda"
)

type Entry struct {
	Table   *eda.Table
	Profile *eda.Profile
}

type Store struct {
	mu    sync.RWMutex
	items map[string]*Entry
}

func New() *Store {
	return &Store{items: make(map[string]*Entry)}
}

func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Store) Put(id string, table *eda.Table, profile *eda.Profile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = &Entry{Table: table, Profile: profile}
}

func (s *Store) Get(id string) (*Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.items[id]
	return e, ok
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
