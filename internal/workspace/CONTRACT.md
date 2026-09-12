# Workspace contract

Revision: 1 · State: specified · Task: workspace-contract

Workspace is the registry of user-owned project locations used by the desktop
workbench. It owns workspace metadata and lifecycle state. Files, documents,
membership, permissions and project-specific schemas belong elsewhere.

## Vocabulary

`Workspace` has a nonzero ID, display name, canonical absolute location, status,
revision, created-at and updated-at timestamps. The application supplies the ID
and UTC millisecond timestamps through the existing capabilities. The domain does
not read the filesystem or a clock.

Names are 1–120 Unicode scalar values, contain no control characters and have no
leading or trailing whitespace. Locations are nonempty, absolute, UTF-8 paths
with no NUL. The platform path adapter resolves separators, symlinks and case
policy before a location enters the domain; the domain stores the resulting
canonical string without reinterpreting platform syntax.

Status is `active` or `archived`. Archiving hides a workspace from the default
active list but never deletes or moves its files. Restoring is reversible.
Forgetting permanently removes registry metadata only and is allowed only for an
archived workspace. “Open” is a runtime application action and does not mutate
this lifecycle state.

## Operations and outcomes

| Operation | Success | Expected refusal or absence |
| --- | --- | --- |
| Register(id, name, location, at) | Active workspace at revision 1 and `workspace.registered` | Invalid value; Conflict when an active record already uses the location or ID |
| Read(id) | Found(snapshot) | Absent is distinct from an unavailable filesystem location |
| List(filter) | Deterministically ordered snapshots | Empty list is successful |
| Rename(id, name, expected, at) | New revision and `workspace.renamed`, or Unchanged for the same name | Absent; stale revision Conflict; invalid name |
| Archive(id, expected, at) | Archived revision and `workspace.archived`, or Unchanged | Absent; stale revision Conflict; already forgotten |
| Restore(id, expected, at) | Active revision and `workspace.restored`, or Unchanged | Absent; stale revision Conflict; active location collision |
| Forget(id, expected) | Metadata removed and `workspace.forgotten` | Absent is a successful no-op; active workspace is refused |

Each successful mutation advances the positive uint64 revision, except an
idempotent Unchanged result. A failed validation, uniqueness check or revision
check does not mutate state. The store enforces active-location uniqueness as
one conditional operation; archived records may retain a location, but Restore
must refuse if another active record now owns it.

Events are returned from transitions and published by the application after the
metadata transaction commits. `workspace.opened` and `workspace.closed` are
transient runtime notifications, not durable domain events in this revision.

## Boundaries and ports

The application owns a `WorkspaceStore` for conditional metadata writes and
queries, plus a path capability that answers whether a canonical location exists
and is a directory when opening. A missing or unreadable path is an application
outcome; it does not corrupt the registry record. The domain never imports host
path APIs, Preferences or Jobs.

The memory adapter is first and must isolate instances, copy snapshots, preserve
ordering and synchronize uniqueness/revision checks. A later SQLite adapter stores
metadata and migration history; it does not store workspace files. Root chooses
the adapter and its lifetime. Persistence failure never becomes an empty registry.

## Explicitly outside revision 1

No nested or multi-root workspaces, project manifests, file watchers, indexing,
recent-open ordering, cloud sync, membership/identity, permissions, deletion or
relocation of files is selected. Whether a workspace must contain an Atelier
manifest is deferred until the first concrete workspace workflow.
