package root

import (
	"errors"
	"strings"
	"time"

	"github.com/0xsj/atelier-wails/pkg/clock"
	"github.com/0xsj/atelier-wails/pkg/id"
	"github.com/0xsj/atelier-wails/pkg/logger"
	"github.com/0xsj/atelier-wails/pkg/provenance"
	"testing"
)

func TestLoggedHostPreservesResult(t *testing.T) {
	for _, expected := range []error{nil, errors.New("private host failure")} {
		m := &logger.Memory{}
		log, e := logger.New(clock.System{}, m, "test", logger.Info)
		if e != nil {
			t.Fatal(e)
		}
		calls := 0
		actual := runLogged(log, func() error { calls++; return expected })
		if actual != expected || calls != 1 {
			t.Fatal("host result changed")
		}
		records := m.Records()
		if len(records) != 2 || records[0].Message != "app.starting" {
			t.Fatal(records)
		}
		if expected != nil && records[1].Error["message"].Scalar() != "internal error" {
			t.Fatal("host diagnostic exposed")
		}
	}
}

type failedSink struct{}

func (failedSink) Enabled(logger.Level) bool { return true }
func (failedSink) Write(logger.Record) error { return errors.New("private sink failure") }
func (failedSink) Flush() error              { return errors.New("private flush failure") }
func TestSinkFailureDoesNotReplaceHostResult(t *testing.T) {
	log, e := logger.New(clock.System{}, failedSink{}, "test", logger.Info)
	if e != nil {
		t.Fatal(e)
	}
	want := errors.New("host")
	calls := 0
	got := runLogged(log, func() error { calls++; return want })
	if calls != 1 || got != want {
		t.Fatal("logging changed host outcome")
	}
}

func TestBootstrapLoggerBindsProcessScope(t *testing.T) {
	c := clock.NewFixed(time.UnixMilli(42))
	first := mustID(t, "01900000-0000-7000-8000-000000000001")
	ids, err := id.NewSequence(first)
	if err != nil {
		t.Fatal(err)
	}
	memory := &logger.Memory{}
	log, err := bootstrapLogger(c, ids, memory, "atelier-test")
	if err != nil {
		t.Fatal(err)
	}
	if err := log.Log(logger.Info, "app.starting", nil); err != nil {
		t.Fatal(err)
	}
	records := memory.Records()
	if len(records) != 1 {
		t.Fatal(records)
	}
	scope := records[0].Scope
	if scope["scope_id"].Scalar() != first.String() ||
		scope["work_id"].Scalar() != first.String() ||
		scope["correlation_id"].Scalar() != first.String() ||
		scope["operation"].Scalar() != "app.startup" ||
		scope["executor_id"].Scalar() != "atelier-test" ||
		scope["started_at_ms"].Scalar() != int64(42) {
		t.Fatal(scope)
	}
}

func TestBootstrapLoggerRejectsIDFailure(t *testing.T) {
	c := clock.NewFixed(time.UnixMilli(42))
	ids, err := id.NewSequence()
	if err != nil {
		t.Fatal(err)
	}
	_, err = bootstrapLogger(c, ids, &logger.Memory{}, "atelier-test")
	if err == nil || !strings.HasPrefix(err.Error(), "ID sequence exhausted") {
		t.Fatal(err)
	}
}

func mustID(t *testing.T, value string) id.ID {
	t.Helper()
	parsed, err := id.Parse(value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

var _ provenance.Clock = clock.System{}
