// Package command owns workspace mutations and their required port.
package command

import (
	"errors"
	"time"

	"github.com/0xsj/atelier-wails/internal/workspace/app/port"
	"github.com/0xsj/atelier-wails/internal/workspace/domain"
	"github.com/0xsj/atelier-wails/pkg/id"
)

type Service struct{ store port.WorkspaceStore }

func New(store port.WorkspaceStore) (*Service, error) {
	if store == nil {
		return nil, errors.New("workspace: nil store")
	}
	return &Service{store: store}, nil
}
func (s *Service) Register(workspaceID id.ID, name, location string, at time.Time) (domain.RegisterResult, error) {
	return s.store.Register(workspaceID, name, location, at)
}
func (s *Service) Rename(workspaceID id.ID, name string, expected domain.Expected, at time.Time) (domain.MutationResult, error) {
	return s.store.Rename(workspaceID, name, expected, at)
}
func (s *Service) Archive(workspaceID id.ID, expected domain.Expected, at time.Time) (domain.MutationResult, error) {
	return s.store.Archive(workspaceID, expected, at)
}
func (s *Service) Restore(workspaceID id.ID, expected domain.Expected, at time.Time) (domain.MutationResult, error) {
	return s.store.Restore(workspaceID, expected, at)
}
func (s *Service) Forget(workspaceID id.ID, expected domain.Expected) (domain.ForgetResult, error) {
	return s.store.Forget(workspaceID, expected)
}
