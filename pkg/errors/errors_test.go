package errors_test

import (
	stderrors "errors"
	"fmt"
	"strings"
	"testing"

	failure "github.com/0xsj/atelier-wails/pkg/errors"
)

var vocabulary = map[string]failure.Kind{
	"internal": failure.Internal, "unauthenticated": failure.Unauthenticated,
	"forbidden": failure.Forbidden, "rate_limited": failure.RateLimited,
	"unavailable": failure.Unavailable, "timeout": failure.Timeout,
	"canceled": failure.Canceled, "not_found": failure.NotFound,
	"invalid": failure.Invalid, "conflict": failure.Conflict,
}

func TestCanonicalKinds(t *testing.T) {
	kinds := failure.Kinds()
	if len(kinds) != len(vocabulary) {
		t.Fatalf("got %d kinds", len(kinds))
	}
	seen := map[string]bool{}
	for _, kind := range kinds {
		name := kind.String()
		want, exists := vocabulary[name]
		if !exists || seen[name] || want != kind {
			t.Fatalf("unexpected/duplicate kind: %q", name)
		}
		seen[name] = true
		if got, ok := failure.ParseKind(name); !ok || got != want {
			t.Fatalf("parse %q = %v, %v", name, got, ok)
		}
	}
}

func TestRejectNoncanonicalNames(t *testing.T) {
	inputs := []string{"", "unknown", "deadline", "time\x00out", "ｔimeout"}
	for name := range vocabulary {
		inputs = append(inputs, strings.ToUpper(name), " "+name, name+"\n")
	}
	for _, input := range inputs {
		if got, ok := failure.ParseKind(input); ok || got != failure.Internal {
			t.Errorf("accepted %q: %v %v", input, got, ok)
		}
	}
}

func TestConstructionAndDiagnostics(t *testing.T) {
	for _, kind := range vocabulary {
		for _, message := range []string{"", "diagnostic", "Δ\nraw\x00text"} {
			f := failure.New(kind, message)
			if f == nil || f.Kind() != kind || f.Diagnostic() != message || f.Error() != message || f.Unwrap() != nil {
				t.Fatalf("construction mismatch: %v", f)
			}
		}
	}
}

func TestClassifiedAndForeignCauses(t *testing.T) {
	for _, cause := range []error{failure.New(failure.Timeout, "inner"), stderrors.New("timeout")} {
		f := failure.New(failure.Conflict, "outer").WithCause(cause)
		if f.Kind() != failure.Conflict || f.Diagnostic() != "outer" || f.Error() != "outer" || f.Unwrap() != cause {
			t.Fatalf("cause altered outer meaning: %v", f)
		}
		if !stderrors.Is(f, cause) {
			t.Fatal("cause is not reachable")
		}
	}
}

func TestOccurrencesAreDistinct(t *testing.T) {
	a, b := failure.New(failure.Invalid, "same"), failure.New(failure.Invalid, "same")
	if a == b || stderrors.Is(a, b) || stderrors.Is(b, a) {
		t.Fatal("equal contents became equal occurrences")
	}
}

func TestEnumerationOwnership(t *testing.T) {
	kinds := failure.Kinds()
	for i := range kinds {
		kinds[i] = failure.Kind(255)
	}
	if len(failure.Kinds()) != 10 {
		t.Fatal("enumeration length changed")
	}
	for _, kind := range failure.Kinds() {
		if kind.String() == "unknown" {
			t.Fatal("caller mutated enumeration")
		}
	}
}

func TestZeroAndInvalidNumericKinds(t *testing.T) {
	var zero failure.Failure
	if zero.Kind() != failure.Internal || zero.Diagnostic() != "" || zero.Error() != "" || zero.Unwrap() != nil {
		t.Fatal("bad zero value")
	}
	for n := 10; n < 256; n++ {
		kind := failure.Kind(n)
		if kind.String() != "unknown" || failure.New(kind, "kept").Kind() != failure.Internal || failure.New(kind, "kept").Diagnostic() != "kept" {
			t.Fatalf("invalid kind %d mishandled", n)
		}
	}
}

type foreignError struct{ code int }

func (e *foreignError) Error() string {
	if e == nil {
		panic("typed-nil cause was inspected")
	}
	return fmt.Sprintf("foreign %d", e.code)
}

func TestCauseReplacementAndNativeInspection(t *testing.T) {
	first := stderrors.New("first")
	base := failure.New(failure.Unavailable, "outer").WithCause(first)
	second := &foreignError{code: 42}
	next := base.WithCause(second)
	if next == base || base.Unwrap() != first || next.Unwrap() != second {
		t.Fatal("replacement mutated original")
	}
	if stderrors.Is(next, first) {
		t.Fatal("old cause retained")
	}
	var typed *foreignError
	if !stderrors.As(next, &typed) || typed != second {
		t.Fatal("typed cause lost")
	}
	annotated := fmt.Errorf("context: %w", next)
	if !stderrors.Is(annotated, second) {
		t.Fatal("native annotation lost identity")
	}
	var outer *failure.Failure
	if !stderrors.As(annotated, &outer) || outer != next {
		t.Fatal("native annotation lost outer failure")
	}
	removed := next.WithCause(nil)
	if removed.Unwrap() != nil || next.Unwrap() != second || removed.Kind() != next.Kind() || removed.Diagnostic() != next.Diagnostic() {
		t.Fatal("nil removal mismatch")
	}
}

func TestTypedNilCauseIsNotInspected(t *testing.T) {
	var typed *foreignError
	var cause error = typed
	f := failure.New(failure.Canceled, "outer").WithCause(cause)
	if f.Unwrap() == nil || f.Unwrap() != cause || f.Kind() != failure.Canceled || f.Error() != "outer" {
		t.Fatal("typed-nil cause changed")
	}
}
