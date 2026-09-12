# native-config: captured process configuration

Contract: [pkg/config/CONTRACT.md](../../pkg/config/CONTRACT.md), revision 1.
Dependencies: native-errors revision 2 and native-secret revision 1, both complete.
State authority: work/manifest.json.

Implement the exact API and V01–V10 behavior. Read AGENTS.md, WORKERS.md and the
contract first. Own pkg/config/ production files, tests, README and
work/handoffs/native-config.md. Contract edits and Rust shared/mod.rs registration
belong to the coordinator. No frontend, application settings, dependency/lockfile,
root boot wiring, dotenv, remote providers or preference behavior in this slice.

Checks from repository root:

```sh
go test -count=1 -v ./pkg/config
go test -race -count=1 ./pkg/config
go vet ./...
gofmt -l pkg/config
```

Map every scenario to tests or explicit source review; record actual test counts
and command results. Keep test secrets synthetic. Hand off limitations and public
API questions; do not modify contracts to hide failures. Coordinator integrates,
checks all native shared tests, and records evidence before marking complete.
