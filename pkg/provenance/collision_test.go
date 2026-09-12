package provenance_test

import (
	p "github.com/0xsj/atelier-wails/pkg/provenance"
	"testing"
)

func TestP23ExplicitCauseDoesNotHideExecutionCollision(t *testing.T) {
	f, _, g := setup()
	parent := must(f.Open(root()))
	g.fixed = ptr(parent.Snapshot().ScopeID)
	event := must(p.NewReference(p.EventReference, ident(90)))
	_, e := f.Child(parent, p.StepSpec{WorkID: ptr(ident(80)), Operation: op(), Executor: actor(), Cause: &event})
	assertType(t, e, "invalid_generated_id")
}
