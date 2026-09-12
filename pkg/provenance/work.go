package provenance

import "github.com/0xsj/atelier-wails/pkg/id"

type Origin string

const (
	Request  Origin = "request"
	Schedule Origin = "schedule"
	Backfill Origin = "backfill"
	Startup  Origin = "startup"
)

func (o Origin) valid() bool {
	return o == Request || o == Schedule || o == Backfill || o == Startup
}

type CorrelationSource string

const (
	Local    CorrelationSource = "local"
	External CorrelationSource = "external"
)

type WorkSnapshot struct {
	WorkID, CorrelationID id.ID
	CorrelationSource     CorrelationSource
	Causation             *Reference
	Origin                *Origin
	Operation             Operation
	Attribution           Attribution
	Depth                 *uint32
}
type WorkContext struct {
	data  WorkSnapshot
	valid bool
}

func copyWork(s WorkSnapshot) WorkSnapshot {
	s.Causation = clone(s.Causation)
	s.Origin = clone(s.Origin)
	s.Depth = clone(s.Depth)
	s.Attribution, _ = NewAttribution(s.Attribution.Snapshot())
	return s
}
func (w WorkContext) Snapshot() WorkSnapshot { return copyWork(w.data) }
func validateWork(s WorkSnapshot, defaults bool) error {
	if !defaults && (s.WorkID.IsZero() || s.CorrelationID.IsZero()) {
		return invalid("invalid_work")
	}
	if s.CorrelationSource != Local && s.CorrelationSource != External {
		return invalid("invalid_work")
	}
	if _, e := NewOperation(s.Operation.name); e != nil {
		return e
	}
	if _, e := NewAttribution(s.Attribution.Snapshot()); e != nil {
		return e
	}
	if s.Origin != nil && !s.Origin.valid() {
		return invalid("invalid_origin")
	}
	if s.Causation != nil {
		if !s.Causation.valid() {
			return invalid("invalid_reference")
		}
		if s.Causation.kind == WorkReference && s.Causation.id == s.WorkID {
			return invalid("invalid_work")
		}
	}
	if s.Depth != nil {
		if s.CorrelationSource == Local && *s.Depth == 0 && s.Causation != nil {
			return invalid("invalid_depth")
		}
		if *s.Depth > 0 && s.Causation == nil {
			return invalid("invalid_depth")
		}
	}
	return nil
}
func RestoreWork(s WorkSnapshot) (WorkContext, error) {
	if e := validateWork(s, false); e != nil {
		return WorkContext{}, e
	}
	return WorkContext{copyWork(s), true}, nil
}
