package provenance

import (
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
	"math"
	"time"
)

type Clock interface{ Now() time.Time }
type Generator interface{ NewID() (id.ID, error) }
type Factory struct {
	clock Clock
	ids   Generator
}
type RootSpec struct {
	WorkID      *id.ID
	Origin      Origin
	Operation   Operation
	Attribution Attribution
	Executor    Actor
}
type StepSpec struct {
	WorkID    *id.ID
	Operation Operation
	Executor  Actor
	Cause     *Reference
}
type WorkSpec struct {
	WorkID    id.ID
	Operation Operation
	Cause     *Reference
}
type ExecutionSpec struct {
	Executor Actor
	Attempt  uint32
}

func NewFactory(c Clock, g Generator) (*Factory, error) {
	if c == nil || g == nil {
		return nil, invalid("invalid_configuration")
	}
	return &Factory{c, g}, nil
}
func generatedError() error {
	return faults.New(faults.Internal, "invalid generated provenance ID").WithType("provenance.invalid_generated_id")
}
func (f *Factory) emit(w WorkSnapshot, a Actor, attempt uint32, previous *id.ID, forbidden ...id.ID) (Scope, error) {
	if !a.named() {
		return Scope{}, invalid("invalid_actor")
	}
	if attempt == 0 {
		return Scope{}, invalid("invalid_attempt")
	}
	if e := validateWork(w, true); e != nil {
		return Scope{}, e
	}
	started, e := normalizedTime(f.clock.Now())
	if e != nil {
		return Scope{}, e
	}
	next, e := f.ids.NewID()
	if e != nil {
		return Scope{}, e
	}
	if next.IsZero() {
		return Scope{}, generatedError()
	}
	for _, v := range forbidden {
		if next == v {
			return Scope{}, generatedError()
		}
	}
	if w.WorkID.IsZero() {
		w.WorkID = next
	}
	if w.CorrelationID.IsZero() {
		w.CorrelationID = next
	}
	scope, e := RestoreScope(ScopeSnapshot{w, next, started, a, attempt, previous})
	if e != nil {
		return Scope{}, generatedError()
	}
	return scope, nil
}
func rootWork(s RootSpec) (WorkSnapshot, error) {
	if !s.Origin.valid() {
		return WorkSnapshot{}, invalid("invalid_origin")
	}
	w := WorkSnapshot{CorrelationSource: Local, Origin: clone(&s.Origin), Operation: s.Operation, Attribution: s.Attribution, Depth: new(uint32)}
	if s.WorkID != nil {
		if s.WorkID.IsZero() {
			return WorkSnapshot{}, invalid("invalid_work")
		}
		w.WorkID = *s.WorkID
	}
	return w, nil
}
func (f *Factory) Open(s RootSpec) (Scope, error) {
	w, e := rootWork(s)
	if e != nil {
		return Scope{}, e
	}
	return f.emit(w, s.Executor, 1, nil)
}
func childWork(parent Scope, workID *id.ID, operation Operation, cause *Reference) (WorkSnapshot, error) {
	if !parent.valid {
		return WorkSnapshot{}, invalid("invalid_scope")
	}
	w := parent.WorkContext().Snapshot()
	w.WorkID = id.ID{}
	if workID != nil {
		if workID.IsZero() || *workID == parent.data.Work.WorkID {
			return WorkSnapshot{}, invalid("invalid_work")
		}
		w.WorkID = *workID
	}
	w.Operation = operation
	if w.Depth != nil {
		if *w.Depth == math.MaxUint32 {
			return WorkSnapshot{}, invalid("depth_exhausted")
		}
		*w.Depth++
	}
	if cause == nil {
		r, _ := NewReference(ScopeReference, parent.data.ScopeID)
		w.Causation = &r
	} else {
		w.Causation = clone(cause)
	}
	if e := validateWork(w, true); e != nil {
		return WorkSnapshot{}, e
	}
	return w, nil
}
func (f *Factory) Child(parent Scope, s StepSpec) (Scope, error) {
	w, e := childWork(parent, s.WorkID, s.Operation, s.Cause)
	if e != nil {
		return Scope{}, e
	}
	forbidden := []id.ID{parent.data.ScopeID}
	if s.WorkID == nil {
		forbidden = append(forbidden, parent.data.Work.WorkID)
	}
	return f.emit(w, s.Executor, 1, nil, forbidden...)
}
func Prepare(parent Scope, s WorkSpec) (WorkContext, error) {
	w, e := childWork(parent, &s.WorkID, s.Operation, s.Cause)
	if e != nil {
		return WorkContext{}, e
	}
	return RestoreWork(w)
}
func (f *Factory) Execute(work WorkContext, s ExecutionSpec) (Scope, error) {
	if !work.valid {
		return Scope{}, invalid("invalid_work")
	}
	return f.emit(work.Snapshot(), s.Executor, s.Attempt, nil)
}
func (f *Factory) Retry(previous Scope, executor Actor) (Scope, error) {
	if !previous.valid {
		return Scope{}, invalid("invalid_scope")
	}
	if previous.data.Attempt == math.MaxUint32 {
		return Scope{}, invalid("attempt_exhausted")
	}
	return f.emit(previous.WorkContext().Snapshot(), executor, previous.data.Attempt+1, clone(&previous.data.ScopeID), previous.data.ScopeID)
}
