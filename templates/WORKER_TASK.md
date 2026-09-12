# <task ID>: <one deliverable>

State: draft | ready | waiting | implemented | complete
Contract: <path> revision <N> · Dependencies: <manifest IDs>

## Read

AGENTS.md, WORKERS.md, the linked contract and only its named prerequisites.
Do not read every sibling implementation or copy sibling imports.

## Scope

Allowed files: <exact paths/globs>. Forbidden edits: contract, root wiring,
dependency/lockfiles, shared integration files and unrelated modules unless named.
Deliver: <one bounded behavior>. Excludes: <adjacent work>.

## Acceptance and checks

Implement every named contract scenario. Use tests exercising observable behavior,
including negative cases. Run <exact commands>. Record actual test counts and
scenario mappings. Compilation alone is not completion.

## Escalation to coordinator

Report a contract contradiction, missing dependency or necessary out-of-scope
change with a concrete proposed revision. Continue independent in-scope work.
Do not change contracts or broaden ownership merely to make checks pass.
This is worker coordination, not a new user approval requirement.

## Handoff

- Contract revision and files changed.
- Scenario ID → test mapping.
- Commands, outcomes and test counts, including any failed checks.
- Remaining limitations and proposed contract changes.
- Integration files the coordinator must connect.

Save at work/handoffs/<task ID>.md. Claim implemented when local work/checks pass;
only the coordinator marks complete after review and integration evidence.
