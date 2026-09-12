package logger_test

import (
	"bytes"
	"errors"
	"github.com/0xsj/atelier-wails/pkg/clock"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
	l "github.com/0xsj/atelier-wails/pkg/logger"
	p "github.com/0xsj/atelier-wails/pkg/provenance"
	"github.com/0xsj/atelier-wails/pkg/secret"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

func must[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}
func logFor(s l.Sink) *l.Logger {
	return must(l.New(clock.NewFixed(time.UnixMilli(1000)), s, "atelier", l.Info))
}
func emit(t *testing.T, log *l.Logger) {
	t.Helper()
	if e := log.Log(l.Info, "test", nil); e != nil {
		t.Fatal(e)
	}
}
func TestConfiguration(t *testing.T) {
	for _, s := range []string{"INFO", " info", "", "trace"} {
		if _, e := l.ParseLevel(s); e == nil {
			t.Fatal(s)
		}
	}
	for _, v := range []l.Level{l.Debug, l.Info, l.Warn, l.Error} {
		if must(l.ParseLevel(v.String())) != v {
			t.Fatal(v)
		}
	}
	if _, e := l.New(nil, &l.Memory{}, "app", l.Info); e == nil {
		t.Fatal("nil clock")
	}
	if _, e := l.NewConsole(nil); e == nil {
		t.Fatal("nil writer")
	}
	if _, e := l.NewSlog(nil, nil); e == nil {
		t.Fatal("nil handler")
	}
}

type countClock struct {
	calls int
	at    time.Time
}

func (c *countClock) Now() time.Time { c.calls++; return c.at }
func TestFilteringAndTime(t *testing.T) {
	c := &countClock{at: time.UnixMilli(0)}
	m := &l.Memory{}
	log := must(l.New(c, m, "app", l.Warn))
	if e := log.Log(l.Info, "filtered", nil); e != nil {
		t.Fatal(e)
	}
	if c.calls != 0 {
		t.Fatal("filtered clock")
	}
	if e := log.Log(l.Warn, "kept", nil); e != nil {
		t.Fatal(e)
	}
	if len(m.Records()) != 1 || c.calls != 1 {
		t.Fatal("floor")
	}
	d := must(l.New(c, l.Discard{}, "app", l.Debug))
	emit(t, d)
	if c.calls != 1 {
		t.Fatal("discard clock")
	}
	c.at = time.Unix(0, -1)
	if e := log.Log(l.Error, "bad time", nil); !faults.IsKind(e, faults.Invalid) {
		t.Fatal(e)
	}
	if len(m.Records()) != 1 {
		t.Fatal("invalid time wrote")
	}
}
func TestBindingAndOwnership(t *testing.T) {
	m := &l.Memory{}
	parent := logFor(m)
	fields := l.Fields{"x": l.Int(1), "token": l.Secret(secret.New("PRIVATE"))}
	child := parent.With(fields)
	fields["x"] = l.Int(9)
	if e := child.Log(l.Info, "child", l.Fields{"x": l.Int(2), "service": l.Text("spoof")}); e != nil {
		t.Fatal(e)
	}
	emit(t, parent)
	r := m.Records()
	if r[0].Fields["x"].Scalar() != int64(2) || r[0].Service != "atelier" || r[0].Fields["token"].Scalar() != "[REDACTED]" || len(r[1].Fields) != 0 {
		t.Fatal(r)
	}
	r[0].Fields["x"] = l.Int(99)
	if m.Records()[0].Fields["x"].Scalar() != int64(2) {
		t.Fatal("alias")
	}
}
func TestConsoleAndSlogSwap(t *testing.T) {
	for _, json := range []bool{false, true} {
		var b bytes.Buffer
		var sink l.Sink
		if json {
			sink = must(l.NewSlog(slog.NewJSONHandler(&b, nil), nil))
		} else {
			sink = must(l.NewConsole(&b))
		}
		log := logFor(sink)
		if e := log.Log(l.Info, "line\n\x1b[31m", l.Fields{"token": l.Secret(secret.New("PRIVATE"))}); e != nil {
			t.Fatal(e)
		}
		if e := log.Flush(); e != nil {
			t.Fatal(e)
		}
		s := b.String()
		if strings.Count(s, "\n") != 1 || strings.Contains(s, "\x1b") || strings.Contains(s, "PRIVATE") || !strings.Contains(s, "REDACTED") {
			t.Fatal(s)
		}
	}
}
func TestSafeErrors(t *testing.T) {
	m := &l.Memory{}
	log := logFor(m)
	raw := errors.New("PRIVATE")
	internal := faults.New(faults.Internal, "PRIVATE").WithType("private.type").WithField("secret", "PRIVATE")
	safe := faults.New(faults.Invalid, "bad input").WithType("input.invalid").WithCause(raw).WithField("name", "required")
	for _, e := range []error{raw, internal, safe} {
		emit(t, log.WithError(e))
	}
	emit(t, log.WithError(safe).WithError(nil))
	r := m.Records()
	if r[0].Error["classified"].Scalar() != false || r[1].Error["message"].Scalar() != "internal error" || len(r[1].ErrorFields) != 0 || r[2].Error["has_cause"].Scalar() != true || r[2].ErrorFields["name"].Scalar() != "required" || len(r[3].Error) != 0 {
		t.Fatal(r)
	}
}
func TestProvenance(t *testing.T) {
	i := must(id.Parse("01900000-0000-7000-8000-000000000001"))
	seq := must(id.NewSequence(i))
	f := must(p.NewFactory(clock.NewFixed(time.UnixMilli(5)), seq))
	scope := must(f.Open(p.RootSpec{Origin: p.Startup, Operation: must(p.NewOperation("app.start")), Executor: must(p.NewActor(p.System, "desktop"))}))
	m := &l.Memory{}
	log := must(logFor(m).WithScope(scope))
	emit(t, log)
	r := m.Records()[0]
	if r.Scope["scope_id"].Scalar() != i.String() || r.Scope["started_at_ms"].Scalar() != int64(5) || r.Time.UnixMilli() != 1000 {
		t.Fatal(r)
	}
	if _, ok := r.Scope["tenant"]; ok {
		t.Fatal("invented tenant")
	}
	if _, e := log.WithScope(p.Scope{}); e == nil {
		t.Fatal("invalid scope")
	}
}

type broken struct{}

func (broken) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }
func (broken) Flush() error              { return io.ErrUnexpectedEOF }
func TestDeliveryFailure(t *testing.T) {
	log := logFor(must(l.NewConsole(broken{})))
	if e := log.Log(l.Info, "test", nil); !errors.Is(e, io.ErrClosedPipe) || !faults.IsKind(e, faults.Unavailable) {
		t.Fatal(e)
	}
	if e := log.Flush(); !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatal(e)
	}
}
func TestConcurrent(t *testing.T) {
	m := &l.Memory{}
	log := logFor(m)
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				if e := log.Log(l.Info, "event", nil); e != nil {
					t.Error(e)
				}
			}
		}()
	}
	wg.Wait()
	if len(m.Records()) != 200 {
		t.Fatal("lost records")
	}
}

