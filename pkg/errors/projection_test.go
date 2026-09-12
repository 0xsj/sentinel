package errors_test

import (
	"context"
	stderrors "errors"
	"fmt"
	failure "github.com/0xsj/atelier-wails/pkg/errors"
	"reflect"
	"testing"
)

func TestPublicProjectionAndFallbacks(t *testing.T) {
	for _, kind := range failure.Kinds() {
		f := failure.New(kind, "safe message").WithType("workspace.name_taken").WithField("name", "taken").WithDetail("sql", "private").WithCause(stderrors.New("private source"))
		got, present := failure.Public(f)
		if !present {
			t.Fatal("failure missing")
		}
		if kind == failure.Internal {
			if !reflect.DeepEqual(got, failure.PublicInfo{Kind: failure.Internal, Message: "internal error"}) {
				t.Fatalf("Internal leaked: %#v", got)
			}
		} else if got.Kind != kind || got.Message != "safe message" || got.Type != "workspace.name_taken" || !reflect.DeepEqual(got.Fields, map[string]string{"name": "taken"}) {
			t.Fatalf("bad public snapshot: %#v", got)
		}
		if failure.DiagnosticTypeOf(f) != "workspace.name_taken" || failure.DetailsOf(f)["sql"] != "private" {
			t.Fatal("private frame lost")
		}
	}
}

func TestEmptyPublicValues(t *testing.T) {
	f := failure.New(failure.Invalid, "").WithType("old").WithType("")
	got, _ := failure.Public(f)
	if got.Message != "request failed" || got.Type != "" || got.Kind != failure.Invalid {
		t.Fatalf("bad fallback: %#v", got)
	}
	if failure.Message(f) != "request failed" || failure.TypeOf(f) != "" {
		t.Fatal("accessor policy mismatch")
	}
}

func TestMetadataAndSnapshotOwnership(t *testing.T) {
	fields := map[string]string{"a": "first", "keep": "yes"}
	details := map[string]string{"a": "private", "keep": "private keep"}
	base := failure.New(failure.Invalid, "safe").WithFields(fields).WithDetails(details)
	fields["a"] = "input edit"
	details["a"] = "input edit"
	next := base.WithField("a", "second").WithDetail("a", "second private")
	if failure.FieldsOf(base)["a"] != "first" || failure.DetailsOf(base)["a"] != "private" {
		t.Fatal("source/input mutated")
	}
	if failure.FieldsOf(next)["a"] != "second" || failure.FieldsOf(next)["keep"] != "yes" || failure.DetailsOf(next)["a"] != "second private" || failure.DetailsOf(next)["keep"] != "private keep" {
		t.Fatal("merge failed")
	}
	public, _ := failure.Public(next)
	public.Fields["a"] = "snapshot edit"
	projected := failure.FieldsOf(next)
	projected["a"] = "accessor edit"
	private := failure.DetailsOf(next)
	private["a"] = "details edit"
	if failure.FieldsOf(next)["a"] != "second" || failure.DetailsOf(next)["a"] != "second private" {
		t.Fatal("returned map aliases stored data")
	}
	nilMerged := next.WithFields(nil).WithDetails(nil)
	if !reflect.DeepEqual(failure.FieldsOf(next), failure.FieldsOf(nilMerged)) || !reflect.DeepEqual(failure.DetailsOf(next), failure.DetailsOf(nilMerged)) {
		t.Fatal("nil merge lost data")
	}
}

func TestConditionIdentityAcrossDerivation(t *testing.T) {
	sentinel := failure.New(failure.Conflict, "name taken")
	other := failure.New(failure.Conflict, "name taken")
	cause := stderrors.New("database")
	derived := sentinel.WithType("workspace.name_taken").WithFields(map[string]string{"name": "taken"}).WithDetails(map[string]string{"table": "private"}).WithCause(cause)
	if !stderrors.Is(derived, sentinel) || !stderrors.Is(sentinel, derived) || stderrors.Is(derived, other) || !stderrors.Is(derived, cause) {
		t.Fatal("condition/source identity mismatch")
	}
	if sentinel.Unwrap() != nil || failure.TypeOf(sentinel) != "" || len(failure.FieldsOf(sentinel)) != 0 || len(failure.DetailsOf(sentinel)) != 0 {
		t.Fatal("template mutated")
	}
	annotated := fmt.Errorf("save: %w", derived)
	kind, ok := failure.KindOf(annotated)
	if !ok || kind != failure.Conflict || !stderrors.Is(annotated, sentinel) || failure.TypeOf(annotated) != "workspace.name_taken" {
		t.Fatal("annotation lost meaning")
	}
}

