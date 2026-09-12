// Package memory provides an isolated, synchronized preference store.
package memory

import (
	"sort"
	"sync"

	"github.com/0xsj/atelier-wails/internal/preferences/domain"
)

type scopeKey struct{ scope string }

type Store struct {
	mu    sync.RWMutex
	items map[scopeKey]map[domain.Key]domain.Entry
}

func New() *Store { return &Store{items: make(map[scopeKey]map[domain.Key]domain.Entry)} }
func (s *Store) Read(scope domain.Scope, key domain.Key) (domain.Entry, bool, error) {
	if !scope.Valid() || !key.Valid() {
		return domain.Entry{}, false, invalidStoreInput()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.items[scopeKey{scope.String()}][key]
	return entry, ok, nil
}
func (s *Store) List(scope domain.Scope) ([]domain.Entry, error) {
	if !scope.Valid() {
		return nil, invalidStoreInput()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := make([]domain.Entry, 0, len(s.items[scopeKey{scope.String()}]))
	for _, entry := range s.items[scopeKey{scope.String()}] {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key.String() < entries[j].Key.String() })
	return entries, nil
}
func (s *Store) Replace(scope domain.Scope, key domain.Key, value domain.Value, expected domain.Expected) (domain.ReplaceResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !scope.Valid() || !key.Valid() || !value.Valid() || !expected.Valid() {
		return domain.ReplaceResult{}, invalidStoreInput()
	}
	k := scopeKey{scope.String()}
	bucket := s.items[k]
	var current *domain.Entry
	if bucket != nil {
		if entry, ok := bucket[key]; ok {
			current = &entry
		}
	}
	result, err := domain.DecideReplace(current, scope, key, value, expected)
	if err != nil || result.Status == domain.Unchanged {
		return result, err
	}
	if bucket == nil {
		bucket = make(map[domain.Key]domain.Entry)
		s.items[k] = bucket
	}
	bucket[key] = result.Entry
	return result, nil
}
func (s *Store) Remove(scope domain.Scope, key domain.Key, expected domain.Expected) (domain.RemoveResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !scope.Valid() || !key.Valid() || !expected.Valid() {
		return domain.RemoveResult{}, invalidStoreInput()
	}
	k := scopeKey{scope.String()}
	bucket := s.items[k]
	var current *domain.Entry
	if bucket != nil {
		if entry, ok := bucket[key]; ok {
			current = &entry
		}
	}
	result, err := domain.DecideRemove(current, scope, key, expected)
	if err != nil || !result.Removed {
		return result, err
	}
	delete(bucket, key)
	if len(bucket) == 0 {
		delete(s.items, k)
	}
	return result, nil
}
func invalidStoreInput() error     { return domainError("preferences.invalid_input") }
func domainError(typ string) error { return &storeFailure{typ: typ} }

type storeFailure struct{ typ string }

func (e *storeFailure) Error() string { return "preferences: invalid input" }
