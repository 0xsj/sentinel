package errors

import (
	"context"
	"maps"
)

// inspect selects one frame, stopping before independent aggregate branches.
func inspect(err error) (*Failure, Kind, bool) {
	for depth := 0; err != nil && depth < 128; depth++ {
		switch e := err.(type) {
		case *Failure:
			if e == nil {
				return nil, Internal, false
			}
			return e, e.kind, e.kind <= Conflict
		case interface{ Unwrap() []error }:
			return nil, Internal, false
		}
		switch err {
		case context.Canceled:
			return nil, Canceled, true
		case context.DeadlineExceeded:
			return nil, Timeout, true
		}
		wrapped, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = wrapped.Unwrap()
	}
	return nil, Internal, false
}

// KindOf separates recognized classification from presentation fallback.
func KindOf(err error) (Kind, bool) { _, kind, ok := inspect(err); return kind, ok }

// DiagnosticTypeOf reads the selected frame's private type, including Internal.
// Use TypeOf/Public at a response boundary; this accessor does not redact it.
func DiagnosticTypeOf(err error) string {
	e, _, _ := inspect(err)
	if e == nil {
		return ""
	}
	return e.typ
}

// IsKind is false for success and unclassified failures, including for Internal.
func IsKind(err error, kind Kind) bool { actual, ok := KindOf(err); return ok && actual == kind }

// DetailsOf copies only the selected private frame, including for Internal.
// It does not flatten inner causes or aggregate branches.
func DetailsOf(err error) map[string]string {
	e, _, _ := inspect(err)
	if e == nil {
		return nil
	}
	return maps.Clone(e.details)
}
