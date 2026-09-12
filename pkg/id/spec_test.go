package id_test

import (
	"bytes"
	stderrors "errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
)

const tick int64 = 0x0123456789ab
const first = "01234567-89ab-7000-8000-000000000000"

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(b []byte) (int, error) { return f(b) }
func zeros(b []byte) (int, error)               { clear(b); return len(b), nil }
func requireID(t *testing.T, s string) id.ID {
	t.Helper()
	v, e := id.Parse(s)
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func next(t *testing.T, g interface{ NewID() (id.ID, error) }) id.ID {
	t.Helper()
	v, e := g.NewID()
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func checkFailure(t *testing.T, e error, k faults.Kind, typ string) {
	t.Helper()
	if !faults.IsKind(e, k) || faults.DiagnosticTypeOf(e) != typ {
		t.Fatalf("unexpected failure: %v / %s", e, faults.DiagnosticTypeOf(e))
	}
}
func generator(t *testing.T, c *testClock, r io.Reader) *id.V7 {
	t.Helper()
	g, e := id.NewV7WithEntropy(c, r)
	if e != nil {
		t.Fatal(e)
	}
	return g
}

func TestI01Parse(t *testing.T) {
	for _, s := range []string{"F81D4FAE-7DEC-41D0-A765-00A0C91E6BF6", "01234567-89ab-f000-8000-000000000000", "01234567-89ab-0000-8000-000000000000"} {
		v := requireID(t, s)
		if v.String() != strings.ToLower(s) {
			t.Fatal("canonical formatting")
		}
	}
	for _, s := range []string{"", first + " ", " " + first, strings.ReplaceAll(first, "-", ""), "{" + first + "}", "urn:uuid:" + first, "01234567_89ab-7000-8000-000000000000", "01234567-89ab-7000-8000-00000000000g", "00000000-0000-0000-0000-000000000000", "01234567-89ab-7000-0000-000000000000", "01234567-89ab-7000-c000-000000000000", "01234567-89ab-7000-f000-000000000000"} {
		v, e := id.Parse(s)
		checkFailure(t, e, faults.Invalid, "id.invalid")
		if !v.IsZero() || e.Error() != "invalid ID" {
			t.Fatalf("parse disclosure or nonzero failure: %v", e)
		}
	}
}
func TestI02ValueAndTime(t *testing.T) {
	v := requireID(t, "017F22E2-79B0-7CC3-98C4-DC0C0C07398F")
	m := map[id.ID]string{v: "found"}
	if m[requireID(t, v.String())] != "found" {
		t.Fatal("value equality")
	}
	stamp, ok := v.Time()
	if !ok || stamp.UnixMilli() != 1645557742000 || stamp.Location() != time.UTC || v.Version() != 7 {
		t.Fatal("RFC time fixture")
	}
	if _, ok := requireID(t, "f81d4fae-7dec-41d0-a765-00a0c91e6bf6").Time(); ok {
		t.Fatal("v4 has no v7 timestamp")
	}
	b, e := v.MarshalText()
	if e != nil || string(b) != v.String() {
		t.Fatal("marshal")
	}
	before := v
	if v.UnmarshalText([]byte("bad")) == nil || v != before {
		t.Fatal("failed unmarshal changed receiver")
	}
	if e := v.UnmarshalText([]byte(first)); e != nil || v.String() != first {
		t.Fatal("unmarshal")
	}
	var zero id.ID
	if !zero.IsZero() || zero.String() != "00000000-0000-0000-0000-000000000000" {
		t.Fatal("zero sentinel")
	}
	if _, e := zero.MarshalText(); e == nil {
		t.Fatal("zero marshal accepted")
	}
}
func TestI03ExactV7Bytes(t *testing.T) {
	c := newTestClock(time.UnixMilli(tick))
	g := generator(t, c, bytes.NewReader([]byte{7, 255, 255, 255, 255, 255, 255, 255, 255, 255}))
	v := next(t, g)
	if v.String() != "01234567-89ab-77ff-bfff-ffffffffffff" {
		t.Fatal(v)
	}
	if got := next(t, generator(t, c, readerFunc(zeros))).String(); got != first {
		t.Fatal(got)
	}
}
func TestI04OrderingAndRollback(t *testing.T) {
	c := newTestClock(time.UnixMilli(tick))
	calls := 0
	g := generator(t, c, readerFunc(func(b []byte) (int, error) {
		calls++
		clear(b)
		if calls == 1 {
			for i := range b {
				b[i] = 255
			}
		}
		return len(b), nil
	}))
	a := next(t, g)
	b := next(t, g)
	c.Set(time.UnixMilli(tick - 100))
	d := next(t, g)
	if !(a.String() < b.String() && b.String() < d.String()) || b.String() != "01234567-89ab-7800-8000-000000000000" {
		t.Fatal("counter ordering")
	}
	stamp, _ := d.Time()
	if stamp.UnixMilli() != tick {
		t.Fatal("rollback lost high water")
	}
	c.Set(time.UnixMilli(tick + 1))
	e := next(t, g)
	if e.String() != "01234567-89ac-7000-8000-000000000000" || e.String() <= d.String() {
		t.Fatal("new tick did not reset counter")
	}
}
func TestI05ExhaustionAndRecovery(t *testing.T) {
	c := newTestClock(time.UnixMilli(tick))
	calls := 0
	g := generator(t, c, readerFunc(func(b []byte) (int, error) { calls++; return zeros(b) }))
	var last id.ID
	for range 4096 {
		v := next(t, g)
		if !last.IsZero() && v.String() <= last.String() {
			t.Fatal("duplicate/regression")
		}
		last = v
	}
	for range 2 {
		v, e := g.NewID()
		checkFailure(t, e, faults.Unavailable, "id.exhausted")
		if !v.IsZero() {
			t.Fatal("failure returned ID")
		}
	}
	if calls != 4096 {
		t.Fatal("exhaustion consumed entropy")
	}
	c.Set(time.UnixMilli(tick - 1))
	_, e := g.NewID()
	checkFailure(t, e, faults.Unavailable, "id.exhausted")
	c.Set(time.UnixMilli(tick + 1))
	if next(t, g).String() <= last.String() {
		t.Fatal("recovery regressed")
	}
}
func TestI06EntropyFailureAtomicity(t *testing.T) {
	c := newTestClock(time.UnixMilli(tick))
	sentinel := stderrors.New("private entropy diagnostic")
	fail := true
	g := generator(t, c, readerFunc(func(b []byte) (int, error) {
		if fail {
			return 0, sentinel
		}
		return zeros(b)
	}))
	v, e := g.NewID()
	checkFailure(t, e, faults.Unavailable, "id.entropy")
	if !v.IsZero() || !stderrors.Is(e, sentinel) {
		t.Fatal("failure/cause")
	}
	pub, _ := faults.Public(e)
	if strings.Contains(pub.Message, "private") {
		t.Fatal("cause disclosed")
	}
	fail = false
	if next(t, g).String() != first {
		t.Fatal("initial failure consumed state")
	}
	c.Set(time.UnixMilli(tick + 100))
	fail = true
	_, e = g.NewID()
	checkFailure(t, e, faults.Unavailable, "id.entropy")
	fail = false
	c.Set(time.UnixMilli(tick))
	if next(t, g).String() != "01234567-89ab-7001-8000-000000000000" {
		t.Fatal("failed new tick committed state")
	}
	short := generator(t, c, bytes.NewReader(make([]byte, 9)))
	_, e = short.NewID()
	if !stderrors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatal("short read cause")
	}
}
func TestI07TimeRange(t *testing.T) {
	c := newTestClock(time.UnixMilli(0))
	calls := 0
	g := generator(t, c, readerFunc(func(b []byte) (int, error) { calls++; return zeros(b) }))
	if next(t, g).String() != "00000000-0000-7000-8000-000000000000" {
		t.Fatal("epoch rejected")
	}
	for _, wall := range []time.Time{time.Unix(0, -1), time.UnixMilli(1 << 48)} {
		c.Set(wall)
		_, e := g.NewID()
		checkFailure(t, e, faults.Invalid, "id.time_range")
	}
	if calls != 1 {
		t.Fatal("bad time consumed entropy")
	}
	c.Set(time.UnixMilli((1 << 48) - 1).Add(999999 * time.Nanosecond))
	v := next(t, g)
	stamp, _ := v.Time()
	if stamp.UnixMilli() != (1<<48)-1 {
		t.Fatal("upper bound")
	}
}
func TestI08Sequence(t *testing.T) {
	a, b := requireID(t, first), requireID(t, "01234567-89ab-7001-8000-000000000000")
	items := []id.ID{a, b}
	s, e := id.NewSequence(items...)
	if e != nil {
		t.Fatal(e)
	}
	items[0] = b
	if next(t, s) != a || next(t, s) != b {
		t.Fatal("sequence ownership/order")
	}
	for range 2 {
		_, e = s.NewID()
		checkFailure(t, e, faults.Unavailable, "id.sequence_exhausted")
	}
	empty, e := id.NewSequence()
	if e != nil {
		t.Fatal(e)
	}
	_, e = empty.NewID()
	checkFailure(t, e, faults.Unavailable, "id.sequence_exhausted")
	if _, e = id.NewSequence(id.ID{}); e == nil {
		t.Fatal("nil fixture")
	}
}
func TestI09SystemAndConcurrent(t *testing.T) {
	c := newTestClock(time.UnixMilli(tick))
	g, e := id.NewV7(c)
	if e != nil {
		t.Fatal(e)
	}
	if next(t, g).Version() != 7 {
		t.Fatal("production adapter")
	}
	g = generator(t, c, readerFunc(zeros))
	ch := make(chan id.ID, 1000)
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				v, e := g.NewID()
				if e != nil {
					t.Error(e)
					return
				}
				ch <- v
			}
		}()
	}
	wg.Wait()
	close(ch)
	seen := map[id.ID]bool{}
	for v := range ch {
		if seen[v] {
			t.Fatal("duplicate")
		}
		seen[v] = true
	}
	if len(seen) != 1000 {
		t.Fatal("lost ID")
	}
	if _, e := id.NewV7(nil); e == nil {
		t.Fatal("nil clock")
	}
	if _, e := id.NewV7WithEntropy(c, nil); e == nil {
		t.Fatal("nil entropy")
	}
}

// Mutable wall time is a private test fixture, not a new clock API.
type testClock struct {
	mu sync.Mutex
	at time.Time
}

func newTestClock(at time.Time) *testClock { return &testClock{at: at} }
func (c *testClock) Now() time.Time        { c.mu.Lock(); defer c.mu.Unlock(); return c.at }
func (c *testClock) Set(at time.Time)      { c.mu.Lock(); defer c.mu.Unlock(); c.at = at }
