package store

import (
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.RWMutex
	urls    map[string]string
	clicks  map[string]int
	expiry  map[string]time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		urls:   make(map[string]string),
		clicks: make(map[string]int),
		expiry: make(map[string]time.Time),
	}
}

func (s *MemoryStore) Close() error { return nil }

func (s *MemoryStore) Set(alias, originalURL string, expiresIn int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[alias] = originalURL
	s.clicks[alias] = 0

	if expiresIn > 0 {
		s.expiry[alias] = time.Now().Add(time.Duration(expiresIn) * time.Second)
	} else {
		delete(s.expiry, alias)
	}
	return nil
}

func (s *MemoryStore) Get(alias string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if exp, ok := s.expiry[alias]; ok && time.Now().After(exp) {
		return "", ErrNotFound
	}

	val, ok := s.urls[alias]
	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

func (s *MemoryStore) Exists(alias string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if exp, ok := s.expiry[alias]; ok && time.Now().After(exp) {
		return false
	}
	_, ok := s.urls[alias]
	return ok
}

func (s *MemoryStore) IncrementClick(alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clicks[alias]++
}

var ErrNotFound = &notFoundError{}

type notFoundError struct{}

func (e *notFoundError) Error() string { return "not found" }