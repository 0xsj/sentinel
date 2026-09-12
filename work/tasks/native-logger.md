# native-logger: console with swappable sinks

Contract: [pkg/logger/CONTRACT.md](../../pkg/logger/CONTRACT.md), revision 1.
Prerequisites: errors revision 2; clock, secret and provenance revision 1.
State authority: work/manifest.json. Own pkg/logger/ and the handoff.
Coordinator owns contract, registration, root logging, manifests and build guides.
Implement L01–L10. Do not copy n2f queues, tracing, colors or remote output.
No frontend changes. Run focused tests then the full native integration gate:

```sh
go test -race -count=1 ./pkg/... ./root
go vet ./...
gofmt -l pkg/logger root
```
