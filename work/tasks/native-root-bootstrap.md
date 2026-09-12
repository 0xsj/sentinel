# native-root-bootstrap: compose the foundation at startup

This coordination task wires the existing clock, ID, provenance and logger
packages into the Wails composition root. It does not add a domain, transport,
frontend behavior or persistence.

The root creates one UUIDv7 generator, opens one `app.startup` provenance scope,
binds the scope to the console logger and uses that bound logger for startup,
failure and shutdown records. ID or provenance bootstrap errors are returned from
`root.Run`; a failed logger delivery never changes the host result.

Checks:

```sh
go test -race -count=1 ./pkg/... ./root
go vet ./...
go build -o /tmp/atelier-wails-bootstrap-check .
test -z "$(gofmt -l pkg root)"
```