func TestOuterTranslationOwnsWholeFrame(t *testing.T) {
	inner := failure.New(failure.Invalid, "inner public").WithType("inner.type").WithField("inner", "bad").WithDetail("inner", "private")
	outer := failure.Wrap(inner, failure.Conflict, "outer public")
	got, _ := failure.Public(outer)
	if got.Kind != failure.Conflict || got.Message != "outer public" || got.Type != "" || len(got.Fields) != 0 || len(failure.DetailsOf(outer)) != 0 {
		t.Fatalf("inherited inner data: %#v", got)
	}
	if !stderrors.Is(outer, inner) {
		t.Fatal("translation lost source")
	}
	if failure.Wrap(nil, failure.Internal, "ignored") != nil {
		t.Fatal("nil wrap became failure")
	}
}

func TestReplacingSourcePreservesMetadata(t *testing.T) {
	first, second := stderrors.New("first"), stderrors.New("second")
	base := failure.New(failure.Conflict, "safe").WithType("domain.condition").WithField("a", "b").WithDetail("x", "y").WithCause(first)
	next := base.WithCause(second)
	before, _ := failure.Public(base)
	after, _ := failure.Public(next)
	if !reflect.DeepEqual(before, after) || !reflect.DeepEqual(failure.DetailsOf(base), failure.DetailsOf(next)) || !stderrors.Is(next, base) || next.Unwrap() != second || base.Unwrap() != first {
		t.Fatal("replacement damaged frame/source")
	}
}

func TestPresenceUnknownAndTypedNil(t *testing.T) {
	if view, ok := failure.Public(nil); ok || !reflect.DeepEqual(view, failure.PublicInfo{}) {
		t.Fatal("nil is failure")
	}
	if _, ok := failure.KindOf(nil); ok || failure.IsKind(nil, failure.Internal) {
		t.Fatal("nil classified")
	}
	if failure.Message(nil) != "" || failure.TypeOf(nil) != "" || len(failure.FieldsOf(nil)) != 0 {
		t.Fatal("nil public accessors")
	}
	var typed *failure.Failure
	for _, err := range []error{typed, stderrors.New("private timeout"), fmt.Errorf("context: %w", typed)} {
		if kind, ok := failure.KindOf(err); ok || kind != failure.Internal {
			t.Fatal("unknown classified")
		}
		view, present := failure.Public(err)
		if !present || !reflect.DeepEqual(view, failure.PublicInfo{Kind: failure.Internal, Message: "internal error"}) {
			t.Fatalf("unknown leaked: %#v", view)
		}
	}
}

func TestCancellationAndDeadlineThroughAnnotation(t *testing.T) {
	for err, kind := range map[error]failure.Kind{context.Canceled: failure.Canceled, context.DeadlineExceeded: failure.Timeout} {
		annotated := fmt.Errorf("operation: %w", err)
		if actual, ok := failure.KindOf(annotated); !ok || actual != kind || !failure.IsKind(annotated, kind) {
			t.Fatal("context classification lost")
		}
		public, _ := failure.Public(annotated)
		if public.Kind != kind || public.Message != "request failed" {
			t.Fatalf("bad context projection: %#v", public)
		}
		outer := failure.New(failure.Conflict, "translated").WithCause(err)
		if !failure.IsKind(outer, failure.Conflict) {
			t.Fatal("source overrode explicit outer kind")
		}
	}
}

func TestAggregateHasNoImplicitSummary(t *testing.T) {
	known, unknown := failure.New(failure.Timeout, "safe"), stderrors.New("private")
	for _, aggregate := range []error{stderrors.Join(known, unknown), stderrors.Join(unknown, known), stderrors.Join(known)} {
		if _, ok := failure.KindOf(fmt.Errorf("batch: %w", aggregate)); ok {
			t.Fatal("aggregate selected branch")
		}
		public, _ := failure.Public(aggregate)
		if public.Kind != failure.Internal || public.Message != "internal error" {
			t.Fatal("aggregate public leak")
		}
		summary := failure.Wrap(aggregate, failure.Unavailable, "batch incomplete")
		if !failure.IsKind(summary, failure.Unavailable) || !stderrors.Is(summary, known) {
			t.Fatal("explicit summary failed")
		}
	}
}

type cycle struct{}

func (e *cycle) Error() string { return "private" }
func (e *cycle) Unwrap() error { return e }
func TestBoundedInspection(t *testing.T) {
	if _, ok := failure.KindOf(&cycle{}); ok {
		t.Fatal("cycle classified")
	}
	var chain error = failure.New(failure.Timeout, "safe")
	for i := 0; i < 127; i++ {
		chain = fmt.Errorf("layer: %w", chain)
	}
	if !failure.IsKind(chain, failure.Timeout) {
		t.Fatal("128th frame not inspected")
	}
	chain = fmt.Errorf("layer: %w", chain)
	if _, ok := failure.KindOf(chain); ok {
		t.Fatal("walk exceeded bound")
	}
	public, _ := failure.Public(&cycle{})
	if public.Message != "internal error" {
		t.Fatal("cycle leaked")
	}
}