func TestTimeBoundsAndSinkFloor(t *testing.T) {
	var b bytes.Buffer
	c := &countClock{at: time.UnixMilli(7).Add(999999 * time.Nanosecond)}
	sink := must(l.NewSlog(slog.NewTextHandler(&b, &slog.HandlerOptions{Level: slog.LevelError}), nil))
	log := must(l.New(c, sink, "app", l.Debug))
	emit(t, log)
	if c.calls != 0 {
		t.Fatal("handler filtering consumed time")
	}
	m := &l.Memory{}
	log = must(l.New(c, m, "app", l.Info))
	emit(t, log)
	if m.Records()[0].Time.UnixNano() != 7000000 {
		t.Fatal("submillisecond time")
	}
	c.at = time.UnixMilli(1 << 48)
	if e := log.Log(l.Info, "bad", nil); faults.DiagnosticTypeOf(e) != "logger.invalid_time" {
		t.Fatal(e)
	}
	if len(m.Records()) != 1 {
		t.Fatal("invalid time emitted")
	}
}
func TestProvenanceOptionalFields(t *testing.T) {
	i := must(id.Parse("01900000-0000-7000-8000-000000000001"))
	j := must(id.Parse("01900000-0000-7000-8000-000000000002"))
	seq := must(id.NewSequence(i, j))
	f := must(p.NewFactory(clock.NewFixed(time.UnixMilli(5)), seq))
	alice := must(p.NewActor(p.User, "alice"))
	bob := must(p.NewActor(p.User, "bob"))
	tenant := "local"
	a := must(p.NewAttribution(p.AttributionSpec{Initiator: &alice, OnBehalfOf: &bob, Tenant: &tenant}))
	scope := must(f.Open(p.RootSpec{Origin: p.Request, Attribution: a, Operation: must(p.NewOperation("app.run")), Executor: must(p.NewActor(p.System, "desktop"))}))
	scope = must(f.Retry(scope, must(p.NewActor(p.Service, "worker"))))
	snap := scope.Snapshot()
	snap.Work.Origin = nil
	snap.Work.Depth = nil
	snap.Work.CorrelationSource = p.External
	cause := must(p.NewReference(p.EventReference, i))
	snap.Work.Causation = &cause
	scope = must(p.RestoreScope(snap))
	m := &l.Memory{}
	log := must(logFor(m).WithScope(scope))
	emit(t, log)
	s := m.Records()[0].Scope
	for k, want := range map[string]any{"attempt": int64(2), "previous_attempt": i.String(), "cause_kind": "event", "cause_id": i.String(), "initiator_id": "alice", "on_behalf_of_id": "bob", "tenant": "local", "executor_id": "worker", "correlation_source": "external"} {
		if s[k].Scalar() != want {
			t.Fatalf("%s: %v", k, s[k])
		}
	}
	if _, ok := s["depth"]; ok {
		t.Fatal("invented depth")
	}
	if _, ok := s["origin"]; ok {
		t.Fatal("invented origin")
	}
}
