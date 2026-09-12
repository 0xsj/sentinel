// Package id separates UUID values from identity generation effects.
//
// Parse accepts the canonical hexadecimal form, either case, for the standard
// UUID variant, independently of version. IDs format lowercase and compare by
// value. Nil is a sentinel, not an accepted parsed identity. Time reports only
// UUIDv7's approximate Unix millisecond timestamp.
//
// V7 consumes a caller-supplied Now capability and cryptographic entropy. It
// orders successful calls within one instance using a seeded 12-bit counter,
// holds the timestamp on rollback, and refuses counter exhaustion. An entropy
// failure retains its cause and commits no generator state. Sequence supplies
// finite explicit test fixtures. Neither generator chooses retries or sleeps.
//
// Generators serialize their own calls. Do not copy a used generator, reenter it
// from a dependency, or assume independent instances establish global ordering.
// CONTRACT.md specifies exact bytes, error types, limits and native differences.
package id
