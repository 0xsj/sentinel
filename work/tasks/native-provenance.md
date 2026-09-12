# native-provenance: work and execution attribution

Contract: [pkg/provenance/CONTRACT.md](../../pkg/provenance/CONTRACT.md), revision 1.
Prerequisites: native-errors revision 2, native-id revision 1, native-clock revision 1.
State authority: work/manifest.json. Own pkg/provenance/ source/tests/README and handoff.
Coordinator owns contract, shared registration and manifest. Implement the listed
P scenarios only; do not import sibling projects, add replay/transport/persistence,
change frontend code or wire dummy runtime consumers. Map scenarios to focused
tests and record full native integration checks before completion.

```sh
go test -race -count=1 -v ./pkg/...
go vet ./...
gofmt -l pkg/provenance
```
