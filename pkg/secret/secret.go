// Package secret makes revealing string contents explicit and redacts supported
// diagnostic formatting. It does not encrypt or wipe memory.
package secret

type Secret struct{ value string }

func New(value string) Secret   { return Secret{value: value} }
func (s Secret) Reveal() string { return s.value }
func (Secret) String() string   { return "[REDACTED]" }
func (Secret) GoString() string { return "[REDACTED]" }
