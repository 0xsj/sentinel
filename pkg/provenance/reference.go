package provenance

import "github.com/0xsj/atelier-wails/pkg/id"

type ReferenceKind string

const (
	ScopeReference ReferenceKind = "scope"
	WorkReference  ReferenceKind = "work"
	EventReference ReferenceKind = "event"
)

type Reference struct {
	kind ReferenceKind
	id   id.ID
}

func validReferenceKind(k ReferenceKind) bool {
	return k == ScopeReference || k == WorkReference || k == EventReference
}
func NewReference(kind ReferenceKind, value id.ID) (Reference, error) {
	r := Reference{kind, value}
	if !r.valid() {
		return Reference{}, invalid("invalid_reference")
	}
	return r, nil
}
func (r Reference) Kind() ReferenceKind { return r.kind }
func (r Reference) ID() id.ID           { return r.id }
func (r Reference) valid() bool         { return validReferenceKind(r.kind) && !r.id.IsZero() }
