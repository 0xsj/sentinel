package memory_test

import (
	"sync"
	"testing"

	"github.com/0xsj/atelier-wails/internal/preferences/domain"
	"github.com/0xsj/atelier-wails/internal/preferences/infra/memory"
)

func TestConditionalOperationsAndOrdering(t *testing.T) {
	store := memory.New()
	scope := domain.Global()
	first, _ := domain.NewKey("z.last")
	second, _ := domain.NewKey("a.first")
	one, _ := domain.Text("one")
	two, _ := domain.Text("two")
	created, err := store.Replace(scope, first, one, domain.ExpectAbsent())
	if err != nil || created.Entry.Revision != 1 {
		t.Fatal(err)
	}
	if _, err := store.Replace(scope, first, two, domain.ExpectAbsent()); err == nil {
		t.Fatal("blind create accepted")
	}
	expected, _ := domain.ExpectRevision(1)
	if _, err := store.Replace(scope, second, two, expected); err == nil {
		t.Fatal("stale create accepted")
	}
	if _, err := store.Replace(scope, first, one, expected); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Replace(scope, first, two, expected); err != nil {
		t.Fatal(err)
	}
	items, err := store.List(scope)
	if err != nil || len(items) != 1 || items[0].Key != first {
		t.Fatal(items, err)
	}
	removed, err := store.Remove(scope, first, mustRevision(2))
	if err != nil || !removed.Removed {
		t.Fatal(err)
	}
	if _, found, err := store.Read(scope, first); err != nil || found {
		t.Fatal("remove persisted")
	}
}

func TestConcurrentCreateHasOneWinner(t *testing.T) {
	store := memory.New()
	key, _ := domain.NewKey("editor.theme")
	value, _ := domain.Text("dark")
	var wg sync.WaitGroup
	twins := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Replace(domain.Global(), key, value, domain.ExpectAbsent())
			twins <- err
		}()
	}
	wg.Wait()
	close(twins)
	var successes, conflicts int
	for err := range twins {
		if err == nil {
			successes++
		} else {
			conflicts++
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}
}

func mustRevision(n uint64) domain.Expected { v, _ := domain.ExpectRevision(n); return v }
