# Synchronous diagnostic logger

Revision: 1 · State: specified · Task: native-logger

Console is the root default at info level. Go uses slog.TextHandler, adapting a
swappable slog.Handler via NewSlog; native slog types stay in the sink adapter API.
Rust uses a Write-based Console behind a Send+Sync Sink trait. Callers depend on
Logger or consumer-owned logging capabilities, not the console implementation.
No global/default logger is installed. Each project is self-contained.

## Record and binding contract

L01: Four levels debug/info/warn/error with exact lowercase parsing and inclusive
floor. Invalid configuration refuses (invalid / logger.invalid_configuration).
Go also checks literal nil dependencies and invalid numeric levels; typed nil,
zero unconstructed Logger, panicking or reentrant adapters are programmer misuse.
Rust's enum and constructors prevent these invalid values.

L02: Filtered and Discard log calls perform no clock read, record creation or sink
write inside Logger. Argument construction happens before the call. Enabled is a
pure adapter policy; both the configured floor and sink must permit an event.

L03: Values are text, signed 64-bit integer, boolean or explicit Secret projection
(always [REDACTED]). There is no Any/object formatter or recursive arbitrary data
encoder. Caller-authored text/messages must be chosen as safe diagnostics; raw
strings do not magically become secrets. Go Value is private with Scalar access;
Rust Value is an owned enum. With returns an independent logger binding. New field
values replace matching old keys; unrelated keys and parent bindings survive.

L04: Records contain emission time, severity, message and service; ordinary fields,
scope, error and public error_fields occupy separate groups. Application fields
cannot overwrite envelope identity. Sink records and Memory snapshots are owned.
Field keys are literal within a group. Console dotted rendering is human-readable,
not a reversible wire schema. Custom sinks receive the structured groups intact.

L05: WithError snapshots classified status, kind, safe public message/type/fields,
and cause presence. No raw diagnostic message, diagnostic-only type, details,
stack or source text is printed. Internal and unknown use the shared fixed public
projection. Unknown remains an actual error; nil/None clears an inherited binding.
Rust ErrorInfo::classified uses explicit classification, unknown accepts foreign
errors; neither searches a registry or formats the raw error. Error bindings do
not themselves emit a record or change its level.

L06: Explicit WithScope preserves scope/work/correlation IDs and source, operation,
start, attempt, executor, and known origin/depth/cause/previous-attempt/attribution.
Unknown optional values stay absent. No new identity, time or permissions arise
from projection. Go validates its potentially zero Scope at this boundary; Rust
Scope is already valid. Snapshot changes cannot affect the projected binding.

L07: Emission samples the injected wall clock exactly once, truncates to integer
UTC milliseconds within inclusive 0..281474976710655, and rejects invalid time
with invalid / logger.invalid_time before writing. Scope.startedAt is independent
of log emission time. A successful Log means the sink accepted the record, not a
durable audit write. Console produces one physical line, escaping newlines and
terminal controls; no color/ANSI is generated. Go native slog text and Rust console
have equivalent record semantics, not identical display bytes. Go handlers may
be replaced (including JSONHandler); custom adapter formatting is its owner's policy.

L08: Log and Flush are synchronous and return delivery failures as unavailable /
logger.sink, preserving the original cause. No recursion, panic-on-error, retry,
queue, counter, timeout, dropping-on-overload, process exit or fallback sink is
hidden in Logger. A partial write may have happened; retries can duplicate output.
A blocked writer blocks the caller. Sinks must obey their ordinary write contracts.

L09: Related/bound loggers serialize clock and write/flush calls. Sink implementations
support concurrent calls even across separately constructed loggers. Callers own
safe concurrent access to fields while constructing bindings/events and own shared
clocks across unrelated logger instances. Flush waits behind ongoing writes and
asks the adapter to flush. It does not close the sink or stop subsequent logging.
Root owns sink lifetime. Go Console borrows its writer and invokes Flush if supplied;
Rust Console owns its writer and calls Write::flush. Neither promises fsync/durability.
Externally shared Go flush callbacks/writers must meet Sink concurrency requirements.

L10: Memory captures and returns owned records; Discard disables logging effects.
Root selects stderr console and records app.starting and app.stopped; Go also
projects a returned native-host error as app.failed and preserves that exact result.
Tauri's existing host entry returns unit; panics/process exits are not converted
into claimed completion records. Diagnostic delivery failures do not prevent the
host call or change its result: root makes one best-effort fixed stderr report.

Native scopes are not fabricated solely for bootstrap logs. Later application
operations supply meaningful provenance. No file rotation, tracing SDK, remote
export, runtime sink reload, persistent history or platform bridge schema is
selected. Console defaults and shutdown are composition-root choices.
