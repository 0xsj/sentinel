package provenance

import (
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"regexp"
)

func invalid(code string) error {
	return faults.New(faults.Invalid, "invalid provenance").WithType("provenance." + code)
}

var identitySyntax = regexp.MustCompile(`^[A-Za-z0-9._:/-]{1,128}$`)

type ActorKind string

const (
	User    ActorKind = "user"
	Service ActorKind = "service"
	System  ActorKind = "system"
)

type Actor struct {
	kind     ActorKind
	identity string
}

func NewActor(kind ActorKind, identity string) (Actor, error) {
	a := Actor{kind, identity}
	if !a.valid() {
		return Actor{}, invalid("invalid_actor")
	}
	return a, nil
}
func Anonymous() Actor           { return Actor{kind: "anonymous"} }
func (a Actor) Kind() ActorKind  { return a.kind }
func (a Actor) Identity() string { return a.identity }
func (a Actor) valid() bool {
	if a.kind == "anonymous" {
		return a.identity == ""
	}
	return (a.kind == User || a.kind == Service || a.kind == System) && identitySyntax.MatchString(a.identity)
}
func (a Actor) named() bool { return a.valid() && a.kind != "anonymous" }
func clone[T any](v *T) *T {
	if v == nil {
		return nil
	}
	out := *v
	return &out
}
