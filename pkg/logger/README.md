# Logger

Implemented synchronous console logging with swappable sinks. Native root uses
stderr at info level and records startup/shutdown. There is no process-global
logger. See [contract](CONTRACT.md) for record, ownership and delivery guarantees.

```go
sink, err := logger.NewConsole(os.Stderr)
if err != nil { return err }
log, err := logger.New(clock.System{}, sink, "atelier", logger.Info)
if err != nil { return err }
err = log.With(logger.Fields{"component": logger.Text("imports")}).Log(
    logger.Info, "import.started", logger.Fields{"files": logger.Int(3)},
)
// Handle err as a diagnostic delivery outcome, separately from import results.
```

To change formatting, supply `logger.NewSlog(slog.NewJSONHandler(writer, nil), nil)`
as the sink. To capture records, inject `&logger.Memory{}` and inspect `Records()`.
`logger.Discard{}` disables logging. NewSlog's optional flush callback belongs to
the adapter; a handler with a stricter level can further filter records.

Bind provenance explicitly with WithScope/with_scope. Bind errors through the
safe projection API; raw diagnostic causes and internal details are not emitted.
Explicit secret values redact; arbitrary caller strings are not auto-sanitized.
Ordinary fields stay separate from envelope, scope and error fields.

Calls and Flush/flush are synchronous. Delivery errors preserve their cause; a
blocked writer blocks the caller. Root reports failures without replacing the
host result. Flush does not close the sink or prevent later logging. Queueing,
rotation, remote export and durable audit history remain future adapters.
