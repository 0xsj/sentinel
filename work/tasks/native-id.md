# native-id: UUID values and generation

Contract: [pkg/id/CONTRACT.md](../../pkg/id/CONTRACT.md), revision 1.
Dependency: native-errors revision 2. State authority: work/manifest.json.

Implement I01–I09. Own pkg/id/ source, tests, README and
work/handoffs/native-id.md. Coordinator owns contract, module registration,
dependency/lockfiles and manifest updates. No frontend, provenance, domain IDs
or root runtime wiring in this slice. Use private test clocks for rollback.

Checks:

```sh
go test -count=1 -v ./pkg/id
go test -race -count=1 ./pkg/id
go vet ./...
gofmt -l pkg/id
```

Coordinator reviews scenario coverage and runs the full native shared suite
before marking complete. Record command results and actual test counts.
