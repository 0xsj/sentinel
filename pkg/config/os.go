package config

import (
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"os"
	"strings"
	"unicode/utf8"
)

func OS() (Lookup, error) {
	values := map[string]string{}
	for _, entry := range os.Environ() {
		k, v, ok := strings.Cut(entry, "=")
		if !ok || !utf8.ValidString(k) || !utf8.ValidString(v) {
			return nil, faults.New(faults.Invalid, "invalid environment source").WithType("config.source")
		}
		values[k] = v
	}
	return Map(values), nil
}
