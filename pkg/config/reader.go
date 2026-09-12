package config

import (
	"cmp"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"github.com/0xsj/atelier-wails/pkg/secret"
	"slices"
	"strconv"
)

type Var struct {
	Key, Value, Source string
	Secret             bool
}
type Reader struct {
	lookup   Lookup
	seen     map[string]bool
	problems map[string]string
	resolved []Var
}

func New(lookup Lookup) *Reader {
	r := &Reader{lookup: lookup, seen: map[string]bool{}, problems: map[string]string{}}
	if lookup == nil {
		r.problems["<source>"] = "invalid_source"
		r.lookup = Map(nil)
	}
	return r
}
func (r *Reader) problem(k, why string) {
	if !keySyntax.MatchString(k) {
		k = "<key>"
	}
	r.problems[k] = why
}
func (r *Reader) read(k string) (string, bool, bool) {
	if !keySyntax.MatchString(k) {
		r.problem(k, "invalid_key")
		return "", false, false
	}
	if r.seen[k] {
		r.problem(k, "duplicate_key")
		return "", false, false
	}
	r.seen[k] = true
	v, p := r.lookup(k)
	return v, p, true
}
func (r *Reader) record(k, v string, p, s bool) {
	source := "default"
	if p {
		source = "environment"
	}
	if s {
		v = "[REDACTED]"
	}
	r.resolved = append(r.resolved, Var{k, v, source, s})
}
func (r *Reader) String(k, fallback string) string {
	v, p, ok := r.read(k)
	if !ok {
		return ""
	}
	if !p {
		v = fallback
	}
	r.record(k, v, p, false)
	return v
}
func (r *Reader) Required(k string) string {
	v, p, ok := r.read(k)
	if !ok {
		return ""
	}
	if !p || v == "" {
		r.problem(k, "required")
		return ""
	}
	r.record(k, v, p, false)
	return v
}
func (r *Reader) Secret(k string) secret.Secret {
	v, p, ok := r.read(k)
	if !ok {
		return secret.New("")
	}
	if !p || v == "" {
		r.problem(k, "required")
		return secret.New("")
	}
	r.record(k, v, p, true)
	return secret.New(v)
}
func (r *Reader) Int(k string, fallback, min, max int64) int64 {
	if min > max || min < -safeInteger || max > safeInteger || fallback < min || fallback > max {
		r.problem(k, "invalid_definition")
		return 0
	}
	s, p, ok := r.read(k)
	if !ok {
		return 0
	}
	v := fallback
	if p {
		var valid bool
		v, valid = integer(s)
		if !valid {
			r.problem(k, "invalid_integer")
			return 0
		}
	}
	if v < min || v > max {
		r.problem(k, "out_of_range")
		return 0
	}
	r.record(k, strconv.FormatInt(v, 10), p, false)
	return v
}
func (r *Reader) Bool(k string, fallback bool) bool {
	s, p, ok := r.read(k)
	if !ok {
		return false
	}
	v := fallback
	if p {
		switch s {
		case "true":
			v = true
		case "false":
			v = false
		default:
			r.problem(k, "invalid_boolean")
			return false
		}
	}
	r.record(k, strconv.FormatBool(v), p, false)
	return v
}
func (r *Reader) Enum(k, fallback string, allowed []string) string {
	seen := map[string]bool{}
	for _, v := range allowed {
		if v == "" || seen[v] {
			r.problem(k, "invalid_definition")
			return ""
		}
		seen[v] = true
	}
	if !seen[fallback] {
		r.problem(k, "invalid_definition")
		return ""
	}
	v, p, ok := r.read(k)
	if !ok {
		return ""
	}
	if !p {
		v = fallback
	}
	if !seen[v] {
		r.problem(k, "invalid_choice")
		return ""
	}
	r.record(k, v, p, false)
	return v
}
func (r *Reader) Err() error {
	if len(r.problems) == 0 {
		return nil
	}
	return faults.New(faults.Invalid, "invalid configuration").WithType("config.invalid").WithFields(r.problems)
}
func (r *Reader) Manifest() ([]Var, error) {
	if e := r.Err(); e != nil {
		return nil, e
	}
	v := slices.Clone(r.resolved)
	slices.SortFunc(v, func(a, b Var) int { return cmp.Compare(a.Key, b.Key) })
	return v, nil
}
