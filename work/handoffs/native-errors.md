# Native errors revision 2 handoff · 2026-09-11

Contract: `pkg/errors/CONTRACT.md`, revision 2.
State: complete after the combined errors/clock/secret integration gate.
Revision 1 evidence remains in [native-errors.v1.md](native-errors.v1.md).

## Changes

Split responsibilities to match n2f; add condition identifiers, public fields,
private details, annotation/classification and independent public projection.
Go preserves sentinel identity across derivation. Rust retains typed domain cases
through Context<E> and explicit Classified frames. Keep original exact diagnostic
formatting and kind parsing/Go numeric normalization. Go inspection is bounded to
128 frames and stops at aggregates; Rust does not infer through arbitrary sources.
Implementations adapt the local n2f-go/n2f-rs error modules into independent Atelier
files, without sibling imports or dependency changes.

Code edits are confined to the errors module. Coordinator updates contract/task
revision, manifests, READMEs, this handoff, and decision/foundation notes. Existing
native registration stays unchanged. No frontend files or lockfiles edited.

## Scenario mapping

All 9 original tests still pass; original E01–E10 mapping remains in the v1 handoff.
The original E07 exclusion on additional APIs is superseded by revision 2, and
Go's explicit condition-identity rule replaces the earlier derivation behavior.

| Scenario | Evidence |
| --- | --- |
| P01–P02 | TestPublicProjectionAndFallbacks |
| P03 | TestEmptyPublicValues |
| P04–P05, G03 | TestMetadataAndSnapshotOwnership |
| P06, G01 | TestConditionIdentityAcrossDerivation |
| P07 | TestOuterTranslationOwnsWholeFrame |
| P08 | TestReplacingSourcePreservesMetadata |
| P09, G02 | TestPresenceUnknownAndTypedNil; TestAggregateHasNoImplicitSummary; TestBoundedInspection |
| P10, G02 | TestCancellationAndDeadlineThroughAnnotation; original foreign-cause test |

## Checks

- go test -count=1 -v ./pkg/errors: 19 tests pass (9 original + 10 projection/identity tests).
- go test -race -count=1 ./pkg/errors: passes.
- go vet ./pkg/errors; gofmt -l pkg/errors: passes, formatting output empty.
- go vet ./...; go build -o /tmp/atelier-wails-errors-check .: passes.

## Review

Coordinator reviewed private field ownership, immutable Go builders, consuming
Rust builders, exact outer-frame selection, safe fallback, source retention and
absence/unknown distinctions against revision 2. Production imports remain stdlib
only and no native host or domain registry is introduced. Tests were authored with
implementation visibility; no claim of independent or blind verification.

The constructor-message semantic change is explicit: non-Internal messages are
now public-safe input. Existing consumers are only leaf tests; no live application
caller requires migration. Bridge DTO/codec design and native runtime UI behavior
are not exercised or claimed by these tests. Rust's unused warnings are expected
until a real consumer arrives; no warning suppressions added.
