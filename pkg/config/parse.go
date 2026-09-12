package config

import (
	"regexp"
	"strconv"
)

const safeInteger int64 = 9007199254740991

var integerSyntax = regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`)
var keySyntax = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

func integer(s string) (int64, bool) {
	if !integerSyntax.MatchString(s) {
		return 0, false
	}
	v, e := strconv.ParseInt(s, 10, 64)
	return v, e == nil && v >= -safeInteger && v <= safeInteger
}
