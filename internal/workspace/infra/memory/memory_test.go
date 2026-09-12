package memory_test

import (
	"testing"
	"time"

	"github.com/0xsj/atelier-wails/internal/workspace/domain"
	"github.com/0xsj/atelier-wails/internal/workspace/infra/memory"
	"github.com/0xsj/atelier-wails/pkg/id"
)

func idFor(t *testing.T, text string) id.ID {
	t.Helper()
	value, err := id.Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func TestRegistryUniquenessLifecycleAndOrdering(t *testing.T) {
	store := memory.New()
	when := time.UnixMilli(100)
	a := idFor(t, "01900000-0000-7000-8000-000000000001")
	b := idFor(t, "01900000-0000-7000-8000-000000000002")
	if _, err := store.Register(a, "A", "/tmp/a", when); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Register(b, "B", "/tmp/a", when); err == nil {
		t.Fatal("active location collision accepted")
	}
	if _, err := store.Register(b, "B", "/tmp/b", when); err != nil {
		t.Fatal(err)
	}
	items, err := store.List(domain.All)
	if err != nil || len(items) != 2 || items[0].ID != a {
		t.Fatal(items, err)
	}
	first, _, _ := store.Read(a)
	archived, err := store.Archive(a, mustRevision(1), when)
	if err != nil || archived.Workspace.Status != domain.Archived {
		t.Fatal(err)
	}
	if _, err := store.Restore(b, mustRevision(1), when); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Forget(a, mustRevision(2)); err != nil {
		t.Fatal(err)
	}
	if _, found, err := store.Read(a); err != nil || found {
		t.Fatal("forgotten workspace persisted")
	}
	if first.Status != domain.Active {
		t.Fatal("read snapshot changed")
	}
}

func mustRevision(n uint64) domain.Expected { v, _ := domain.ExpectRevision(n); return v }
