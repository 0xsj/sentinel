# native-leaves-integrate: coordinator gate

State: waiting · Depends on native-errors, native-clock, native-secret
This task owns integration, not redefinition of the three contracts. Read all
three handoffs and inspect actual changed files before integrating.

Check scenario coverage, no new dependencies, no framework imports, diagnostic
versus public data, fixed-clock semantics and explicit secret exposure. Investigate
any undeclared edits. Do not treat test counts as a substitute for scenario coverage.

Go discovers the independent packages without root registration. Do not add
unused root imports or manufacture consumers to mark integration complete.
Root/host behavior remains unchanged.

Run from the repository root:

```sh
go test -count=1 ./pkg/errors ./pkg/clock ./pkg/secret
go test -race -count=1 ./pkg/errors ./pkg/clock ./pkg/secret
go vet ./...
npm --prefix frontend run build
go build -o /tmp/atelier-wails-worker-check .
```

Verify every leaf's tests actually ran (a zero-test filtered pass is insufficient).
Save evidence in `work/handoffs/native-leaves-integrate.md`; coordinator may then
mark reviewed leaf tasks and this task complete. Dependencies may be integrated
from implemented handoffs together for this gate; ordinary later tasks require
complete prerequisites. Native window smoke checking is required if startup or
host behavior changes beyond module registration.
