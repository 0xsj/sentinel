// Package command owns preference mutations and their required port.
package command

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
func (s *Service) Replace(scope domain.Scope, key domain.Key, value domain.Value, expected domain.Expected) (domain.ReplaceResult, error) {
	if !scope.Valid() || !key.Valid() || !value.Valid() || !expected.Valid() {
		return domain.ReplaceResult{}, errors.New("preferences: invalid command")
	}
	return s.store.Replace(scope, key, value, expected)
}
func (s *Service) Remove(scope domain.Scope, key domain.Key, expected domain.Expected) (domain.RemoveResult, error) {
	if !scope.Valid() || !key.Valid() || !expected.Valid() {
		return domain.RemoveResult{}, errors.New("preferences: invalid command")
	}
	return s.store.Remove(scope, key, expected)
}
