package config_test

import (
	"github.com/0xsj/atelier-wails/pkg/config"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"strings"
	"testing"
)

func TestV01Presence(t *testing.T) {
	r := config.New(config.Map(map[string]string{"EMPTY": "", "SET": "yes"}))
	if r.String("ABSENT", "default") != "default" || r.String("EMPTY", "default") != "" || r.Required("SET") != "yes" {
		t.Fatal("presence")
	}
	if r.Err() != nil {
		t.Fatal(r.Err())
	}
	r.Secret("MISSING")
	if !faults.IsKind(r.Err(), faults.Invalid) {
		t.Fatal("secret required")
	}
}
func TestV02Parsers(t *testing.T) {
	for _, s := range []string{"", " 2", "2 ", "+2", "02", "1e2", "1.5", "9007199254740992"} {
		r := config.New(config.Map(map[string]string{"N": s}))
		r.Int("N", 2, -9007199254740991, 9007199254740991)
		if r.Err() == nil {
			t.Fatal(s)
		}
	}
	r := config.New(config.Map(map[string]string{"N": "-3", "B": "false", "E": "json"}))
	if r.Int("N", 2, -3, 3) != -3 || r.Bool("B", true) || r.Enum("E", "console", []string{"console", "json"}) != "json" || r.Err() != nil {
		t.Fatal("valid parsers")
	}
	bad := config.New(config.Map(map[string]string{"B": "TRUE", "E": " json"}))
	bad.Bool("B", false)
	bad.Enum("E", "console", []string{"console", "json"})
	if len(faults.FieldsOf(bad.Err())) != 2 {
		t.Fatal(bad.Err())
	}
}
func TestV03Definitions(t *testing.T) {
	calls := 0
	r := config.New(func(string) (string, bool) { calls++; return "2", true })
	r.Int("N", 20, 0, 10)
	r.Enum("E", "bad", []string{"ok"})
	if calls != 0 || len(faults.FieldsOf(r.Err())) != 2 {
		t.Fatal(calls, r.Err())
	}
}
func TestV04Failures(t *testing.T) {
	r := config.New(config.Map(map[string]string{"N": "private-SENTINEL", "EMPTY": ""}))
	r.Int("N", 1, 0, 10)
	r.Required("EMPTY")
	if len(faults.FieldsOf(r.Err())) != 2 || strings.Contains(r.Err().Error(), "private-SENTINEL") {
		t.Fatal(r.Err())
	}
	if _, e := r.Manifest(); e == nil {
		t.Fatal("invalid manifest")
	}
}
func TestV05Manifest(t *testing.T) {
	m := map[string]string{"TOKEN": "private-SENTINEL", "NAME": "visible"}
	r := config.New(config.Map(m))
	m["TOKEN"] = "changed"
	if r.Secret("TOKEN").Reveal() != "private-SENTINEL" {
		t.Fatal("input alias")
	}
	r.Required("NAME")
	r.Int("COUNT", 3, 0, 10)
	v, e := r.Manifest()
	if e != nil || len(v) != 3 || v[0].Key != "COUNT" || v[0].Source != "default" || v[1].Value != "visible" || v[2].Value != "[REDACTED]" {
		t.Fatal(v, e)
	}
	v[2].Value = "oops"
	again, _ := r.Manifest()
	if again[2].Value != "[REDACTED]" {
		t.Fatal("output alias")
	}
}
func TestV06Duplicate(t *testing.T) {
	r := config.New(config.Map(nil))
	r.String("A", "a")
	r.String("A", "b")
	if faults.FieldsOf(r.Err())["A"] != "duplicate_key" {
		t.Fatal(r.Err())
	}
}
