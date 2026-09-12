# native-id handoff

Contract revision: 1. Implemented and coordinator-integrated on 2026-09-12.

Local adaptation of n2f's value, UUIDv7 state machine and finite sequence;
no sibling runtime imports. Files: `pkg/id/`, task and manifest, foundation and
worker documentation; Rust additionally shared/mod.rs and Cargo.toml/Cargo.lock.
Clock mutability exists only in private tests; the existing clock API is unchanged.

## Scenario coverage

I01–I09 map to correspondingly numbered tests in spec_test.go:
strict parsing and redaction; equality/canonical formatting/timestamp;
exact byte fixtures; ordering and rollback; exhaustion and recovery; original
entropy causes and state atomicity; time bounds; owned finite sequences;
production OS entropy and concurrent/shared generator access.
Two additional tests verify existing Fixed clock compatibility, intentional
fixture duplicates, and exactly 2049 successes from the highest fresh seed.
Go additionally verifies text decoding leaves the receiver unchanged on failure,
zero-value rejection, nil dependencies and short entropy reads.

## Verification

`go test -race -count=1 -v ./pkg/...`: 55 tests passed, including 11 ID tests.
`go vet ./...`: passed. `gofmt -l pkg/id`: empty.
The sandbox denied Go cache access on the final pass; the test and vet commands
were rerun with approved cache access and passed.

Coordinator source review: values have private storage; production imports are
limited to standard runtime facilities, local errors and Rust getrandom. Only
successful generation commits state. Invalid time and exhausted counters return
before entropy reads; no sleep, retry, insecure fallback or future tick invention.
Both languages use identical seeded-counter and suffix fixtures. This is local
implementation-visible verification, not independent certification.

No frontend or host code changed. Desktop packaging and visual launch were not
rerun for this leaf. No cross-platform execution was performed. Per-generator
ordering is not a global uniqueness or event-time guarantee. Root wiring,
provenance, domain wrappers, database codecs and bridge schemas remain future work.
