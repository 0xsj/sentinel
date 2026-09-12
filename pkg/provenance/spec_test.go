package provenance_test

import (
	"fmt"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
	p "github.com/0xsj/atelier-wails/pkg/provenance"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func must[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}
func ident(n int) id.ID { return must(id.Parse(fmt.Sprintf("01900000-0000-7000-8000-%012x", n))) }
func ptr[T any](v T) *T { return &v }

type wall struct {
	value time.Time
	calls int
}

func (w *wall) Now() time.Time { w.calls++; return w.value }

type ids struct {
	n, calls int
	fail     error
	fixed    *id.ID
}

func (g *ids) NewID() (id.ID, error) {
	g.calls++
	if g.fail != nil {
		return id.ID{}, g.fail
	}
	if g.fixed != nil {
		return *g.fixed, nil
	}
	g.n++
	return ident(g.n), nil
}
func setup() (*p.Factory, *wall, *ids) {
	w := &wall{value: time.UnixMilli(1000)}
	g := &ids{}
	return must(p.NewFactory(w, g)), w, g
}
func actor() p.Actor   { return must(p.NewActor(p.Service, "api")) }
func op() p.Operation  { return must(p.NewOperation("demo.export")) }
func root() p.RootSpec { return p.RootSpec{Origin: p.Startup, Operation: op(), Executor: actor()} }
func assertType(t *testing.T, e error, typ string) {
	t.Helper()
	if faults.DiagnosticTypeOf(e) != "provenance."+typ {
		t.Fatalf("wanted %s, got %v", typ, e)
	}
}

func TestP01P02Attribution(t *testing.T) {
	for _, s := range []string{"", "a b", "é", strings.Repeat("a", 129)} {
		_, e := p.NewActor(p.User, s)
		assertType(t, e, "invalid_actor")
	}
	a := must(p.NewActor(p.User, "alice"))
	b := must(p.NewActor(p.User, "bob"))
	tenant := "acme"
	attr := must(p.NewAttribution(p.AttributionSpec{Initiator: &a, OnBehalfOf: &b, Tenant: &tenant}))
	tenant = "changed"
	if *attr.Snapshot().Tenant != "acme" {
		t.Fatal("input ownership")
	}
	for _, spec := range []p.AttributionSpec{{OnBehalfOf: &b}, {Initiator: ptr(p.Anonymous()), OnBehalfOf: &b}, {Initiator: &a, OnBehalfOf: &a}} {
		_, e := p.NewAttribution(spec)
		assertType(t, e, "invalid_attribution")
	}
	if (p.Attribution{}).Snapshot().Initiator != nil || p.Anonymous().Kind() != "anonymous" {
		t.Fatal("unknown versus anonymous")
	}
}
func TestP03P04P05P06P07P08P09Transitions(t *testing.T) {
	f, w, g := setup()
	spec := root()
	alice := must(p.NewActor(p.User, "alice"))
	spec.Attribution = must(p.NewAttribution(p.AttributionSpec{Initiator: &alice, Tenant: ptr("acme")}))
	a := must(f.Open(spec))
	s := a.Snapshot()
	if s.ScopeID != ident(1) || s.Work.WorkID != ident(1) || s.Work.CorrelationID != ident(1) || s.Attempt != 1 || *s.Work.Depth != 0 || s.StartedAt.UnixMilli() != 1000 || w.calls != 1 || g.calls != 1 {
		t.Fatal(s)
	}
	c := must(f.Child(a, p.StepSpec{Operation: op(), Executor: actor()})).Snapshot()
	if c.ScopeID == s.ScopeID || c.Work.WorkID != c.ScopeID || c.Work.CorrelationID != s.Work.CorrelationID || c.Work.Causation.ID() != s.ScopeID || *c.Work.Depth != 1 || *c.Work.Attribution.Snapshot().Tenant != "acme" {
		t.Fatal(c)
	}
	event := must(p.NewReference(p.EventReference, ident(99)))
	before := g.calls
	wc := must(p.Prepare(a, p.WorkSpec{WorkID: ident(50), Operation: op(), Cause: &event}))
	if g.calls != before || w.calls != before {
		t.Fatal("prepare consumed effects")
	}
	b := must(f.Execute(wc, p.ExecutionSpec{Executor: actor(), Attempt: 1}))
	w.value = time.UnixMilli(2000)
	worker := must(p.NewActor(p.Service, "worker"))
	retry := must(f.Retry(b, worker))
	r := retry.Snapshot()
	if r.ScopeID == b.Snapshot().ScopeID || r.Work.WorkID != ident(50) || r.Work.Causation.ID() != ident(99) || r.Attempt != 2 || *r.PreviousAttempt != b.Snapshot().ScopeID || r.Executor.Identity() != "worker" || r.StartedAt.UnixMilli() != 2000 {
		t.Fatal(r)
	}
	other := must(f.Retry(b, actor())).Snapshot()
	if other.Attempt != 2 || other.ScopeID == r.ScopeID {
		t.Fatal("concurrent attempt ordinal")
	}
	resumed := must(f.Execute(wc, p.ExecutionSpec{Executor: actor(), Attempt: 3})).Snapshot()
	if resumed.PreviousAttempt != nil || resumed.Work.WorkID != ident(50) {
		t.Fatal(resumed)
	}
	explicit := root()
	explicit.WorkID = ptr(ident(80))
	x := must(f.Open(explicit)).Snapshot()
	if x.Work.WorkID != ident(80) || x.Work.CorrelationID != x.ScopeID {
		t.Fatal(x)
	}
}
func TestP10P11OverflowAndTime(t *testing.T) {
	f, w, g := setup()
	a := must(f.Open(root()))
	snap := a.Snapshot()
	snap.Attempt = math.MaxUint32
	maxAttempt := must(p.RestoreScope(snap))
	before := g.calls
	_, e := f.Retry(maxAttempt, actor())
	assertType(t, e, "attempt_exhausted")
	if g.calls != before || w.calls != before {
		t.Fatal("overflow effects")
	}
	snap = a.Snapshot()
	snap.Work.Depth = ptr(uint32(math.MaxUint32))
	snap.Work.Causation = ptr(must(p.NewReference(p.EventReference, ident(90))))
	maxDepth := must(p.RestoreScope(snap))
	_, e = f.Child(maxDepth, p.StepSpec{Operation: op(), Executor: actor()})
	assertType(t, e, "depth_exhausted")
	w.value = time.UnixMilli(1)
	child := must(f.Child(a, p.StepSpec{Operation: op(), Executor: actor()}))
	if !child.Snapshot().StartedAt.Before(a.Snapshot().StartedAt) {
		t.Fatal("clock correction lost")
	}
}
func TestP12P23EffectsAndInvalidValues(t *testing.T) {
	f, w, g := setup()
	bad := root()
	bad.Executor = p.Anonymous()
	_, e := f.Open(bad)
	assertType(t, e, "invalid_actor")
	if w.calls != 0 || g.calls != 0 {
		t.Fatal("validation effects")
	}
	w.value = time.UnixMilli(-1)
	_, e = f.Open(root())
	assertType(t, e, "invalid_time")
	if g.calls != 0 {
		t.Fatal("invalid time minted")
	}
	w.value = time.UnixMilli(1000)
	failure := faults.New(faults.Unavailable, "entropy").WithType("id.entropy")
	g.fail = failure
	_, e = f.Open(root())
	if e != failure {
		t.Fatal("generator error changed")
	}
	g.fail = nil
	g.fixed = ptr(id.ID{})
	_, e = f.Open(root())
	assertType(t, e, "invalid_generated_id")
	g.fixed = nil
	a := must(f.Open(root()))
	g.fixed = ptr(a.Snapshot().ScopeID)
	_, e = f.Child(a, p.StepSpec{Operation: op(), Executor: actor()})
	assertType(t, e, "invalid_generated_id")
	_, e = p.Prepare(p.Scope{}, p.WorkSpec{})
	assertType(t, e, "invalid_scope")
}
func TestP18P19Snapshots(t *testing.T) {
	f, _, _ := setup()
	spec := root()
	tenant := "acme"
	spec.Attribution = must(p.NewAttribution(p.AttributionSpec{Tenant: &tenant}))
	a := must(f.Open(spec))
	original := a.Snapshot()
	s := a.Snapshot()
	*s.Work.Depth = 9
	*s.Work.Origin = p.Schedule
	attr := s.Work.Attribution.Snapshot()
	*attr.Tenant = "other"
	if !reflect.DeepEqual(a.Snapshot(), original) {
		t.Fatal("snapshot alias")
	}
	bad := original
	bad.Attempt = 0
	_, e := p.RestoreScope(bad)
	assertType(t, e, "invalid_attempt")
	bad = original
	bad.PreviousAttempt = ptr(bad.ScopeID)
	_, e = p.RestoreScope(bad)
	assertType(t, e, "invalid_scope")
	w := original.Work
	w.Causation = ptr(must(p.NewReference(p.WorkReference, w.WorkID)))
	_, e = p.RestoreWork(w)
	if e == nil {
		t.Fatal("self cause accepted")
	}
}
