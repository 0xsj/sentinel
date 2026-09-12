// Package memory provides an isolated, synchronized workspace registry.
package memory

import (
	"sort"
	"sync"
	"time"

	"github.com/0xsj/atelier-wails/internal/workspace/domain"
	"github.com/0xsj/atelier-wails/pkg/id"
)

type Store struct {
	mu    sync.RWMutex
	items map[id.ID]domain.Workspace
}

func New() *Store { return &Store{items: make(map[id.ID]domain.Workspace)} }
func (s *Store) Register(workspaceID id.ID, name, location string, at time.Time) (domain.RegisterResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[workspaceID]; exists {
		return domain.RegisterResult{}, conflict()
	}
	for _, item := range s.items {
		if item.Status == domain.Active && item.Location == location {
			return domain.RegisterResult{}, conflict()
		}
	}
	result, err := domain.Register(workspaceID, name, location, at)
	if err != nil {
		return domain.RegisterResult{}, err
	}
	s.items[workspaceID] = result.Workspace
	return result, nil
}
func (s *Store) Read(workspaceID id.ID) (domain.Workspace, bool, error) {
	if workspaceID.IsZero() {
		return domain.Workspace{}, false, invalid()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.items[workspaceID]
	return value, ok, nil
}
func (s *Store) List(filter domain.ListFilter) ([]domain.Workspace, error) {
	if !filter.Valid() {
		return nil, invalid()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]domain.Workspace, 0, len(s.items))
	for _, item := range s.items {
		if filter == domain.ActiveOnly && item.Status != domain.Active {
			continue
		}
		values = append(values, item)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].ID.String() < values[j].ID.String() })
	return values, nil
}
func (s *Store) Rename(workspaceID id.ID, name string, expected domain.Expected, at time.Time) (domain.MutationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[workspaceID]
	if !ok {
		return domain.MutationResult{}, nil
	}
	result, err := domain.Rename(current, name, expected, at)
	if err != nil || result.Status == domain.Unchanged {
		return result, err
	}
	s.items[workspaceID] = result.Workspace
	return result, nil
}
func (s *Store) Archive(workspaceID id.ID, expected domain.Expected, at time.Time) (domain.MutationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[workspaceID]
	if !ok {
		return domain.MutationResult{}, nil
	}
	result, err := domain.Archive(current, expected, at)
	if err != nil || result.Status == domain.Unchanged {
		return result, err
	}
	s.items[workspaceID] = result.Workspace
	return result, nil
}
func (s *Store) Restore(workspaceID id.ID, expected domain.Expected, at time.Time) (domain.MutationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[workspaceID]
	if !ok {
		return domain.MutationResult{}, nil
	}
	if current.Status == domain.Archived {
		for otherID, item := range s.items {
			if otherID != workspaceID && item.Status == domain.Active && item.Location == current.Location {
				return domain.MutationResult{}, conflict()
			}
		}
	}
	result, err := domain.Restore(current, expected, at)
	if err != nil || result.Status == domain.Unchanged {
		return result, err
	}
	s.items[workspaceID] = result.Workspace
	return result, nil
}
func (s *Store) Forget(workspaceID id.ID, expected domain.Expected) (domain.ForgetResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[workspaceID]
	if !ok {
		return domain.ForgetResult{}, nil
	}
	result, err := domain.Forget(current, expected)
	if err != nil || !result.Removed {
		return result, err
	}
	delete(s.items, workspaceID)
	return result, nil
}
func invalid() error  { return &storeFailure{"workspace: invalid input"} }
func conflict() error { return &storeFailure{"workspace: conflict"} }

type storeFailure struct{ message string }

func (e *storeFailure) Error() string { return e.message }
