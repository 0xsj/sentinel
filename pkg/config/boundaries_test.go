package config_test

import (
	"fmt"
	"github.com/0xsj/atelier-wails/pkg/config"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"reflect"
	"strings"
	"testing"
)

func TestRequiredAndSecretPresence(t *testing.T) {
	for _, present := range []bool{false, true} {
		values := map[string]string{}
		if present {
			values["TEXT"] = ""
			values["TOKEN"] = ""
		}
		r := config.New(config.Map(values))
		r.Required("TEXT")
		r.Secret("TOKEN")
		if !reflect.DeepEqual(faults.FieldsOf(r.Err()), map[string]string{"TEXT": "required", "TOKEN": "required"}) {
			t.Fatal("missing/empty accepted")
		}
	}
	r := config.New(config.Map(map[string]string{"TEXT": " ", "TOKEN": " \n"}))
	if r.Required("TEXT") != " " || r.Secret("TOKEN").Reveal() != " \n" || r.Err() != nil {
		t.Fatal("whitespace trimmed")
	}
}
func TestIntegerBoundariesAndCanonicalization(t *testing.T) {
	for _, raw := range []string{"-9007199254740991", "9007199254740991", "0", "-0"} {
		r := config.New(config.Map(map[string]string{"N": raw}))
		r.Int("N", 0, -9007199254740991, 9007199254740991)
		v, e := r.Manifest()
		if e != nil {
			t.Fatal("valid boundary refused")
		}
		want := raw
		if raw == "-0" {
			want = "0"
		}
		if v[0].Value != want {
			t.Fatal("noncanonical integer")
		}
	}
	for _, raw := range []string{"-9007199254740992", "999999999999999999999999", "-01", "--1", "١", "1\n"} {
		r := config.New(config.Map(map[string]string{"N": raw}))
		r.Int("N", 0, -9007199254740991, 9007199254740991)
		if faults.FieldsOf(r.Err())["N"] != "invalid_integer" {
			t.Fatal("invalid integer accepted/misclassified")
		}
	}
	r := config.New(config.Map(map[string]string{"N": "11"}))
	r.Int("N", 5, 0, 10)
	if faults.FieldsOf(r.Err())["N"] != "out_of_range" {
		t.Fatal("bounds not enforced")
	}
}
func TestDefinitionAndLookupPrecedence(t *testing.T) {
	calls := 0
	r := config.New(func(string) (string, bool) { calls++; return "ok", true })
	r.Int("BOUNDS", 0, 1, -1)
	r.Int("SAFE", 0, -9007199254740992, 10)
	r.Enum("EMPTY", "ok", nil)
	r.Enum("DUP", "ok", []string{"ok", "ok"})
	r.Enum("BLANK", "ok", []string{"ok", ""})
	r.String("private-key-value", "fallback")
	if calls != 0 {
		t.Fatal("invalid definition/key accessed source")
	}
	if len(faults.FieldsOf(r.Err())) != 6 || faults.FieldsOf(r.Err())["<key>"] != "invalid_key" {
		t.Fatal("problem classification missing")
	}
	r.String("ONCE", "fallback")
	r.String("ONCE", "fallback")
	if calls != 1 || faults.FieldsOf(r.Err())["ONCE"] != "duplicate_key" {
		t.Fatal("duplicate accessed source")
	}
	r.Int("BOUNDS", 0, 0, 10)
	if faults.FieldsOf(r.Err())["BOUNDS"] == "" {
		t.Fatal("earlier error cleared")
	}
}
func TestBooleanAndEnumExactness(t *testing.T) {
	for _, raw := range []string{"", "True", "FALSE", "1", "0", " true", "false\n"} {
		r := config.New(config.Map(map[string]string{"B": raw}))
		r.Bool("B", true)
		if faults.FieldsOf(r.Err())["B"] != "invalid_boolean" {
			t.Fatal("boolean accepted")
		}
	}
	r := config.New(config.Map(nil))
	if !r.Bool("B", true) || r.Enum("E", "one", []string{"one", "two"}) != "one" {
		t.Fatal("defaults lost")
	}
	manifest, e := r.Manifest()
	if e != nil || manifest[0].Value != "true" || manifest[0].Source != "default" {
		t.Fatal("bad bool default manifest")
	}
}
func TestProblemsAreRedactedAndOwned(t *testing.T) {
	r := config.New(config.Map(map[string]string{"N": "dummy-private-value", "E": "dummy-private-value"}))
	r.Int("N", 0, 0, 1)
	r.Enum("E", "yes", []string{"yes"})
	r.Required("dummy-private-key")
	e := r.Err()
	view, _ := faults.Public(e)
	text := fmt.Sprintf("%v %#v %#v", e, view, faults.DetailsOf(e))
	if strings.Contains(text, "dummy-private") || view.Type != "config.invalid" || view.Message != "invalid configuration" || view.Kind != faults.Invalid {
		t.Fatal("configuration disclosure or classification mismatch")
	}
	view.Fields["N"] = "changed"
	if faults.FieldsOf(r.Err())["N"] != "invalid_integer" {
		t.Fatal("error output aliases problems")
	}
	if v, e := r.Manifest(); e == nil || v != nil {
		t.Fatal("invalid reader produced manifest")
	}
}
func TestEmptyAndNilSources(t *testing.T) {
	r := config.New(config.Map(nil))
	v, e := r.Manifest()
	if e != nil || len(v) != 0 {
		t.Fatal("empty valid manifest failed")
	}
	r = config.New(nil)
	r.String("A", "fallback")
	if faults.FieldsOf(r.Err())["<source>"] != "invalid_source" {
		t.Fatal("nil source accepted")
	}
	if _, e := r.Manifest(); e == nil {
		t.Fatal("nil source manifest succeeded")
	}
}
func TestManifestOnlyContainsDeclaredSettings(t *testing.T) {
	r := config.New(config.Map(map[string]string{"UNUSED": "dummy-hidden", "TOKEN": "dummy-secret", "B": "true"}))
	s := r.Secret("TOKEN")
	r.Bool("B", false)
	r.String("A", "visible")
	v, e := r.Manifest()
	if e != nil || len(v) != 3 || v[0].Key != "A" || v[1].Key != "B" || v[2].Key != "TOKEN" || !v[2].Secret || v[2].Source != "environment" {
		t.Fatal("manifest shape")
	}
	if strings.Contains(fmt.Sprint(v), "dummy-") || fmt.Sprint(s) != "[REDACTED]" {
		t.Fatal("manifest or secret disclosed")
	}
}
