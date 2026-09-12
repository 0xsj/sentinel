# native-errors handoff · 2026-09-11

Contract: `pkg/errors/CONTRACT.md`, revision 1, unchanged.
State: implemented; individual native integration verified. The combined
errors/clock/secret integration task remains waiting for the other two leaves.
No frontend files or dependency/lockfiles changed in this task.

## Files

`pkg/errors/errors.go`, `pkg/errors/errors_test.go`, module README and this handoff.

Coordinator metadata: manifest task owner/state and the foundation status note.

## Scenario evidence

| Scenario | Test or review |
| --- | --- |
| E01 | TestCanonicalKinds |
| E02 | TestRejectNoncanonicalNames |
| E03 | TestConstructionAndDiagnostics |
| E04–E05 | TestClassifiedAndForeignCauses |
| E06 | TestOccurrencesAreDistinct |
| E07 | Source/API review: private fields, declared exports only, stdlib-only imports |
| E08 | TestEnumerationOwnership |
| E09 | TestZeroAndInvalidNumericKinds; TestCauseReplacementAndNativeInspection |
| E10 | TestCauseReplacementAndNativeInspection; TestTypedNilCauseIsNotInspected |

## Checks

- `go test -count=1 -v ./pkg/errors`: 9 tests pass.
- `go test -race -count=1 ./pkg/errors`: passes.
- `go vet ./pkg/errors`: passes.
- `gofmt -l pkg/errors`: empty output.
- Coordinator: `go vet ./...` and `go build -o /tmp/atelier-wails-errors-check .`: pass.

## Review and limitations

Reviewed actual production code against the declared API and contract: fixed
classification vocabulary, exact diagnostic preservation, source replacement,
private fields, standard-library-only dependencies, no retry/serialization/public
projection, no hidden IO/time/identity generation, no peer or host imports.
The tests were authored with implementation visibility by the implementing
coordinator; this is not a blind or independently authored oracle.

No contract changes required. Causes remain native error objects and diagnostics
remain private-facing data. Go retains foreign error references; Rust owns a boxed
Send + Sync source. Native bridge mapping is still a separate future contract.
Host behavior is unchanged, so no frontend build or native UI smoke was rerun.
