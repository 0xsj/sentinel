package id_test

import (
	"github.com/0xsj/atelier-wails/pkg/clock"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
	"testing"
	"time"
)

func TestExistingClockAndRepeatedFixtures(t *testing.T) {
	c := clock.NewFixed(time.UnixMilli(tick))
	g, err := id.NewV7WithEntropy(c, readerFunc(zeros))
	if err != nil {
		t.Fatal(err)
	}
	a := next(t, g)
	if a.String() != first {
		t.Fatal(a)
	}
	s, err := id.NewSequence(a, a)
	if err != nil {
		t.Fatal(err)
	}
	if next(t, s) != a || next(t, s) != a {
		t.Fatal("fixture duplicates lost")
	}
}

func TestHighSeedCapacity(t *testing.T) {
	c := newTestClock(time.UnixMilli(tick))
	calls := 0
	g := generator(t, c, readerFunc(func(b []byte) (int, error) {
		calls++
		for i := range b {
			b[i] = 255
		}
		return len(b), nil
	}))
	for range 2049 {
		next(t, g)
	}
	_, err := g.NewID()
	checkFailure(t, err, faults.Unavailable, "id.exhausted")
	if calls != 2049 {
		t.Fatal("exhaustion read entropy")
	}
}
