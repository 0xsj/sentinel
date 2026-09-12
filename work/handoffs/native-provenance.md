# native-provenance handoff

Contract revision 1. Implemented and coordinator-integrated 2026-09-12.

Files: `pkg/provenance/` (actor, attribution, operation, reference, work, scope,
factory, tests and docs), work task/manifest/build order, FOUNDATION and WORKERS.
Rust also registers the module in shared/mod.rs. No dependency or lockfile changes.

## Behavior and evidence

P01–P02: attribution tests and vocabulary boundaries validate construction,
unknown/anonymous distinctions, delegation and owned data.
P03–P09: transition tests exercise Open, Child, pure Prepare, Execute, Retry,
explicit work IDs, changed executor, retained attribution and repeated attempt ordinals.
P10–P12: bounds/effects tests cover checked overflow, clock correction, invalid
time, pre-effect validation and unchanged generator failures/diagnostic causes.
P18–P19: snapshot and restore tests cover ownership, self-causes, inconsistent
attempt/depth/executor combinations and preserved unknown upstream ancestry.
P23: collision tests verify that explicit event causation cannot hide parent scope
reuse; retry collision and failed generation are refused without hidden retries.
Existing-clock/sequence tests integrate the actual preceding packages, including
exhaustion. Time tests cover millisecond truncation and the upper time boundary.

`go test -race -count=1 -v ./pkg/...`: 67 tests passed, including 12 provenance tests.
`go vet ./...`: passed. `gofmt -l pkg/provenance`: empty.

Coordinator source review: Go has private value storage and copies all optional
pointer fields; Rust exposes cloned/owned snapshots. Pure preparation has no
factory/effect dependency. Factories import local errors/IDs and standard runtime
facilities only. Parent values are unchanged by successful or failed transitions.
Replay, link sets and incoming propagation were deliberately omitted from the
local n2f adaptation; no unreachable stubs were copied. No sibling imports.

Checks are local implementation-visible verification, not independent certification.
No desktop packaging, visual launch or cross-platform execution was repeated.
No frontend, host or runtime composition code changed. This module records claims;
it does not establish permission, commit, durable history or global uniqueness.
Go callers serialize factory/dependency access; Rust uses exclusive mutable access.
Logger projection is implemented; the root bootstrap is recorded separately in
`work/handoffs/native-root-bootstrap.md`. Domain runtime consumers remain future work.
