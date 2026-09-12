# native-logger handoff

Contract revision 1. Implemented and coordinator-integrated 2026-09-12.

Files: `pkg/logger/` production code, tests and documentation; root logging/entry
and root tests; task/manifest/build guides. Rust additionally shared/mod.rs.
No frontend, dependency or lockfile changes in this slice.

## Coverage

- L01–L02: configuration, filtering/time, time-bounds/sink-floor tests.
- L03–L04: binding/ownership tests cover scalar values, secret redaction,
  parent isolation, replacement, owned capture and protected envelope service.
- L05: safe-error tests cover classified, internal, unknown, causes and clearing.
- L06: provenance and optional-field tests cover identity, separate timestamps,
  delegation, tenant, prior attempt, cause and absence of unknown origin/depth.
- L07: console tests verify escaping and single physical lines; Go also exercises
  replacement with the actual slog.JSONHandler. Time tests cover normalization
  and rejected timestamps.
- L08–L09: delivery failure and concurrent tests cover observable original causes,
  flush failure and 200 concurrent records. Source review confirms synchronized
  write/flush on a shared sink and no queue, hidden retry or close semantics.
- L10: native-root tests verify host invocation, lifecycle records, preserved Go
  host error identity and continued host operation when sink/flush fail.

## Verification

- `go test -race -count=1 ./pkg/... ./root`: passed, 79 tests (10 logger, 2 root).
- `go vet ./...`: passed.
- `gofmt -l pkg/logger root`: empty.
- `go build -o /tmp/atelier-wails-logger-check .`: passed after resume.

The final test/vet pass initially encountered sandbox-denied Go cache access;
reruns with approved cache access passed. No code-test failure was hidden.

Coordinator source review: console adapters own formatting; Go slog types occur
only at the adapter seam. Fields are restricted to owned scalars, with no arbitrary
object display fallback. Error projection calls shared public disclosure policy.
Provenance projection reads validated snapshots without generating IDs or time.
Go related loggers serialize shared clock/writes and slog adapter write/flush;
Rust uses shared state and writer mutexes. Memory returns copies in both languages.
Root binds process lifecycle events to one generated bootstrap provenance scope;
it does not fabricate user, workspace or durable-operation provenance.

Verification is local and implementation-visible, not independent certification.
Native GUI launch, packaging and cross-platform release checks were not repeated.
Go uses native slog text, Rust uses quoted key/value console text; record meaning
is mirrored, exact display bytes are not. Custom sinks own their encoding policy.
Synchronous logging can block on a writer. It supplies diagnostics, not durable
history. Caller/root decides how delivery failures affect diagnostics separately
from an operation's outcome.
