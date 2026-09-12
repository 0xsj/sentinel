package config_test

import (
	"github.com/0xsj/atelier-wails/pkg/config"
	"testing"
)

func TestV01OSIsCapturedOnce(t *testing.T) {
	t.Setenv("ATELIER_CONFIG_SNAPSHOT_FIXTURE", "before")
	lookup, e := config.OS()
	if e != nil {
		t.Fatal(e)
	}
	t.Setenv("ATELIER_CONFIG_SNAPSHOT_FIXTURE", "after")
	if value, present := lookup("ATELIER_CONFIG_SNAPSHOT_FIXTURE"); !present || value != "before" {
		t.Fatal(value, present)
	}
}
