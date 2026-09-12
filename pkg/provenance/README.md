# Execution provenance

Immutable logical work and individual execution scopes. Adapted locally from
n2f with replay, input links and incoming transport hints deferred.
See [contract](CONTRACT.md) for exact transitions, effects and limits.

```text
Open import request → Scope A / Work A
  Prepare export → Work J (not yet executing)
  Execute Work J → Scope B / attempt 1
  Retry Scope B → Scope C / Work J / attempt 2
```

A named executor is required; the initiator and tenant may be absent. Callers
supply existing clocks and ID generators. Keep provenance alongside operations
and records that need it; no runtime singleton or automatic logging is installed.
Domain-owned artifact lineage and durable audit storage remain separate concerns.
