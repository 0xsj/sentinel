# native-clock: implement the clock leaf

State: ready · Contract: [pkg/clock/CONTRACT.md](../../pkg/clock/CONTRACT.md) revision 1
Dependencies: none · Manifest: ../manifest.json

## Read and scope

Read AGENTS.md, WORKERS.md and the linked contract. Implement the exact public
surface and all acceptance scenarios. Allowed edits: `pkg/clock/` production files,
tests and README.md, plus `work/handoffs/native-clock.md`. The CONTRACT.md is read-only
for workers. Do not edit dependencies/lockfiles, other leaves, root/host wiring,
manifest state or shared module registration. No new external packages.

Private representation and test organization are your decisions within the
contract. A public API contradiction goes to the coordinator with a proposed
revision; do not weaken the contract or implement guessed adjacent behavior.

## Checks from repository root

```sh
go test -count=1 ./pkg/clock
go test -race -count=1 ./pkg/clock
go vet ./pkg/clock
gofmt -l pkg/clock
```

Formatting output must be empty. Record nonzero test counts and scenario-to-test
mapping. Where a scenario requires source/API review, report that separately.
Tests authored with implementation access are ordinary regression evidence.

## Handoff

Write `work/handoffs/native-clock.md`: contract revision, files changed, each scenario's
test/review evidence, exact commands/results/counts, limitations and any proposed
contract changes. Report implemented, not integrated/complete. The coordinator
owns shared wiring and final source/dependency review.
