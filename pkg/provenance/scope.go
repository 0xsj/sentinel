package provenance

import (
	"github.com/0xsj/atelier-wails/pkg/id"
	"time"
)

const maxMillis int64 = 281474976710655

type ScopeSnapshot struct {
	Work            WorkSnapshot
	ScopeID         id.ID
	StartedAt       time.Time
	Executor        Actor
	Attempt         uint32
	PreviousAttempt *id.ID
}
type Scope struct {
	data  ScopeSnapshot
	valid bool
}

func normalizedTime(t time.Time) (time.Time, error) {
	if t.Before(time.UnixMilli(0)) || !t.Before(time.UnixMilli(maxMillis).Add(time.Millisecond)) {
		return time.Time{}, invalid("invalid_time")
	}
	return time.UnixMilli(t.UnixMilli()).UTC(), nil
}
func (s Scope) Snapshot() ScopeSnapshot {
	v := s.data
	v.Work = copyWork(v.Work)
	v.PreviousAttempt = clone(v.PreviousAttempt)
	return v
}
func (s Scope) WorkContext() WorkContext { return WorkContext{copyWork(s.data.Work), s.valid} }
func RestoreScope(s ScopeSnapshot) (Scope, error) {
	if s.ScopeID.IsZero() {
		return Scope{}, invalid("invalid_scope")
	}
	if e := validateWork(s.Work, false); e != nil {
		return Scope{}, e
	}
	if !s.Executor.named() {
		return Scope{}, invalid("invalid_actor")
	}
	if s.Attempt == 0 {
		return Scope{}, invalid("invalid_attempt")
	}
	if s.PreviousAttempt != nil && (s.PreviousAttempt.IsZero() || *s.PreviousAttempt == s.ScopeID || s.Attempt == 1) {
		return Scope{}, invalid("invalid_scope")
	}
	if s.Work.Causation != nil && s.Work.Causation.kind == ScopeReference && s.Work.Causation.id == s.ScopeID {
		return Scope{}, invalid("invalid_scope")
	}
	t, e := normalizedTime(s.StartedAt)
	if e != nil {
		return Scope{}, e
	}
	s.StartedAt = t
	s.Work = copyWork(s.Work)
	s.PreviousAttempt = clone(s.PreviousAttempt)
	return Scope{s, true}, nil
}
