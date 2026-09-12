package errors

import "maps"

// Failure carries immutable meaning and a private diagnostic frame.
// Use public accessors for responses; Error() is diagnostic text.
type Failure struct {
	kind      Kind
	message   string
	typ       string
	fields    map[string]string
	details   map[string]string
	cause     error
	condition *Failure
}

// New constructs a real failure, including when its message is empty.
func New(kind Kind, message string) *Failure {
	if kind > Conflict {
		kind = Internal
	}
	return &Failure{kind: kind, message: message}
}

// Wrap deliberately classifies a cause. Nil remains a nil error interface.
func Wrap(err error, kind Kind, message string) error {
	if err == nil {
		return nil
	}
	return New(kind, message).WithCause(err)
}

// Error is diagnostic; use Public at caller-facing boundaries.
func (e *Failure) Error() string      { return e.message }
func (e *Failure) Kind() Kind         { return e.kind }
func (e *Failure) Diagnostic() string { return e.message }

// Unwrap exposes the diagnostic cause to standard errors.Is and errors.As.
func (e *Failure) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}
func (e *Failure) identity() *Failure {
	if e.condition != nil {
		return e.condition
	}
	return e
}

// Is preserves condition identity across derivations; equal text is not identity.
func (e *Failure) Is(target error) bool {
	other, ok := target.(*Failure)
	return ok && e != nil && other != nil && e.identity() == other.identity()
}
func (e *Failure) clone() *Failure {
	copy := *e
	copy.condition = e.identity()
	copy.fields = maps.Clone(e.fields)
	copy.details = maps.Clone(e.details)
	return &copy
}

// WithType derives a public condition identifier.
func (e *Failure) WithType(value string) *Failure { copy := e.clone(); copy.typ = value; return copy }

// WithField derives a public problem. New values replace the same key.
func (e *Failure) WithField(key, value string) *Failure {
	return e.WithFields(map[string]string{key: value})
}

// WithFields copies and merges public problems without aliasing the input.
func (e *Failure) WithFields(values map[string]string) *Failure {
	copy := e.clone()
	if len(values) != 0 {
		if copy.fields == nil {
			copy.fields = make(map[string]string, len(values))
		}
		maps.Copy(copy.fields, values)
	}
	return copy
}

// WithDetail derives a private diagnostic value.
func (e *Failure) WithDetail(key, value string) *Failure {
	return e.WithDetails(map[string]string{key: value})
}

// WithDetails copies and merges private metadata; new values win.
func (e *Failure) WithDetails(values map[string]string) *Failure {
	copy := e.clone()
	if len(values) != 0 {
		if copy.details == nil {
			copy.details = make(map[string]string, len(values))
		}
		maps.Copy(copy.details, values)
	}
	return copy
}

// WithCause constructs an occurrence with a replacement source. A nil cause still
// produces a real failure; use Wrap for nil propagation.
func (e *Failure) WithCause(cause error) *Failure { copy := e.clone(); copy.cause = cause; return copy }
