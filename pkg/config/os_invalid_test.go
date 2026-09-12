//go:build !windows

package config_test

import (
	"fmt"
	"github.com/0xsj/atelier-wails/pkg/config"
	faults "github.com/0xsj/atelier-wails/pkg/errors"
	"strings"
	"testing"
)

func TestInvalidNativeSourceIsRedacted(t *testing.T) {
	t.Setenv("ATELIER_INVALID_SOURCE_FIXTURE", "dummy-private-\xff")
	lookup, e := config.OS()
	if e == nil || lookup != nil {
		t.Fatal("invalid UTF-8 source accepted")
	}
	view, _ := faults.Public(e)
	if view.Kind != faults.Invalid || view.Type != "config.source" || strings.Contains(fmt.Sprint(view, e), "dummy-private") || len(view.Fields) != 0 {
		t.Fatal("source error leaked")
	}
}
