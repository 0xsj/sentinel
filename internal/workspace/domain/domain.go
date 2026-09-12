// Package domain contains pure workspace metadata and lifecycle transitions.
package domain

import (
	"math"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
)

type Status uint8

const (
	Active Status = iota + 1
	Archived
)

func (s Status) String() string {
	if s == Active {
		return "active"
	}
	if s == Archived {
		return "archived"
	}
	return "invalid"
}

type ListFilter uint8

const (
	All ListFilter = iota + 1
	ActiveOnly
)

func (f ListFilter) Valid() bool { return f == All || f == ActiveOnly }

type Expected struct {
	exists   bool
	revision uint64
}

func ExpectAbsent() Expected { return Expected{} }
func ExpectRevision(revision uint64) (Expected, error) {
	if revision == 0 {
		return Expected{}, invalid("workspace.invalid_revision")
	}
	return Expected{exists: true, revision: revision}, nil
}
func (e Expected) IsAbsent() bool           { return !e.exists }
func (e Expected) Revision() (uint64, bool) { return e.revision, e.exists }
func (e Expected) Valid() bool              { return !e.exists || e.revision != 0 }

type Workspace struct {
	ID        id.ID
	Name      string
	Location  string
	Status    Status
	Revision  uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(idValue id.ID, name, location string, at time.Time) (Workspace, error) {
	if idValue.IsZero() {
		return Workspace{}, invalid("workspace.invalid_id")
	}
	if !validName(name) {
		return Workspace{}, invalid("workspace.invalid_name")
	}
	if !validLocation(location) {
		return Workspace{}, invalid("workspace.invalid_location")
	}
	at = normalizeTime(at)
	return Workspace{ID: idValue, Name: name, Location: location, Status: Active, Revision: 1, CreatedAt: at, UpdatedAt: at}, nil
}
func (w Workspace) Valid() bool {
	return !w.ID.IsZero() && validName(w.Name) && validLocation(w.Location) && (w.Status == Active || w.Status == Archived) && w.Revision != 0
}

type Event struct {
	Name      string
	Workspace id.ID
	Revision  uint64
	NameValue string
	Location  string
	Status    Status
}
type ChangeStatus uint8

const (
	Changed ChangeStatus = iota + 1
	Unchanged
)

type MutationResult struct {
	Status    ChangeStatus
	Workspace Workspace
	Event     *Event
}
type RegisterResult struct {
	Workspace Workspace
	Event     Event
}
type ForgetResult struct {
	Removed  bool
	Revision uint64
	Event    *Event
}

func Register(idValue id.ID, name, location string, at time.Time) (RegisterResult, error) {
	w, err := New(idValue, name, location, at)
	if err != nil {
		return RegisterResult{}, err
	}
	return RegisterResult{Workspace: w, Event: Event{Name: "workspace.registered", Workspace: idValue, Revision: 1, NameValue: name, Location: location, Status: Active}}, nil
}
func Rename(current Workspace, name string, expected Expected, at time.Time) (MutationResult, error) {
	if err := validateCurrent(current, expected); err != nil {
		return MutationResult{}, err
	}
	if !validName(name) {
		return MutationResult{}, invalid("workspace.invalid_name")
	}
	if current.Name == name {
		return MutationResult{Status: Unchanged, Workspace: current}, nil
	}
	revision, err := nextRevision(current.Revision)
	if err != nil {
		return MutationResult{}, err
	}
	current.Name, current.Revision, current.UpdatedAt = name, revision, normalizeTime(at)
	return MutationResult{Status: Changed, Workspace: current, Event: &Event{Name: "workspace.renamed", Workspace: current.ID, Revision: revision, NameValue: name, Location: current.Location, Status: current.Status}}, nil
}
func Archive(current Workspace, expected Expected, at time.Time) (MutationResult, error) {
	if err := validateCurrent(current, expected); err != nil {
		return MutationResult{}, err
	}
	if current.Status == Archived {
		return MutationResult{Status: Unchanged, Workspace: current}, nil
	}
	revision, err := nextRevision(current.Revision)
	if err != nil {
		return MutationResult{}, err
	}
	current.Status, current.Revision, current.UpdatedAt = Archived, revision, normalizeTime(at)
	return MutationResult{Status: Changed, Workspace: current, Event: &Event{Name: "workspace.archived", Workspace: current.ID, Revision: revision, NameValue: current.Name, Location: current.Location, Status: Archived}}, nil
}
func Restore(current Workspace, expected Expected, at time.Time) (MutationResult, error) {
	if err := validateCurrent(current, expected); err != nil {
		return MutationResult{}, err
	}
	if current.Status == Active {
		return MutationResult{Status: Unchanged, Workspace: current}, nil
	}
	revision, err := nextRevision(current.Revision)
	if err != nil {
		return MutationResult{}, err
	}
	current.Status, current.Revision, current.UpdatedAt = Active, revision, normalizeTime(at)
	return MutationResult{Status: Changed, Workspace: current, Event: &Event{Name: "workspace.restored", Workspace: current.ID, Revision: revision, NameValue: current.Name, Location: current.Location, Status: Active}}, nil
}
func Forget(current Workspace, expected Expected) (ForgetResult, error) {
	if err := validateCurrent(current, expected); err != nil {
		return ForgetResult{}, err
	}
	if current.Status != Archived {
		return ForgetResult{}, invalid("workspace.active_forget")
	}
	revision, err := nextRevision(current.Revision)
	if err != nil {
		return ForgetResult{}, err
	}
	return ForgetResult{Removed: true, Revision: revision, Event: &Event{Name: "workspace.forgotten", Workspace: current.ID, Revision: revision, NameValue: current.Name, Location: current.Location, Status: Archived}}, nil
}

func validateCurrent(current Workspace, expected Expected) error {
	if !current.Valid() {
		return invalid("workspace.invalid_record")
	}
	if !expected.Valid() {
		return invalid("workspace.invalid_revision")
	}
	if !expected.exists || expected.revision != current.Revision {
		return conflict()
	}
	return nil
}
func validName(value string) bool {
	if !utf8.ValidString(value) || value == "" || utf8.RuneCountInString(value) > 120 {
		return false
	}
	runes := []rune(value)
	if unicode.IsSpace(runes[0]) || unicode.IsSpace(runes[len(runes)-1]) {
		return false
	}
	for _, r := range runes {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func validLocation(value string) bool {
	if !utf8.ValidString(value) || value == "" || strings.IndexByte(value, 0) >= 0 || !absolutePath(value) {
		return false
	}
	return true
}
func absolutePath(value string) bool {
	if strings.HasPrefix(value, "/") || strings.HasPrefix(value, `\\`) {
		return true
	}
	return len(value) >= 3 && ((value[0] >= 'a' && value[0] <= 'z') || (value[0] >= 'A' && value[0] <= 'Z')) && value[1] == ':' && (value[2] == '/' || value[2] == '\\')
}
func normalizeTime(value time.Time) time.Time { return value.UTC().Truncate(time.Millisecond) }
func nextRevision(current uint64) (uint64, error) {
	if current == math.MaxUint64 {
		return 0, invalid("workspace.revision_exhausted")
	}
	return current + 1, nil
}
func invalid(typ string) error { return faults.New(faults.Invalid, "invalid workspace").WithType(typ) }
func conflict() error {
	return faults.New(faults.Conflict, "workspace changed").WithType("workspace.conflict")
}
