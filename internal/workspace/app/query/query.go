// Package query owns workspace reads and deterministic projections.
package query

import (
	"errors"

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
func (s *Service) Read(workspaceID id.ID) (domain.Workspace, bool, error) {
	return s.store.Read(workspaceID)
}
func (s *Service) List(filter domain.ListFilter) ([]domain.Workspace, error) {
	return s.store.List(filter)
}
