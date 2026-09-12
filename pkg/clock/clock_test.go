package clock_test

import (
	"github.com/0xsj/atelier-wails/pkg/clock"
	"sync"
	"testing"
	"time"
)

func TestFixedPreservesInstantAndNormalizesUTC(t *testing.T) {
	at := time.Date(2026, 9, 11, 12, 34, 56, 123456789, time.FixedZone("test", 3*60*60))
	var c clock.Clock = clock.NewFixed(at)
	for i := 0; i < 100; i++ {
		got := c.Now()
		if !got.Equal(at) || got.Location() != time.UTC || got.Nanosecond() != 123456789 {
			t.Fatalf("instant changed: %v", got)
		}
	}
}
func TestFixedInstancesAreIndependent(t *testing.T) {
	a, b := clock.NewFixed(time.Unix(1, 2)), clock.NewFixed(time.Unix(3, 4))
	returned := a.Now()
	returned = returned.Add(time.Hour)
	if !a.Now().Equal(time.Unix(1, 2)) || !b.Now().Equal(time.Unix(3, 4)) || returned.Equal(a.Now()) {
		t.Fatal("instance or returned-value aliasing")
	}
}
func TestBeforeEpoch(t *testing.T) {
	at := time.Unix(-123456, 987654321)
	got := clock.NewFixed(at).Now()
	if !got.Equal(at) || got.Location() != time.UTC {
		t.Fatal("pre-epoch time changed")
	}
}
func TestZeroFixed(t *testing.T) {
	var c clock.Fixed
	if !c.Now().IsZero() || c.Now().Location() != time.UTC {
		t.Fatal("zero clock changed")
	}
}
func TestSystemClock(t *testing.T) {
	var c clock.Clock = clock.System{}
	got := c.Now()
	if got.IsZero() || got.Location() != time.UTC {
		t.Fatal("invalid system wall time")
	}
}
func TestConcurrentFixedReads(t *testing.T) {
	at := time.Unix(123, 456)
	c := clock.NewFixed(at)
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if !c.Now().Equal(at) {
					t.Error("concurrent clock drift")
				}
			}
		}()
	}
	wg.Wait()
}
