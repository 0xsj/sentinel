# Jobs contract

Revision: 1 · State: specified · Task: jobs-contract

Jobs tracks a long-running local operation whose progress and terminal outcome
should be visible to a workbench user. It records operation state; it does not
spawn processes, schedule threads or decide whether an underlying write
committed. It uses opaque workspace IDs and provenance IDs without importing
those contexts.

## Vocabulary

`Job` has a nonzero job ID, operation name, optional workspace ID, required work
ID, status, revision, timestamps, optional progress, current attempt/scope and an
optional safe failure summary. Operation names use the provenance syntax. A job's
work ID remains stable across attempts; each execution has a new provenance scope.

Statuses are `pending`, `running`, `succeeded`, `failed` and `canceled`.
Cancellation request is a separate boolean on pending/running state. Progress is
optional integer percent 0–100 and may only increase; successful completion sets
it to 100. Failure summaries contain only a shared error category and optional
stable type, never raw causes, stacks or arbitrary output.

## State transitions

```text
pending ── Start ──> running ── Complete ──> succeeded
   │                    │  └──── Fail ─────> failed
   │                    └────── Cancel ────> canceled
   └─ RequestCancel ────────────────────────> canceled
failed/canceled ── Requeue ──> pending
```

`RequestCancel` on running sets the request flag; the worker must acknowledge
with `Cancel`. Completion or failure remains valid after a request because the
operation may have committed before cancellation was observed. A pending job can
be canceled before start. Terminal jobs reject further transitions except an
explicit Requeue. Requeue clears terminal outcome/progress and requires a new
execution attempt; it never starts work automatically.

## Operations and outcomes

| Operation | Success | Expected refusal |
| --- | --- | --- |
| Schedule(job, work, operation, workspace?, at) | Pending job at revision 1 and `job.scheduled` | Invalid IDs/operation; duplicate job ID Conflict |
| Start(job, scope, attempt, at) | Running job and `job.started` | Missing/terminal job; stale revision; scope/work mismatch; attempt not previous+1 |
| ReportProgress(job, percent, expected) | Updated progress and `job.progress`, or Unchanged | Non-running job; regression; percent >100; stale revision |
| RequestCancel(job, expected, at) | Flagged running job or canceled pending job and event; terminal Unchanged | Stale revision or missing job |
| Complete(job, expected, at) | Succeeded and `job.completed` | Non-running job or stale revision |
| Fail(job, summary, expected, at) | Failed and `job.failed` | Non-running job, invalid summary or stale revision |
| Cancel(job, expected, at) | Canceled running job and `job.canceled` | Cancellation was not requested; non-running/stale job |
| Requeue(job, expected, at) | Pending job, cleared outcome and `job.requeued` | Non-terminal job or stale revision |

Every accepted mutation advances a positive uint64 revision. Events are values
returned from the transition and published by the application after the store
transaction commits. Reads and idempotent Unchanged outcomes emit nothing.

## Ports and adapters

The application owns a conditional `JobStore` and a worker capability. The worker
creates provenance scopes through the existing Factory and passes the validated
scope/work/attempt facts to Start. The domain does not call the worker, logger or
event dispatcher. The first memory adapter must make conditional transitions
atomic, return deterministic lists and copy snapshots. SQLite persistence and a
durable event/outbox record are later decisions, made only when restart recovery
or cross-process subscribers are required.

## Explicitly outside revision 1

No scheduler, priority, dependency graph, queue lease, automatic retry/backoff,
process output capture, restart recovery policy, retention period, remote worker,
distributed lock, event bus or exactly-once delivery is selected. A cancellation
request never proves that an operation failed or that its side effects rolled
back.
