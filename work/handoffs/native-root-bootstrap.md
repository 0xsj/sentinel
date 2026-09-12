# native-root-bootstrap handoff

Coordinator integration completed 2026-09-12. The Wails root now composes
`clock.System`, `id.V7`, `provenance.Factory` and the console `logger`.

`root.Run` opens an `app.startup` scope with a generated UUIDv7 used as scope,
work and local correlation identity. The bound logger emits lifecycle records
with operation, executor and start time. Host errors remain the exact returned
error; only their safe logger projection is emitted. Bootstrap generation or
configuration failures stop `root.Run` before the host starts, while sink and
flush failures are reported without replacing host results.

Verification: `go test -race -count=1 ./pkg/... ./root` passed; `go vet ./...`
passed; `go build -o /tmp/atelier-wails-bootstrap-check .` passed; `gofmt -l
pkg root` was empty. Root tests cover generated scope fields, ID sequence
failure, host result preservation and sink failure behavior.

No frontend, domain, persistence or transport code changed. The generated scope
is process bootstrap context, not a claim about a user, workspace or durable
operation. Later domain operations create their own scopes.
