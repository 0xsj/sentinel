package logger

import (
	stderrors "errors"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	p "github.com/0xsj/atelier-wails/pkg/provenance"
)

// WithError snapshots safe metadata only. nil clears an inherited error binding.
func (l *Logger) WithError(err error) *Logger {
	n := *l
	n.failure = nil
	n.errorFields = nil
	if err == nil {
		return &n
	}
	pub, _ := faults.Public(err)
	k, known := faults.KindOf(err)
	n.failure = Fields{"classified": Bool(known), "kind": Text(k.String()), "message": Text(pub.Message), "has_cause": Bool(stderrors.Unwrap(err) != nil)}
	if pub.Type != "" {
		n.failure["type"] = Text(pub.Type)
	}
	n.errorFields = Fields{}
	for k, v := range pub.Fields {
		n.errorFields[k] = Text(v)
	}
	return &n
}

// WithScope records validated known claims, without sampling time or generating IDs.
func (l *Logger) WithScope(scope p.Scope) (*Logger, error) {
	s := scope.Snapshot()
	if _, e := p.RestoreScope(s); e != nil {
		return nil, e
	}
	w := s.Work
	f := Fields{"scope_id": Text(s.ScopeID.String()), "work_id": Text(w.WorkID.String()), "correlation_id": Text(w.CorrelationID.String()), "correlation_source": Text(string(w.CorrelationSource)), "operation": Text(w.Operation.String()), "started_at_ms": Int(s.StartedAt.UnixMilli()), "attempt": Int(int64(s.Attempt)), "executor_kind": Text(string(s.Executor.Kind())), "executor_id": Text(s.Executor.Identity())}
	if w.Depth != nil {
		f["depth"] = Int(int64(*w.Depth))
	}
	if w.Origin != nil {
		f["origin"] = Text(string(*w.Origin))
	}
	if w.Causation != nil {
		f["cause_kind"] = Text(string(w.Causation.Kind()))
		f["cause_id"] = Text(w.Causation.ID().String())
	}
	if s.PreviousAttempt != nil {
		f["previous_attempt"] = Text(s.PreviousAttempt.String())
	}
	a := w.Attribution.Snapshot()
	if a.Initiator != nil {
		f["initiator_kind"] = Text(string(a.Initiator.Kind()))
		f["initiator_id"] = Text(a.Initiator.Identity())
	}
	if a.OnBehalfOf != nil {
		f["on_behalf_of_kind"] = Text(string(a.OnBehalfOf.Kind()))
		f["on_behalf_of_id"] = Text(a.OnBehalfOf.Identity())
	}
	if a.Tenant != nil {
		f["tenant"] = Text(*a.Tenant)
	}
	n := *l
	n.scope = f
	return &n, nil
}
