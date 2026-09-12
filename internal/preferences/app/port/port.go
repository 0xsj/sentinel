// Package port contains application-owned preference storage capabilities.
package port

import "github.com/0xsj/atelier-wails/internal/preferences/domain"

type PreferenceStore interface {
	Read(domain.Scope, domain.Key) (domain.Entry, bool, error)
	List(domain.Scope) ([]domain.Entry, error)
	Replace(domain.Scope, domain.Key, domain.Value, domain.Expected) (domain.ReplaceResult, error)
	Remove(domain.Scope, domain.Key, domain.Expected) (domain.RemoveResult, error)
}
