// Package service orchestrates parsing, profiling and storage.
package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"prismeda/internal/eda"
	"prismeda/internal/store"
)

var ErrNotFound = errors.New("dataset not found")

type Service struct {
	store      *store.Store
	maxBytes   int64
	samplePath string
}

func New(maxBytes int64, samplePath string) *Service {
	return &Service{store: store.New(), maxBytes: maxBytes, samplePath: samplePath}
}

func (s *Service) MaxBytes() int64 { return s.maxBytes }

func (s *Service) Ingest(name string, r io.Reader) (*eda.Profile, error) {
	table, err := eda.Parse(r, name, s.maxBytes)
	if err != nil {
		return nil, err
	}
	id := store.NewID()
	profile := eda.BuildProfile(id, table)
	s.store.Put(id, table, profile)
	return profile, nil
}

func (s *Service) LoadSample() (*eda.Profile, error) {
	f, err := os.Open(s.samplePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return s.Ingest(filepath.Base(s.samplePath), f)
}

func (s *Service) entry(id string) (*store.Entry, error) {
	e, ok := s.store.Get(id)
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *Service) Profile(id string) (*eda.Profile, error) {
	e, err := s.entry(id)
	if err != nil {
		return nil, err
	}
	return e.Profile, nil
}

func (s *Service) Preview(id string, limit int) (*eda.Preview, error) {
	e, err := s.entry(id)
	if err != nil {
		return nil, err
	}
	return e.Table.Preview(limit), nil
}

func (s *Service) Scatter(id, x, y string, limit int) ([]eda.Point, error) {
	e, err := s.entry(id)
	if err != nil {
		return nil, err
	}
	return e.Table.Scatter(x, y, limit)
}

func (s *Service) Count() int { return s.store.Count() }
