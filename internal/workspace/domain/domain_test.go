package domain_test

import (
	"testing"
	"time"

	"github.com/0xsj/atelier-wails/internal/workspace/domain"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
)

func workspaceID(t *testing.T) id.ID {
	t.Helper()
	v, e := id.Parse("01900000-0000-7000-8000-000000000001")
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func at() time.Time                     { return time.UnixMilli(1234).UTC() }
func revision(n uint64) domain.Expected { v, _ := domain.ExpectRevision(n); return v }

func TestRegisterValidationAndLifecycle(t *testing.T) {
	w, err := domain.Register(workspaceID(t), "Atelier", "/tmp/atelier", at())
	if err != nil || w.Workspace.Revision != 1 || w.Workspace.Status != domain.Active || w.Event.Name != "workspace.registered" {
		t.Fatal("register")
	}
	archived, err := domain.Archive(w.Workspace, revision(1), at())
	if err != nil || archived.Workspace.Status != domain.Archived || archived.Workspace.Revision != 2 {
		t.Fatal("archive")
	}
	if same, err := domain.Archive(archived.Workspace, revision(2), at()); err != nil || same.Status != domain.Unchanged || same.Event != nil {
		t.Fatal("idempotent archive")
	}
	restored, err := domain.Restore(archived.Workspace, revision(2), at())
	if err != nil || restored.Workspace.Status != domain.Active || restored.Workspace.Revision != 3 {
		t.Fatal("restore")
	}
	if _, err := domain.Forget(restored.Workspace, revision(3)); !faults.IsKind(err, faults.Invalid) {
		t.Fatal("active forget")
	}
	forgotten, err := domain.Forget(archived.Workspace, revision(2))
	if err != nil || !forgotten.Removed || forgotten.Revision != 3 || forgotten.Event == nil {
		t.Fatal("forget")
	}
}

func TestWorkspaceValidation(t *testing.T) {
	for _, path := range []string{"", "relative", "C:/atelier\x00"} {
		if _, err := domain.Register(workspaceID(t), "Atelier", path, at()); err == nil {
			t.Fatalf("accepted path %q", path)
		}
	}
	for _, name := range []string{"", " leading", "trailing ", "bad\nname"} {
		if _, err := domain.Register(workspaceID(t), name, "/tmp/atelier", at()); err == nil {
			t.Fatalf("accepted name %q", name)
		}
	}
}
