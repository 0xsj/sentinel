package secret_test

import (
	"fmt"
	"github.com/0xsj/atelier-wails/pkg/secret"
	"testing"
)

func TestRevealPreservesInput(t *testing.T) {
	for _, input := range []string{"", "dummy-secret", "Δ\nvalue", "  whitespace\x00  "} {
		if secret.New(input).Reveal() != input {
			t.Fatal("reveal changed input")
		}
	}
}
func TestSupportedFormatsAreRedacted(t *testing.T) {
	for _, input := range []string{"", "dummy-secret", "Δ\nvalue", "  whitespace\x00  "} {
		s := secret.New(input)
		outputs := []string{fmt.Sprint(s), fmt.Sprintf("%s", s), fmt.Sprintf("%v", s), fmt.Sprintf("%+v", s), fmt.Sprintf("%#v", s), s.String(), s.GoString()}
		for _, got := range outputs {
			if got != "[REDACTED]" {
				t.Fatal("supported format was not redacted")
			}
		}
	}
}
func TestIndependentSecrets(t *testing.T) {
	a, b := secret.New("first dummy"), secret.New("second dummy")
	if a.Reveal() != "first dummy" || b.Reveal() != "second dummy" || fmt.Sprint(a) != fmt.Sprint(b) {
		t.Fatal("independence or redaction mismatch")
	}
}
func TestZeroValue(t *testing.T) {
	var s secret.Secret
	if s.Reveal() != "" || fmt.Sprint(s) != "[REDACTED]" || fmt.Sprintf("%#v", s) != "[REDACTED]" {
		t.Fatal("zero secret mismatch")
	}
}
