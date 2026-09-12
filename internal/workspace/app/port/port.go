// Package port contains application-owned workspace capabilities.
package port

import (
	"time"

	"github.com/0xsj/atelier-wails/internal/workspace/domain"
	"github.com/0xsj/atelier-wails/pkg/id"
)

type WorkspaceStore interface {
	Register(id.ID, string, string, time.Time) (domain.RegisterResult, error)
	Read(id.ID) (domain.Workspace, bool, error)
	List(domain.ListFilter) ([]domain.Workspace, error)
	Rename(id.ID, string, domain.Expected, time.Time) (domain.MutationResult, error)
	Archive(id.ID, domain.Expected, time.Time) (domain.MutationResult, error)
	Restore(id.ID, domain.Expected, time.Time) (domain.MutationResult, error)
	Forget(id.ID, domain.Expected) (domain.ForgetResult, error)
}
