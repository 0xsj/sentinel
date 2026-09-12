# Clock: first contract

Revision: 1 · State: specified · Task: native-clock

## Purpose and boundary

Expose controllable wall time to future application/adapters. Scheduled consumers
are timestamped diagnostics and later tracked work; no current domain must import
this module. Standard library only. Pure transitions receive a timestamp value;
they do not look up a clock internally. This leaf performs only the explicit time
read, never sleeps, starts threads or schedules work.

### Go API

```go
// package clock
type Clock interface { Now() time.Time }
type System struct{}
func (System) Now() time.Time

type Fixed struct { /* private timestamp */ }
func NewFixed(at time.Time) Fixed
func (Fixed) Now() time.Time
```

System returns time.Now().UTC(); Fixed normalizes the supplied timestamp to UTC.
Preserve the instant and nanosecond precision. The zero Fixed is Go's zero time in
UTC. Returned time.Time is a value; there is no Set/Advance. Verify C01 with a
non-UTC input using Equal plus a UTC location check. Do not promise preserved
monotonic readings after UTC normalization. Verify zero-value behavior separately.


## Acceptance scenarios

| ID | Given / when | Required observable outcome |
| --- | --- | --- |
| C01 | Fixed timestamp queried repeatedly | Exactly the supplied instant on every read; no drift or advancement |
| C02 | Two differently configured fixed clocks | Independent results; no global clock override |
| C03 | Timestamp before Unix epoch | Accepted and preserved, without unsigned epoch assumptions |
| C04 | System clock called | Returns a native timestamp without application initialization; structural smoke check, no equality-to-now assertion |
| C05 | Read fixed clock from concurrent callers | Same configured instant; no data race/mutation |
| C06 | Inspect API and production imports | Only the declared clock surface; no sleeps, timer/timeout, monotonic duration or native-host APIs |

No ordering between two wall-clock reads is promised. Monotonic measurement,
manual advance, schedulers, calendars/time zones, persistence formats and timestamp
serialization are excluded. Never use real sleeping as a test strategy. Production
consumers own the clock instance; this package does not register itself in root.
