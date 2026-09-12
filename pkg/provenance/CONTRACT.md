# Execution provenance

Revision: 1 · State: specified · Task: native-provenance

Local adaptation of n2f's core; independent implementation in each Atelier build.
Shared dependencies: errors revision 2 and IDs revision 1. Factories consume
caller-supplied wall-time and ID generation capabilities; existing clocks satisfy
the time seam. No native host, domain, logger, context storage or transport imports.

## Values

Actor is explicitly anonymous or named user/service/system. Named identity and
optional tenant match `[A-Za-z0-9._:/-]{1,128}`. They are opaque labels, not display
names, authentication credentials or resource ownership. Operation matches
`[a-z][a-z0-9_.:-]{0,127}`. Reference is a nonzero UUID tagged scope/work/event.

Attribution owns optional initiator, on-behalf-of and tenant. Absent initiator
means unknown; explicit anonymous means known anonymous. On-behalf-of requires a
named initiator and a different named represented principal. It does not denote
the target of the operation. Tenant may remain absent in a local workbench.

WorkContext owns work ID, correlation ID, correlation source (local/external),
optional primary causation, optional origin (request/schedule/backfill/startup),
operation, attribution and optional uint32 depth. External marks a correlation
claim in restored data; this module does not admit public frontend hints.

Scope adds scope ID, started-at wall time, named executor, positive uint32 attempt,
and optional previous-attempt scope ID. Scope and work expose owned snapshots.
An execution start is normalized to integer UTC milliseconds in inclusive
0..281474976710655; a clock correction can make a child start before its parent.
ID timestamps are not used to derive execution time.

## Transitions and observable scenarios

| Scenario | Contract |
| --- | --- |
| P01–P02 | Validate actor/operation/reference/attribution without echoing rejected input. Preserve unknown versus anonymous and own all input/output data. |
| P03 | Open creates a scope, local correlation equal to its scope ID, depth 0, no cause, attempt 1; optional explicit work ID otherwise defaults to scope ID. Operation, origin and executor are explicit. |
| P04–P05 | Child creates new logical work and execution, preserves correlation/source/origin/attribution, increments known depth, keeps unknown depth unknown. Default cause is parent scope; caller may supply another typed cause. Executor and operation are explicit. Explicit child work ID cannot equal parent work ID. |
| P06 | Prepare is pure: requires a new explicit work ID, applies Child's work rules, consumes no clock/IDs, and records no start, executor or attempt. |
| P07 | Execute starts a prepared/restored work context with a new scope ID, current start and caller-owned positive attempt; absent prior-attempt evidence stays absent. |
| P08–P09 | Retry keeps the entire WorkContext, starts a new scope, increments attempt, records previous scope and takes an explicit executor. Retrying the same scope twice can yield two attempt-2 scopes. No automatic execution, delivery counter, retry safety or idempotency is provided. |
| P10 | Depth/attempt overflow returns invalid / provenance.depth_exhausted or provenance.attempt_exhausted before effects; no saturation or arbitrary depth budget. |
| P11 | Invalid time is rejected before generation. Read exactly one clock sample per otherwise valid execution; truncate submilliseconds, retain wall-clock rollback. |
| P12 | Invalid caller values fail before effects where independently decidable. ID generation failures preserve their classification, metadata and original diagnostic cause. Go preserves the same error instance; Rust moves the Failure unchanged. |
| P18–P19 | Restore validates combinations and copies/owns snapshots; it never repairs invalid input into a fresh root. Reject direct work self-cause, scope self-cause, attempt zero, prior ID equal to scope ID, or prior ID on attempt 1. Known positive depth requires a cause; local depth zero forbids a cause. Unknown depth/origin remains unknown. |
| P23 | Detect locally decidable generated-ID collisions: child scope cannot equal parent scope (even with an explicit different cause/work ID); generated child work cannot equal parent work; retry scope cannot equal previous scope; generated scope/work cannot be its own typed cause. Go also rejects a generated zero ID. Return internal / provenance.invalid_generated_id without retrying. |

Scenario numbering retains n2f core references; omitted numbers are deliberately
out of scope. Each successful execution consumes one generated ID. Validation
that depends on that ID occurs after generation; consumed IDs/time cannot be
rolled back. Factories do not track every prior ID: freshness across unrelated
calls requires a correct generator and persistence uniqueness constraints.

Invalid values use Invalid with diagnostic type `provenance.invalid_actor`,
`invalid_attribution`, `invalid_operation`, `invalid_reference`, `invalid_origin`,
`invalid_work`, `invalid_scope`, `invalid_attempt`, `invalid_depth`, or `invalid_time`
(each with the `provenance.` prefix). Messages are fixed, input-free text.
Go additionally rejects literal nil factory dependencies with
`provenance.invalid_configuration`; typed nil, reentrant callbacks and dependency
panics are caller programming errors. Zero Go scope/work/operation/actor/reference
values are rejected at consuming boundaries; zero Attribution is valid unknown.
Rust constructors prevent invalid zero values and use Option for absence.

## Ownership and boundaries

Factory owns no global state. Go callers serialize factory calls and shared
injected dependencies; do not assume synchronization from the ID generator alone.
Rust requires `&mut Factory`; callers synchronize shared dependencies when needed.
Pure snapshots may be shared safely; Go snapshots copy pointer-bearing fields,
Rust snapshots own cloned values. Consumers own their execution lifetime.

Reference existence, authenticity, graph cycles, authorization, committed outcomes,
artifact lineage and durable history are not proven by construction or restoration.
Stored input needs a versioned codec and trust policy at its owning adapter;
Restore is not a wire decoder. This slice has no replay, multi-input link sets,
frontend hint admission, tracing, ambient context, persistence, runtime root wiring,
logger projections or agent identity subsystem.
