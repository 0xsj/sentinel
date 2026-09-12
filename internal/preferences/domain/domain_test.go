package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/0xsj/atelier-wails/internal/preferences/domain"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
)

func testID(t *testing.T) id.ID {
	t.Helper()
	v, err := id.Parse("01900000-0000-7000-8000-000000000001")
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestValuesAndScopes(t *testing.T) {
	global := domain.Global()
	if !global.Valid() || global.String() != "global" {
		t.Fatal("global scope")
	}
	workspace, err := domain.ForWorkspace(testID(t))
	if err != nil || !workspace.Valid() || workspace.String() == "global" {
		t.Fatal("workspace scope")
	}
	if _, err := domain.ForWorkspace(id.ID{}); err == nil {
		t.Fatal("zero workspace accepted")
	}
	if _, err := domain.NewKey("Bad"); err == nil {
		t.Fatal("invalid key accepted")
	}
	key, err := domain.NewKey("editor.theme")
	if err != nil || !key.Valid() {
		t.Fatal(err)
	}
	if _, err := domain.Text(strings.Repeat("x", 4097)); err == nil {
		t.Fatal("long text accepted")
	}
	empty, err := domain.Text("")
	if err != nil || !empty.Valid() {
		t.Fatal("empty text rejected")
	}
}

func TestCompareReplaceAndRemove(t *testing.T) {
	scope := domain.Global()
	key, _ := domain.NewKey("editor.theme")
	value, _ := domain.Text("dark")
	create, err := domain.DecideReplace(nil, scope, key, value, domain.ExpectAbsent())
	if err != nil || create.Status != domain.Changed || create.Entry.Revision != 1 || create.Event == nil {
		t.Fatal("create")
	}
	if create.Event.Name != "preference.changed" {
		t.Fatal("event")
	}
	stale, _ := domain.ExpectRevision(9)
	if _, err := domain.DecideReplace(&create.Entry, scope, key, value, stale); !faults.IsKind(err, faults.Conflict) {
		t.Fatal("stale replace")
	}
	expected, _ := domain.ExpectRevision(1)
	equal, err := domain.DecideReplace(&create.Entry, scope, key, value, expected)
	if err != nil || equal.Status != domain.Unchanged || equal.Event != nil {
		t.Fatal("equal replace")
	}
	nextValue, _ := domain.Text("light")
	changed, err := domain.DecideReplace(&create.Entry, scope, key, nextValue, expected)
	if err != nil || changed.Entry.Revision != 2 || changed.Event == nil {
		t.Fatal("changed replace")
	}
	removed, err := domain.DecideRemove(&changed.Entry, scope, key, mustRevision(2))
	if err != nil || !removed.Removed || removed.Revision != 3 || removed.Event == nil {
		t.Fatal("remove")
	}
	missing, err := domain.DecideRemove(nil, scope, key, domain.ExpectAbsent())
	if err != nil || missing.Removed {
		t.Fatal("missing remove")
	}
}

func mustRevision(n uint64) domain.Expected { v, _ := domain.ExpectRevision(n); return v }

func TestFailureProjectionDoesNotEchoValues(t *testing.T) {
	_, err := domain.NewKey("private key")
	if !errors.Is(err, err) || faults.DiagnosticTypeOf(err) != "preferences.invalid_key" || strings.Contains(err.Error(), "private") {
		t.Fatal(err)
	}
}
