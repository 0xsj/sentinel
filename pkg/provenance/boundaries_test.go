package provenance_test

import (
	"github.com/0xsj/atelier-wails/pkg/clock"
	"github.com/0xsj/atelier-wails/pkg/id"
	p "github.com/0xsj/atelier-wails/pkg/provenance"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestVocabularyBoundaries(t *testing.T) {
	for _, name := range []string{"", "A", "1run", "run work", "é", strings.Repeat("a", 129)} {
		_, e := p.NewOperation(name)
		assertType(t, e, "invalid_operation")
	}
	name := "a" + strings.Repeat("b", 127)
	if must(p.NewOperation(name)).String() != name {
		t.Fatal("operation altered")
	}
	for _, k := range []p.ActorKind{p.User, p.Service, p.System} {
		a := must(p.NewActor(k, "local:desktop"))
		if a.Kind() != k {
			t.Fatal(a)
		}
	}
	for _, k := range []p.ReferenceKind{p.ScopeReference, p.WorkReference, p.EventReference} {
		r := must(p.NewReference(k, ident(9)))
		if r.Kind() != k || r.ID() != ident(9) {
			t.Fatal(r)
		}
	}
	_, e := p.NewReference("bad", ident(9))
	assertType(t, e, "invalid_reference")
	_, e = p.NewReference(p.WorkReference, id.ID{})
	assertType(t, e, "invalid_reference")
	_, e = p.NewAttribution(p.AttributionSpec{Tenant: ptr("")})
	assertType(t, e, "invalid_attribution")
}
func TestUnknownRestoredAncestry(t *testing.T) {
	f, _, _ := setup()
	a := must(f.Open(root()))
	snap := a.Snapshot()
	snap.Work.CorrelationSource = p.External
	snap.Work.Depth = nil
	snap.Work.Origin = nil
	restored := must(p.RestoreScope(snap))
	child := must(f.Child(restored, p.StepSpec{Operation: op(), Executor: actor()}))
	w := child.Snapshot().Work
	if w.Depth != nil || w.Origin != nil || w.CorrelationSource != p.External || w.CorrelationID != snap.Work.CorrelationID {
		t.Fatal("invented ancestry")
	}
}
func TestTimeNormalizationAndBounds(t *testing.T) {
	f, w, g := setup()
	for _, origin := range []p.Origin{p.Request, p.Schedule, p.Backfill, p.Startup} {
		s := root()
		s.Origin = origin
		w.value = time.Unix(0, 1999999).In(time.FixedZone("offset", 3600))
		got := must(f.Open(s)).Snapshot().StartedAt
		if got.UnixNano() != 1000000 || got.Location() != time.UTC {
			t.Fatal(got)
		}
	}
	before := g.calls
	w.value = time.UnixMilli(1 << 48)
	_, e := f.Open(root())
	assertType(t, e, "invalid_time")
	if g.calls != before {
		t.Fatal("generated on invalid time")
	}
	w.value = time.UnixMilli((1 << 48) - 1).Add(999999 * time.Nanosecond)
	if must(f.Open(root())).Snapshot().StartedAt.UnixMilli() != (1<<48)-1 {
		t.Fatal("upper bound")
	}
}
func TestRestoreRejectsCombinations(t *testing.T) {
	f, _, _ := setup()
	a := must(f.Open(root()))
	for _, mutate := range []func(*p.ScopeSnapshot){
		func(s *p.ScopeSnapshot) {
			s.Work.Causation = ptr(must(p.NewReference(p.ScopeReference, s.ScopeID)))
			s.Work.Depth = ptr(uint32(1))
		},
		func(s *p.ScopeSnapshot) { s.Work.Depth = ptr(uint32(1)) },
		func(s *p.ScopeSnapshot) { s.PreviousAttempt = ptr(ident(99)) },
		func(s *p.ScopeSnapshot) { s.Executor = p.Anonymous() },
		func(s *p.ScopeSnapshot) { s.Work.CorrelationID = id.ID{} },
	} {
		s := a.Snapshot()
		mutate(&s)
		if _, e := p.RestoreScope(s); e == nil {
			t.Fatal("invalid snapshot accepted")
		}
	}
	retry := must(f.Retry(a, actor()))
	input := retry.Snapshot()
	restored := must(p.RestoreScope(input))
	before := restored.Snapshot()
	*input.PreviousAttempt = ident(99)
	*input.Work.Depth = 9
	out := restored.Snapshot()
	*out.PreviousAttempt = ident(98)
	if !reflect.DeepEqual(before, restored.Snapshot()) {
		t.Fatal("restoration alias")
	}
}
func TestRefusalsHaveNoEffectsAndRetryPreservesWork(t *testing.T) {
	f, w, g := setup()
	a := must(f.Open(root()))
	before := g.calls
	_, e := f.Child(a, p.StepSpec{WorkID: ptr(a.Snapshot().Work.WorkID), Operation: op(), Executor: actor()})
	assertType(t, e, "invalid_work")
	_, e = f.Execute(a.WorkContext(), p.ExecutionSpec{Executor: actor()})
	assertType(t, e, "invalid_attempt")
	if g.calls != before || w.calls != before {
		t.Fatal("invalid inputs consumed effects")
	}
	retry := must(f.Retry(a, must(p.NewActor(p.System, "desktop"))))
	if !reflect.DeepEqual(a.Snapshot().Work, retry.Snapshot().Work) {
		t.Fatal("retry changed work")
	}
	g.fixed = ptr(a.Snapshot().ScopeID)
	_, e = f.Retry(a, actor())
	assertType(t, e, "invalid_generated_id")
}
func TestExistingClockAndSequenceIntegration(t *testing.T) {
	seq := must(id.NewSequence(ident(1), ident(2)))
	f := must(p.NewFactory(clock.NewFixed(time.UnixMilli(1000)), seq))
	a := must(f.Open(root()))
	b := must(f.Retry(a, actor()))
	if b.Snapshot().ScopeID != ident(2) {
		t.Fatal("sequence not consumed")
	}
	_, e := f.Retry(b, actor())
	if e == nil {
		t.Fatal("exhaustion hidden")
	}
	_, e = p.NewFactory(nil, seq)
	assertType(t, e, "invalid_configuration")
	_, e = p.NewFactory(clock.System{}, nil)
	assertType(t, e, "invalid_configuration")
}
