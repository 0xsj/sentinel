# native-errors: implement the errors leaf

State: ready · Contract: [pkg/errors/CONTRACT.md](../../pkg/errors/CONTRACT.md) revision 2
Dependencies: none · Manifest: ../manifest.json

## Read and scope

Read AGENTS.md, WORKERS.md and the linked contract. Implement the exact public
surface and all acceptance scenarios. Allowed edits: `pkg/errors/` production files,
tests and README.md, plus `work/handoffs/native-errors.md`. The CONTRACT.md is read-only
for workers. Do not edit dependencies/lockfiles, other leaves, root/host wiring,
manifest state or shared module registration. No new external packages.

Private representation and test organization are your decisions within the
contract. A public API contradiction goes to the coordinator with a proposed
revision; do not weaken the contract or implement guessed adjacent behavior.

## Checks from repository root

```sh
go test -count=1 ./pkg/errors
go test -race -count=1 ./pkg/errors
go vet ./pkg/errors
gofmt -l pkg/errors
```

Formatting output must be empty. Record nonzero test counts and scenario-to-test
mapping. Where a scenario requires source/API review, report that separately.
Tests authored with implementation access are ordinary regression evidence.

## Handoff

Write `work/handoffs/native-errors.md`: contract revision, files changed, each scenario's
test/review evidence, exact commands/results/counts, limitations and any proposed
contract changes. Report implemented, not integrated/complete. The coordinator
owns shared wiring and final source/dependency review.
