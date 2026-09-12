# Preferences contract

Revision: 1 · State: specified · Task: preferences-contract

Preferences stores user-chosen overrides for one Atelier installation. It is a
bounded context, not process configuration, credentials, workspace files or a
remote synchronization service. The first application slice exercises global
preferences; workspace-scoped values are part of the value contract and use an
opaque workspace ID without importing the Workspace context.

## Vocabulary

`Scope` is either `global` or `workspace(<nonzero ID>)`. A workspace scope does
not prove that the referenced workspace exists; the owning application decides
whether that is a valid target. `Key` is lowercase ASCII matching
`[a-z][a-z0-9_.-]{0,127}`. Keys are case-sensitive and never trimmed.

`Value` is one of UTF-8 text (at most 4096 bytes), boolean or signed 64-bit
integer. Empty text is valid. Arrays, objects, floating point, binary values and
secrets are outside this contract. Secret settings use the existing Secret
capability at a separate boundary and are never stored as ordinary preferences.

`Entry` contains scope, key, value and a positive uint64 revision. A new entry
starts at revision 1. A successful replacement increments the revision by one;
overflow is refused. Stored values and returned snapshots are independently owned.

## Operations and outcomes

The application owns a narrow store port; the domain owns validation and
transitions. Every write uses compare-and-replace semantics:

| Operation | Success | Expected refusal or absence |
| --- | --- | --- |
| Read(scope, key) | Found(entry) | Absent is distinct from an empty value |
| List(scope) | Entries sorted by key | Empty list is a successful empty result |
| Replace(scope, key, value, expected) | Changed(entry, event) or Unchanged(entry) when the value is equal | Invalid input; Conflict when expected revision does not match current presence/revision |
| Remove(scope, key, expected) | Removed(revision, event) | Absent is a successful no-op; Conflict for a stale expected revision |

`expected` is either absent-for-create or an exact existing revision. There is no
blind overwrite operation. A failed compare-and-replace does not mutate state or
advance a revision. A repeated equal replacement does not emit an event.

`Changed` and `Removed` events are values returned by the successful transition.
The application publishes them after the store transaction commits. The domain
does not call a dispatcher. Event names are `preference.changed` and
`preference.removed`; each carries scope, key and revision, and `changed` carries
the new non-secret value.

## Invariants and failures

- Scope, key and value are validated before a store effect.
- A scope cannot contain duplicate keys.
- A key's value type is determined by the value supplied; schema/default
  definitions belong to the consuming feature and are not persisted here.
- `invalid_scope`, `invalid_key`, `invalid_value`, `revision_exhausted` and
  `conflict` are stable domain categories. Messages do not echo raw values.
- A store failure remains an infrastructure failure and is not converted into
  absence or success.

## Ports and adapters

The application owns `PreferenceStore` operations for read, list, conditional
replace and conditional remove. The memory adapter is the first implementation
and must synchronize a whole conditional operation, copy inputs/outputs and
provide deterministic key ordering. A later SQLite adapter owns its schema,
migrations and transaction behavior. Persistence is not implied by this domain
contract; an explicitly selected ephemeral mode remains valid.

The desktop transport exposes only validated plain values and public failure
categories. It never serializes arbitrary maps or diagnostic causes. Defaults are
merged by the application/query boundary and are not silently written on read.

## Explicitly outside revision 1

No environment lookup, secret vault, profile/account system, cloud sync, schema
registry, arbitrary JSON, preference history, undo log, event bus or automatic
migration policy is selected. Whether workspace-scoped values are enabled in the
first UI is an application decision; their identity and validation are already
defined here.
