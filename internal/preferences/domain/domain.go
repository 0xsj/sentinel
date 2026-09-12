// Package domain contains pure preference values and compare-and-replace rules.
package domain

import (
	"math"
	"unicode/utf8"

	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/id"
)

type ScopeKind uint8

const (
	GlobalScope ScopeKind = iota + 1
	WorkspaceScope
)

type Scope struct {
	kind      ScopeKind
	workspace id.ID
}

func Global() Scope { return Scope{kind: GlobalScope} }

func ForWorkspace(workspace id.ID) (Scope, error) {
	if workspace.IsZero() {
		return Scope{}, invalid("preferences.invalid_scope")
	}
	return Scope{kind: WorkspaceScope, workspace: workspace}, nil
}

func (s Scope) Kind() ScopeKind { return s.kind }
func (s Scope) WorkspaceID() (id.ID, bool) {
	return s.workspace, s.kind == WorkspaceScope && !s.workspace.IsZero()
}
func (s Scope) Valid() bool {
	return s.kind == GlobalScope || (s.kind == WorkspaceScope && !s.workspace.IsZero())
}
func (s Scope) String() string {
	if s.kind == GlobalScope {
		return "global"
	}
	if s.kind == WorkspaceScope && !s.workspace.IsZero() {
		return "workspace:" + s.workspace.String()
	}
	return "invalid"
}

type Key string

func NewKey(value string) (Key, error) {
	if len(value) == 0 || len(value) > 128 || value[0] < 'a' || value[0] > 'z' {
		return "", invalid("preferences.invalid_key")
	}
	for i := 1; i < len(value); i++ {
		b := value[i]
		if !((b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '.' || b == '_' || b == '-') {
			return "", invalid("preferences.invalid_key")
		}
	}
	return Key(value), nil
}
func (k Key) String() string { return string(k) }
func (k Key) Valid() bool {
	_, err := NewKey(string(k))
	return err == nil
}

type ValueKind uint8

const (
	TextValue ValueKind = iota + 1
	BoolValue
	IntValue
)

type Value struct {
	kind ValueKind
	text string
	flag bool
	intv int64
}

func Text(value string) (Value, error) {
	if !utf8.ValidString(value) || len(value) > 4096 {
		return Value{}, invalid("preferences.invalid_value")
	}
	return Value{kind: TextValue, text: value}, nil
}
func Bool(value bool) Value          { return Value{kind: BoolValue, flag: value} }
func Int(value int64) Value          { return Value{kind: IntValue, intv: value} }
func (v Value) Kind() ValueKind      { return v.kind }
func (v Value) Text() (string, bool) { return v.text, v.kind == TextValue }
func (v Value) Bool() (bool, bool)   { return v.flag, v.kind == BoolValue }
func (v Value) Int() (int64, bool)   { return v.intv, v.kind == IntValue }
func (v Value) Valid() bool {
	switch v.kind {
	case TextValue:
		return utf8.ValidString(v.text) && len(v.text) <= 4096
	case BoolValue, IntValue:
		return true
	default:
		return false
	}
}
func (v Value) Equal(other Value) bool {
	return v.kind == other.kind && v.text == other.text && v.flag == other.flag && v.intv == other.intv
}

type Expected struct {
	exists   bool
	revision uint64
}

func ExpectAbsent() Expected { return Expected{} }
func ExpectRevision(revision uint64) (Expected, error) {
	if revision == 0 {
		return Expected{}, invalid("preferences.invalid_revision")
	}
	return Expected{exists: true, revision: revision}, nil
}
func (e Expected) IsAbsent() bool           { return !e.exists }
func (e Expected) Revision() (uint64, bool) { return e.revision, e.exists }
func (e Expected) Valid() bool              { return !e.exists || e.revision != 0 }

type Entry struct {
	Scope    Scope
	Key      Key
	Value    Value
	Revision uint64
}

func (e Entry) Valid() bool {
	return e.Scope.Valid() && e.Key.Valid() && e.Value.Valid() && e.Revision != 0
}

type Event struct {
	Name     string
	Scope    Scope
	Key      Key
	Value    Value
	HasValue bool
	Revision uint64
}

type ChangeStatus uint8

const (
	Changed ChangeStatus = iota + 1
	Unchanged
)

type ReplaceResult struct {
	Status ChangeStatus
	Entry  Entry
	Event  *Event
}

type RemoveResult struct {
	Removed  bool
	Revision uint64
	Event    *Event
}

func DecideReplace(current *Entry, scope Scope, key Key, value Value, expected Expected) (ReplaceResult, error) {
	if !scope.Valid() {
		return ReplaceResult{}, invalid("preferences.invalid_scope")
	}
	if !key.Valid() {
		return ReplaceResult{}, invalid("preferences.invalid_key")
	}
	if !value.Valid() {
		return ReplaceResult{}, invalid("preferences.invalid_value")
	}
	if !expected.Valid() {
		return ReplaceResult{}, invalid("preferences.invalid_revision")
	}
	if current == nil {
		if expected.exists {
			return ReplaceResult{}, conflict()
		}
		entry := Entry{Scope: scope, Key: key, Value: value, Revision: 1}
		return ReplaceResult{Status: Changed, Entry: entry, Event: &Event{Name: "preference.changed", Scope: scope, Key: key, Value: value, HasValue: true, Revision: 1}}, nil
	}
	if !current.Valid() || current.Scope != scope || current.Key != key {
		return ReplaceResult{}, invalid("preferences.invalid_entry")
	}
	if !expected.exists || expected.revision != current.Revision {
		return ReplaceResult{}, conflict()
	}
	if current.Value.Equal(value) {
		return ReplaceResult{Status: Unchanged, Entry: *current}, nil
	}
	revision, err := nextRevision(current.Revision)
	if err != nil {
		return ReplaceResult{}, err
	}
	entry := Entry{Scope: scope, Key: key, Value: value, Revision: revision}
	return ReplaceResult{Status: Changed, Entry: entry, Event: &Event{Name: "preference.changed", Scope: scope, Key: key, Value: value, HasValue: true, Revision: revision}}, nil
}

func DecideRemove(current *Entry, scope Scope, key Key, expected Expected) (RemoveResult, error) {
	if !scope.Valid() {
		return RemoveResult{}, invalid("preferences.invalid_scope")
	}
	if !key.Valid() {
		return RemoveResult{}, invalid("preferences.invalid_key")
	}
	if !expected.Valid() {
		return RemoveResult{}, invalid("preferences.invalid_revision")
	}
	if current == nil {
		return RemoveResult{}, nil
	}
	if !current.Valid() || current.Scope != scope || current.Key != key {
		return RemoveResult{}, invalid("preferences.invalid_entry")
	}
	if !expected.exists || expected.revision != current.Revision {
		return RemoveResult{}, conflict()
	}
	revision, err := nextRevision(current.Revision)
	if err != nil {
		return RemoveResult{}, err
	}
	return RemoveResult{Removed: true, Revision: revision, Event: &Event{Name: "preference.removed", Scope: scope, Key: key, Revision: revision}}, nil
}

func nextRevision(current uint64) (uint64, error) {
	if current == math.MaxUint64 {
		return 0, invalid("preferences.revision_exhausted")
	}
	return current + 1, nil
}

func invalid(typ string) error { return faults.New(faults.Invalid, "invalid preference").WithType(typ) }
func conflict() error {
	return faults.New(faults.Conflict, "preference changed").WithType("preferences.conflict")
}
