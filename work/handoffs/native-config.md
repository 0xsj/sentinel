# Native config handoff · 2026-09-11

Contract: `pkg/config/CONTRACT.md`, revision 1. State: complete after native integration
and coordinator source review. Dependencies: errors revision 2 and secret revision 1.

## Changes and scope

Added captured map/OS sources, exact typed readers, aggregated safe problems and
owned redacted manifests. Adapted n2f's environment module locally under config;
all imports target standard library or Atelier's own errors/secret modules.
Rust additionally registers shared::config. No new dependencies or lockfile edits.
No frontend source edits, startup behavior changes or invented application settings.
Coordinator owns contract/task/manifest/status updates.

## Evidence

| Scenario | Tests / review |
| --- | --- |
| V01 | TestV01Presence; TestRequiredAndSecretPresence |
| V02 | TestV02Parsers; TestBooleanAndEnumExactness |
| V03 | TestV03Definitions; TestDefinitionAndLookupPrecedence |
| V04 | TestV04Failures; TestProblemsAreRedactedAndOwned |
| V05 | TestV05Manifest; TestManifestOnlyContainsDeclaredSettings |
| V06 | TestV06Duplicate; TestV01OSIsCapturedOnce |
| V07 | TestDefinitionAndLookupPrecedence |
| V08 | TestInvalidNativeSourceIsRedacted; TestV01OSIsCapturedOnce |
| V09 | TestIntegerBoundariesAndCanonicalization |
| V10 | TestEmptyAndNilSources; TestProblemsAreRedactedAndOwned |

- go test -count=1 -v ./pkg/config: 15 tests pass.
- go test -race -count=1 ./pkg/config: passes.
- go vet ./pkg/config and go vet ./...: pass.
- gofmt -l pkg/config: empty output.
- go test -count=1 ./pkg/errors ./pkg/clock ./pkg/secret ./pkg/config: all four pass (44 tests across focused counts).

## Review and limitations

Reviewed lookup/definition ordering, absent versus empty behavior, exact parser
bounds, canonical output, snapshot ownership, copied failure metadata and redaction.
Native source capture errors contain neither raw values nor raw keys. No arbitrary
callback panic recovery is promised. V08 invalid-byte tests were run on macOS;
Go's invalid native text fixture is excluded on Windows and Rust's is Unix-only.
No claim of cross-platform invalid-text evidence. Reader callbacks are trusted and
readers are not concurrent mutable settings stores.

These tests have implementation visibility, including adapted n2f tests. They are
not an independently authored oracle. Root consumers must check Err/check before
using read values; the API does not enforce that through typestate. No root settings
are selected in this slice. Native UI/frontend builds were not rerun because host
startup and frontend code are unchanged; Go vet and Rust Cargo integration pass.
