// Package query owns preference reads and deterministic projections.
package query

import (
	"errors"

	"github.com/0xsj/atelier-wails/internal/preferences/app/port"
	"github.com/0xsj/atelier-wails/internal/preferences/domain"
)

type Service struct{ store port.PreferenceStore }

func New(store port.PreferenceStore) (*Service, error) {
	if store == nil {
		return nil, errors.New("preferences: nil store")
	}
	return &Service{store: store}, nil
}
func (s *Service) Read(scope domain.Scope, key domain.Key) (domain.Entry, bool, error) {
	if !scope.Valid() || !key.Valid() {
		return domain.Entry{}, false, errors.New("preferences: invalid query")
	}
	return s.store.Read(scope, key)
}
func (s *Service) List(scope domain.Scope) ([]domain.Entry, error) {
	if !scope.Valid() {
		return nil, errors.New("preferences: invalid query")
	}
	return s.store.List(scope)
}
