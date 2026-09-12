package provenance

import "regexp"

var operationSyntax = regexp.MustCompile(`^[a-z][a-z0-9_.:-]{0,127}$`)

type Operation struct{ name string }

func NewOperation(name string) (Operation, error) {
	if !operationSyntax.MatchString(name) {
		return Operation{}, invalid("invalid_operation")
	}
	return Operation{name}, nil
}
func (o Operation) String() string { return o.name }
